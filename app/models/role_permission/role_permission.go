// Package role_permission 角色权限关系模型
package role_permission

import (
	"GoHub-Service/app/models"
)

// RolePermission 角色权限关系
type RolePermission struct {
	models.BaseModel
	RoleID       uint64 `gorm:"not null;index:idx_role_permission" json:"role_id"`
	PermissionID uint64 `gorm:"not null;index:idx_role_permission" json:"permission_id"`
	models.CommonTimestampsField
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}
