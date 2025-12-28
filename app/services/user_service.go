// Package services 用户业务逻辑服务
package services

import (
	"GoHub-Service/app/models/user"
	apperrors "GoHub-Service/pkg/errors"
	"GoHub-Service/pkg/paginator"

	"github.com/gin-gonic/gin"
)

// UserService 用户服务
type UserService struct{}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{}
}

// GetByID 根据ID获取用户
func (s *UserService) GetByID(id string) (*user.User, error) {
	userModel := user.Get(id)
	if userModel.ID == 0 {
		return nil, apperrors.NotFoundError("用户").WithDetails(map[string]interface{}{
			"user_id": id,
		})
	}
	return &userModel, nil
}

// List 获取用户列表
func (s *UserService) List(c *gin.Context, perPage int) ([]user.User, *paginator.Paging, error) {
	data, pager := user.Paginate(c, perPage)
	return data, &pager, nil
}
