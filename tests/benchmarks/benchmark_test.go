// Package benchmarks 性能基准测试
package benchmarks

import (
	"context"
	"testing"

	"github.com/Oyxiaoxi/GoHub-Service/app/models/comment"
	"github.com/Oyxiaoxi/GoHub-Service/app/models/topic"
	"github.com/Oyxiaoxi/GoHub-Service/app/models/user"
	"github.com/Oyxiaoxi/GoHub-Service/pkg/mapper"
	"github.com/Oyxiaoxi/GoHub-Service/pkg/resource"
	"github.com/Oyxiaoxi/GoHub-Service/pkg/testutil"
)

// BenchmarkMapperBatchConvert 测试 Mapper 批量转换性能
func BenchmarkMapperBatchConvert(b *testing.B) {
	users := testutil.MockUsers(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapper.MapSlice(users, func(u *user.User) map[string]interface{} {
			return map[string]interface{}{
				"id":   u.ID,
				"name": u.Name,
			}
		})
	}
}

// BenchmarkMapperSingleConvert 测试 Mapper 单个转换性能
func BenchmarkMapperSingleConvert(b *testing.B) {
	user := testutil.MockUserFactory()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapper.Map(user, func(u *user.User) map[string]interface{} {
			return map[string]interface{}{
				"id":   u.ID,
				"name": u.Name,
			}
		})
	}
}

// BenchmarkResourceManager 测试 ResourceManager 性能
func BenchmarkResourceManager(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rm := resource.NewResourceManager()
		defer rm.Release()

		// 模拟资源管理
		rm.Track("test-resource", func() {})
	}
}

// BenchmarkBatchCreate 测试批量创建性能
func BenchmarkBatchCreate(b *testing.B) {
	helper := testutil.SetupTestEnvironment(&testing.T{})
	defer helper.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tx := helper.DB.Begin()
		users := testutil.MockUsers(10)
		for _, u := range users {
			_ = tx.Create(u).Error
		}
		tx.Rollback()
	}
}

// BenchmarkMapperComplexConversion 测试复杂对象转换性能
func BenchmarkMapperComplexConversion(b *testing.B) {
	topics := make([]*topic.Topic, 100)
	for i := range topics {
		topics[i] = testutil.MockTopicFactory()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapper.MapSlice(topics, func(t *topic.Topic) map[string]interface{} {
			return map[string]interface{}{
				"id":            t.ID,
				"title":         t.Title,
				"body":          t.Body,
				"user_id":       t.UserID,
				"category_id":   t.CategoryID,
				"view_count":    t.ViewCount,
				"like_count":    t.LikeCount,
				"comment_count": t.CommentCount,
				"created_at":    t.CreatedAt,
			}
		})
	}
}

// BenchmarkMapperWithPreload 测试带预加载的转换性能
func BenchmarkMapperWithPreload(b *testing.B) {
	comments := make([]*comment.Comment, 100)
	for i := range comments {
		comments[i] = testutil.MockCommentFactory()
		comments[i].User = *testutil.MockUserFactory()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapper.MapSlice(comments, func(c *comment.Comment) map[string]interface{} {
			return map[string]interface{}{
				"id":         c.ID,
				"content":    c.Content,
				"user_id":    c.UserID,
				"topic_id":   c.TopicID,
				"like_count": c.LikeCount,
				"user": map[string]interface{}{
					"id":   c.User.ID,
					"name": c.User.Name,
				},
			}
		})
	}
}

// BenchmarkResourceManagerConcurrent 测试并发资源管理性能
func BenchmarkResourceManagerConcurrent(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rm := resource.NewResourceManager()
			rm.Track("concurrent-resource", func() {})
			rm.Release()
		}
	})
}

// BenchmarkMapperParallel 测试并行 Mapper 性能
func BenchmarkMapperParallel(b *testing.B) {
	users := testutil.MockUsers(100)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = mapper.MapSlice(users, func(u *user.User) map[string]interface{} {
				return map[string]interface{}{
					"id":   u.ID,
					"name": u.Name,
				}
			})
		}
	})
}

// BenchmarkMapperMemoryAllocation 测试内存分配
func BenchmarkMapperMemoryAllocation(b *testing.B) {
	users := testutil.MockUsers(1000)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = mapper.MapSlice(users, func(u *user.User) map[string]interface{} {
			return map[string]interface{}{
				"id":   u.ID,
				"name": u.Name,
			}
		})
	}
}
