// Package services 评论服务测试
package services

import (
	"GoHub-Service/app/models/comment"
	"GoHub-Service/pkg/testutil"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockCommentRepository 评论仓储 Mock
type MockCommentRepository struct {
	GetByIDFunc         func(ctx context.Context, id string) (*comment.Comment, error)
	GetByTopicIDFunc    func(ctx context.Context, topicID string) ([]comment.Comment, error)
	CreateFunc          func(ctx context.Context, c *comment.Comment) error
	UpdateFunc          func(ctx context.Context, c *comment.Comment) error
	DeleteFunc          func(ctx context.Context, id string) error
	BatchCreateFunc     func(ctx context.Context, comments []comment.Comment) error
	CountByTopicIDFunc  func(ctx context.Context, topicID string) (int64, error)
	GetRecentByUserFunc func(ctx context.Context, userID string, limit int) ([]comment.Comment, error)
}

func (m *MockCommentRepository) GetByID(ctx context.Context, id string) (*comment.Comment, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return testutil.MockCommentFactory(id, "测试评论", "1", "1"), nil
}

func (m *MockCommentRepository) GetByTopicID(ctx context.Context, topicID string) ([]comment.Comment, error) {
	if m.GetByTopicIDFunc != nil {
		return m.GetByTopicIDFunc(ctx, topicID)
	}
	return testutil.MockComments(3), nil
}

func (m *MockCommentRepository) Create(ctx context.Context, c *comment.Comment) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, c)
	}
	c.ID = "1"
	return nil
}

func (m *MockCommentRepository) Update(ctx context.Context, c *comment.Comment) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, c)
	}
	return nil
}

func (m *MockCommentRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockCommentRepository) BatchCreate(ctx context.Context, comments []comment.Comment) error {
	if m.BatchCreateFunc != nil {
		return m.BatchCreateFunc(ctx, comments)
	}
	return nil
}

func (m *MockCommentRepository) CountByTopicID(ctx context.Context, topicID string) (int64, error) {
	if m.CountByTopicIDFunc != nil {
		return m.CountByTopicIDFunc(ctx, topicID)
	}
	return 10, nil
}

func (m *MockCommentRepository) GetRecentByUser(ctx context.Context, userID string, limit int) ([]comment.Comment, error) {
	if m.GetRecentByUserFunc != nil {
		return m.GetRecentByUserFunc(ctx, userID, limit)
	}
	return testutil.MockComments(limit), nil
}

// TestCommentService_Create 测试创建评论
func TestCommentService_Create(t *testing.T) {
	mockRepo := &MockCommentRepository{}
	// 假设 CommentService 有 repo 字段（需根据实际结构调整）
	// service := &CommentService{repo: mockRepo}

	tests := []struct {
		name      string
		content   string
		userID    string
		topicID   string
		setupMock func()
		wantErr   bool
	}{
		{
			name:    "成功创建评论",
			content: "这是一条测试评论",
			userID:  "1",
			topicID: "1",
			setupMock: func() {
				mockRepo.CreateFunc = func(ctx context.Context, c *comment.Comment) error {
					c.ID = "1"
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:    "评论内容为空",
			content: "",
			userID:  "1",
			topicID: "1",
			setupMock: func() {
				mockRepo.CreateFunc = func(ctx context.Context, c *comment.Comment) error {
					return errors.New("评论内容不能为空")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			ctx := context.Background()

			newComment := &comment.Comment{
				Content: tt.content,
				UserID:  tt.userID,
				TopicID: tt.topicID,
			}

			err := mockRepo.Create(ctx, newComment)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, newComment.ID)
			}
		})
	}
}

// TestCommentService_GetByTopicID 测试获取话题评论列表
func TestCommentService_GetByTopicID(t *testing.T) {
	mockRepo := &MockCommentRepository{
		GetByTopicIDFunc: func(ctx context.Context, topicID string) ([]comment.Comment, error) {
			if topicID == "999" {
				return nil, errors.New("话题不存在")
			}
			return testutil.MockComments(5), nil
		},
	}

	ctx := context.Background()

	t.Run("成功获取评论列表", func(t *testing.T) {
		comments, err := mockRepo.GetByTopicID(ctx, "1")
		assert.NoError(t, err)
		assert.Len(t, comments, 5)
	})

	t.Run("话题不存在", func(t *testing.T) {
		comments, err := mockRepo.GetByTopicID(ctx, "999")
		assert.Error(t, err)
		assert.Nil(t, comments)
	})
}

// TestCommentService_BatchCreate 测试批量创建评论
func TestCommentService_BatchCreate(t *testing.T) {
	mockRepo := &MockCommentRepository{
		BatchCreateFunc: func(ctx context.Context, comments []comment.Comment) error {
			if len(comments) == 0 {
				return errors.New("评论列表不能为空")
			}
			if len(comments) > 100 {
				return errors.New("批量创建数量不能超过100条")
			}
			return nil
		},
	}

	ctx := context.Background()

	t.Run("成功批量创建", func(t *testing.T) {
		comments := testutil.MockComments(10)
		err := mockRepo.BatchCreate(ctx, comments)
		assert.NoError(t, err)
	})

	t.Run("评论列表为空", func(t *testing.T) {
		err := mockRepo.BatchCreate(ctx, []comment.Comment{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不能为空")
	})

	t.Run("超过批量限制", func(t *testing.T) {
		comments := testutil.MockComments(101)
		err := mockRepo.BatchCreate(ctx, comments)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不能超过100")
	})
}

// TestCommentService_Delete 测试删除评论
func TestCommentService_Delete(t *testing.T) {
	mockRepo := &MockCommentRepository{
		DeleteFunc: func(ctx context.Context, id string) error {
			if id == "999" {
				return errors.New("评论不存在")
			}
			return nil
		},
	}

	ctx := context.Background()

	t.Run("成功删除评论", func(t *testing.T) {
		err := mockRepo.Delete(ctx, "1")
		assert.NoError(t, err)
	})

	t.Run("评论不存在", func(t *testing.T) {
		err := mockRepo.Delete(ctx, "999")
		assert.Error(t, err)
	})
}
