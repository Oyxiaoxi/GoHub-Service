package repositories

import (
	"GoHub-Service/app/models/category"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockCategoryRepository 用于测试的 Mock 实现
type MockCategoryRepository struct {
	GetByIDFunc     func(id string) (*category.Category, error)
	CreateFunc      func(c *category.Category) error
	UpdateFunc      func(c *category.Category) error
	DeleteFunc      func(id string) error
	BatchCreateFunc func(categories []category.Category) error
	BatchDeleteFunc func(ids []string) error
}

func (m *MockCategoryRepository) GetByID(id string) (*category.Category, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *MockCategoryRepository) List(c interface{}, perPage int) ([]category.Category, interface{}, error) {
	return []category.Category{}, nil, nil
}

func (m *MockCategoryRepository) Create(cat *category.Category) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(cat)
	}
	return nil
}

func (m *MockCategoryRepository) Update(cat *category.Category) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(cat)
	}
	return nil
}

func (m *MockCategoryRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func (m *MockCategoryRepository) BatchCreate(categories []category.Category) error {
	if m.BatchCreateFunc != nil {
		return m.BatchCreateFunc(categories)
	}
	return nil
}

func (m *MockCategoryRepository) BatchDelete(ids []string) error {
	if m.BatchDeleteFunc != nil {
		return m.BatchDeleteFunc(ids)
	}
	return nil
}

func (m *MockCategoryRepository) GetAllCached() ([]category.Category, error) {
	return []category.Category{}, nil
}

func (m *MockCategoryRepository) SetListCache(categories []category.Category) error {
	return nil
}

func (m *MockCategoryRepository) FlushCache() error {
	return nil
}

// TestMockCategoryRepository_GetByID 测试分类 GetByID
func TestMockCategoryRepository_GetByID(t *testing.T) {
	mock := &MockCategoryRepository{
		GetByIDFunc: func(id string) (*category.Category, error) {
			if id == "1" {
				return &category.Category{Name: "Tech"}, nil
			}
			return nil, ErrNotFound
		},
	}

	result, err := mock.GetByID("1")
	assert.NoError(t, err)
	assert.Equal(t, "Tech", result.Name)

	result, err = mock.GetByID("999")
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestMockCategoryRepository_Create 测试分类创建
func TestMockCategoryRepository_Create(t *testing.T) {
	mock := &MockCategoryRepository{
		CreateFunc: func(cat *category.Category) error {
			if cat.Name == "" {
				return ErrCreateFailed
			}
			cat.ID = 1
			return nil
		},
	}

	newCat := &category.Category{Name: "New Category"}
	err := mock.Create(newCat)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), newCat.ID)

	invalidCat := &category.Category{}
	err = mock.Create(invalidCat)
	assert.Error(t, err)
}

// TestMockCategoryRepository_BatchCreate 测试批量创建
func TestMockCategoryRepository_BatchCreate(t *testing.T) {
	mock := &MockCategoryRepository{
		BatchCreateFunc: func(categories []category.Category) error {
			if len(categories) == 0 {
				return ErrCreateFailed
			}
			return nil
		},
	}

	cats := []category.Category{{Name: "Cat1"}, {Name: "Cat2"}}
	err := mock.BatchCreate(cats)
	assert.NoError(t, err)

	err = mock.BatchCreate([]category.Category{})
	assert.Error(t, err)
}
