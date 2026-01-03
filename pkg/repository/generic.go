// Package repository 通用 Repository 泛型基类
package repository

import (
	"context"
	"fmt"

	"GoHub-Service/pkg/database"
	"GoHub-Service/pkg/paginator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Model 模型接口，所有模型必须实现
type Model interface {
	GetStringID() string
	TableName() string
}

// GenericRepository 泛型仓储，提供通用 CRUD 操作
type GenericRepository[T Model] struct {
	db *gorm.DB
}

// NewGenericRepository 创建泛型仓储实例
func NewGenericRepository[T Model]() *GenericRepository[T] {
	return &GenericRepository[T]{
		db: database.DB,
	}
}

// GetByID 根据 ID 获取单条记录
func (r *GenericRepository[T]) GetByID(ctx context.Context, id string) (*T, error) {
	var model T
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

// GetByIDWithPreload 根据 ID 获取单条记录（带预加载）
func (r *GenericRepository[T]) GetByIDWithPreload(ctx context.Context, id string, preloads ...string) (*T, error) {
	var model T
	query := r.db.WithContext(ctx)
	
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	
	err := query.Where("id = ?", id).First(&model).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

// List 获取列表（带分页）
func (r *GenericRepository[T]) List(ctx context.Context, c *gin.Context, perPage int) ([]T, *paginator.Paging, error) {
	var models []T
	var model T
	
	query := r.db.WithContext(ctx).Model(&model)
	paging := paginator.Paginate(c, query, &models, perPage)
	
	return models, paging, nil
}

// ListWithCondition 获取列表（带条件和分页）
func (r *GenericRepository[T]) ListWithCondition(ctx context.Context, c *gin.Context, condition interface{}, args []interface{}, perPage int) ([]T, *paginator.Paging, error) {
	var models []T
	var model T
	
	query := r.db.WithContext(ctx).Model(&model).Where(condition, args...)
	paging := paginator.Paginate(c, query, &models, perPage)
	
	return models, paging, nil
}

// Create 创建记录
func (r *GenericRepository[T]) Create(ctx context.Context, model *T) error {
	return r.db.WithContext(ctx).Create(model).Error
}

// Update 更新记录
func (r *GenericRepository[T]) Update(ctx context.Context, model *T) error {
	return r.db.WithContext(ctx).Save(model).Error
}

// UpdateFields 更新指定字段
func (r *GenericRepository[T]) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	var model T
	return r.db.WithContext(ctx).Model(&model).Where("id = ?", id).Updates(fields).Error
}

// Delete 删除记录
func (r *GenericRepository[T]) Delete(ctx context.Context, id string) error {
	var model T
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model).Error
}

// BatchCreate 批量创建
func (r *GenericRepository[T]) BatchCreate(ctx context.Context, models []T) error {
	if len(models) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&models).Error
}

// BatchCreateInChunks 分块批量创建（大数据量优化）
func (r *GenericRepository[T]) BatchCreateInChunks(ctx context.Context, models []T, chunkSize int) error {
	if len(models) == 0 {
		return nil
	}

	if chunkSize <= 0 {
		chunkSize = 100
	}

	for i := 0; i < len(models); i += chunkSize {
		end := i + chunkSize
		if end > len(models) {
			end = len(models)
		}
		
		if err := r.db.WithContext(ctx).Create(models[i:end]).Error; err != nil {
			return err
		}
	}

	return nil
}

// BatchDelete 批量删除
func (r *GenericRepository[T]) BatchDelete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	
	var model T
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model).Error
}

// Count 计数
func (r *GenericRepository[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	var model T
	err := r.db.WithContext(ctx).Model(&model).Count(&count).Error
	return count, err
}

// CountWithCondition 条件计数
func (r *GenericRepository[T]) CountWithCondition(ctx context.Context, condition interface{}, args ...interface{}) (int64, error) {
	var count int64
	var model T
	err := r.db.WithContext(ctx).Model(&model).Where(condition, args...).Count(&count).Error
	return count, err
}

// Exists 检查记录是否存在
func (r *GenericRepository[T]) Exists(ctx context.Context, id string) (bool, error) {
	var count int64
	var model T
	err := r.db.WithContext(ctx).Model(&model).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// Transaction 执行事务
func (r *GenericRepository[T]) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

// GetDB 获取数据库实例（用于复杂查询）
func (r *GenericRepository[T]) GetDB() *gorm.DB {
	return r.db
}

// WithContext 返回带 context 的 DB
func (r *GenericRepository[T]) WithContext(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// Increment 增加字段值
func (r *GenericRepository[T]) Increment(ctx context.Context, id string, field string, value int) error {
	var model T
	return r.db.WithContext(ctx).Model(&model).Where("id = ?", id).
		Update(field, gorm.Expr(fmt.Sprintf("%s + ?", field), value)).Error
}

// Decrement 减少字段值
func (r *GenericRepository[T]) Decrement(ctx context.Context, id string, field string, value int) error {
	var model T
	return r.db.WithContext(ctx).Model(&model).Where("id = ?", id).
		Update(field, gorm.Expr(fmt.Sprintf("%s - ?", field), value)).Error
}

// FindBy 根据字段查找单条记录
func (r *GenericRepository[T]) FindBy(ctx context.Context, field string, value interface{}) (*T, error) {
	var model T
	err := r.db.WithContext(ctx).Where(fmt.Sprintf("%s = ?", field), value).First(&model).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

// FindManyBy 根据字段查找多条记录
func (r *GenericRepository[T]) FindManyBy(ctx context.Context, field string, value interface{}) ([]T, error) {
	var models []T
	err := r.db.WithContext(ctx).Where(fmt.Sprintf("%s = ?", field), value).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}
