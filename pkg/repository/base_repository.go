// Package repository 数据访问层基础接口
package repository

import (
	"GoHub-Service/pkg/paginator"

	"github.com/gin-gonic/gin"
)

// BaseRepository 基础仓储接口
type BaseRepository interface {
	// GetByID 根据ID获取单个实体
	GetByID(id string) (interface{}, error)
	
	// Create 创建实体
	Create(entity interface{}) error
	
	// Update 更新实体
	Update(entity interface{}) error
	
	// Delete 删除实体
	Delete(id string) error
	
	// List 获取分页列表
	List(c *gin.Context, perPage int) (interface{}, *paginator.Paging, error)
}

// CacheableRepository 支持缓存的仓储接口
type CacheableRepository interface {
	BaseRepository
	
	// GetFromCache 从缓存获取
	GetFromCache(key string) (interface{}, error)
	
	// SetCache 设置缓存
	SetCache(key string, value interface{}, ttl int) error
	
	// DeleteCache 删除缓存
	DeleteCache(key string) error
	
	// FlushCache 清空相关缓存
	FlushCache() error
}
