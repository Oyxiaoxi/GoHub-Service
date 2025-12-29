// Package repositories User仓储实现
package repositories

import (
	"GoHub-Service/app/cache"
	"GoHub-Service/app/models/user"
	"GoHub-Service/pkg/database"
	apperrors "GoHub-Service/pkg/errors"
	"GoHub-Service/pkg/paginator"
	"GoHub-Service/pkg/redis"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"

	"github.com/gin-gonic/gin"
)

// UserRepository 定义用户的查询、批处理以及缓存读写接口.
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

// userRepository 基于 GORM + Redis 的用户仓储实现.
type userRepository struct {
	cacheTTL     int
	cacheKeyUser string
}

// BatchCreate 批量创建用户，事务保证批次全成或全失败.
func (r *userRepository) BatchCreate(users []user.User) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&users).Error; err != nil {
			return err
		}
		return nil
	})
}

// BatchDelete 批量删除用户，避免局部删除留下孤儿记录.
func (r *userRepository) BatchDelete(ids []string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id IN ?", ids).Delete(&user.User{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// NewUserRepository 返回默认的用户仓储实现，使用分级缓存策略.
func NewUserRepository() UserRepository {
	tier := cache.GetEntityTier("user")
	return &userRepository{
		cacheTTL:     int(tier.TTL.Seconds()),
		cacheKeyUser: "user:%s",
	}
}

// GetByID 优先读取缓存，未命中则回源数据库并写回缓存.
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

// List 获取用户列表，沿用模型层分页实现.
func (r *userRepository) List(c *gin.Context, perPage int) ([]user.User, *paginator.Paging, error) {
	data, pager := user.Paginate(c, perPage)
	return data, &pager, nil
}

// GetFromCache 从缓存获取用户，缓存未命中返回可视为软错误供上层回源.
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

// SetCache 将用户信息写入缓存，调用方负责选择写回时机.
func (r *userRepository) SetCache(u *user.User) error {
	key := fmt.Sprintf(r.cacheKeyUser, fmt.Sprintf("%d", u.ID))
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}

	redis.Redis.Set(key, string(data), time.Duration(r.cacheTTL)*time.Second)
	return nil
}

// DeleteCache 删除指定用户缓存，供更新/删除后调用.
func (r *userRepository) DeleteCache(id string) error {
	key := fmt.Sprintf(r.cacheKeyUser, id)
	redis.Redis.Del(key)
	return nil
}
