// Package services 分类业务逻辑服务
package services

import (
	"GoHub-Service/app/models/category"
	apperrors "GoHub-Service/pkg/errors"
	"GoHub-Service/pkg/paginator"

	"github.com/gin-gonic/gin"
)

// CategoryService 分类服务
type CategoryService struct{}

// NewCategoryService 创建分类服务实例
func NewCategoryService() *CategoryService {
	return &CategoryService{}
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
	categoryModel := category.Get(id)
	if categoryModel.ID == 0 {
		return nil, apperrors.NotFoundError("分类").WithDetails(map[string]interface{}{
			"category_id": id,
		})
	}
	return &categoryModel, nil
}

// List 获取分类列表
func (s *CategoryService) List(c *gin.Context, perPage int) ([]category.Category, *paginator.Paging, error) {
	data, pager := category.Paginate(c, perPage)
	return data, &pager, nil
}

// Create 创建分类
func (s *CategoryService) Create(dto CategoryCreateDTO) (*category.Category, error) {
	categoryModel := category.Category{
		Name:        dto.Name,
		Description: dto.Description,
	}

	categoryModel.Create()
	if categoryModel.ID == 0 {
		return nil, apperrors.DatabaseError("创建分类", nil)
	}

	return &categoryModel, nil
}

// Update 更新分类
func (s *CategoryService) Update(id string, dto CategoryUpdateDTO) (*category.Category, error) {
	categoryModel := category.Get(id)
	if categoryModel.ID == 0 {
		return nil, apperrors.NotFoundError("分类").WithDetails(map[string]interface{}{
			"category_id": id,
		})
	}

	categoryModel.Name = dto.Name
	categoryModel.Description = dto.Description

	rowsAffected := categoryModel.Save()
	if rowsAffected == 0 {
		return nil, apperrors.DatabaseError("更新分类", nil)
	}

	return &categoryModel, nil
}

// Delete 删除分类
func (s *CategoryService) Delete(id string) error {
	categoryModel := category.Get(id)
	if categoryModel.ID == 0 {
		return apperrors.NotFoundError("分类").WithDetails(map[string]interface{}{
			"category_id": id,
		})
	}

	rowsAffected := categoryModel.Delete()
	if rowsAffected == 0 {
		return apperrors.DatabaseError("删除分类", nil)
	}

	return nil
}
