// Package role 角色模型
package role

import (
	"GoHub-Service/app/models"
	"GoHub-Service/pkg/database"
)

type Role struct {
	models.BaseModel
	Name        string `gorm:"type:varchar(50);not null;uniqueIndex" json:"name"`
	DisplayName string `gorm:"type:varchar(100);not null" json:"display_name"`
	Description string `gorm:"type:varchar(255)" json:"description"`
	models.CommonTimestampsField
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// Create 创建角色
func (r *Role) Create() {
	database.DB.Create(&r)
}

// Save 保存角色
func (r *Role) Save() (rowsAffected int64) {
	result := database.DB.Save(&r)
	return result.RowsAffected
}

// Delete 删除角色
func (r *Role) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&r)
	return result.RowsAffected
}

// GetByName 根据名称获取角色
func GetByName(name string) Role {
	var role Role
	database.DB.Where("name = ?", name).First(&role)
	return role
}

// Get 根据ID获取角色
func Get(id string) Role {
	var role Role
	database.DB.Where("id = ?", id).First(&role)
	return role
}

// All 获取所有角色
func All() []Role {
	var roles []Role
	database.DB.Find(&roles)
	return roles
}

// IsExist 检查角色名是否存在
func IsExist(name string) bool {
	var count int64
	database.DB.Model(&Role{}).Where("name = ?", name).Count(&count)
	return count > 0
}
