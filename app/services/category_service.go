// Package services 分类业务逻辑服务
package services

import (
	"GoHub-Service/app/models/category"
	"GoHub-Service/app/repositories"
	"GoHub-Service/pkg/paginator"

	"github.com/gin-gonic/gin"
)

// CategoryService 分类服务
type CategoryService struct{
	repo repositories.CategoryRepository
}

// NewCategoryService 创建分类服务实例
func NewCategoryService() *CategoryService {
	return &CategoryService{
		repo: repositories.NewCategoryRepository(),
	}
}

// CategoryCreateDTO 创建分类数据传输对象
type CategoryCreateDTO struct {
	Name        string
	Description string
}

// CategoryUpdateDTO 更新分类数据传输对象
type CategoryUpdateDTO struct {
	Name        string
	Description string
}

// GetByID 根据ID获取分类
func (s *CategoryService) GetByID(id string) (*category.Category, error) {
	return s.repo.GetByID(id)
}

// List 获取分类列表
func (s *CategoryService) List(c *gin.Context, perPage int) ([]category.Category, *paginator.Paging, error) {
	return s.repo.List(c, perPage)
}

// Create 创建分类
func (s *CategoryService) Create(dto CategoryCreateDTO) (*category.Category, error) {
	categoryModel := &category.Category{
		Name:        dto.Name,
		Description: dto.Description,
	}

	if err := s.repo.Create(categoryModel); err != nil {
		return nil, err
	}

	return categoryModel, nil
}

// Update 更新分类
func (s *CategoryService) Update(id string, dto CategoryUpdateDTO) (*category.Category, error) {
	categoryModel, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	categoryModel.Name = dto.Name
	categoryModel.Description = dto.Description

	if err := s.repo.Update(categoryModel); err != nil {
		return nil, err
	}

	return categoryModel, nil
}

// Delete 删除分类
func (s *CategoryService) Delete(id string) error {
	return s.repo.Delete(id)
}
