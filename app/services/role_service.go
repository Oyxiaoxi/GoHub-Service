package services

import (
	"fmt"

	"GoHub-Service/app/models/role"
	"GoHub-Service/app/repositories"
	apperrors "GoHub-Service/pkg/errors"
)

// RoleService 角色业务逻辑
type RoleService struct {
	repo repositories.RoleRepository
}

// NewRoleService 创建角色服务
func NewRoleService() *RoleService {
	return &RoleService{
		repo: repositories.NewRoleRepository(),
	}
}

// RoleCreateDTO 创建角色请求DTO
type RoleCreateDTO struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	DisplayName string `json:"display_name" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"max=255"`
}

// RoleUpdateDTO 更新角色请求DTO
type RoleUpdateDTO struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=50"`
	DisplayName string `json:"display_name" binding:"omitempty,min=1,max=100"`
	Description string `json:"description" binding:"max=255"`
}

// RoleResponseDTO 角色响应DTO
type RoleResponseDTO struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// toRoleResponseDTO 转换为响应DTO
func toRoleResponseDTO(r *role.Role) RoleResponseDTO {
	return RoleResponseDTO{
		ID:          r.ID,
		Name:        r.Name,
		DisplayName: r.DisplayName,
		Description: r.Description,
		CreatedAt:   r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   r.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(dto RoleCreateDTO) (*RoleResponseDTO, error) {
	// 检查角色是否存在
	existingRole, _ := s.repo.GetByName(dto.Name)
	if existingRole != nil {
		return nil, apperrors.ConflictError("角色")
	}

	newRole := &role.Role{
		Name:        dto.Name,
		DisplayName: dto.DisplayName,
		Description: dto.Description,
	}

	if err := s.repo.Create(newRole); err != nil {
		return nil, apperrors.DatabaseCreateError("角色", err)
	}

	resp := toRoleResponseDTO(newRole)
	return &resp, nil
}

// GetRoleByID 根据ID获取角色
func (s *RoleService) GetRoleByID(id uint64) (*RoleResponseDTO, error) {
	role, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperrors.NotFoundErrorWithCode(apperrors.CodeRoleNotFound, "角色")
	}

	resp := toRoleResponseDTO(role)
	return &resp, nil
}

// GetAllRoles 获取所有角色
func (s *RoleService) GetAllRoles() ([]RoleResponseDTO, error) {
	roles, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("获取角色列表失败: %v", err)
	}

	responses := make([]RoleResponseDTO, len(roles))
	for i, r := range roles {
		responses[i] = toRoleResponseDTO(&r)
	}

	return responses, nil
}

// GetRolesPaginated 分页获取角色
func (s *RoleService) GetRolesPaginated(page, perPage int) ([]RoleResponseDTO, int64, error) {
	roles, count, err := s.repo.Paginate(page, perPage)
	if err != nil {
		return nil, 0, fmt.Errorf("获取角色列表失败: %v", err)
	}

	responses := make([]RoleResponseDTO, len(roles))
	for i, r := range roles {
		responses[i] = toRoleResponseDTO(&r)
	}

	return responses, count, nil
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(id uint64, dto RoleUpdateDTO) (*RoleResponseDTO, error) {
	role, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperrors.NotFoundErrorWithCode(apperrors.CodeRoleNotFound, "角色")
	}

	// 检查新名称是否被占用
	if dto.Name != "" && dto.Name != role.Name {
		existing, _ := s.repo.GetByName(dto.Name)
		if existing != nil {
			return nil, apperrors.ConflictError("角色名称")
		}
		role.Name = dto.Name
	}

	if dto.DisplayName != "" {
		role.DisplayName = dto.DisplayName
	}
	if dto.Description != "" {
		role.Description = dto.Description
	}

	if err := s.repo.Update(role); err != nil {
		return nil, apperrors.DatabaseUpdateError("角色", err)
	}

	resp := toRoleResponseDTO(role)
	return &resp, nil
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(id uint64) error {
	// 检查角色是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		return apperrors.NotFoundErrorWithCode(apperrors.CodeRoleNotFound, "角色")
	}

	if err := s.repo.Delete(id); err != nil {
		return apperrors.DatabaseDeleteError("角色", err)
	}

	return nil
}
