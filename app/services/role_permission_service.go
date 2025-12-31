package services

import (
	"GoHub-Service/app/repositories"
	"fmt"
)

// RolePermissionService 角色权限关系业务逻辑
type RolePermissionService struct {
	repo repositories.RolePermissionRepository
}

// NewRolePermissionService 创建角色权限服务
func NewRolePermissionService() *RolePermissionService {
	return &RolePermissionService{
		repo: repositories.NewRolePermissionRepository(),
	}
}

// RolePermissionDTO 角色权限响应DTO
type RolePermissionDTO struct {
	RoleID       uint64 `json:"role_id"`
	PermissionID uint64 `json:"permission_id"`
}

// PermissionDetailDTO 权限详情DTO（用于获取角色权限时返回）
type PermissionDetailDTO struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Group       string `json:"group"`
}

// AssignPermissions 为角色分配权限
func (s *RolePermissionService) AssignPermissions(roleID uint64, permissionIDs []uint64) error {
	// 验证权限是否存在
	permSvc := NewPermissionService()
	for _, permID := range permissionIDs {
		if _, err := permSvc.GetPermissionByID(permID); err != nil {
			return fmt.Errorf("权限ID %d 不存在", permID)
		}
	}

	// 分配权限
	if err := s.repo.AssignMultiple(roleID, permissionIDs); err != nil {
		return fmt.Errorf("分配权限失败: %v", err)
	}

	return nil
}

// RevokePermission 撤销角色的单个权限
func (s *RolePermissionService) RevokePermission(roleID, permissionID uint64) error {
	if err := s.repo.Revoke(roleID, permissionID); err != nil {
		return fmt.Errorf("撤销权限失败: %v", err)
	}
	return nil
}

// RevokeAllPermissions 撤销角色的所有权限
func (s *RolePermissionService) RevokeAllPermissions(roleID uint64) error {
	if err := s.repo.RevokeAll(roleID); err != nil {
		return fmt.Errorf("撤销权限失败: %v", err)
	}
	return nil
}

// GetRolePermissions 获取角色的权限ID列表
func (s *RolePermissionService) GetRolePermissions(roleID uint64) ([]uint64, error) {
	permissionIDs, err := s.repo.GetPermissionsByRoleID(roleID)
	if err != nil {
		return nil, fmt.Errorf("获取角色权限失败: %v", err)
	}
	return permissionIDs, nil
}

// GetRolePermissionDetails 获取角色的权限详情
func (s *RolePermissionService) GetRolePermissionDetails(roleID uint64) ([]PermissionDetailDTO, error) {
	permissionIDs, err := s.repo.GetPermissionsByRoleID(roleID)
	if err != nil {
		return nil, fmt.Errorf("获取角色权限失败: %v", err)
	}

	// 获取权限详情
	permSvc := NewPermissionService()
	details := make([]PermissionDetailDTO, 0, len(permissionIDs))

	for _, permID := range permissionIDs {
		perm, err := permSvc.GetPermissionByID(permID)
		if err != nil {
			continue
		}

		details = append(details, PermissionDetailDTO{
			ID:          perm.ID,
			Name:        perm.Name,
			DisplayName: perm.DisplayName,
			Description: perm.Description,
			Group:       perm.Group,
		})
	}

	return details, nil
}

// HasPermission 检查角色是否拥有权限
func (s *RolePermissionService) HasPermission(roleID, permissionID uint64) (bool, error) {
	has, err := s.repo.HasPermission(roleID, permissionID)
	if err != nil {
		return false, fmt.Errorf("检查权限失败: %v", err)
	}
	return has, nil
}

// CountRolePermissions 获取角色权限总数
func (s *RolePermissionService) CountRolePermissions(roleID uint64) (int64, error) {
	count, err := s.repo.Count(roleID)
	if err != nil {
		return 0, fmt.Errorf("获取权限数量失败: %v", err)
	}
	return count, nil
}
