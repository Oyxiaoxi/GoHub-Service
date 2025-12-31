// Package services 分类服务测试
package services

import (
	"GoHub-Service/app/models/category"
	"GoHub-Service/app/repositories"
	"GoHub-Service/pkg/testutil"
	"errors"
	"testing"
)

// MockCategoryRepository 分类仓储Mock
type MockCategoryRepository struct {
	GetByIDFunc func(id string) (*category.Category, error)
	ListFunc    func(c interface{}, perPage int) ([]category.Category, interface{}, error)
	CreateFunc  func(c *category.Category) error
	UpdateFunc  func(c *category.Category) error
	DeleteFunc  func(id string) error
}

func (m *MockCategoryRepository) GetByID(id string) (*category.Category, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return testutil.MockCategoryFactory("1", "测试分类", "测试描述"), nil
}

func (m *MockCategoryRepository) List(c interface{}, perPage int) ([]category.Category, interface{}, error) {
	if m.ListFunc != nil {
		return m.ListFunc(c, perPage)
	}
	return testutil.MockCategories(), nil, nil
}

func (m *MockCategoryRepository) Create(c *category.Category) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(c)
	}
	return nil
}

func (m *MockCategoryRepository) Update(c *category.Category) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(c)
	}
	return nil
}

func (m *MockCategoryRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

// 确保MockCategoryRepository实现了CategoryRepository接口
var _ repositories.CategoryRepository = (*MockCategoryRepository)(nil)

// TestCategoryService_GetByID 测试获取分类
func TestCategoryService_GetByID(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		mockFunc  func(id string) (*category.Category, error)
		wantErr   bool
		checkFunc func(t *testing.T, result *CategoryResponseDTO)
	}{
		{
			name: "成功获取分类",
			id:   "1",
			mockFunc: func(id string) (*category.Category, error) {
				return testutil.MockCategoryFactory("1", "技术", "技术讨论"), nil
			},
			wantErr: false,
			checkFunc: func(t *testing.T, result *CategoryResponseDTO) {
				testutil.AssertNotNil(t, result, "结果不应为nil")
				testutil.AssertEqual(t, "1", result.ID, "分类ID应该匹配")
				testutil.AssertEqual(t, "技术", result.Name, "分类名称应该匹配")
				testutil.AssertEqual(t, "技术讨论", result.Description, "分类描述应该匹配")
			},
		},
		{
			name: "分类不存在",
			id:   "999",
			mockFunc: func(id string) (*category.Category, error) {
				return nil, errors.New("分类不存在")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建Mock仓储
			mockRepo := &MockCategoryRepository{
				GetByIDFunc: tt.mockFunc,
			}

			// 创建服务实例
			service := &CategoryService{repo: mockRepo}

			// 执行测试
			result, err := service.GetByID(tt.id)

			// 验证错误
			if tt.wantErr {
				testutil.AssertNotNil(t, err, "应该返回错误")
			} else {
				testutil.AssertNil(t, err, "不应该返回错误")
				if tt.checkFunc != nil {
					tt.checkFunc(t, result)
				}
			}
		})
	}
}

// TestCategoryService_Create 测试创建分类
func TestCategoryService_Create(t *testing.T) {
	tests := []struct {
		name      string
		dto       CategoryCreateDTO
		mockFunc  func(c *category.Category) error
		wantErr   bool
		checkFunc func(t *testing.T, result *CategoryResponseDTO)
	}{
		{
			name: "成功创建分类",
			dto: CategoryCreateDTO{
				Name:        "新分类",
				Description: "新分类描述",
			},
			mockFunc: func(c *category.Category) error {
				c.ID = "10"
				return nil
			},
			wantErr: false,
			checkFunc: func(t *testing.T, result *CategoryResponseDTO) {
				testutil.AssertNotNil(t, result, "结果不应为nil")
				testutil.AssertEqual(t, "新分类", result.Name, "分类名称应该匹配")
				testutil.AssertEqual(t, "新分类描述", result.Description, "分类描述应该匹配")
			},
		},
		{
			name: "创建失败-数据库错误",
			dto: CategoryCreateDTO{
				Name:        "错误分类",
				Description: "错误描述",
			},
			mockFunc: func(c *category.Category) error {
				return errors.New("数据库连接失败")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockCategoryRepository{
				CreateFunc: tt.mockFunc,
			}

			service := &CategoryService{repo: mockRepo}
			result, err := service.Create(tt.dto)

			if tt.wantErr {
				testutil.AssertNotNil(t, err, "应该返回错误")
			} else {
				testutil.AssertNil(t, err, "不应该返回错误")
				if tt.checkFunc != nil {
					tt.checkFunc(t, result)
				}
			}
		})
	}
}

// TestCategoryService_Update 测试更新分类
func TestCategoryService_Update(t *testing.T) {
	name := "更新后的名称"
	description := "更新后的描述"

	tests := []struct {
		name       string
		id         string
		dto        CategoryUpdateDTO
		getByIDFunc func(id string) (*category.Category, error)
		updateFunc func(c *category.Category) error
		wantErr    bool
		checkFunc  func(t *testing.T, result *CategoryResponseDTO)
	}{
		{
			name: "成功更新分类",
			id:   "1",
			dto: CategoryUpdateDTO{
				Name:        &name,
				Description: &description,
			},
			getByIDFunc: func(id string) (*category.Category, error) {
				return testutil.MockCategoryFactory("1", "旧名称", "旧描述"), nil
			},
			updateFunc: func(c *category.Category) error {
				return nil
			},
			wantErr: false,
			checkFunc: func(t *testing.T, result *CategoryResponseDTO) {
				testutil.AssertNotNil(t, result, "结果不应为nil")
				testutil.AssertEqual(t, "更新后的名称", result.Name, "名称应该已更新")
				testutil.AssertEqual(t, "更新后的描述", result.Description, "描述应该已更新")
			},
		},
		{
			name: "分类不存在",
			id:   "999",
			dto: CategoryUpdateDTO{
				Name: &name,
			},
			getByIDFunc: func(id string) (*category.Category, error) {
				return nil, errors.New("分类不存在")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockCategoryRepository{
				GetByIDFunc: tt.getByIDFunc,
				UpdateFunc:  tt.updateFunc,
			}

			service := &CategoryService{repo: mockRepo}
			result, err := service.Update(tt.id, tt.dto)

			if tt.wantErr {
				testutil.AssertNotNil(t, err, "应该返回错误")
			} else {
				testutil.AssertNil(t, err, "不应该返回错误")
				if tt.checkFunc != nil {
					tt.checkFunc(t, result)
				}
			}
		})
	}
}

// TestCategoryService_Delete 测试删除分类
func TestCategoryService_Delete(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		mockFunc func(id string) error
		wantErr  bool
	}{
		{
			name: "成功删除分类",
			id:   "1",
			mockFunc: func(id string) error {
				return nil
			},
			wantErr: false,
		},
		{
			name: "删除失败-分类不存在",
			id:   "999",
			mockFunc: func(id string) error {
				return errors.New("分类不存在")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockCategoryRepository{
				DeleteFunc: tt.mockFunc,
			}

			service := &CategoryService{repo: mockRepo}
			err := service.Delete(tt.id)

			if tt.wantErr {
				testutil.AssertNotNil(t, err, "应该返回错误")
			} else {
				testutil.AssertNil(t, err, "不应该返回错误")
			}
		})
	}
}

// TestCategoryService_toResponseDTO 测试DTO转换
func TestCategoryService_toResponseDTO(t *testing.T) {
	service := &CategoryService{}
	cat := testutil.MockCategoryFactory("1", "测试分类", "测试描述")

	dto := service.toResponseDTO(cat)

	testutil.AssertNotNil(t, dto, "DTO不应为nil")
	testutil.AssertEqual(t, "1", dto.ID, "ID应该匹配")
	testutil.AssertEqual(t, "测试分类", dto.Name, "名称应该匹配")
	testutil.AssertEqual(t, "测试描述", dto.Description, "描述应该匹配")
}

// TestCategoryService_toResponseDTOList 测试DTO列表转换
func TestCategoryService_toResponseDTOList(t *testing.T) {
	service := &CategoryService{}
	categories := testutil.MockCategories()

	dtos := service.toResponseDTOList(categories)

	testutil.AssertEqual(t, 3, len(dtos), "DTO列表长度应该匹配")
	testutil.AssertEqual(t, "1", dtos[0].ID, "第一个分类ID应该匹配")
	testutil.AssertEqual(t, "技术", dtos[0].Name, "第一个分类名称应该匹配")
}
