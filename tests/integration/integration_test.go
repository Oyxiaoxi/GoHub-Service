// Package integration 集成测试示例
package integration_test

import (
	"GoHub-Service/app/models/category"
	"GoHub-Service/app/models/topic"
	"GoHub-Service/app/models/user"
	"GoHub-Service/pkg/testutil"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 这是一个集成测试示例，展示如何使用集成测试框架

// TestUserTopicCategoryIntegration 测试用户-话题-分类的完整流程
func TestUserTopicCategoryIntegration(t *testing.T) {
	// 初始化测试环境
	testutil.SetupTestEnvironment(t, nil)

	// 使用测试数据运行测试
	testutil.RunWithTestData(t, seedUserTopicCategoryData, func(db *gorm.DB) {
		ctx := testutil.MockContext()

		// 1. 查询用户
		var u user.User
		err := db.WithContext(ctx).First(&u, "email = ?", "test@example.com").Error
		assert.NoError(t, err)
		assert.Equal(t, "测试用户", u.Name)

		// 2. 查询分类
		var cat category.Category
		err = db.WithContext(ctx).First(&cat, "name = ?", "技术分享").Error
		assert.NoError(t, err)

		// 3. 创建话题
		newTopic := &topic.Topic{
			Title:      "集成测试话题",
			Body:       "这是一个集成测试创建的话题",
			UserID:     u.ID,
			CategoryID: cat.ID,
		}
		err = db.WithContext(ctx).Create(newTopic).Error
		assert.NoError(t, err)
		assert.NotEmpty(t, newTopic.ID)

		// 4. 验证话题已创建
		var createdTopic topic.Topic
		err = db.WithContext(ctx).Preload("User").Preload("Category").First(&createdTopic, newTopic.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, "集成测试话题", createdTopic.Title)
		assert.Equal(t, u.ID, createdTopic.UserID)
		assert.Equal(t, cat.ID, createdTopic.CategoryID)

		// 5. 统计用户话题数
		var topicCount int64
		err = db.WithContext(ctx).Model(&topic.Topic{}).Where("user_id = ?", u.ID).Count(&topicCount).Error
		assert.NoError(t, err)
		assert.Greater(t, topicCount, int64(0))

		// 6. 断言表记录数
		testutil.AssertDBCount(t, db, "topics", topicCount)
		testutil.AssertRecordExists(t, db, "topics", "title = ?", "集成测试话题")
	})
}

// seedUserTopicCategoryData 填充测试数据
func seedUserTopicCategoryData(db *gorm.DB) error {
	ctx := context.Background()

	// 创建测试用户
	testUser := &user.User{
		Name:  "测试用户",
		Email: "test@example.com",
		Phone: "13800138000",
	}
	if err := db.WithContext(ctx).Create(testUser).Error; err != nil {
		return err
	}

	// 创建测试分类
	testCategory := &category.Category{
		Name:        "技术分享",
		Description: "技术相关的话题",
	}
	if err := db.WithContext(ctx).Create(testCategory).Error; err != nil {
		return err
	}

	// 创建初始话题
	testTopic := &topic.Topic{
		Title:      "初始测试话题",
		Body:       "这是一个预先创建的测试话题",
		UserID:     testUser.ID,
		CategoryID: testCategory.ID,
	}
	if err := db.WithContext(ctx).Create(testTopic).Error; err != nil {
		return err
	}

	return nil
}

// TestCommentIntegration 测试评论完整流程
func TestCommentIntegration(t *testing.T) {
	t.Skip("需要完整的模型导入和数据库配置")

	testutil.SetupTestEnvironment(t, nil)

	testutil.RunWithTestData(t, nil, func(db *gorm.DB) {
		helper := testutil.NewTestHelper(t, db)

		// 测试逻辑...
		helper.AssertNoError(nil, "测试评论创建")
	})
}

// TestBatchOperationIntegration 测试批量操作
func TestBatchOperationIntegration(t *testing.T) {
	t.Skip("需要完整的模型导入和数据库配置")

	testutil.SetupTestEnvironment(t, &testutil.TestConfig{
		DBPath:      "./test_batch.db",
		LogLevel:    "warn",
		EnableCache: false,
	})

	tx := testutil.BeginTransaction(t)

	// 批量创建测试数据
	users := testutil.MockUsers(10)
	for i := range users {
		if err := tx.Create(&users[i]).Error; err != nil {
			t.Fatalf("批量创建用户失败: %v", err)
		}
	}

	// 验证批量创建结果
	testutil.AssertDBCount(t, tx, "users", 10)
}
