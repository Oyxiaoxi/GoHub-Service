// Package services 用户业务逻辑服务
package services

import (
	"GoHub-Service/app/models/user"
	"GoHub-Service/app/repositories"
	"GoHub-Service/pkg/paginator"
	"time"

	"github.com/gin-gonic/gin"
)

// UserService 用户服务
type UserService struct{
	repo repositories.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		repo: repositories.NewUserRepository(),
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

// toResponseDTO 将User模型转换为响应DTO
func (s *UserService) toResponseDTO(u *user.User) *UserResponseDTO {
	return &UserResponseDTO{
		ID:        u.GetStringID(),
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// toResponseDTOList 将User模型列表转换为响应DTO列表
func (s *UserService) toResponseDTOList(users []user.User) []UserResponseDTO {
	dtos := make([]UserResponseDTO, len(users))
	for i, u := range users {
		dtos[i] = UserResponseDTO{
			ID:        u.GetStringID(),
			Name:      u.Name,
			Email:     u.Email,
			Phone:     u.Phone,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}
	}
	return dtos
}

// GetByID 根据ID获取用户
func (s *UserService) GetByID(id string) (*UserResponseDTO, error) {
	u, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.toResponseDTO(u), nil
}

// List 获取用户列表
func (s *UserService) List(c *gin.Context, perPage int) (*UserListResponseDTO, error) {
	users, paging, err := s.repo.List(c, perPage)
	if err != nil {
		return nil, err
	}
	
	return &UserListResponseDTO{
		Users:  s.toResponseDTOList(users),
		Paging: paging,
	}, nil
}
