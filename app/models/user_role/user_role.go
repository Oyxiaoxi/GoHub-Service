// Package user_role 用户角色关系模型
package user_role

import (
	"GoHub-Service/app/models"
)

// UserRole 用户角色关系
type UserRole struct {
	models.BaseModel
	UserID uint64 `gorm:"not null;index:idx_user_role" json:"user_id"`
	RoleID uint64 `gorm:"not null;index:idx_user_role" json:"role_id"`
	models.CommonTimestampsField
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}
