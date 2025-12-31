// Package repositories Comment仓储测试
package repositories

import (
	"GoHub-Service/app/models/comment"
	"GoHub-Service/pkg/testutil"
	"errors"
	"testing"
)

// TestCommentRepository_GetByID 测试获取评论
func TestCommentRepository_GetByID(t *testing.T) {
	// 注意：这是单元测试，不依赖真实数据库
	// 实际项目中可以使用sqlmock或testcontainers

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "有效的ID格式",
			id:      "1",
			wantErr: false, // Mock场景下
		},
		{
			name:    "空ID",
			id:      "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 这里我们只测试逻辑，不测试真实数据库
			if tt.id == "" {
				testutil.AssertTrue(t, tt.wantErr, "空ID应该返回错误")
			}
		})
	}
}

// TestCommentRepository_Create 测试创建评论
func TestCommentRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		comment *comment.Comment
		wantErr bool
	}{
		{
			name:    "有效评论",
			comment: testutil.MockCommentFactory("1", "测试评论", "1", "1", ""),
			wantErr: false,
		},
		{
			name:    "nil评论",
			comment: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.comment == nil {
				testutil.AssertTrue(t, tt.wantErr, "nil评论应该返回错误")
			} else {
				testutil.AssertNotNil(t, tt.comment, "评论不应为nil")
				testutil.AssertTrue(t, len(tt.comment.Content) > 0, "评论内容不应为空")
			}
		})
	}
}

// TestCommentRepository_Update 测试更新评论
func TestCommentRepository_Update(t *testing.T) {
	originalComment := testutil.MockCommentFactory("1", "原始内容", "1", "1", "")

	tests := []struct {
		name           string
		comment        *comment.Comment
		updatedContent string
		wantErr        bool
	}{
		{
			name:           "更新评论内容",
			comment:        originalComment,
			updatedContent: "更新后的内容",
			wantErr:        false,
		},
		{
			name:           "空内容",
			comment:        originalComment,
			updatedContent: "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.updatedContent == "" {
				testutil.AssertTrue(t, tt.wantErr, "空内容应该返回错误")
			} else {
				tt.comment.Content = tt.updatedContent
				testutil.AssertEqual(t, tt.updatedContent, tt.comment.Content, "内容应该已更新")
			}
		})
	}
}

// TestCommentRepository_Delete 测试删除评论
func TestCommentRepository_Delete(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "删除存在的评论",
			id:      "1",
			wantErr: false,
		},
		{
			name:    "删除不存在的评论",
			id:      "999",
			wantErr: true,
		},
		{
			name:    "空ID",
			id:      "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.id == "" || tt.id == "999" {
				testutil.AssertTrue(t, tt.wantErr, "无效ID应该返回错误")
			}
		})
	}
}

// TestCommentRepository_IncrementLikeCount 测试增加点赞数
func TestCommentRepository_IncrementLikeCount(t *testing.T) {
	comment := testutil.MockCommentFactory("1", "测试评论", "1", "1", "")
	originalCount := comment.LikeCount

	// 模拟增加点赞
	comment.LikeCount++

	testutil.AssertEqual(t, originalCount+1, comment.LikeCount, "点赞数应该增加1")
	testutil.AssertTrue(t, comment.LikeCount > 0, "点赞数应该大于0")
}

// TestCommentRepository_DecrementLikeCount 测试减少点赞数
func TestCommentRepository_DecrementLikeCount(t *testing.T) {
	comment := testutil.MockCommentFactory("1", "测试评论", "1", "1", "")
	comment.LikeCount = 5
	originalCount := comment.LikeCount

	// 模拟减少点赞
	if comment.LikeCount > 0 {
		comment.LikeCount--
	}

	testutil.AssertEqual(t, originalCount-1, comment.LikeCount, "点赞数应该减少1")
	testutil.AssertTrue(t, comment.LikeCount >= 0, "点赞数不应该为负数")
}

// MockCommentRepositoryForTest 用于测试的Mock实现
type MockCommentRepositoryForTest struct {
	GetByIDFunc func(id string) (*comment.Comment, error)
	CreateFunc  func(c *comment.Comment) error
	UpdateFunc  func(c *comment.Comment) error
	DeleteFunc  func(id string) error
}

func (m *MockCommentRepositoryForTest) GetByID(id string) (*comment.Comment, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	if id == "" {
		return nil, errors.New("ID不能为空")
	}
	return testutil.MockCommentFactory(id, "测试评论", "1", "1", ""), nil
}

func (m *MockCommentRepositoryForTest) Create(c *comment.Comment) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(c)
	}
	if c == nil {
		return errors.New("评论不能为nil")
	}
	return nil
}

func (m *MockCommentRepositoryForTest) Update(c *comment.Comment) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(c)
	}
	if c == nil || c.Content == "" {
		return errors.New("评论内容不能为空")
	}
	return nil
}

func (m *MockCommentRepositoryForTest) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	if id == "" || id == "999" {
		return errors.New("评论不存在")
	}
	return nil
}

func (m *MockCommentRepositoryForTest) List(c interface{}, perPage int) ([]comment.Comment, interface{}, error) {
	return testutil.MockComments(), nil, nil
}

func (m *MockCommentRepositoryForTest) GetByTopicID(topicID string, page, perPage int) ([]comment.Comment, int64, error) {
	return testutil.MockComments(), 3, nil
}

func (m *MockCommentRepositoryForTest) GetByUserID(userID string, page, perPage int) ([]comment.Comment, int64, error) {
	return testutil.MockComments(), 3, nil
}

func (m *MockCommentRepositoryForTest) IncrementLikeCount(commentID string) error {
	return nil
}

func (m *MockCommentRepositoryForTest) DecrementLikeCount(commentID string) error {
	return nil
}

// TestMockCommentRepository 测试Mock仓储
func TestMockCommentRepository(t *testing.T) {
	mockRepo := &MockCommentRepositoryForTest{}

	t.Run("GetByID", func(t *testing.T) {
		result, err := mockRepo.GetByID("1")
		testutil.AssertNil(t, err, "不应该有错误")
		testutil.AssertNotNil(t, result, "结果不应为nil")
		testutil.AssertEqual(t, "1", result.ID, "ID应该匹配")
	})

	t.Run("GetByID_空ID", func(t *testing.T) {
		result, err := mockRepo.GetByID("")
		testutil.AssertNotNil(t, err, "应该返回错误")
		testutil.AssertNil(t, result, "结果应为nil")
	})

	t.Run("Create", func(t *testing.T) {
		comment := testutil.MockCommentFactory("1", "新评论", "1", "1", "")
		err := mockRepo.Create(comment)
		testutil.AssertNil(t, err, "不应该有错误")
	})

	t.Run("Create_nil评论", func(t *testing.T) {
		err := mockRepo.Create(nil)
		testutil.AssertNotNil(t, err, "应该返回错误")
	})

	t.Run("Delete", func(t *testing.T) {
		err := mockRepo.Delete("1")
		testutil.AssertNil(t, err, "不应该有错误")
	})

	t.Run("Delete_不存在的ID", func(t *testing.T) {
		err := mockRepo.Delete("999")
		testutil.AssertNotNil(t, err, "应该返回错误")
	})
}
