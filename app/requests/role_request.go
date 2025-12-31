package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// RoleCreateRequest 创建角色请求
type RoleCreateRequest struct {
	Name        string `json:"name" form:"name"`
	DisplayName string `json:"display_name" form:"display_name"`
	Description string `json:"description" form:"description"`
}

// ValidateRoleCreate 验证创建角色请求
func ValidateRoleCreate(c *gin.Context, req *RoleCreateRequest) bool {
	return Validate(c, req, func(data interface{}, c *gin.Context) map[string][]string {
		return validate(data, govalidator.MapData{
			"name":         []string{"required", "min:1", "max:50"},
			"display_name": []string{"required", "min:1", "max:100"},
			"description":  []string{"max:255"},
		}, govalidator.MapData{
			"name":         []string{"required:角色名称不能为空", "min:角色名称至少1个字符", "max:角色名称最多50个字符"},
			"display_name": []string{"required:显示名称不能为空", "min:显示名称至少1个字符", "max:显示名称最多100个字符"},
			"description":  []string{"max:描述最多255个字符"},
		})
	})
}

// RoleUpdateRequest 更新角色请求
type RoleUpdateRequest struct {
	Name        string `json:"name" form:"name"`
	DisplayName string `json:"display_name" form:"display_name"`
	Description string `json:"description" form:"description"`
}

// ValidateRoleUpdate 验证更新角色请求
func ValidateRoleUpdate(c *gin.Context, req *RoleUpdateRequest) bool {
	return Validate(c, req, func(data interface{}, c *gin.Context) map[string][]string {
		return validate(data, govalidator.MapData{
			"name":         []string{"min:1", "max:50"},
			"display_name": []string{"min:1", "max:100"},
			"description":  []string{"max:255"},
		}, govalidator.MapData{
			"name":         []string{"min:角色名称至少1个字符", "max:角色名称最多50个字符"},
			"display_name": []string{"min:显示名称至少1个字符", "max:显示名称最多100个字符"},
			"description":  []string{"max:描述最多255个字符"},
		})
	})
}

// AssignPermissionsRequest 分配权限请求
type AssignPermissionsRequest struct {
	PermissionIDs []uint64 `json:"permission_ids" form:"permission_ids"`
}

// ValidateAssignPermissions 验证分配权限请求
func ValidateAssignPermissions(c *gin.Context, req *AssignPermissionsRequest) bool {
	return Validate(c, req, func(data interface{}, c *gin.Context) map[string][]string {
		return validate(data, govalidator.MapData{
			"permission_ids": []string{"required"},
		}, govalidator.MapData{
			"permission_ids": []string{"required:权限ID列表不能为空"},
		})
	})
}
