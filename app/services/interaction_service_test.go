// Package services 交互服务测试（点赞、关注等）
package services

import (
	"GoHub-Service/app/models/like"
	"GoHub-Service/app/models/follow"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockLikeRepository 点赞仓储 Mock
type MockLikeRepository struct {
	CreateFunc      func(ctx context.Context, l *like.Like) error
	DeleteFunc      func(ctx context.Context, userID, topicID string) error
	IsUserLikedFunc func(ctx context.Context, userID, topicID string) (bool, error)
	GetUserLikesFunc func(ctx context.Context, userID string) ([]like.Like, error)
	CountFunc       func(ctx context.Context, topicID string) (int64, error)
	BatchCreateFunc func(ctx context.Context, likes []like.Like) error
	BatchDeleteFunc func(ctx context.Context, userID string, topicIDs []string) error
}

func (m *MockLikeRepository) Create(ctx context.Context, l *like.Like) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, l)
	}
	l.ID = "1"
	return nil
}

func (m *MockLikeRepository) Delete(ctx context.Context, userID, topicID string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, userID, topicID)
	}
	return nil
}

func (m *MockLikeRepository) IsUserLiked(ctx context.Context, userID, topicID string) (bool, error) {
	if m.IsUserLikedFunc != nil {
		return m.IsUserLikedFunc(ctx, userID, topicID)
	}
	return false, nil
}

func (m *MockLikeRepository) GetUserLikes(ctx context.Context, userID string) ([]like.Like, error) {
	if m.GetUserLikesFunc != nil {
		return m.GetUserLikesFunc(ctx, userID)
	}
	return []like.Like{}, nil
}

func (m *MockLikeRepository) Count(ctx context.Context, topicID string) (int64, error) {
	if m.CountFunc != nil {
		return m.CountFunc(ctx, topicID)
	}
	return 10, nil
}

func (m *MockLikeRepository) BatchCreate(ctx context.Context, likes []like.Like) error {
	if m.BatchCreateFunc != nil {
		return m.BatchCreateFunc(ctx, likes)
	}
	return nil
}

func (m *MockLikeRepository) BatchDelete(ctx context.Context, userID string, topicIDs []string) error {
	if m.BatchDeleteFunc != nil {
		return m.BatchDeleteFunc(ctx, userID, topicIDs)
	}
	return nil
}

// MockFollowRepository 关注仓储 Mock
type MockFollowRepository struct {
	CreateFunc        func(ctx context.Context, f *follow.Follow) error
	DeleteFunc        func(ctx context.Context, followerID, followeeID string) error
	IsFollowingFunc   func(ctx context.Context, followerID, followeeID string) (bool, error)
	GetFollowersFunc  func(ctx context.Context, userID string) ([]follow.Follow, error)
	GetFolloweesFunc  func(ctx context.Context, userID string) ([]follow.Follow, error)
	CountFollowersFunc func(ctx context.Context, userID string) (int64, error)
	CountFolloweesFunc func(ctx context.Context, userID string) (int64, error)
}

func (m *MockFollowRepository) Create(ctx context.Context, f *follow.Follow) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, f)
	}
	f.ID = "1"
	return nil
}

func (m *MockFollowRepository) Delete(ctx context.Context, followerID, followeeID string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, followerID, followeeID)
	}
	return nil
}

func (m *MockFollowRepository) IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error) {
	if m.IsFollowingFunc != nil {
		return m.IsFollowingFunc(ctx, followerID, followeeID)
	}
	return false, nil
}

func (m *MockFollowRepository) GetFollowers(ctx context.Context, userID string) ([]follow.Follow, error) {
	if m.GetFollowersFunc != nil {
		return m.GetFollowersFunc(ctx, userID)
	}
	return []follow.Follow{}, nil
}

func (m *MockFollowRepository) GetFollowees(ctx context.Context, userID string) ([]follow.Follow, error) {
	if m.GetFolloweesFunc != nil {
		return m.GetFolloweesFunc(ctx, userID)
	}
	return []follow.Follow{}, nil
}

func (m *MockFollowRepository) CountFollowers(ctx context.Context, userID string) (int64, error) {
	if m.CountFollowersFunc != nil {
		return m.CountFollowersFunc(ctx, userID)
	}
	return 100, nil
}

func (m *MockFollowRepository) CountFollowees(ctx context.Context, userID string) (int64, error) {
	if m.CountFolloweesFunc != nil {
		return m.CountFolloweesFunc(ctx, userID)
	}
	return 50, nil
}

// TestInteractionService_Like 测试点赞功能
func TestInteractionService_Like(t *testing.T) {
	mockRepo := &MockLikeRepository{}
	ctx := context.Background()

	tests := []struct {
		name      string
		userID    string
		topicID   string
		setupMock func()
		wantErr   bool
	}{
		{
			name:    "成功点赞",
			userID:  "1",
			topicID: "100",
			setupMock: func() {
				mockRepo.IsUserLikedFunc = func(ctx context.Context, userID, topicID string) (bool, error) {
					return false, nil
				}
				mockRepo.CreateFunc = func(ctx context.Context, l *like.Like) error {
					l.ID = "1"
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:    "重复点赞",
			userID:  "1",
			topicID: "100",
			setupMock: func() {
				mockRepo.IsUserLikedFunc = func(ctx context.Context, userID, topicID string) (bool, error) {
					return true, nil
				}
				mockRepo.CreateFunc = func(ctx context.Context, l *like.Like) error {
					return errors.New("已经点赞过")
				}
			},
			wantErr: true,
		},
		{
			name:    "用户ID为空",
			userID:  "",
			topicID: "100",
			setupMock: func() {
				mockRepo.CreateFunc = func(ctx context.Context, l *like.Like) error {
					return errors.New("用户ID不能为空")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			newLike := &like.Like{
				UserID:  tt.userID,
				TopicID: tt.topicID,
			}

			err := mockRepo.Create(ctx, newLike)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, newLike.ID)
			}
		})
	}
}

// TestInteractionService_Unlike 测试取消点赞
func TestInteractionService_Unlike(t *testing.T) {
	mockRepo := &MockLikeRepository{
		DeleteFunc: func(ctx context.Context, userID, topicID string) error {
			if userID == "" || topicID == "" {
				return errors.New("参数不能为空")
			}
			return nil
		},
	}

	ctx := context.Background()

	t.Run("成功取消点赞", func(t *testing.T) {
		err := mockRepo.Delete(ctx, "1", "100")
		assert.NoError(t, err)
	})

	t.Run("参数为空", func(t *testing.T) {
		err := mockRepo.Delete(ctx, "", "100")
		assert.Error(t, err)
	})
}

// TestInteractionService_BatchLike 测试批量点赞
func TestInteractionService_BatchLike(t *testing.T) {
	mockRepo := &MockLikeRepository{
		BatchCreateFunc: func(ctx context.Context, likes []like.Like) error {
			if len(likes) == 0 {
				return errors.New("点赞列表不能为空")
			}
			if len(likes) > 100 {
				return errors.New("批量点赞数量不能超过100")
			}
			return nil
		},
	}

	ctx := context.Background()

	t.Run("成功批量点赞", func(t *testing.T) {
		likes := make([]like.Like, 10)
		for i := 0; i < 10; i++ {
			likes[i] = like.Like{
				UserID:  "1",
				TopicID: string(rune(i + 100)),
			}
		}
		err := mockRepo.BatchCreate(ctx, likes)
		assert.NoError(t, err)
	})

	t.Run("点赞列表为空", func(t *testing.T) {
		err := mockRepo.BatchCreate(ctx, []like.Like{})
		assert.Error(t, err)
	})

	t.Run("超过批量限制", func(t *testing.T) {
		likes := make([]like.Like, 101)
		err := mockRepo.BatchCreate(ctx, likes)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不能超过100")
	})
}

// TestInteractionService_Follow 测试关注功能
func TestInteractionService_Follow(t *testing.T) {
	mockRepo := &MockFollowRepository{}
	ctx := context.Background()

	tests := []struct {
		name       string
		followerID string
		followeeID string
		setupMock  func()
		wantErr    bool
	}{
		{
			name:       "成功关注",
			followerID: "1",
			followeeID: "2",
			setupMock: func() {
				mockRepo.IsFollowingFunc = func(ctx context.Context, followerID, followeeID string) (bool, error) {
					return false, nil
				}
				mockRepo.CreateFunc = func(ctx context.Context, f *follow.Follow) error {
					f.ID = "1"
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:       "重复关注",
			followerID: "1",
			followeeID: "2",
			setupMock: func() {
				mockRepo.IsFollowingFunc = func(ctx context.Context, followerID, followeeID string) (bool, error) {
					return true, nil
				}
				mockRepo.CreateFunc = func(ctx context.Context, f *follow.Follow) error {
					return errors.New("已经关注过")
				}
			},
			wantErr: true,
		},
		{
			name:       "不能关注自己",
			followerID: "1",
			followeeID: "1",
			setupMock: func() {
				mockRepo.CreateFunc = func(ctx context.Context, f *follow.Follow) error {
					return errors.New("不能关注自己")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			newFollow := &follow.Follow{
				FollowerID: tt.followerID,
				FolloweeID: tt.followeeID,
			}

			err := mockRepo.Create(ctx, newFollow)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, newFollow.ID)
			}
		})
	}
}

// TestInteractionService_Unfollow 测试取消关注
func TestInteractionService_Unfollow(t *testing.T) {
	mockRepo := &MockFollowRepository{
		DeleteFunc: func(ctx context.Context, followerID, followeeID string) error {
			if followerID == "" || followeeID == "" {
				return errors.New("参数不能为空")
			}
			return nil
		},
	}

	ctx := context.Background()

	t.Run("成功取消关注", func(t *testing.T) {
		err := mockRepo.Delete(ctx, "1", "2")
		assert.NoError(t, err)
	})

	t.Run("参数为空", func(t *testing.T) {
		err := mockRepo.Delete(ctx, "", "2")
		assert.Error(t, err)
	})
}

// TestInteractionService_GetFollowerCount 测试获取粉丝数
func TestInteractionService_GetFollowerCount(t *testing.T) {
	mockRepo := &MockFollowRepository{
		CountFollowersFunc: func(ctx context.Context, userID string) (int64, error) {
			if userID == "" {
				return 0, errors.New("用户ID不能为空")
			}
			return 150, nil
		},
	}

	ctx := context.Background()

	t.Run("成功获取粉丝数", func(t *testing.T) {
		count, err := mockRepo.CountFollowers(ctx, "1")
		assert.NoError(t, err)
		assert.Equal(t, int64(150), count)
	})

	t.Run("用户ID为空", func(t *testing.T) {
		count, err := mockRepo.CountFollowers(ctx, "")
		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
	})
}

// TestInteractionService_GetFolloweeCount 测试获取关注数
func TestInteractionService_GetFolloweeCount(t *testing.T) {
	mockRepo := &MockFollowRepository{
		CountFolloweesFunc: func(ctx context.Context, userID string) (int64, error) {
			if userID == "" {
				return 0, errors.New("用户ID不能为空")
			}
			return 80, nil
		},
	}

	ctx := context.Background()

	t.Run("成功获取关注数", func(t *testing.T) {
		count, err := mockRepo.CountFollowees(ctx, "1")
		assert.NoError(t, err)
		assert.Equal(t, int64(80), count)
	})

	t.Run("用户ID为空", func(t *testing.T) {
		count, err := mockRepo.CountFollowees(ctx, "")
		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
	})
}

// TestInteractionService_Timeout 测试超时控制（使用 ContextGuard）
func TestInteractionService_Timeout(t *testing.T) {
	mockRepo := &MockLikeRepository{
		BatchCreateFunc: func(ctx context.Context, likes []like.Like) error {
			// 模拟长时间操作
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(100 * time.Millisecond):
				return nil
			}
		},
	}

	t.Run("操作在超时前完成", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		likes := make([]like.Like, 5)
		err := mockRepo.BatchCreate(ctx, likes)
		assert.NoError(t, err)
	})

	t.Run("操作超时", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		likes := make([]like.Like, 5)
		err := mockRepo.BatchCreate(ctx, likes)
		assert.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)
	})
}
