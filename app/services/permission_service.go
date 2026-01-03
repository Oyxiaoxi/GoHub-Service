package services

import (
	"fmt"

	"GoHub-Service/app/models/permission"
	"GoHub-Service/app/repositories"
	"GoHub-Service/pkg/mapper"
)

// PermissionService 权限业务逻辑
type PermissionService struct {
	repo   repositories.PermissionRepository
	mapper mapper.Mapper[permission.Permission, PermissionResponseDTO]
}

// NewPermissionService 创建权限服务
func NewPermissionService() *PermissionService {
	converter := func(p *permission.Permission) *PermissionResponseDTO {
		return &PermissionResponseDTO{
			ID:          p.ID,
			Name:        p.Name,
			DisplayName: p.DisplayName,
			Description: p.Description,
			Group:       p.Group,
			CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &PermissionService{
		repo:   repositories.NewPermissionRepository(),
		mapper: mapper.NewSimpleMapper(converter),
	}
}

// PermissionCreateDTO 创建权限请求DTO
type PermissionCreateDTO struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	DisplayName string `json:"display_name" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"max=255"`
	Group       string `json:"group" binding:"max=50"`
}

// PermissionUpdateDTO 更新权限请求DTO
type PermissionUpdateDTO struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=100"`
	DisplayName string `json:"display_name" binding:"omitempty,min=1,max=100"`
	Description string `json:"description" binding:"max=255"`
	Group       string `json:"group" binding:"max=50"`
}

// PermissionResponseDTO 权限响应DTO
type PermissionResponseDTO struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Group       string `json:"group"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// toPermissionResponseDTO 使用Mapper转换为响应DTO
func (s *PermissionService) toPermissionResponseDTO(p *permission.Permission) *PermissionResponseDTO {
	return s.mapper.ToDTO(p)
}

// toPermissionResponseDTOList 批量转换
func (s *PermissionService) toPermissionResponseDTOList(perms []permission.Permission) []PermissionResponseDTO {
	return s.mapper.ToDTOList(perms)
}

// CreatePermission 创建权限
func (s *PermissionService) CreatePermission(dto PermissionCreateDTO) (*PermissionResponseDTO, error) {
	// 检查权限是否存在
	existing, _ := s.repo.GetByName(dto.Name)
	if existing != nil {
		return nil, fmt.Errorf("权限已存在")
	}

	newPerm := &permission.Permission{
		Name:        dto.Name,
		DisplayName: dto.DisplayName,
		Description: dto.Description,
		Group:       dto.Group,
	}

	if err := s.repo.Create(newPerm); err != nil {
		return nil, fmt.Errorf("创建权限失败: %v", err)
	}

	return s.toPermissionResponseDTO(newPerm), nil
}

// GetPermissionByID 根据ID获取权限
func (s *PermissionService) GetPermissionByID(id uint64) (*PermissionResponseDTO, error) {
	perm, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("权限不存在")
	}

	return s.toPermissionResponseDTO(perm), nil
}

// GetAllPermissions 获取所有权限
func (s *PermissionService) GetAllPermissions() ([]PermissionResponseDTO, error) {
	perms, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("获取权限列表失败: %v", err)
	}

	return s.toPermissionResponseDTOList(perms), nil
}

// GetPermissionsPaginated 分页获取权限
func (s *PermissionService) GetPermissionsPaginated(page, perPage int) ([]PermissionResponseDTO, int64, error) {
	perms, count, err := s.repo.Paginate(page, perPage)
	if err != nil {
		return nil, 0, fmt.Errorf("获取权限列表失败: %v", err)
	}

	return s.toPermissionResponseDTOList(perms), count, nil
}

// GetPermissionsByGroup 按分组获取权限
func (s *PermissionService) GetPermissionsByGroup(group string) ([]PermissionResponseDTO, error) {
	perms, err := s.repo.GetByGroup(group)
	if err != nil {
		return nil, fmt.Errorf("获取权限列表失败: %v", err)
	}

	return s.toPermissionResponseDTOList(perms), nil
}

// UpdatePermission 更新权限
func (s *PermissionService) UpdatePermission(id uint64, dto PermissionUpdateDTO) (*PermissionResponseDTO, error) {
	perm, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("权限不存在")
	}

	// 检查新名称是否被占用
	if dto.Name != "" && dto.Name != perm.Name {
		existing, _ := s.repo.GetByName(dto.Name)
		if existing != nil {
			return nil, fmt.Errorf("权限名称已被使用")
		}
		perm.Name = dto.Name
	}

	if dto.DisplayName != "" {
		perm.DisplayName = dto.DisplayName
	}
	if dto.Description != "" {
		perm.Description = dto.Description
	}
	if dto.Group != "" {
		perm.Group = dto.Group
	}

	if err := s.repo.Update(perm); err != nil {
		return nil, fmt.Errorf("更新权限失败: %v", err)
	}

	return s.toPermissionResponseDTO(perm), nil
}

// DeletePermission 删除权限
func (s *PermissionService) DeletePermission(id uint64) error {
	// 检查权限是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("权限不存在")
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("删除权限失败: %v", err)
	}

	return nil
}
