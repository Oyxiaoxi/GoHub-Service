package repositories

import (
	"context"
	"testing"

	"github.com/Oyxiaoxi/GoHub-Service/app/models/message"
	"github.com/Oyxiaoxi/GoHub-Service/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupMessageRepositoryTest(t *testing.T) (*MessageRepository, *gorm.DB, func()) {
	helper := testutil.SetupTestEnvironment(t)
	repo := NewMessageRepository()
	return repo, helper.DB, helper.Close
}

func TestMessageRepository_Create(t *testing.T) {
	repo, db, cleanup := setupMessageRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("创建消息成功", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		testMessage := testutil.MockMessageFactory()
		err := repo.Create(ctx, tx, testMessage)

		assert.NoError(t, err)
		assert.NotZero(t, testMessage.ID)
		assert.NotZero(t, testMessage.CreatedAt)
	})

	t.Run("创建消息缺少必填字段", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		invalidMessage := &message.Message{}
		err := repo.Create(ctx, tx, invalidMessage)

		assert.Error(t, err)
	})
}

func TestMessageRepository_ListConversation(t *testing.T) {
	repo, db, cleanup := setupMessageRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("获取对话消息列表", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		// 创建测试消息
		user1 := testutil.MockUserFactory()
		user2 := testutil.MockUserFactory()
		_ = tx.Create(user1).Error
		_ = tx.Create(user2).Error

		for i := 0; i < 5; i++ {
			msg := testutil.MockMessageFactory()
			msg.FromUserID = user1.ID
			msg.ToUserID = user2.ID
			_ = repo.Create(ctx, tx, msg)
		}

		messages, total, err := repo.ListConversation(ctx, tx, user1.ID, user2.ID, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, messages, 5)
	})
}

func TestMessageRepository_MarkConversationRead(t *testing.T) {
	repo, db, cleanup := setupMessageRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("标记对话已读", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		user1 := testutil.MockUserFactory()
		user2 := testutil.MockUserFactory()
		_ = tx.Create(user1).Error
		_ = tx.Create(user2).Error

		// 创建未读消息
		msg := testutil.MockMessageFactory()
		msg.FromUserID = user1.ID
		msg.ToUserID = user2.ID
		msg.ReadAt = nil
		_ = repo.Create(ctx, tx, msg)

		err := repo.MarkConversationRead(ctx, tx, user1.ID, user2.ID)

		assert.NoError(t, err)

		// 验证已读
		var updatedMsg message.Message
		tx.Where("id = ?", msg.ID).First(&updatedMsg)
		assert.NotNil(t, updatedMsg.ReadAt)
	})
}

func TestMessageRepository_CountUnread(t *testing.T) {
	repo, db, cleanup := setupMessageRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("统计未读消息数", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		user1 := testutil.MockUserFactory()
		user2 := testutil.MockUserFactory()
		_ = tx.Create(user1).Error
		_ = tx.Create(user2).Error

		// 创建 3 条未读消息
		for i := 0; i < 3; i++ {
			msg := testutil.MockMessageFactory()
			msg.FromUserID = user1.ID
			msg.ToUserID = user2.ID
			msg.ReadAt = nil
			_ = repo.Create(ctx, tx, msg)
		}

		count, err := repo.CountUnread(ctx, tx, user2.ID)

		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})
}

func TestMessageRepository_Delete(t *testing.T) {
	repo, db, cleanup := setupMessageRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("删除消息", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		testMessage := testutil.MockMessageFactory()
		_ = repo.Create(ctx, tx, testMessage)

		err := repo.Delete(ctx, tx, testMessage.ID)

		assert.NoError(t, err)

		// 验证删除
		var deletedMsg message.Message
		result := tx.Where("id = ?", testMessage.ID).First(&deletedMsg)
		assert.Error(t, result.Error)
	})
}

func TestMessageRepository_GetByID(t *testing.T) {
	repo, db, cleanup := setupMessageRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("根据ID获取消息", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		testMessage := testutil.MockMessageFactory()
		_ = repo.Create(ctx, tx, testMessage)

		result, err := repo.GetByID(ctx, tx, testMessage.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, testMessage.ID, result.ID)
		assert.Equal(t, testMessage.Content, result.Content)
	})

	t.Run("获取不存在的消息", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		result, err := repo.GetByID(ctx, tx, 99999)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestMessageRepository_BatchCreate(t *testing.T) {
	repo, db, cleanup := setupMessageRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("批量创建消息", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		messages := make([]*message.Message, 10)
		for i := range messages {
			messages[i] = testutil.MockMessageFactory()
		}

		err := repo.BatchCreate(ctx, tx, messages)

		assert.NoError(t, err)
		for _, msg := range messages {
			assert.NotZero(t, msg.ID)
		}
	})
}
