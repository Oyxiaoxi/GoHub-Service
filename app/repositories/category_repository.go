// Package repositories Category仓储实现
package repositories

import (
	"GoHub-Service/app/models/category"
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

// CategoryRepository Category仓储接口
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

// categoryRepository Category仓储实现
type categoryRepository struct {
	cacheTTL         int
	cacheKeyCategory string
	cacheKeyList     string
}

// NewCategoryRepository 创建Category仓储实例
func NewCategoryRepository() CategoryRepository {
	return &categoryRepository{
		cacheTTL:         7200,
		cacheKeyCategory: "category:%s",
		cacheKeyList:     "category:list",
	}
}

// BatchCreate 批量创建分类（事务包裹）
func (r *categoryRepository) BatchCreate(categories []category.Category) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&categories).Error; err != nil {
			return err
		}
		return nil
	})
}

// BatchDelete 批量删除分类（事务包裹）
func (r *categoryRepository) BatchDelete(ids []string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id IN ?", ids).Delete(&category.Category{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetByID 根据ID获取分类
func (r *categoryRepository) GetByID(id string) (*category.Category, error) {
	categoryModel := category.Get(id)
	if categoryModel.ID == 0 {
		return nil, apperrors.NotFoundError("分类").WithDetails(map[string]interface{}{
			"category_id": id,
		})
	}
	return &categoryModel, nil
}

// List 获取分类列表
func (r *categoryRepository) List(c *gin.Context, perPage int) ([]category.Category, *paginator.Paging, error) {
	data, pager := category.Paginate(c, perPage)
	return data, &pager, nil
}

// Create 创建分类
func (r *categoryRepository) Create(cat *category.Category) error {
	cat.Create()
	if cat.ID == 0 {
		return apperrors.DatabaseError("创建分类", nil)
	}

	// 清除列表缓存
	_ = r.FlushCache()

	return nil
}

// Update 更新分类
func (r *categoryRepository) Update(cat *category.Category) error {
	rowsAffected := cat.Save()
	if rowsAffected == 0 {
		return apperrors.DatabaseError("更新分类", nil)
	}

	// 清除缓存
	_ = r.FlushCache()

	return nil
}

// Delete 删除分类
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


// GetAllCached 获取所有分类（缓存）
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

// SetListCache 设置列表缓存
func (r *categoryRepository) SetListCache(categories []category.Category) error {
	data, err := json.Marshal(categories)
	if err != nil {
		return err
	}

	redis.Redis.Set(r.cacheKeyList, string(data), time.Duration(r.cacheTTL)*time.Second)
	return nil
}

// FlushCache 清空缓存
func (r *categoryRepository) FlushCache() error {
	_ = redis.Redis.Del(r.cacheKeyList)
	pattern := fmt.Sprintf("%s*", "category:")
	_ = pattern // 简化实现
	return nil
}
