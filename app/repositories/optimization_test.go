package repositories

import (
	"context"

	"GoHub-Service/app/models/comment"
	"GoHub-Service/pkg/database"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// BenchmarkQueryOptimization 查询优化性能测试
func BenchmarkQueryOptimization(b *testing.B) {
	// 初始化数据库连接（在实际测试中需要配置）
	// setupTestDB()

	b.Run("优化前的查询", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var comments []comment.Comment
			// 优化前：查询所有字段
			database.DB.
				Preload("User").
				Preload("Topic").
				Limit(50).
				Find(&comments)
		}
	})

	b.Run("优化后的查询", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var comments []comment.Comment
			// 优化后：只查询必要字段
			database.DB.Model(&comment.Comment{}).
				Select("id", "topic_id", "user_id", "content", "like_count", "created_at").
				Preload("User", database.PreloadUser("user_basic")).
				Preload("Topic", database.PreloadTopic("topic_basic")).
				Limit(50).
				Find(&comments)
		}
	})
}

// BenchmarkBatchOperations 批量操作性能测试
func BenchmarkBatchOperations(b *testing.B) {
	repo := NewCommentRepository()

	// 准备测试数据
	prepareComments := func(count int) []comment.Comment {
		comments := make([]comment.Comment, count)
		for i := 0; i < count; i++ {
			comments[i] = comment.Comment{
				Content: "Test comment",
				TopicID: "1",
				UserID:  "1",
			}
		}
		return comments
	}

	b.Run("循环单条插入100条", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			comments := prepareComments(100)
			for j := range comments {
				repo.Create(context.Background(), &comments[j])
			}
		}
	})

	b.Run("批量插入100条", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			comments := prepareComments(100)
			repo.BatchCreate(context.Background(), comments)
		}
	})
}

// TestQueryPerformance 查询性能对比测试
func TestQueryPerformance(t *testing.T) {
	// 跳过集成测试（需要数据库连接）
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	t.Run("单条查询性能对比", func(t *testing.T) {
		// 优化前
		start := time.Now()
		var comment1 comment.Comment
		database.DB.Preload("User").Preload("Topic").First(&comment1, 1)
		duration1 := time.Since(start)

		// 优化后
		start = time.Now()
		var comment2 comment.Comment
		database.DB.
			Select("id", "topic_id", "user_id", "content", "like_count", "created_at").
			Preload("User", database.PreloadUser("user_basic")).
			Preload("Topic", database.PreloadTopic("topic_basic")).
			First(&comment2, 1)
		duration2 := time.Since(start)

		t.Logf("优化前耗时: %v", duration1)
		t.Logf("优化后耗时: %v", duration2)
		t.Logf("性能提升: %.2f%%", float64(duration1-duration2)/float64(duration1)*100)

		// 验证优化后的查询速度更快
		assert.Less(t, duration2, duration1, "优化后的查询应该更快")
	})

	t.Run("列表查询性能对比", func(t *testing.T) {
		// 优化前
		start := time.Now()
		var comments1 []comment.Comment
		database.DB.Preload("User").Preload("Topic").Limit(50).Find(&comments1)
		duration1 := time.Since(start)

		// 优化后
		start = time.Now()
		var comments2 []comment.Comment
		database.DB.Model(&comment.Comment{}).
			Select("id", "topic_id", "user_id", "content", "like_count", "created_at").
			Preload("User", database.PreloadUser("user_basic")).
			Preload("Topic", database.PreloadTopic("topic_basic")).
			Limit(50).
			Find(&comments2)
		duration2 := time.Since(start)

		t.Logf("优化前耗时: %v", duration1)
		t.Logf("优化后耗时: %v", duration2)
		t.Logf("性能提升: %.2f%%", float64(duration1-duration2)/float64(duration1)*100)

		assert.Less(t, duration2, duration1, "优化后的查询应该更快")
		assert.Equal(t, len(comments1), len(comments2), "查询结果数量应该相同")
	})
}

// TestBatchOperations 批量操作功能测试
func TestBatchOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	repo := NewCommentRepository()

	t.Run("批量创建评论", func(t *testing.T) {
		comments := []comment.Comment{
			{Content: "Test 1", TopicID: "1", UserID: "1"},
			{Content: "Test 2", TopicID: "1", UserID: "1"},
			{Content: "Test 3", TopicID: "1", UserID: "1"},
		}

		err := repo.BatchCreate(context.Background(), comments)
		assert.NoError(t, err)

		// 验证所有评论都已创建
		for _, c := range comments {
			assert.NotZero(t, c.ID, "评论应该有ID")
		}
	})

	t.Run("批量删除评论", func(t *testing.T) {
		// 先创建一些测试数据
		comments := []comment.Comment{
			{Content: "To Delete 1", TopicID: "1", UserID: "1"},
			{Content: "To Delete 2", TopicID: "1", UserID: "1"},
			{Content: "To Delete 3", TopicID: "1", UserID: "1"},
		}
		err := repo.BatchCreate(context.Background(), comments)
		assert.NoError(t, err)

		// 批量删除
		ids := make([]string, len(comments))
		for i, c := range comments {
			ids[i] = c.GetStringID()
		}

		err = repo.BatchDelete(context.Background(), ids)
		assert.NoError(t, err)

		// 验证已删除
		for _, id := range ids {
			c, err := repo.GetByID(context.Background(), id)
			assert.NoError(t, err)
			assert.Nil(t, c, "评论应该已被删除")
		}
	})

	t.Run("空列表批量操作", func(t *testing.T) {
		// 测试空列表的处理
		err := repo.BatchCreate([]comment.Comment{})
		assert.NoError(t, err, "空列表应该正常处理")

		err = repo.BatchDelete([]string{})
		assert.NoError(t, err, "空列表应该正常处理")
	})
}

// TestQueryHelper 查询助手测试
func TestQueryHelper(t *testing.T) {
	t.Run("字段选择器", func(t *testing.T) {
		// 测试预定义的字段选择器
		userBasic := database.SelectFields["user_basic"]
		assert.Contains(t, userBasic, "id")
		assert.Contains(t, userBasic, "name")
		assert.Contains(t, userBasic, "email")
		assert.Contains(t, userBasic, "avatar")

		commentList := database.SelectFields["comment_list"]
		assert.Contains(t, commentList, "id")
		assert.Contains(t, commentList, "content")
		assert.Contains(t, commentList, "user_id")
		assert.Contains(t, commentList, "like_count")
	})

	t.Run("Preload助手函数", func(t *testing.T) {
		// 测试 Preload 助手函数不会返回 nil
		userPreload := database.PreloadUser("user_basic")
		assert.NotNil(t, userPreload)

		topicPreload := database.PreloadTopic("topic_basic")
		assert.NotNil(t, topicPreload)

		categoryPreload := database.PreloadCategory("category_basic")
		assert.NotNil(t, categoryPreload)
	})
}

// 运行基准测试命令：
// go test -bench=. -benchmem -benchtime=10s ./app/repositories/

// 运行性能测试：
// go test -v -run TestQueryPerformance ./app/repositories/

// 查看测试覆盖率：
// go test -cover ./app/repositories/
