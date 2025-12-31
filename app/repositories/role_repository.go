package repositories

import (
	"GoHub-Service/app/models/role"
	"GoHub-Service/pkg/database"
)

// RoleRepository 角色仓储
type RoleRepository interface {
	// Create 创建角色
	Create(r *role.Role) error
	// GetByID 根据ID获取角色
	GetByID(id uint64) (*role.Role, error)
	// GetByName 根据名称获取角色
	GetByName(name string) (*role.Role, error)
	// GetAll 获取所有角色
	GetAll() ([]role.Role, error)
	// Update 更新角色
	Update(r *role.Role) error
	// Delete 删除角色
	Delete(id uint64) error
	// Count 获取角色总数
	Count() (int64, error)
	// Paginate 分页获取角色
	Paginate(page, perPage int) ([]role.Role, int64, error)
}

type RoleRepositoryImpl struct{}

// NewRoleRepository 创建角色仓储
func NewRoleRepository() RoleRepository {
	return &RoleRepositoryImpl{}
}

// Create 创建角色
func (r *RoleRepositoryImpl) Create(role *role.Role) error {
	if result := database.DB.Create(role); result.Error != nil {
		return result.Error
	}
	return nil
}

// GetByID 根据ID获取角色
func (r *RoleRepositoryImpl) GetByID(id uint64) (*role.Role, error) {
	var role role.Role
	if result := database.DB.Where("id = ?", id).First(&role); result.Error != nil {
		return nil, result.Error
	}
	return &role, nil
}

// GetByName 根据名称获取角色
func (r *RoleRepositoryImpl) GetByName(name string) (*role.Role, error) {
	var role role.Role
	if result := database.DB.Where("name = ?", name).First(&role); result.Error != nil {
		return nil, result.Error
	}
	return &role, nil
}

// GetAll 获取所有角色
func (r *RoleRepositoryImpl) GetAll() ([]role.Role, error) {
	var roles []role.Role
	if result := database.DB.Find(&roles); result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

// Update 更新角色
func (r *RoleRepositoryImpl) Update(role *role.Role) error {
	if result := database.DB.Save(role); result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete 删除角色
func (r *RoleRepositoryImpl) Delete(id uint64) error {
	if result := database.DB.Delete(&role.Role{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}

// Count 获取角色总数
func (r *RoleRepositoryImpl) Count() (int64, error) {
	var count int64
	if result := database.DB.Model(&role.Role{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// Paginate 分页获取角色
func (r *RoleRepositoryImpl) Paginate(page, perPage int) ([]role.Role, int64, error) {
	var roles []role.Role
	var count int64

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	db := database.DB
	if result := db.Model(&role.Role{}).Count(&count); result.Error != nil {
		return nil, 0, result.Error
	}

	if result := db.Offset(offset).Limit(perPage).Find(&roles); result.Error != nil {
		return nil, 0, result.Error
	}

	return roles, count, nil
}
