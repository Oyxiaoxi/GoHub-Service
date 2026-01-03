// Package services 通知服务测试
package services

import (
	"GoHub-Service/app/models/notification"
	"GoHub-Service/pkg/testutil"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockNotificationRepository 通知仓储 Mock
type MockNotificationRepository struct {
	GetByIDFunc          func(ctx context.Context, id string) (*notification.Notification, error)
	GetByUserIDFunc      func(ctx context.Context, userID string, limit int) ([]notification.Notification, error)
	CreateFunc           func(ctx context.Context, n *notification.Notification) error
	BatchCreateFunc      func(ctx context.Context, notifications []notification.Notification) error
	MarkAsReadFunc       func(ctx context.Context, id string) error
	MarkAllAsReadFunc    func(ctx context.Context, userID string) error
	DeleteFunc           func(ctx context.Context, id string) error
	CountUnreadFunc      func(ctx context.Context, userID string) (int64, error)
	GetUnreadFunc        func(ctx context.Context, userID string) ([]notification.Notification, error)
}

func (m *MockNotificationRepository) GetByID(ctx context.Context, id string) (*notification.Notification, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return testutil.MockNotificationFactory(id, "测试通知", "1"), nil
}

func (m *MockNotificationRepository) GetByUserID(ctx context.Context, userID string, limit int) ([]notification.Notification, error) {
	if m.GetByUserIDFunc != nil {
		return m.GetByUserIDFunc(ctx, userID, limit)
	}
	return testutil.MockNotifications(limit), nil
}

func (m *MockNotificationRepository) Create(ctx context.Context, n *notification.Notification) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, n)
	}
	n.ID = "1"
	return nil
}

func (m *MockNotificationRepository) BatchCreate(ctx context.Context, notifications []notification.Notification) error {
	if m.BatchCreateFunc != nil {
		return m.BatchCreateFunc(ctx, notifications)
	}
	return nil
}

func (m *MockNotificationRepository) MarkAsRead(ctx context.Context, id string) error {
	if m.MarkAsReadFunc != nil {
		return m.MarkAsReadFunc(ctx, id)
	}
	return nil
}

func (m *MockNotificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	if m.MarkAllAsReadFunc != nil {
		return m.MarkAllAsReadFunc(ctx, userID)
	}
	return nil
}

func (m *MockNotificationRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockNotificationRepository) CountUnread(ctx context.Context, userID string) (int64, error) {
	if m.CountUnreadFunc != nil {
		return m.CountUnreadFunc(ctx, userID)
	}
	return 5, nil
}

func (m *MockNotificationRepository) GetUnread(ctx context.Context, userID string) ([]notification.Notification, error) {
	if m.GetUnreadFunc != nil {
		return m.GetUnreadFunc(ctx, userID)
	}
	return testutil.MockNotifications(5), nil
}

// TestNotificationService_Create 测试创建通知
func TestNotificationService_Create(t *testing.T) {
	mockRepo := &MockNotificationRepository{}
	ctx := context.Background()

	tests := []struct {
		name      string
		content   string
		userID    string
		setupMock func()
		wantErr   bool
	}{
		{
			name:    "成功创建通知",
			content: "您有一条新消息",
			userID:  "1",
			setupMock: func() {
				mockRepo.CreateFunc = func(ctx context.Context, n *notification.Notification) error {
					n.ID = "1"
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:    "通知内容为空",
			content: "",
			userID:  "1",
			setupMock: func() {
				mockRepo.CreateFunc = func(ctx context.Context, n *notification.Notification) error {
					return errors.New("通知内容不能为空")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			newNotification := &notification.Notification{
				Content: tt.content,
				UserID:  tt.userID,
				Type:    "system",
			}

			err := mockRepo.Create(ctx, newNotification)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, newNotification.ID)
			}
		})
	}
}

// TestNotificationService_BatchCreate 测试批量创建通知
func TestNotificationService_BatchCreate(t *testing.T) {
	mockRepo := &MockNotificationRepository{
		BatchCreateFunc: func(ctx context.Context, notifications []notification.Notification) error {
			if len(notifications) == 0 {
				return errors.New("通知列表不能为空")
			}
			if len(notifications) > 50 {
				return errors.New("批量创建数量不能超过50条")
			}
			return nil
		},
	}

	ctx := context.Background()

	t.Run("成功批量创建", func(t *testing.T) {
		notifications := testutil.MockNotifications(20)
		err := mockRepo.BatchCreate(ctx, notifications)
		assert.NoError(t, err)
	})

	t.Run("通知列表为空", func(t *testing.T) {
		err := mockRepo.BatchCreate(ctx, []notification.Notification{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不能为空")
	})

	t.Run("超过批量限制", func(t *testing.T) {
		notifications := testutil.MockNotifications(51)
		err := mockRepo.BatchCreate(ctx, notifications)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不能超过50")
	})
}

// TestNotificationService_GetUnread 测试获取未读通知
func TestNotificationService_GetUnread(t *testing.T) {
	mockRepo := &MockNotificationRepository{
		GetUnreadFunc: func(ctx context.Context, userID string) ([]notification.Notification, error) {
			if userID == "" {
				return nil, errors.New("用户ID不能为空")
			}
			return testutil.MockNotifications(3), nil
		},
	}

	ctx := context.Background()

	t.Run("成功获取未读通知", func(t *testing.T) {
		notifications, err := mockRepo.GetUnread(ctx, "1")
		assert.NoError(t, err)
		assert.Len(t, notifications, 3)
	})

	t.Run("用户ID为空", func(t *testing.T) {
		notifications, err := mockRepo.GetUnread(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, notifications)
	})
}

// TestNotificationService_MarkAsRead 测试标记已读
func TestNotificationService_MarkAsRead(t *testing.T) {
	mockRepo := &MockNotificationRepository{
		MarkAsReadFunc: func(ctx context.Context, id string) error {
			if id == "999" {
				return errors.New("通知不存在")
			}
			return nil
		},
	}

	ctx := context.Background()

	t.Run("成功标记已读", func(t *testing.T) {
		err := mockRepo.MarkAsRead(ctx, "1")
		assert.NoError(t, err)
	})

	t.Run("通知不存在", func(t *testing.T) {
		err := mockRepo.MarkAsRead(ctx, "999")
		assert.Error(t, err)
	})
}

// TestNotificationService_MarkAllAsRead 测试全部标记已读
func TestNotificationService_MarkAllAsRead(t *testing.T) {
	mockRepo := &MockNotificationRepository{
		MarkAllAsReadFunc: func(ctx context.Context, userID string) error {
			if userID == "" {
				return errors.New("用户ID不能为空")
			}
			return nil
		},
	}

	ctx := context.Background()

	t.Run("成功全部标记已读", func(t *testing.T) {
		err := mockRepo.MarkAllAsRead(ctx, "1")
		assert.NoError(t, err)
	})

	t.Run("用户ID为空", func(t *testing.T) {
		err := mockRepo.MarkAllAsRead(ctx, "")
		assert.Error(t, err)
	})
}

// TestNotificationService_CountUnread 测试统计未读数量
func TestNotificationService_CountUnread(t *testing.T) {
	mockRepo := &MockNotificationRepository{
		CountUnreadFunc: func(ctx context.Context, userID string) (int64, error) {
			if userID == "" {
				return 0, errors.New("用户ID不能为空")
			}
			return 8, nil
		},
	}

	ctx := context.Background()

	t.Run("成功统计未读数量", func(t *testing.T) {
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
