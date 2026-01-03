// Package database 数据库查询助手
package database

import (
	"gorm.io/gorm"
)

// QueryHelper 查询助手，提供常用的查询优化方法
type QueryHelper struct {
	db *gorm.DB
}

// NewQueryHelper 创建查询助手实例
func NewQueryHelper(db *gorm.DB) *QueryHelper {
	return &QueryHelper{db: db}
}

// SelectFields 常用字段选择器，避免查询所有字段
var SelectFields = map[string][]string{
	// 用户基本字段
	"user_basic": {"id", "name", "email", "avatar"},
	"user_list":  {"id", "name", "email", "avatar", "created_at"},
	"user_full":  {"id", "name", "email", "avatar", "phone", "introduction", "created_at", "updated_at"},

	// 话题字段
	"topic_basic": {"id", "title", "user_id", "category_id"},
	"topic_list":  {"id", "title", "body", "user_id", "category_id", "like_count", "favorite_count", "view_count", "created_at"},
	"topic_full":  {"id", "title", "body", "user_id", "category_id", "like_count", "favorite_count", "view_count", "created_at", "updated_at"},

	// 评论字段
	"comment_basic": {"id", "content", "user_id", "topic_id"},
	"comment_list":  {"id", "topic_id", "user_id", "content", "parent_id", "like_count", "created_at"},
	"comment_full":  {"id", "topic_id", "user_id", "content", "parent_id", "like_count", "created_at", "updated_at"},

	// 分类字段
	"category_basic": {"id", "name"},
	"category_list":  {"id", "name", "description"},
	"category_full":  {"id", "name", "description", "created_at", "updated_at"},
}

// ApplySelect 应用字段选择
func (h *QueryHelper) ApplySelect(db *gorm.DB, selectType string) *gorm.DB {
	if fields, ok := SelectFields[selectType]; ok {
		return db.Select(fields)
	}
	return db
}

// PreloadUser 预加载用户关联（优化版）
func PreloadUser(selectType string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if fields, ok := SelectFields[selectType]; ok {
			return db.Select(fields)
		}
		return db.Select(SelectFields["user_basic"])
	}
}

// PreloadTopic 预加载话题关联（优化版）
func PreloadTopic(selectType string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if fields, ok := SelectFields[selectType]; ok {
			return db.Select(fields)
		}
		return db.Select(SelectFields["topic_basic"])
	}
}

// PreloadCategory 预加载分类关联（优化版）
func PreloadCategory(selectType string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if fields, ok := SelectFields[selectType]; ok {
			return db.Select(fields)
		}
		return db.Select(SelectFields["category_basic"])
	}
}

// PreloadComment 预加载评论关联（优化版）
func PreloadComment(selectType string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if fields, ok := SelectFields[selectType]; ok {
			return db.Select(fields)
		}
		return db.Select(SelectFields["comment_basic"])
	}
}

// BatchSize 批量操作的默认大小
const BatchSize = 100

// BatchInsert 批量插入助手
func BatchInsert(db *gorm.DB, records interface{}, batchSize int) error {
	if batchSize <= 0 {
		batchSize = BatchSize
	}
	return db.CreateInBatches(records, batchSize).Error
}

// OptimizeQuery 优化查询（添加常用的优化选项）
func OptimizeQuery(db *gorm.DB) *gorm.DB {
	// 可以在这里添加通用的查询优化
	// 例如：设置超时、添加索引提示等
	return db
}
