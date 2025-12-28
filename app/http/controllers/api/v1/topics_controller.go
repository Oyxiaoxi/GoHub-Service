package v1

import (
	"GoHub-Service/app/http/middlewares"
	"GoHub-Service/app/requests"
	"GoHub-Service/app/services"
	apperrors "GoHub-Service/pkg/errors"
	"GoHub-Service/pkg/auth"
	"GoHub-Service/pkg/logger"
	"GoHub-Service/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TopicsController struct {
	BaseAPIController
	topicService *services.TopicService
}

// NewTopicsController 创建TopicsController实例
func NewTopicsController() *TopicsController {
	return &TopicsController{
		topicService: services.NewTopicService(),
	}
}

// Index godoc
// @Summary 获取话题列表
// @Description 分页获取所有话题
// @Tags topics
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param per_page query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{} "data: 话题列表, pager: 分页信息"
// @Failure 500 {object} map[string]interface{}
// @Router /topics [get]
// Index 话题列表
// @Summary 获取话题列表
// @Description 获取话题列表，支持分页
// @Tags 话题管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param per_page query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{} "成功"
// @Router /api/v1/topics [get]
func (ctrl *TopicsController) Index(c *gin.Context) {
	request := requests.PaginationRequest{}
	if ok := requests.Validate(c, &request, requests.Pagination); !ok {
		return
	}

	data, pager, err := ctrl.topicService.List(c, 10)
	if err != nil {
		logger.LogErrorWithContext(c, err, "获取话题列表失败")
		response.Abort500(c, "获取列表失败")
		return
	}

	response.JSON(c, gin.H{
		"data":  data,
		"pager": pager,
	})
}

// Show godoc
// @Summary 获取话题详情
// @Description 根据ID获取话题详细信息
// @Tags topics
// @Accept json
// @Produce json
// @Param id path string true "话题ID"
// @Success 200 {object} topic.Topic
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /topics/{id} [get]
// Show 话题详情
// @Summary 获取话题详情
// @Description 根据ID获取话题详情
// @Tags 话题管理
// @Accept json
// @Produce json
// @Param id path string true "话题ID"
// @Success 200 {object} topic.Topic "成功"
// @Failure 404 {object} map[string]interface{} "话题不存在"
// @Router /api/v1/topics/{id} [get]
func (ctrl *TopicsController) Show(c *gin.Context) {
	topicModel, err := ctrl.topicService.GetByID(c.Param("id"))
	if err != nil {
		if apperrors.IsAppError(err) {
			appErr := apperrors.GetAppError(err)
			appErr.WithRequestID(middlewares.GetRequestID(c))
			response.Abort404(c)
			return
		}
		logger.LogErrorWithContext(c, err, "获取话题失败")
		response.Abort500(c)
		return
	}
	response.Data(c, topicModel)
}

// Store godoc
// @Summary 创建话题
// @Description 创建新的话题
// @Tags topics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param topic body requests.TopicRequest true "话题信息"
// @Success 201 {object} topic.Topic
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /topics [post]
// Store 创建话题
// @Summary 创建新话题
// @Description 创建一个新的话题
// @Tags 话题管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param topic body requests.TopicRequest true "话题信息"
// @Success 201 {object} topic.Topic "创建成功"
// @Failure 422 {object} map[string]interface{} "验证失败"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/v1/topics [post]
func (ctrl *TopicsController) Store(c *gin.Context) {
	request := requests.TopicRequest{}
	if ok := requests.Validate(c, &request, requests.TopicSave); !ok {
		return
	}

	dto := services.TopicCreateDTO{
		Title:      request.Title,
		Body:       request.Body,
		CategoryID: request.CategoryID,
		UserID:     auth.CurrentUID(c),
	}

	topicModel, err := ctrl.topicService.Create(dto)
	if err != nil {
		logger.LogErrorWithContext(c, err, "创建话题失败",
			zap.String("title", request.Title),
			zap.String("user_id", auth.CurrentUID(c)),
		)
		response.Abort500(c, "创建失败，请稍后尝试~")
		return
	}

	response.Created(c, topicModel)
}

// Update godoc
// @Summary 更新话题
// @Description 更新话题信息（需要所有权）
// @Tags topics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "话题ID"
// @Param topic body requests.TopicRequest true "话题信息"
// @Success 200 {object} topic.Topic
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /topics/{id} [put]
func (ctrl *TopicsController) Update(c *gin.Context) {
	// 验证请求
	request := requests.TopicRequest{}
	if ok := requests.Validate(c, &request, requests.TopicSave); !ok {
		return
	}

	// 检查所有权
	topicID := c.Param("id")
	currentUserID := auth.CurrentUID(c)
	
	isOwner, err := ctrl.topicService.CheckOwnership(topicID, currentUserID)
	if err != nil {
		if apperrors.IsAppError(err) {
			response.Abort404(c)
			return
		}
		logger.LogErrorWithContext(c, err, "检查话题所有权失败")
		response.Abort500(c)
		return
	}
	
	if !isOwner {
		response.Abort403(c, "无权限操作")
		return
	}

	// 更新话题
	dto := services.TopicUpdateDTO{
		Title:      request.Title,
		Body:       request.Body,
		CategoryID: request.CategoryID,
	}

	topicModel, err := ctrl.topicService.Update(topicID, dto)
	if err != nil {
		logger.LogErrorWithContext(c, err, "更新话题失败",
			zap.String("topic_id", topicID),
		)
		response.Abort500(c, "更新失败，请稍后尝试~")
		return
	}

	response.Data(c, topicModel)
}

// Delete godoc
// @Summary 删除话题
// @Description 删除话题（需要所有权）
// @Tags topics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "话题ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /topics/{id} [delete]
func (ctrl *TopicsController) Delete(c *gin.Context) {
	// 检查所有权
	topicID := c.Param("id")
	currentUserID := auth.CurrentUID(c)
	
	isOwner, err := ctrl.topicService.CheckOwnership(topicID, currentUserID)
	if err != nil {
		if apperrors.IsAppError(err) {
			response.Abort404(c)
			return
		}
		logger.LogErrorWithContext(c, err, "检查话题所有权失败")
		response.Abort500(c)
		return
	}
	
	if !isOwner {
		response.Abort403(c, "无权限操作")
		return
	}

	// 删除话题
	err = ctrl.topicService.Delete(topicID)
	if err != nil {
		logger.LogErrorWithContext(c, err, "删除话题失败",
			zap.String("topic_id", topicID),
		)
		response.Abort500(c, "删除失败，请稍后尝试~")
		return
	}

	response.Success(c)
}
