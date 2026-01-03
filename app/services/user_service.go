// Package services 用户业务逻辑服务
package services

import (
	"context"
	"fmt"
	"time"

	"GoHub-Service/app/models/user"
	"GoHub-Service/app/repositories"
	apperrors "GoHub-Service/pkg/errors"
	"GoHub-Service/pkg/mapper"
	"GoHub-Service/pkg/paginator"
	"GoHub-Service/pkg/resource"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UpdateProfile 更新当前用户信息，并在成功时刷新缓存，确保后续读取不命中旧值.
func (s *UserService) UpdateProfile(currentUser *user.User, name, city, introduction string) (*user.User, error) {
	currentUser.Name = name
	currentUser.City = city
	currentUser.Introduction = introduction
	rowsAffected := currentUser.Save()
	if rowsAffected > 0 {
		// 更新缓存
		_ = s.repo.SetCache(currentUser)
		return currentUser, nil
	}
	return nil, fmt.Errorf("用户信息更新失败")
}

// UserService 负责用户的读写操作及缓存一致性处理.
type UserService struct {
	repo   repositories.UserRepository
	mapper mapper.Mapper[user.User, UserResponseDTO] // 使用泛型Mapper消除DTO转换重复
	logger *zap.Logger
}

// NewUserService 构造用户服务，默认使用数据库+缓存仓储实现.
func NewUserService() *UserService {
	// 定义DTO转换函数（只需一次）
	converter := func(u *user.User) *UserResponseDTO {
		return &UserResponseDTO{
			ID:        u.GetStringID(),
			Name:      u.Name,
			Email:     u.Email,
			Phone:     u.Phone,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}
	}

	return &UserService{
		repo:   repositories.NewUserRepository(),
		mapper: mapper.NewSimpleMapper(converter),
		logger: zap.L(),
	}
}

// UserCreateDTO 创建用户请求DTO
type UserCreateDTO struct {
	Name     string `json:"name" binding:"required,min=3,max=255"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone" binding:"required"`
}

// UserUpdateDTO 更新用户请求DTO
type UserUpdateDTO struct {
	Name  *string `json:"name,omitempty" binding:"omitempty,min=3,max=255"`
	Email *string `json:"email,omitempty" binding:"omitempty,email"`
	Phone *string `json:"phone,omitempty"`
}

// UserResponseDTO 用户响应DTO
type UserResponseDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserListResponseDTO 用户列表响应DTO
type UserListResponseDTO struct {
	Users  []UserResponseDTO `json:"users"`
	Paging *paginator.Paging `json:"paging"`
}

// toResponseDTO 使用Mapper将User模型转换为响应DTO
// 优化：使用泛型Mapper消除重复代码
func (s *UserService) toResponseDTO(u *user.User) *UserResponseDTO {
	return s.mapper.ToDTO(u)
}

// toResponseDTOList 使用Mapper将User模型列表转换为响应DTO列表
// 优化：使用泛型Mapper消除重复代码，自动优化内存拷贝
func (s *UserService) toResponseDTOList(users []user.User) []UserResponseDTO {
	return s.mapper.ToDTOList(users)
}

// GetByID 获取单个用户，优先命中缓存，失败时返回包装后的业务错误.
// 优化：使用ContextGuard确保超时控制
func (s *UserService) GetByID(id string) (*UserResponseDTO, *apperrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	guard := resource.NewContextGuard(ctx, cancel, s.logger)
	defer guard.Release()

	u, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperrors.WrapError(err, "获取用户失败")
	}

	guard.Cancel() // 操作完成，提前取消
	return s.toResponseDTO(u), nil
}

// List 分页获取用户列表并附带分页元信息.
func (s *UserService) List(c *gin.Context, perPage int) (*UserListResponseDTO, *apperrors.AppError) {
	users, paging, err := s.repo.List(c, perPage)
	if err != nil {
		return nil, apperrors.WrapError(err, "获取用户列表失败")
	}
	return &UserListResponseDTO{
		Users:  s.toResponseDTOList(users),
		Paging: paging,
	}, nil
}
