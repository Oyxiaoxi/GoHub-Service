package repositories

import (
	"GoHub-Service/app/models/role_permission"
	"GoHub-Service/pkg/database"
)

// RolePermissionRepository 角色权限关系仓储
type RolePermissionRepository interface {
	// Assign 为角色分配权限
	Assign(roleID, permissionID uint64) error
	// Revoke 撤销角色的权限
	Revoke(roleID, permissionID uint64) error
	// GetPermissionsByRoleID 获取角色的所有权限
	GetPermissionsByRoleID(roleID uint64) ([]uint64, error)
	// GetRolesByPermissionID 获取拥有该权限的所有角色
	GetRolesByPermissionID(permissionID uint64) ([]uint64, error)
	// RevokeAll 撤销角色的所有权限
	RevokeAll(roleID uint64) error
	// AssignMultiple 批量分配权限
	AssignMultiple(roleID uint64, permissionIDs []uint64) error
	// HasPermission 检查角色是否有权限
	HasPermission(roleID, permissionID uint64) (bool, error)
	// Count 获取角色权限关系总数
	Count(roleID uint64) (int64, error)
}

type RolePermissionRepositoryImpl struct{}

// NewRolePermissionRepository 创建角色权限仓储
func NewRolePermissionRepository() RolePermissionRepository {
	return &RolePermissionRepositoryImpl{}
}

// Assign 为角色分配权限
func (r *RolePermissionRepositoryImpl) Assign(roleID, permissionID uint64) error {
	rp := &role_permission.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	if result := database.DB.Create(rp); result.Error != nil {
		return result.Error
	}
	return nil
}

// Revoke 撤销角色的权限
func (r *RolePermissionRepositoryImpl) Revoke(roleID, permissionID uint64) error {
	if result := database.DB.Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&role_permission.RolePermission{}); result.Error != nil {
		return result.Error
	}
	return nil
}

// GetPermissionsByRoleID 获取角色的所有权限
func (r *RolePermissionRepositoryImpl) GetPermissionsByRoleID(roleID uint64) ([]uint64, error) {
	var permissionIDs []uint64
	if result := database.DB.Model(&role_permission.RolePermission{}).
		Where("role_id = ?", roleID).
		Pluck("permission_id", &permissionIDs); result.Error != nil {
		return nil, result.Error
	}
	return permissionIDs, nil
}

// GetRolesByPermissionID 获取拥有该权限的所有角色
func (r *RolePermissionRepositoryImpl) GetRolesByPermissionID(permissionID uint64) ([]uint64, error) {
	var roleIDs []uint64
	if result := database.DB.Model(&role_permission.RolePermission{}).
		Where("permission_id = ?", permissionID).
		Pluck("role_id", &roleIDs); result.Error != nil {
		return nil, result.Error
	}
	return roleIDs, nil
}

// RevokeAll 撤销角色的所有权限
func (r *RolePermissionRepositoryImpl) RevokeAll(roleID uint64) error {
	if result := database.DB.Where("role_id = ?", roleID).Delete(&role_permission.RolePermission{}); result.Error != nil {
		return result.Error
	}
	return nil
}

// AssignMultiple 批量分配权限
func (r *RolePermissionRepositoryImpl) AssignMultiple(roleID uint64, permissionIDs []uint64) error {
	// 先删除旧的权限
	if err := r.RevokeAll(roleID); err != nil {
		return err
	}

	// 批量插入新权限
	if len(permissionIDs) == 0 {
		return nil
	}

	rps := make([]role_permission.RolePermission, len(permissionIDs))
	for i, permissionID := range permissionIDs {
		rps[i] = role_permission.RolePermission{
			RoleID:       roleID,
			PermissionID: permissionID,
		}
	}

	if result := database.DB.CreateInBatches(rps, 100); result.Error != nil {
		return result.Error
	}
	return nil
}

// HasPermission 检查角色是否有权限
func (r *RolePermissionRepositoryImpl) HasPermission(roleID, permissionID uint64) (bool, error) {
	var count int64
	if result := database.DB.Model(&role_permission.RolePermission{}).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Count(&count); result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// Count 获取角色权限关系总数
func (r *RolePermissionRepositoryImpl) Count(roleID uint64) (int64, error) {
	var count int64
	if result := database.DB.Model(&role_permission.RolePermission{}).
		Where("role_id = ?", roleID).
		Count(&count); result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
