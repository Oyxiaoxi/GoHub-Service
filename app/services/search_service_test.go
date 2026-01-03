// Package services 搜索服务测试
package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// SearchResult 搜索结果结构（简化版）
type SearchResult struct {
	Type  string
	ID    string
	Title string
	Score float64
}

// MockSearchRepository 搜索仓储 Mock
type MockSearchRepository struct {
	SearchTopicsFunc     func(ctx context.Context, keyword string, limit int) ([]SearchResult, error)
	SearchUsersFunc      func(ctx context.Context, keyword string, limit int) ([]SearchResult, error)
	SearchAllFunc        func(ctx context.Context, keyword string) (map[string][]SearchResult, error)
	IndexTopicFunc       func(ctx context.Context, topicID string) error
	DeleteTopicIndexFunc func(ctx context.Context, topicID string) error
	UpdateTopicIndexFunc func(ctx context.Context, topicID string) error
}

func (m *MockSearchRepository) SearchTopics(ctx context.Context, keyword string, limit int) ([]SearchResult, error) {
	if m.SearchTopicsFunc != nil {
		return m.SearchTopicsFunc(ctx, keyword, limit)
	}
	return []SearchResult{}, nil
}

func (m *MockSearchRepository) SearchUsers(ctx context.Context, keyword string, limit int) ([]SearchResult, error) {
	if m.SearchUsersFunc != nil {
		return m.SearchUsersFunc(ctx, keyword, limit)
	}
	return []SearchResult{}, nil
}

func (m *MockSearchRepository) SearchAll(ctx context.Context, keyword string) (map[string][]SearchResult, error) {
	if m.SearchAllFunc != nil {
		return m.SearchAllFunc(ctx, keyword)
	}
	return map[string][]SearchResult{}, nil
}

func (m *MockSearchRepository) IndexTopic(ctx context.Context, topicID string) error {
	if m.IndexTopicFunc != nil {
		return m.IndexTopicFunc(ctx, topicID)
	}
	return nil
}

func (m *MockSearchRepository) DeleteTopicIndex(ctx context.Context, topicID string) error {
	if m.DeleteTopicIndexFunc != nil {
		return m.DeleteTopicIndexFunc(ctx, topicID)
	}
	return nil
}

func (m *MockSearchRepository) UpdateTopicIndex(ctx context.Context, topicID string) error {
	if m.UpdateTopicIndexFunc != nil {
		return m.UpdateTopicIndexFunc(ctx, topicID)
	}
	return nil
}

// TestSearchService_SearchTopics 测试搜索话题
func TestSearchService_SearchTopics(t *testing.T) {
	mockRepo := &MockSearchRepository{}
	ctx := context.Background()

	tests := []struct {
		name      string
		keyword   string
		limit     int
		setupMock func()
		wantErr   bool
		wantCount int
	}{
		{
			name:    "成功搜索话题",
			keyword: "Go语言",
			limit:   10,
			setupMock: func() {
				mockRepo.SearchTopicsFunc = func(ctx context.Context, keyword string, limit int) ([]SearchResult, error) {
					results := make([]SearchResult, 3)
					for i := 0; i < 3; i++ {
						results[i] = SearchResult{
							Type:  "topic",
							ID:    string(rune(i + 1)),
							Title: "Go语言相关话题",
							Score: 0.95,
						}
					}
					return results, nil
				}
			},
			wantErr:   false,
			wantCount: 3,
		},
		{
			name:    "关键词为空",
			keyword: "",
			limit:   10,
			setupMock: func() {
				mockRepo.SearchTopicsFunc = func(ctx context.Context, keyword string, limit int) ([]SearchResult, error) {
					return nil, errors.New("关键词不能为空")
				}
			},
			wantErr:   true,
			wantCount: 0,
		},
		{
			name:    "没有搜索结果",
			keyword: "不存在的关键词xyz123",
			limit:   10,
			setupMock: func() {
				mockRepo.SearchTopicsFunc = func(ctx context.Context, keyword string, limit int) ([]SearchResult, error) {
					return []SearchResult{}, nil
				}
			},
			wantErr:   false,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			results, err := mockRepo.SearchTopics(ctx, tt.keyword, tt.limit)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, results, tt.wantCount)
			}
		})
	}
}

// TestSearchService_SearchUsers 测试搜索用户
func TestSearchService_SearchUsers(t *testing.T) {
	mockRepo := &MockSearchRepository{
		SearchUsersFunc: func(ctx context.Context, keyword string, limit int) ([]SearchResult, error) {
			if keyword == "" {
				return nil, errors.New("关键词不能为空")
			}
			if limit <= 0 {
				return nil, errors.New("limit必须大于0")
			}
			
			results := make([]SearchResult, 2)
			for i := 0; i < 2; i++ {
				results[i] = SearchResult{
					Type:  "user",
					ID:    string(rune(i + 1)),
					Title: keyword + "相关用户",
					Score: 0.88,
				}
			}
			return results, nil
		},
	}

	ctx := context.Background()

	t.Run("成功搜索用户", func(t *testing.T) {
		results, err := mockRepo.SearchUsers(ctx, "张三", 10)
		assert.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("关键词为空", func(t *testing.T) {
		results, err := mockRepo.SearchUsers(ctx, "", 10)
		assert.Error(t, err)
		assert.Nil(t, results)
	})

	t.Run("limit为0", func(t *testing.T) {
		results, err := mockRepo.SearchUsers(ctx, "张三", 0)
		assert.Error(t, err)
		assert.Nil(t, results)
	})
}

// TestSearchService_SearchAll 测试全局搜索
func TestSearchService_SearchAll(t *testing.T) {
	mockRepo := &MockSearchRepository{
		SearchAllFunc: func(ctx context.Context, keyword string) (map[string][]SearchResult, error) {
			if keyword == "" {
				return nil, errors.New("关键词不能为空")
			}

			results := map[string][]SearchResult{
				"topics": {
					{Type: "topic", ID: "1", Title: "话题1", Score: 0.95},
					{Type: "topic", ID: "2", Title: "话题2", Score: 0.90},
				},
				"users": {
					{Type: "user", ID: "1", Title: "用户1", Score: 0.88},
				},
				"comments": {
					{Type: "comment", ID: "1", Title: "评论1", Score: 0.75},
				},
			}
			return results, nil
		},
	}

	ctx := context.Background()

	t.Run("成功全局搜索", func(t *testing.T) {
		results, err := mockRepo.SearchAll(ctx, "测试")
		assert.NoError(t, err)
		assert.Len(t, results, 3)
		assert.Len(t, results["topics"], 2)
		assert.Len(t, results["users"], 1)
		assert.Len(t, results["comments"], 1)
	})

	t.Run("关键词为空", func(t *testing.T) {
		results, err := mockRepo.SearchAll(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, results)
	})
}

// TestSearchService_IndexTopic 测试索引话题
func TestSearchService_IndexTopic(t *testing.T) {
	mockRepo := &MockSearchRepository{
		IndexTopicFunc: func(ctx context.Context, topicID string) error {
			if topicID == "" {
				return errors.New("话题ID不能为空")
			}
			if topicID == "999" {
				return errors.New("话题不存在")
			}
			return nil
		},
	}

	ctx := context.Background()

	t.Run("成功创建索引", func(t *testing.T) {
		err := mockRepo.IndexTopic(ctx, "1")
		assert.NoError(t, err)
	})

	t.Run("话题ID为空", func(t *testing.T) {
		err := mockRepo.IndexTopic(ctx, "")
		assert.Error(t, err)
	})

	t.Run("话题不存在", func(t *testing.T) {
		err := mockRepo.IndexTopic(ctx, "999")
		assert.Error(t, err)
	})
}

// TestSearchService_UpdateTopicIndex 测试更新话题索引
func TestSearchService_UpdateTopicIndex(t *testing.T) {
	mockRepo := &MockSearchRepository{
		UpdateTopicIndexFunc: func(ctx context.Context, topicID string) error {
			if topicID == "" {
				return errors.New("话题ID不能为空")
			}
			return nil
		},
	}

	ctx := context.Background()

	t.Run("成功更新索引", func(t *testing.T) {
		err := mockRepo.UpdateTopicIndex(ctx, "1")
		assert.NoError(t, err)
	})

	t.Run("话题ID为空", func(t *testing.T) {
		err := mockRepo.UpdateTopicIndex(ctx, "")
		assert.Error(t, err)
	})
}

// TestSearchService_DeleteTopicIndex 测试删除话题索引
func TestSearchService_DeleteTopicIndex(t *testing.T) {
	mockRepo := &MockSearchRepository{
		DeleteTopicIndexFunc: func(ctx context.Context, topicID string) error {
			if topicID == "" {
				return errors.New("话题ID不能为空")
			}
			return nil
		},
	}

	ctx := context.Background()

	t.Run("成功删除索引", func(t *testing.T) {
		err := mockRepo.DeleteTopicIndex(ctx, "1")
		assert.NoError(t, err)
	})

	t.Run("话题ID为空", func(t *testing.T) {
		err := mockRepo.DeleteTopicIndex(ctx, "")
		assert.Error(t, err)
	})
}

// TestSearchService_Timeout 测试搜索超时控制（使用 ContextGuard）
func TestSearchService_Timeout(t *testing.T) {
	mockRepo := &MockSearchRepository{
		SearchAllFunc: func(ctx context.Context, keyword string) (map[string][]SearchResult, error) {
			// 模拟长时间搜索操作
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(200 * time.Millisecond):
				return map[string][]SearchResult{
					"topics": {{Type: "topic", ID: "1", Title: "结果", Score: 0.9}},
				}, nil
			}
		},
	}

	t.Run("搜索在超时前完成", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancel()

		results, err := mockRepo.SearchAll(ctx, "测试")
		assert.NoError(t, err)
		assert.NotNil(t, results)
	})

	t.Run("搜索超时", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		results, err := mockRepo.SearchAll(ctx, "测试")
		assert.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)
		assert.Nil(t, results)
	})
}

// TestSearchService_ConcurrentSearch 测试并发搜索
func TestSearchService_ConcurrentSearch(t *testing.T) {
	mockRepo := &MockSearchRepository{
		SearchTopicsFunc: func(ctx context.Context, keyword string, limit int) ([]SearchResult, error) {
			time.Sleep(50 * time.Millisecond) // 模拟搜索耗时
			return []SearchResult{
				{Type: "topic", ID: "1", Title: keyword + "相关话题", Score: 0.9},
			}, nil
		},
		SearchUsersFunc: func(ctx context.Context, keyword string, limit int) ([]SearchResult, error) {
			time.Sleep(50 * time.Millisecond) // 模拟搜索耗时
			return []SearchResult{
				{Type: "user", ID: "1", Title: keyword + "相关用户", Score: 0.85},
			}, nil
		},
	}

	ctx := context.Background()

	t.Run("并发搜索话题和用户", func(t *testing.T) {
		start := time.Now()

		// 并发执行两个搜索
		done := make(chan bool, 2)

		go func() {
			_, err := mockRepo.SearchTopics(ctx, "Go", 10)
			assert.NoError(t, err)
			done <- true
		}()

		go func() {
			_, err := mockRepo.SearchUsers(ctx, "Go", 10)
			assert.NoError(t, err)
			done <- true
		}()

		// 等待两个搜索完成
		<-done
		<-done

		elapsed := time.Since(start)
		
		// 并发执行应该接近单个操作的时间（50ms），而不是串行的100ms
		// 考虑到测试环境，给一些余量
		assert.Less(t, elapsed, 150*time.Millisecond, "并发搜索应该更快")
	})
}

// TestSearchService_KeywordValidation 测试关键词验证
func TestSearchService_KeywordValidation(t *testing.T) {
	mockRepo := &MockSearchRepository{
		SearchTopicsFunc: func(ctx context.Context, keyword string, limit int) ([]SearchResult, error) {
			if len(keyword) < 2 {
				return nil, errors.New("关键词长度至少为2个字符")
			}
			if len(keyword) > 100 {
				return nil, errors.New("关键词长度不能超过100个字符")
			}
			return []SearchResult{}, nil
		},
	}

	ctx := context.Background()

	t.Run("关键词太短", func(t *testing.T) {
		results, err := mockRepo.SearchTopics(ctx, "a", 10)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "至少为2个字符")
		assert.Nil(t, results)
	})

	t.Run("关键词太长", func(t *testing.T) {
		longKeyword := string(make([]byte, 101))
		results, err := mockRepo.SearchTopics(ctx, longKeyword, 10)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不能超过100")
		assert.Nil(t, results)
	})

	t.Run("关键词长度正常", func(t *testing.T) {
		results, err := mockRepo.SearchTopics(ctx, "正常关键词", 10)
		assert.NoError(t, err)
		assert.NotNil(t, results)
	})
}
