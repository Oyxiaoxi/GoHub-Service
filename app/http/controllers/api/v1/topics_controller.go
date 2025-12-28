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
