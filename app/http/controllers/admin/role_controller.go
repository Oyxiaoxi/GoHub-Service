package admin

import (
	"GoHub-Service/app/requests"
	"GoHub-Service/app/services"
	"GoHub-Service/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// RoleController 角色管理控制器
type RoleController struct{}

// Index 角色列表
// @Summary 获取角色列表
// @Description 分页获取所有角色
// @Tags Role
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param per_page query int false "每页数量" default(20)
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/roles [get]
func (ctrl *RoleController) Index(c *gin.Context) {
	page := cast.ToInt(c.DefaultQuery("page", "1"))
	perPage := cast.ToInt(c.DefaultQuery("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	svc := services.NewRoleService()
	roles, total, err := svc.GetRolesPaginated(page, perPage)
	if err != nil {
		response.ApiError(c, http.StatusInternalServerError, response.CodeServerError, err.Error())
		return
	}

	response.ApiSuccess(c, gin.H{
		"data": roles,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"per_page":  perPage,
			"last_page": (total + int64(perPage) - 1) / int64(perPage),
		},
	})
}

// Store 创建角色
// @Summary 创建角色
// @Description 创建新的角色
// @Tags Role
// @Accept json
// @Produce json
// @Param request body requests.RoleCreateRequest true "角色信息"
// @Success 201 {object} map[string]interface{}
// @Router /api/v1/admin/roles [post]
func (ctrl *RoleController) Store(c *gin.Context) {
	var req requests.RoleCreateRequest
	if !requests.ValidateRoleCreate(c, &req) {
		return
	}

	svc := services.NewRoleService()
	dto := services.RoleCreateDTO{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
	}

	role, err := svc.CreateRole(dto)
	if err != nil {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, err.Error())
		return
	}

	response.ApiResponse(c, http.StatusCreated, response.CodeSuccess, "角色创建成功", role)
}

// Show 角色详情
// @Summary 获取角色详情
// @Description 根据ID获取角色详情
// @Tags Role
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/roles/{id} [get]
func (ctrl *RoleController) Show(c *gin.Context) {
	id := cast.ToUint64(c.Param("id"))
	if id == 0 {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, "角色ID无效")
		return
	}

	svc := services.NewRoleService()
	role, err := svc.GetRoleByID(id)
	if err != nil {
		response.ApiError(c, http.StatusNotFound, response.CodeNotFound, err.Error())
		return
	}

	response.ApiSuccess(c, role)
}

// Update 更新角色
// @Summary 更新角色
// @Description 更新指定角色的信息
// @Tags Role
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param request body requests.RoleUpdateRequest true "角色信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/roles/{id} [put]
func (ctrl *RoleController) Update(c *gin.Context) {
	id := cast.ToUint64(c.Param("id"))
	if id == 0 {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, "角色ID无效")
		return
	}

	var req requests.RoleUpdateRequest
	if !requests.ValidateRoleUpdate(c, &req) {
		return
	}

	svc := services.NewRoleService()
	dto := services.RoleUpdateDTO{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
	}

	role, err := svc.UpdateRole(id, dto)
	if err != nil {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, err.Error())
		return
	}

	response.ApiSuccess(c, role)
}

// Delete 删除角色
// @Summary 删除角色
// @Description 删除指定的角色
// @Tags Role
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 204
// @Router /api/v1/admin/roles/{id} [delete]
func (ctrl *RoleController) Delete(c *gin.Context) {
	id := cast.ToUint64(c.Param("id"))
	if id == 0 {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, "角色ID无效")
		return
	}

	svc := services.NewRoleService()
	err := svc.DeleteRole(id)
	if err != nil {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// GetPermissions 获取角色权限
// @Summary 获取角色权限
// @Description 获取指定角色拥有的所有权限
// @Tags Role
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/roles/{id}/permissions [get]
func (ctrl *RoleController) GetPermissions(c *gin.Context) {
	id := cast.ToUint64(c.Param("id"))
	if id == 0 {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, "角色ID无效")
		return
	}

	// 验证角色是否存在
	svc := services.NewRoleService()
	_, err := svc.GetRoleByID(id)
	if err != nil {
		response.ApiError(c, http.StatusNotFound, response.CodeNotFound, err.Error())
		return
	}

	// 获取角色权限
	rpRepo := services.NewRolePermissionService()
	permissions, err := rpRepo.GetRolePermissions(id)
	if err != nil {
		response.ApiError(c, http.StatusInternalServerError, response.CodeServerError, err.Error())
		return
	}

	response.ApiSuccess(c, gin.H{
		"role_id": id,
		"permissions": permissions,
	})
}

// AssignPermissions 分配权限
// @Summary 分配权限到角色
// @Description 为指定角色分配权限
// @Tags Role
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param request body requests.AssignPermissionsRequest true "权限ID列表"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/roles/{id}/permissions [post]
func (ctrl *RoleController) AssignPermissions(c *gin.Context) {
	id := cast.ToUint64(c.Param("id"))
	if id == 0 {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, "角色ID无效")
		return
	}

	var req requests.AssignPermissionsRequest
	if !requests.ValidateAssignPermissions(c, &req) {
		return
	}

	// 验证角色是否存在
	svc := services.NewRoleService()
	_, err := svc.GetRoleByID(id)
	if err != nil {
		response.ApiError(c, http.StatusNotFound, response.CodeNotFound, err.Error())
		return
	}

	// 分配权限
	rpSvc := services.NewRolePermissionService()
	err = rpSvc.AssignPermissions(id, req.PermissionIDs)
	if err != nil {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, err.Error())
		return
	}

	response.ApiSuccessWithMessage(c, "权限分配成功")
}
