package repositories

import (
	"GoHub-Service/app/models/permission"
	"GoHub-Service/pkg/database"
)

// PermissionRepository 权限仓储
type PermissionRepository interface {
	// Create 创建权限
	Create(p *permission.Permission) error
	// GetByID 根据ID获取权限
	GetByID(id uint64) (*permission.Permission, error)
	// GetByName 根据名称获取权限
	GetByName(name string) (*permission.Permission, error)
	// GetAll 获取所有权限
	GetAll() ([]permission.Permission, error)
	// GetByGroup 按分组获取权限
	GetByGroup(group string) ([]permission.Permission, error)
	// Update 更新权限
	Update(p *permission.Permission) error
	// Delete 删除权限
	Delete(id uint64) error
	// Count 获取权限总数
	Count() (int64, error)
	// Paginate 分页获取权限
	Paginate(page, perPage int) ([]permission.Permission, int64, error)
}

type PermissionRepositoryImpl struct{}

// NewPermissionRepository 创建权限仓储
func NewPermissionRepository() PermissionRepository {
	return &PermissionRepositoryImpl{}
}

// Create 创建权限
func (p *PermissionRepositoryImpl) Create(perm *permission.Permission) error {
	if result := database.DB.Create(perm); result.Error != nil {
		return result.Error
	}
	return nil
}

// GetByID 根据ID获取权限
func (p *PermissionRepositoryImpl) GetByID(id uint64) (*permission.Permission, error) {
	var perm permission.Permission
	if result := database.DB.Where("id = ?", id).First(&perm); result.Error != nil {
		return nil, result.Error
	}
	return &perm, nil
}

// GetByName 根据名称获取权限
func (p *PermissionRepositoryImpl) GetByName(name string) (*permission.Permission, error) {
	var perm permission.Permission
	if result := database.DB.Where("name = ?", name).First(&perm); result.Error != nil {
		return nil, result.Error
	}
	return &perm, nil
}

// GetAll 获取所有权限
func (p *PermissionRepositoryImpl) GetAll() ([]permission.Permission, error) {
	var permissions []permission.Permission
	if result := database.DB.Find(&permissions); result.Error != nil {
		return nil, result.Error
	}
	return permissions, nil
}

// GetByGroup 按分组获取权限
func (p *PermissionRepositoryImpl) GetByGroup(group string) ([]permission.Permission, error) {
	var permissions []permission.Permission
	if result := database.DB.Where("group = ?", group).Find(&permissions); result.Error != nil {
		return nil, result.Error
	}
	return permissions, nil
}

// Update 更新权限
func (p *PermissionRepositoryImpl) Update(perm *permission.Permission) error {
	if result := database.DB.Save(perm); result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete 删除权限
func (p *PermissionRepositoryImpl) Delete(id uint64) error {
	if result := database.DB.Delete(&permission.Permission{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}

// Count 获取权限总数
func (p *PermissionRepositoryImpl) Count() (int64, error) {
	var count int64
	if result := database.DB.Model(&permission.Permission{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// Paginate 分页获取权限
func (p *PermissionRepositoryImpl) Paginate(page, perPage int) ([]permission.Permission, int64, error) {
	var permissions []permission.Permission
	var count int64

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	db := database.DB
	if result := db.Model(&permission.Permission{}).Count(&count); result.Error != nil {
		return nil, 0, result.Error
	}

	if result := db.Offset(offset).Limit(perPage).Find(&permissions); result.Error != nil {
		return nil, 0, result.Error
	}

	return permissions, count, nil
}
