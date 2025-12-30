package migrations

import (
	"database/sql"
	"GoHub-Service/app/models/permission"
	"GoHub-Service/app/models/role"
	"GoHub-Service/app/models/role_permission"
	"GoHub-Service/app/models/user_role"
	"GoHub-Service/pkg/migrate"

	"gorm.io/gorm"
)

func init() {

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		// Create roles table
		migrator.AutoMigrate(&role.Role{})
		
		// Create permissions table
		migrator.AutoMigrate(&permission.Permission{})
		
		// Create role_permissions table
		migrator.AutoMigrate(&role_permission.RolePermission{})
		
		// Create user_roles table
		migrator.AutoMigrate(&user_role.UserRole{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		// Delete tables in reverse order of creation
		migrator.DropTable(&user_role.UserRole{})
		migrator.DropTable(&role_permission.RolePermission{})
		migrator.DropTable(&permission.Permission{})
		migrator.DropTable(&role.Role{})
	}

	migrate.Add("2025_12_30_230000_add_rbac_tables", up, down)
}
