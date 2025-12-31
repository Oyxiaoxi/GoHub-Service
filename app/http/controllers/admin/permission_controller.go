package admin

import (
	"GoHub-Service/app/requests"
	"GoHub-Service/app/services"
	"GoHub-Service/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// PermissionController 权限管理控制器
type PermissionController struct{}

// Index 权限列表
// @Summary 获取权限列表
// @Description 分页获取所有权限
// @Tags Permission
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param per_page query int false "每页数量" default(20)
// @Param group query string false "权限分组"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/permissions [get]
func (ctrl *PermissionController) Index(c *gin.Context) {
	page := cast.ToInt(c.DefaultQuery("page", "1"))
	perPage := cast.ToInt(c.DefaultQuery("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	svc := services.NewPermissionService()

	// 如果指定了分组，则获取指定分组的权限
	if group := c.Query("group"); group != "" {
		permissions, err := svc.GetPermissionsByGroup(group)
		if err != nil {
			response.ApiError(c, http.StatusInternalServerError, response.CodeServerError, err.Error())
			return
		}

		response.ApiSuccess(c, gin.H{
			"data": permissions,
			"pagination": gin.H{
				"total":     int64(len(permissions)),
				"page":      page,
				"per_page":  perPage,
				"last_page": 1,
			},
		})
		return
	}

	// 获取所有权限（分页）
	permissions, total, err := svc.GetPermissionsPaginated(page, perPage)
	if err != nil {
		response.ApiError(c, http.StatusInternalServerError, response.CodeServerError, err.Error())
		return
	}

	response.ApiSuccess(c, gin.H{
		"data": permissions,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"per_page":  perPage,
			"last_page": (total + int64(perPage) - 1) / int64(perPage),
		},
	})
}

// Store 创建权限
// @Summary 创建权限
// @Description 创建新的权限
// @Tags Permission
// @Accept json
// @Produce json
// @Param request body requests.PermissionCreateRequest true "权限信息"
// @Success 201 {object} map[string]interface{}
// @Router /api/v1/admin/permissions [post]
func (ctrl *PermissionController) Store(c *gin.Context) {
	var req requests.PermissionCreateRequest
	if !requests.ValidatePermissionCreate(c, &req) {
		return
	}

	svc := services.NewPermissionService()
	dto := services.PermissionCreateDTO{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Group:       req.Group,
	}

	permission, err := svc.CreatePermission(dto)
	if err != nil {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, err.Error())
		return
	}

	response.ApiResponse(c, http.StatusCreated, response.CodeSuccess, "权限创建成功", permission)
}

// Show 权限详情
// @Summary 获取权限详情
// @Description 根据ID获取权限详情
// @Tags Permission
// @Accept json
// @Produce json
// @Param id path int true "权限ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/permissions/{id} [get]
func (ctrl *PermissionController) Show(c *gin.Context) {
	id := cast.ToUint64(c.Param("id"))
	if id == 0 {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, "权限ID无效")
		return
	}

	svc := services.NewPermissionService()
	permission, err := svc.GetPermissionByID(id)
	if err != nil {
		response.ApiError(c, http.StatusNotFound, response.CodeNotFound, err.Error())
		return
	}

	response.ApiSuccess(c, permission)
}

// Update 更新权限
// @Summary 更新权限
// @Description 更新指定权限的信息
// @Tags Permission
// @Accept json
// @Produce json
// @Param id path int true "权限ID"
// @Param request body requests.PermissionUpdateRequest true "权限信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/permissions/{id} [put]
func (ctrl *PermissionController) Update(c *gin.Context) {
	id := cast.ToUint64(c.Param("id"))
	if id == 0 {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, "权限ID无效")
		return
	}

	var req requests.PermissionUpdateRequest
	if !requests.ValidatePermissionUpdate(c, &req) {
		return
	}

	svc := services.NewPermissionService()
	dto := services.PermissionUpdateDTO{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Group:       req.Group,
	}

	permission, err := svc.UpdatePermission(id, dto)
	if err != nil {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, err.Error())
		return
	}

	response.ApiSuccess(c, permission)
}

// Delete 删除权限
// @Summary 删除权限
// @Description 删除指定的权限
// @Tags Permission
// @Accept json
// @Produce json
// @Param id path int true "权限ID"
// @Success 204
// @Router /api/v1/admin/permissions/{id} [delete]
func (ctrl *PermissionController) Delete(c *gin.Context) {
	id := cast.ToUint64(c.Param("id"))
	if id == 0 {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, "权限ID无效")
		return
	}

	svc := services.NewPermissionService()
	err := svc.DeletePermission(id)
	if err != nil {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
