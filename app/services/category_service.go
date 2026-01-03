// Package services 分类业务逻辑服务
package services

import (
	"time"

	"GoHub-Service/app/models/category"
	"GoHub-Service/app/repositories"
	apperrors "GoHub-Service/pkg/errors"
	"GoHub-Service/pkg/mapper"
	"GoHub-Service/pkg/paginator"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CategoryService 负责分类的读写和缓存穿透保护等业务流.
type CategoryService struct {
	repo   repositories.CategoryRepository
	mapper mapper.Mapper[category.Category, CategoryResponseDTO] // 使用泛型Mapper消除DTO转换重复
	logger *zap.Logger
}

// NewCategoryService 构造分类服务，默认使用数据库+缓存仓储实现.
func NewCategoryService() *CategoryService {
	// 定义DTO转换函数（只需一次）
	converter := func(c *category.Category) *CategoryResponseDTO {
		return &CategoryResponseDTO{
			ID:          c.GetStringID(),
			Name:        c.Name,
			Description: c.Description,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
		}
	}

	return &CategoryService{
		repo:   repositories.NewCategoryRepository(),
		mapper: mapper.NewSimpleMapper(converter),
		logger: zap.L(),
	}
}

// CategoryCreateDTO 创建分类数据传输对象
type CategoryCreateDTO struct {
	Name        string `json:"name" binding:"required,min=2,max=255"`
	Description string `json:"description"`
}

// CategoryUpdateDTO 更新分类数据传输对象
type CategoryUpdateDTO struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=2,max=255"`
	Description *string `json:"description,omitempty"`
}

// CategoryResponseDTO 分类响应DTO
type CategoryResponseDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CategoryListResponseDTO 分类列表响应DTO
type CategoryListResponseDTO struct {
	Categories []CategoryResponseDTO `json:"categories"`
	Paging     *paginator.Paging     `json:"paging"`
}

// toResponseDTO 使用Mapper将Category模型转换为响应DTO
// 优化：使用泛型Mapper消除重复代码
func (s *CategoryService) toResponseDTO(c *category.Category) *CategoryResponseDTO {
	return s.mapper.ToDTO(c)
}

// toResponseDTOList 使用Mapper将Category模型列表转换为响应DTO列表
// 优化：使用泛型Mapper消除重复代码，自动优化内存拷贝
func (s *CategoryService) toResponseDTOList(categories []category.Category) []CategoryResponseDTO {
	return s.mapper.ToDTOList(categories)
}

// GetByID 拉取单条分类，包装仓储错误便于向上层透传.
func (s *CategoryService) GetByID(id string) (*CategoryResponseDTO, error) {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperrors.WrapError(err, "获取分类失败")
	}
	return s.toResponseDTO(c), nil
}

// List 分页查询分类列表并返回分页信息.
func (s *CategoryService) List(c *gin.Context, perPage int) (*CategoryListResponseDTO, error) {
	categories, paging, err := s.repo.List(c, perPage)
	if err != nil {
		return nil, apperrors.WrapError(err, "获取分类列表失败")
	}
	return &CategoryListResponseDTO{
		Categories: s.toResponseDTOList(categories),
		Paging:     paging,
	}, nil
}

// Create 创建分类
func (s *CategoryService) Create(dto CategoryCreateDTO) (*CategoryResponseDTO, error) {
	categoryModel := &category.Category{
		Name:        dto.Name,
		Description: dto.Description,
	}

	if err := s.repo.Create(categoryModel); err != nil {
		return nil, apperrors.WrapError(err, "创建分类失败")
	}

	return s.toResponseDTO(categoryModel), nil
}

// Update 更新分类
func (s *CategoryService) Update(id string, dto CategoryUpdateDTO) (*CategoryResponseDTO, error) {
	categoryModel, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperrors.WrapError(err, "获取分类失败")
	}

	// 只更新非空字段
	if dto.Name != nil {
		categoryModel.Name = *dto.Name
	}
	if dto.Description != nil {
		categoryModel.Description = *dto.Description
	}

	if err := s.repo.Update(categoryModel); err != nil {
		return nil, apperrors.WrapError(err, "更新分类失败")
	}

	return s.toResponseDTO(categoryModel), nil
}

// Delete 删除分类
func (s *CategoryService) Delete(id string) error {
	return s.repo.Delete(id)
}
