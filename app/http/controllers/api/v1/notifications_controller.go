package v1

import (
	"GoHub-Service/app/services"
	"GoHub-Service/pkg/auth"
	"GoHub-Service/pkg/logger"
	"GoHub-Service/pkg/response"

	"github.com/gin-gonic/gin"
)

// NotificationsController 通知相关接口
type NotificationsController struct {
	BaseAPIController
	service *services.NotificationService
}

// NewNotificationsController 创建实例
func NewNotificationsController() *NotificationsController {
	return &NotificationsController{service: services.NewNotificationService()}
}

// Index 获取当前用户通知列表
// @Summary 获取通知列表
// @Description 分页获取当前用户的通知列表
// @Tags 通知管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param per_page query int false "每页数量" default(15)
// @Success 200 {object} response.Response "成功"
// @Failure 401 {object} response.Response "未授权"
// @Router /notifications [get]
func (ctrl *NotificationsController) Index(c *gin.Context) {
	userID := auth.CurrentUID(c)
	list, paging, err := ctrl.service.List(c, userID, 15)
	if err != nil {
		logger.LogErrorWithContext(c, err, "获取通知列表失败")
		response.ApiError(c, 500, err.Code, err.Message)
		return
	}
	response.JSON(c, gin.H{"data": ctrl.service.ToResponseList(list), "pager": paging})
}

// Read 标记单条已读
// @Summary 标记通知已读
// @Description 标记指定通知为已读状态
// @Tags 通知管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "通知ID"
// @Success 200 {object} response.Response "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "通知不存在"
// @Router /notifications/{id}/read [put]
func (ctrl *NotificationsController) Read(c *gin.Context) {
	userID := auth.CurrentUID(c)
	id := c.Param("id")

	if err := ctrl.service.MarkRead(id, userID); err != nil {
		logger.LogErrorWithContext(c, err, "标记通知已读失败")
		response.ApiError(c, 500, err.Code, err.Message)
		return
	}
	response.Success(c)
}

// ReadAll 全部已读
// @Summary 标记全部通知已读
// @Description 标记当前用户的所有通知为已读状态
// @Tags 通知管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response "成功"
// @Failure 401 {object} response.Response "未授权"
// @Router /notifications/read-all [put]
func (ctrl *NotificationsController) ReadAll(c *gin.Context) {
	userID := auth.CurrentUID(c)

	if err := ctrl.service.MarkAllRead(userID); err != nil {
		logger.LogErrorWithContext(c, err, "标记全部通知已读失败")
		response.ApiError(c, 500, err.Code, err.Message)
		return
	}
	response.Success(c)
}
