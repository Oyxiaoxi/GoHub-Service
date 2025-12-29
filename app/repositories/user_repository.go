// Package repositories User仓储实现
package repositories

import (
	"GoHub-Service/app/models/user"
	apperrors "GoHub-Service/pkg/errors"
	"GoHub-Service/pkg/paginator"
	"GoHub-Service/pkg/redis"
	"encoding/json"
	"fmt"
	"time"
	"GoHub-Service/pkg/database"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// UserRepository User仓储接口
type UserRepository interface {
	GetByID(id string) (*user.User, error)
	List(c *gin.Context, perPage int) ([]user.User, *paginator.Paging, error)
	BatchCreate(users []user.User) error
	BatchDelete(ids []string) error
	// 缓存方法
	GetFromCache(id string) (*user.User, error)
	SetCache(user *user.User) error
	DeleteCache(id string) error
}

// userRepository User仓储实现
type userRepository struct {
	cacheTTL     int
	cacheKeyUser string
}

// BatchCreate 批量创建用户（事务包裹）
func (r *userRepository) BatchCreate(users []user.User) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&users).Error; err != nil {
			return err
		}
		return nil
	})
}

// BatchDelete 批量删除用户（事务包裹）
func (r *userRepository) BatchDelete(ids []string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id IN ?", ids).Delete(&user.User{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// NewUserRepository 创建User仓储实例
func NewUserRepository() UserRepository {
	return &userRepository{
		cacheTTL:     1800,      // 30分钟
		cacheKeyUser: "user:%s", // user:id
	}
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(id string) (*user.User, error) {
	// 尝试从缓存获取
	if cached, err := r.GetFromCache(id); err == nil && cached != nil {
		return cached, nil
	}

	// 从数据库获取
	userModel := user.Get(id)
	if userModel.ID == 0 {
		return nil, apperrors.NotFoundError("用户").WithDetails(map[string]interface{}{
			"user_id": id,
		})
	}

	// 设置缓存
	_ = r.SetCache(&userModel)

	return &userModel, nil
}

// List 获取用户列表
func (r *userRepository) List(c *gin.Context, perPage int) ([]user.User, *paginator.Paging, error) {
	data, pager := user.Paginate(c, perPage)
	return data, &pager, nil
}

// GetFromCache 从缓存获取用户
func (r *userRepository) GetFromCache(id string) (*user.User, error) {
	key := fmt.Sprintf(r.cacheKeyUser, id)
	val := redis.Redis.Get(key)
	if val == "" {
		return nil, fmt.Errorf("cache miss")
	}

	var userModel user.User
	if err := json.Unmarshal([]byte(val), &userModel); err != nil {
		return nil, err
	}

	return &userModel, nil
}

// SetCache 设置用户缓存
func (r *userRepository) SetCache(u *user.User) error {
	key := fmt.Sprintf(r.cacheKeyUser, fmt.Sprintf("%d", u.ID))
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}

	redis.Redis.Set(key, string(data), time.Duration(r.cacheTTL)*time.Second)
	return nil
}

// DeleteCache 删除用户缓存
func (r *userRepository) DeleteCache(id string) error {
	key := fmt.Sprintf(r.cacheKeyUser, id)
	redis.Redis.Del(key)
	return nil
}
