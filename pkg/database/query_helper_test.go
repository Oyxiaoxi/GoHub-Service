package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSelectFields 测试字段选择器
func TestSelectFields(t *testing.T) {
	t.Run("user_basic字段", func(t *testing.T) {
		fields := SelectFields["user_basic"]
		assert.NotNil(t, fields)
		assert.Contains(t, fields, "id")
		assert.Contains(t, fields, "name")
		assert.Contains(t, fields, "email")
		assert.Contains(t, fields, "avatar")
		assert.Equal(t, 4, len(fields))
	})

	t.Run("topic_list字段", func(t *testing.T) {
		fields := SelectFields["topic_list"]
		assert.NotNil(t, fields)
		assert.Contains(t, fields, "id")
		assert.Contains(t, fields, "title")
		assert.Contains(t, fields, "body")
		assert.Contains(t, fields, "user_id")
		assert.Contains(t, fields, "category_id")
		assert.Contains(t, fields, "like_count")
		assert.Contains(t, fields, "favorite_count")
		assert.Contains(t, fields, "view_count")
		assert.Contains(t, fields, "created_at")
	})

	t.Run("comment_full字段", func(t *testing.T) {
		fields := SelectFields["comment_full"]
		assert.NotNil(t, fields)
		assert.Contains(t, fields, "id")
		assert.Contains(t, fields, "topic_id")
		assert.Contains(t, fields, "user_id")
		assert.Contains(t, fields, "content")
		assert.Contains(t, fields, "parent_id")
		assert.Contains(t, fields, "like_count")
		assert.Contains(t, fields, "created_at")
		assert.Contains(t, fields, "updated_at")
	})
}

// TestPreloadHelpers 测试预加载助手函数
func TestPreloadHelpers(t *testing.T) {
	t.Run("PreloadUser", func(t *testing.T) {
		fn := PreloadUser("user_basic")
		assert.NotNil(t, fn, "PreloadUser应该返回一个函数")
	})

	t.Run("PreloadTopic", func(t *testing.T) {
		fn := PreloadTopic("topic_list")
		assert.NotNil(t, fn, "PreloadTopic应该返回一个函数")
	})

	t.Run("PreloadCategory", func(t *testing.T) {
		fn := PreloadCategory("category_basic")
		assert.NotNil(t, fn, "PreloadCategory应该返回一个函数")
	})

	t.Run("PreloadComment", func(t *testing.T) {
		fn := PreloadComment("comment_list")
		assert.NotNil(t, fn, "PreloadComment应该返回一个函数")
	})
}

// TestBatchSize 测试批量操作大小常量
func TestBatchSize(t *testing.T) {
	assert.Equal(t, 100, BatchSize, "默认批量大小应该是100")
}

// TestNewQueryHelper 测试查询助手创建
func TestNewQueryHelper(t *testing.T) {
	// 不需要真实的数据库连接来测试创建
	helper := NewQueryHelper(nil)
	assert.NotNil(t, helper, "NewQueryHelper应该返回非nil实例")
}
