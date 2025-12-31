package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// PermissionCreateRequest 创建权限请求
type PermissionCreateRequest struct {
	Name        string `json:"name" form:"name"`
	DisplayName string `json:"display_name" form:"display_name"`
	Description string `json:"description" form:"description"`
	Group       string `json:"group" form:"group"`
}

// ValidatePermissionCreate 验证创建权限请求
func ValidatePermissionCreate(c *gin.Context, req *PermissionCreateRequest) bool {
	return Validate(c, req, func(data interface{}, c *gin.Context) map[string][]string {
		return validate(data, govalidator.MapData{
			"name":         []string{"required", "min:1", "max:100"},
			"display_name": []string{"required", "min:1", "max:100"},
			"description":  []string{"max:255"},
			"group":        []string{"max:50"},
		}, govalidator.MapData{
			"name":         []string{"required:权限名称不能为空", "min:权限名称至少1个字符", "max:权限名称最多100个字符"},
			"display_name": []string{"required:显示名称不能为空", "min:显示名称至少1个字符", "max:显示名称最多100个字符"},
			"description":  []string{"max:描述最多255个字符"},
			"group":        []string{"max:分组最多50个字符"},
		})
	})
}

// PermissionUpdateRequest 更新权限请求
type PermissionUpdateRequest struct {
	Name        string `json:"name" form:"name"`
	DisplayName string `json:"display_name" form:"display_name"`
	Description string `json:"description" form:"description"`
	Group       string `json:"group" form:"group"`
}

// ValidatePermissionUpdate 验证更新权限请求
func ValidatePermissionUpdate(c *gin.Context, req *PermissionUpdateRequest) bool {
	return Validate(c, req, func(data interface{}, c *gin.Context) map[string][]string {
		return validate(data, govalidator.MapData{
			"name":         []string{"min:1", "max:100"},
			"display_name": []string{"min:1", "max:100"},
			"description":  []string{"max:255"},
			"group":        []string{"max:50"},
		}, govalidator.MapData{
			"name":         []string{"min:权限名称至少1个字符", "max:权限名称最多100个字符"},
			"display_name": []string{"min:显示名称至少1个字符", "max:显示名称最多100个字符"},
			"description":  []string{"max:描述最多255个字符"},
			"group":        []string{"max:分组最多50个字符"},
		})
	})
}
