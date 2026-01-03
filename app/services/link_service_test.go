// Package services 链接服务测试
package services

import (
	"GoHub-Service/app/models/link"
	"GoHub-Service/pkg/testutil"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockLinkRepository 链接仓储 Mock
type MockLinkRepository struct {
	GetByIDFunc     func(ctx context.Context, id string) (*link.Link, error)
	GetAllFunc      func(ctx context.Context) ([]link.Link, error)
	CreateFunc      func(ctx context.Context, l *link.Link) error
	UpdateFunc      func(ctx context.Context, l *link.Link) error
	DeleteFunc      func(ctx context.Context, id string) error
	UpdateSortFunc  func(ctx context.Context, id string, sort int) error
	GetActivesFunc  func(ctx context.Context) ([]link.Link, error)
}

func (m *MockLinkRepository) GetByID(ctx context.Context, id string) (*link.Link, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return testutil.MockLinkFactory(id, "测试链接", "https://example.com"), nil
}

func (m *MockLinkRepository) GetAll(ctx context.Context) ([]link.Link, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return testutil.MockLinks(5), nil
}

func (m *MockLinkRepository) Create(ctx context.Context, l *link.Link) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, l)
	}
	l.ID = "1"
	return nil
}

func (m *MockLinkRepository) Update(ctx context.Context, l *link.Link) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, l)
	}
	return nil
}

func (m *MockLinkRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockLinkRepository) UpdateSort(ctx context.Context, id string, sort int) error {
	if m.UpdateSortFunc != nil {
		return m.UpdateSortFunc(ctx, id, sort)
	}
	return nil
}

func (m *MockLinkRepository) GetActives(ctx context.Context) ([]link.Link, error) {
	if m.GetActivesFunc != nil {
		return m.GetActivesFunc(ctx)
	}
	return testutil.MockLinks(3), nil
}

// TestLinkService_Create 测试创建链接
func TestLinkService_Create(t *testing.T) {
	mockRepo := &MockLinkRepository{}
	ctx := context.Background()

	tests := []struct {
		name      string
		title     string
		url       string
		setupMock func()
		wantErr   bool
	}{
		{
			name:  "成功创建链接",
			title: "GitHub",
			url:   "https://github.com",
			setupMock: func() {
				mockRepo.CreateFunc = func(ctx context.Context, l *link.Link) error {
					l.ID = "1"
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:  "标题为空",
			title: "",
			url:   "https://github.com",
			setupMock: func() {
				mockRepo.CreateFunc = func(ctx context.Context, l *link.Link) error {
					return errors.New("标题不能为空")
				}
			},
			wantErr: true,
		},
		{
			name:  "URL格式错误",
			title: "测试链接",
			url:   "invalid-url",
			setupMock: func() {
				mockRepo.CreateFunc = func(ctx context.Context, l *link.Link) error {
					return errors.New("URL格式不正确")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			newLink := &link.Link{
				Title: tt.title,
				URL:   tt.url,
			}

			err := mockRepo.Create(ctx, newLink)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, newLink.ID)
			}
		})
	}
}

// TestLinkService_GetAll 测试获取所有链接
func TestLinkService_GetAll(t *testing.T) {
	mockRepo := &MockLinkRepository{
		GetAllFunc: func(ctx context.Context) ([]link.Link, error) {
			return testutil.MockLinks(5), nil
		},
	}

	ctx := context.Background()

	t.Run("成功获取所有链接", func(t *testing.T) {
		links, err := mockRepo.GetAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, links, 5)
	})
}

// TestLinkService_GetActives 测试获取激活的链接
func TestLinkService_GetActives(t *testing.T) {
	mockRepo := &MockLinkRepository{
		GetActivesFunc: func(ctx context.Context) ([]link.Link, error) {
			return testutil.MockLinks(3), nil
		},
	}

	ctx := context.Background()

	t.Run("成功获取激活链接", func(t *testing.T) {
		links, err := mockRepo.GetActives(ctx)
		assert.NoError(t, err)
		assert.Len(t, links, 3)
	})
}

// TestLinkService_Update 测试更新链接
func TestLinkService_Update(t *testing.T) {
	mockRepo := &MockLinkRepository{
		UpdateFunc: func(ctx context.Context, l *link.Link) error {
			if l.ID == "" {
				return errors.New("链接ID不能为空")
			}
			if l.Title == "" {
				return errors.New("标题不能为空")
			}
			return nil
		},
	}

	ctx := context.Background()

	t.Run("成功更新链接", func(t *testing.T) {
		link := testutil.MockLinkFactory("1", "新标题", "https://newurl.com")
		err := mockRepo.Update(ctx, link)
		assert.NoError(t, err)
	})

	t.Run("链接ID为空", func(t *testing.T) {
		link := &link.Link{
			Title: "测试",
			URL:   "https://test.com",
		}
		err := mockRepo.Update(ctx, link)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID不能为空")
	})
}

// TestLinkService_UpdateSort 测试更新排序
func TestLinkService_UpdateSort(t *testing.T) {
	mockRepo := &MockLinkRepository{
		UpdateSortFunc: func(ctx context.Context, id string, sort int) error {
			if id == "999" {
				return errors.New("链接不存在")
			}
			if sort < 0 {
				return errors.New("排序值不能为负数")
			}
			return nil
		},
	}

	ctx := context.Background()

	t.Run("成功更新排序", func(t *testing.T) {
		err := mockRepo.UpdateSort(ctx, "1", 10)
		assert.NoError(t, err)
	})

	t.Run("链接不存在", func(t *testing.T) {
		err := mockRepo.UpdateSort(ctx, "999", 10)
		assert.Error(t, err)
	})

	t.Run("排序值为负数", func(t *testing.T) {
		err := mockRepo.UpdateSort(ctx, "1", -1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不能为负数")
	})
}

// TestLinkService_Delete 测试删除链接
func TestLinkService_Delete(t *testing.T) {
	mockRepo := &MockLinkRepository{
		DeleteFunc: func(ctx context.Context, id string) error {
			if id == "999" {
				return errors.New("链接不存在")
			}
			return nil
		},
	}

	ctx := context.Background()

	t.Run("成功删除链接", func(t *testing.T) {
		err := mockRepo.Delete(ctx, "1")
		assert.NoError(t, err)
	})

	t.Run("链接不存在", func(t *testing.T) {
		err := mockRepo.Delete(ctx, "999")
		assert.Error(t, err)
	})
}

// TestLinkService_GetByID 测试根据ID获取链接
func TestLinkService_GetByID(t *testing.T) {
	mockRepo := &MockLinkRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*link.Link, error) {
			if id == "999" {
				return nil, errors.New("链接不存在")
			}
			return testutil.MockLinkFactory(id, "测试链接", "https://example.com"), nil
		},
	}

	ctx := context.Background()

	t.Run("成功获取链接", func(t *testing.T) {
		link, err := mockRepo.GetByID(ctx, "1")
		assert.NoError(t, err)
		assert.NotNil(t, link)
		assert.Equal(t, "1", link.ID)
	})

	t.Run("链接不存在", func(t *testing.T) {
		link, err := mockRepo.GetByID(ctx, "999")
		assert.Error(t, err)
		assert.Nil(t, link)
	})
}
