// Package permission 权限模型
package permission

import (
	"GoHub-Service/app/models"
	"GoHub-Service/pkg/database"
)

type Permission struct {
	models.BaseModel
	Name        string `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	DisplayName string `gorm:"type:varchar(100);not null" json:"display_name"`
	Description string `gorm:"type:varchar(255)" json:"description"`
	Group       string `gorm:"type:varchar(50)" json:"group"` // 权限分组
	models.CommonTimestampsField
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// Create 创建权限
func (p *Permission) Create() {
	database.DB.Create(&p)
}

// Save 保存权限
func (p *Permission) Save() (rowsAffected int64) {
	result := database.DB.Save(&p)
	return result.RowsAffected
}

// Delete 删除权限
func (p *Permission) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&p)
	return result.RowsAffected
}

// GetByName 根据名称获取权限
func GetByName(name string) Permission {
	var permission Permission
	database.DB.Where("name = ?", name).First(&permission)
	return permission
}

// All 获取所有权限
func All() []Permission {
	var permissions []Permission
	database.DB.Find(&permissions)
	return permissions
}

// GetByGroup 根据分组获取权限
func GetByGroup(group string) []Permission {
	var permissions []Permission
	database.DB.Where("group = ?", group).Find(&permissions)
	return permissions
}
