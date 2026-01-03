// Package services 消息服务测试
package services

import (
	"GoHub-Service/app/models/message"
	"GoHub-Service/pkg/testutil"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockMessageRepository 消息仓储 Mock
type MockMessageRepository struct {
	GetByIDFunc            func(ctx context.Context, id string) (*message.Message, error)
	GetConversationFunc    func(ctx context.Context, userID1, userID2 string, limit int) ([]message.Message, error)
	GetUserMessagesFunc    func(ctx context.Context, userID string, limit int) ([]message.Message, error)
	CreateFunc             func(ctx context.Context, m *message.Message) error
	MarkAsReadFunc         func(ctx context.Context, id string) error
	MarkConversationReadFunc func(ctx context.Context, userID1, userID2 string) error
	DeleteFunc             func(ctx context.Context, id string) error
	CountUnreadFunc        func(ctx context.Context, userID string) (int64, error)
	GetUnreadFunc          func(ctx context.Context, userID string) ([]message.Message, error)
}

func (m *MockMessageRepository) GetByID(ctx context.Context, id string) (*message.Message, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return testutil.MockMessageFactory(id, "测试消息", "1", "2"), nil
}

func (m *MockMessageRepository) GetConversation(ctx context.Context, userID1, userID2 string, limit int) ([]message.Message, error) {
	if m.GetConversationFunc != nil {
		return m.GetConversationFunc(ctx, userID1, userID2, limit)
	}
	return testutil.MockMessages(limit), nil
}

func (m *MockMessageRepository) GetUserMessages(ctx context.Context, userID string, limit int) ([]message.Message, error) {
	if m.GetUserMessagesFunc != nil {
		return m.GetUserMessagesFunc(ctx, userID, limit)
	}
	return testutil.MockMessages(limit), nil
}

func (m *MockMessageRepository) Create(ctx context.Context, msg *message.Message) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, msg)
	}
	msg.ID = "1"
	return nil
}

func (m *MockMessageRepository) MarkAsRead(ctx context.Context, id string) error {
	if m.MarkAsReadFunc != nil {
		return m.MarkAsReadFunc(ctx, id)
	}
	return nil
}

func (m *MockMessageRepository) MarkConversationRead(ctx context.Context, userID1, userID2 string) error {
	if m.MarkConversationReadFunc != nil {
		return m.MarkConversationReadFunc(ctx, userID1, userID2)
	}
	return nil
}

func (m *MockMessageRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockMessageRepository) CountUnread(ctx context.Context, userID string) (int64, error) {
	if m.CountUnreadFunc != nil {
		return m.CountUnreadFunc(ctx, userID)
	}
	return 5, nil
}

func (m *MockMessageRepository) GetUnread(ctx context.Context, userID string) ([]message.Message, error) {
	if m.GetUnreadFunc != nil {
		return m.GetUnreadFunc(ctx, userID)
	}
	return testutil.MockMessages(5), nil
}

// TestMessageService_Create 测试创建消息
func TestMessageService_Create(t *testing.T) {
	mockRepo := &MockMessageRepository{}
	ctx := context.Background()

	tests := []struct {
		name       string
		content    string
		fromUserID string
		toUserID   string
		setupMock  func()
		wantErr    bool
	}{
		{
			name:       "成功创建消息",
			content:    "你好，这是一条测试消息",
			fromUserID: "1",
			toUserID:   "2",
			setupMock: func() {
				mockRepo.CreateFunc = func(ctx context.Context, m *message.Message) error {
					m.ID = "1"
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:       "消息内容为空",
			content:    "",
			fromUserID: "1",
			toUserID:   "2",
			setupMock: func() {
				mockRepo.CreateFunc = func(ctx context.Context, m *message.Message) error {
					return errors.New("消息内容不能为空")
				}
			},
			wantErr: true,
		},
		{
			name:       "发送者ID为空",
			content:    "测试消息",
			fromUserID: "",
			toUserID:   "2",
			setupMock: func() {
				mockRepo.CreateFunc = func(ctx context.Context, m *message.Message) error {
					return errors.New("发送者ID不能为空")
				}
			},
			wantErr: true,
		},
		{
			name:       "接收者ID为空",
			content:    "测试消息",
			fromUserID: "1",
			toUserID:   "",
			setupMock: func() {
				mockRepo.CreateFunc = func(ctx context.Context, m *message.Message) error {
					return errors.New("接收者ID不能为空")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			newMessage := &message.Message{
				Content:    tt.content,
				FromUserID: tt.fromUserID,
				ToUserID:   tt.toUserID,
			}

			err := mockRepo.Create(ctx, newMessage)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, newMessage.ID)
			}
		})
	}
}

// TestMessageService_GetConversation 测试获取会话消息
func TestMessageService_GetConversation(t *testing.T) {
	mockRepo := &MockMessageRepository{
		GetConversationFunc: func(ctx context.Context, userID1, userID2 string, limit int) ([]message.Message, error) {
			if userID1 == "" || userID2 == "" {
				return nil, errors.New("用户ID不能为空")
			}
			if limit <= 0 {
				return nil, errors.New("limit必须大于0")
			}
			return testutil.MockMessages(limit), nil
		},
	}

	ctx := context.Background()

	t.Run("成功获取会话消息", func(t *testing.T) {
		messages, err := mockRepo.GetConversation(ctx, "1", "2", 10)
		assert.NoError(t, err)
		assert.Len(t, messages, 10)
	})

	t.Run("用户ID为空", func(t *testing.T) {
		messages, err := mockRepo.GetConversation(ctx, "", "2", 10)
		assert.Error(t, err)
		assert.Nil(t, messages)
	})

	t.Run("limit为0", func(t *testing.T) {
		messages, err := mockRepo.GetConversation(ctx, "1", "2", 0)
		assert.Error(t, err)
		assert.Nil(t, messages)
	})
}

// TestMessageService_GetUserMessages 测试获取用户消息列表
func TestMessageService_GetUserMessages(t *testing.T) {
	mockRepo := &MockMessageRepository{
		GetUserMessagesFunc: func(ctx context.Context, userID string, limit int) ([]message.Message, error) {
			if userID == "" {
				return nil, errors.New("用户ID不能为空")
			}
			return testutil.MockMessages(limit), nil
		},
	}

	ctx := context.Background()

	t.Run("成功获取用户消息", func(t *testing.T) {
		messages, err := mockRepo.GetUserMessages(ctx, "1", 20)
		assert.NoError(t, err)
		assert.Len(t, messages, 20)
	})

	t.Run("用户ID为空", func(t *testing.T) {
		messages, err := mockRepo.GetUserMessages(ctx, "", 20)
		assert.Error(t, err)
		assert.Nil(t, messages)
	})
}

// TestMessageService_MarkAsRead 测试标记单条消息已读
func TestMessageService_MarkAsRead(t *testing.T) {
	mockRepo := &MockMessageRepository{
		MarkAsReadFunc: func(ctx context.Context, id string) error {
			if id == "999" {
				return errors.New("消息不存在")
			}
			return nil
		},
	}

	ctx := context.Background()

	t.Run("成功标记已读", func(t *testing.T) {
		err := mockRepo.MarkAsRead(ctx, "1")
		assert.NoError(t, err)
	})

	t.Run("消息不存在", func(t *testing.T) {
		err := mockRepo.MarkAsRead(ctx, "999")
		assert.Error(t, err)
	})
}

// TestMessageService_MarkConversationRead 测试标记会话已读
func TestMessageService_MarkConversationRead(t *testing.T) {
	mockRepo := &MockMessageRepository{
		MarkConversationReadFunc: func(ctx context.Context, userID1, userID2 string) error {
			if userID1 == "" || userID2 == "" {
				return errors.New("用户ID不能为空")
			}
			return nil
		},
	}

	ctx := context.Background()

	t.Run("成功标记会话已读", func(t *testing.T) {
		err := mockRepo.MarkConversationRead(ctx, "1", "2")
		assert.NoError(t, err)
	})

	t.Run("用户ID为空", func(t *testing.T) {
		err := mockRepo.MarkConversationRead(ctx, "", "2")
		assert.Error(t, err)
	})
}

// TestMessageService_GetUnread 测试获取未读消息
func TestMessageService_GetUnread(t *testing.T) {
	mockRepo := &MockMessageRepository{
		GetUnreadFunc: func(ctx context.Context, userID string) ([]message.Message, error) {
			if userID == "" {
				return nil, errors.New("用户ID不能为空")
			}
			return testutil.MockMessages(3), nil
		},
	}

	ctx := context.Background()

	t.Run("成功获取未读消息", func(t *testing.T) {
		messages, err := mockRepo.GetUnread(ctx, "1")
		assert.NoError(t, err)
		assert.Len(t, messages, 3)
	})

	t.Run("用户ID为空", func(t *testing.T) {
		messages, err := mockRepo.GetUnread(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, messages)
	})
}

// TestMessageService_CountUnread 测试统计未读消息数
func TestMessageService_CountUnread(t *testing.T) {
	mockRepo := &MockMessageRepository{
		CountUnreadFunc: func(ctx context.Context, userID string) (int64, error) {
			if userID == "" {
				return 0, errors.New("用户ID不能为空")
			}
			return 8, nil
		},
	}

	ctx := context.Background()

	t.Run("成功统计未读消息", func(t *testing.T) {
		count, err := mockRepo.CountUnread(ctx, "1")
		assert.NoError(t, err)
		assert.Equal(t, int64(8), count)
	})

	t.Run("用户ID为空", func(t *testing.T) {
		count, err := mockRepo.CountUnread(ctx, "")
		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
	})
}

// TestMessageService_Delete 测试删除消息
func TestMessageService_Delete(t *testing.T) {
	mockRepo := &MockMessageRepository{
		DeleteFunc: func(ctx context.Context, id string) error {
			if id == "999" {
				return errors.New("消息不存在")
			}
			return nil
		},
	}

	ctx := context.Background()

	t.Run("成功删除消息", func(t *testing.T) {
		err := mockRepo.Delete(ctx, "1")
		assert.NoError(t, err)
	})

	t.Run("消息不存在", func(t *testing.T) {
		err := mockRepo.Delete(ctx, "999")
		assert.Error(t, err)
	})
}
