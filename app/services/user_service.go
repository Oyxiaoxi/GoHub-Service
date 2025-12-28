// Package services 用户业务逻辑服务
package services

import (
	"GoHub-Service/app/models/user"
	"GoHub-Service/app/repositories"
	"GoHub-Service/pkg/paginator"

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

// GetByID 根据ID获取用户
func (s *UserService) GetByID(id string) (*user.User, error) {
	return s.repo.GetByID(id)
}

// List 获取用户列表
func (s *UserService) List(c *gin.Context, perPage int) ([]user.User, *paginator.Paging, error) {
	return s.repo.List(c, perPage)
}
