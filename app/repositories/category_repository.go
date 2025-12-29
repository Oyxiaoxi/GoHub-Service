// Package repositories Category仓储实现
package repositories

import (
	"GoHub-Service/app/models/category"
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

// CategoryRepository 聚合了分类的 CRUD、批处理和简单缓存操作的接口定义.
type CategoryRepository interface {
	GetByID(id string) (*category.Category, error)
	List(c *gin.Context, perPage int) ([]category.Category, *paginator.Paging, error)
	Create(category *category.Category) error
	Update(category *category.Category) error
	Delete(id string) error
	BatchCreate(categories []category.Category) error
	BatchDelete(ids []string) error
	// 缓存方法
	GetAllCached() ([]category.Category, error)
	SetListCache(categories []category.Category) error
	FlushCache() error
}

// categoryRepository 基于 GORM + Redis 的 Category 仓储实现.
type categoryRepository struct {
	cacheTTL         int
	cacheKeyCategory string
	cacheKeyList     string
}

// NewCategoryRepository 返回默认的分类仓储实现，带基础缓存配置.
func NewCategoryRepository() CategoryRepository {
	return &categoryRepository{
		cacheTTL:         7200,
		cacheKeyCategory: "category:%s",
		cacheKeyList:     "category:list",
	}
}

// BatchCreate 批量创建分类（事务包裹），确保批次成功或全部回滚.
func (r *categoryRepository) BatchCreate(categories []category.Category) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&categories).Error; err != nil {
			return err
		}
		return nil
	})
}

// BatchDelete 批量删除分类（事务包裹），避免部分删除留下脏数据.
func (r *categoryRepository) BatchDelete(ids []string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id IN ?", ids).Delete(&category.Category{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetByID 根据ID获取分类，不存在时返回业务级 NotFound.
func (r *categoryRepository) GetByID(id string) (*category.Category, error) {
	categoryModel := category.Get(id)
	if categoryModel.ID == 0 {
		return nil, apperrors.NotFoundError("分类").WithDetails(map[string]interface{}{
			"category_id": id,
		})
	}
	return &categoryModel, nil
}

// List 获取分类列表，复用模型层分页能力.
func (r *categoryRepository) List(c *gin.Context, perPage int) ([]category.Category, *paginator.Paging, error) {
	data, pager := category.Paginate(c, perPage)
	return data, &pager, nil
}

// Create 创建分类并清理列表缓存，避免返回旧数据.
func (r *categoryRepository) Create(cat *category.Category) error {
	cat.Create()
	if cat.ID == 0 {
		return apperrors.DatabaseError("创建分类", nil)
	}

	// 清除列表缓存
	_ = r.FlushCache()

	return nil
}

// Update 更新分类并清理相关缓存.
func (r *categoryRepository) Update(cat *category.Category) error {
	rowsAffected := cat.Save()
	if rowsAffected == 0 {
		return apperrors.DatabaseError("更新分类", nil)
	}

	// 清除缓存
	_ = r.FlushCache()

	return nil
}

// Delete 删除分类并清理相关缓存.
func (r *categoryRepository) Delete(id string) error {
	categoryModel := category.Get(id)
	if categoryModel.ID == 0 {
		return apperrors.NotFoundError("分类").WithDetails(map[string]interface{}{
			"category_id": id,
		})
	}

	rowsAffected := categoryModel.Delete()
	if rowsAffected == 0 {
		return apperrors.DatabaseError("删除分类", nil)
	}

	// 清除缓存
	_ = r.FlushCache()

	return nil
}

// GetAllCached 尝试读取缓存的分类列表，未命中时返回空切片且不视为错误（由上层决定降级策略）。
func (r *categoryRepository) GetAllCached() ([]category.Category, error) {
	val := redis.Redis.Get(r.cacheKeyList)
	if val != "" {
		var categories []category.Category
		if err := json.Unmarshal([]byte(val), &categories); err == nil {
			return categories, nil
		}
	}

	// 从数据库获取所有分类
	var categories []category.Category
	// 这里简化处理，实际需要实现 category.All() 方法
	return categories, nil
}

// SetListCache 将分类列表写入缓存，调用方决定何时刷新.
func (r *categoryRepository) SetListCache(categories []category.Category) error {
	data, err := json.Marshal(categories)
	if err != nil {
		return err
	}

	redis.Redis.Set(r.cacheKeyList, string(data), time.Duration(r.cacheTTL)*time.Second)
	return nil
}

// FlushCache 清空分类相关缓存键，粗粒度策略即可满足当前读多写少场景.
func (r *categoryRepository) FlushCache() error {
	_ = redis.Redis.Del(r.cacheKeyList)
	pattern := fmt.Sprintf("%s*", "category:")
	_ = pattern // 简化实现
	return nil
}
