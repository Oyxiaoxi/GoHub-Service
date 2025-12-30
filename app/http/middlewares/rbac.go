package middlewares

import (
	"GoHub-Service/app/models/role"
	"GoHub-Service/pkg/auth"
	"GoHub-Service/pkg/database"
	"GoHub-Service/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// RequireRole 中间件：要求用户具有指定角色
func RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户ID
		userID := auth.CurrentUID(c)
		if userID == "" {
			response.Abort403(c, "未登录")
			return
		}

		// 检查用户是否拥有该角色
		if !hasRole(userID, roleName) {
			response.Abort403(c, "权限不足")
			return
		}

		c.Next()
	}
}

// RequireAnyRole 中间件：要求用户具有任意一个指定角色
func RequireAnyRole(roleNames ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := auth.CurrentUID(c)
		if userID == "" {
			response.Abort403(c, "未登录")
			return
		}

		for _, roleName := range roleNames {
			if hasRole(userID, roleName) {
				c.Next()
				return
			}
		}

		response.Abort403(c, "权限不足")
	}
}

// RequirePermission 中间件：要求用户具有指定权限
func RequirePermission(permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := auth.CurrentUID(c)
		if userID == "" {
			response.Abort403(c, "未登录")
			return
		}

		if !hasPermission(userID, permissionName) {
			response.Abort403(c, "权限不足")
			return
		}

		c.Next()
	}
}

// RequireAdmin 中间件：要求用户是管理员
func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin")
}

// hasRole 检查用户是否拥有指定角色
func hasRole(userID, roleName string) bool {
	var count int64
	database.DB.Table("user_roles").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.name = ?", userID, roleName).
		Count(&count)
	return count > 0
}

// hasPermission 检查用户是否拥有指定权限（通过角色）
func hasPermission(userID, permissionName string) bool {
	var count int64
	database.DB.Table("user_roles").
		Joins("JOIN role_permissions ON role_permissions.role_id = user_roles.role_id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("user_roles.user_id = ? AND permissions.name = ?", userID, permissionName).
		Count(&count)
	return count > 0
}

// GetUserRoles 获取用户的所有角色
func GetUserRoles(userID string) []role.Role {
	var roles []role.Role
	database.DB.Table("roles").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles)
	return roles
}

// AssignRole 给用户分配角色
func AssignRole(userID, roleID string) error {
	type UserRole struct {
		UserID uint64 `gorm:"not null"`
		RoleID uint64 `gorm:"not null"`
	}
	userRole := UserRole{
		UserID: cast.ToUint64(userID),
		RoleID: cast.ToUint64(roleID),
	}
	return database.DB.Table("user_roles").Create(&userRole).Error
}

// RemoveRole 移除用户的角色
func RemoveRole(userID, roleID string) error {
	return database.DB.Table("user_roles").
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(nil).Error
}
