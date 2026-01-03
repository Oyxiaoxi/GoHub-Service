package v1

import (
	"GoHub-Service/app/requests"
	"GoHub-Service/app/services"
	"GoHub-Service/pkg/auth"
	ctx "GoHub-Service/pkg/ctx"
	"GoHub-Service/pkg/logger"
	"GoHub-Service/pkg/response"

	"github.com/gin-gonic/gin"
)

type CommentsController struct {
	BaseAPIController
	commentService *services.CommentService
}

// NewCommentsController 创建CommentsController实例
func NewCommentsController() *CommentsController {
	return &CommentsController{
		commentService: services.NewCommentService(),
	}
}

// Index 评论列表
// @Summary 获取评论列表
// @Description 获取评论列表，支持分页
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param per_page query int false "每页数量" default(15)
// @Success 200 {object} map[string]interface{} "成功"
// @Router /api/v1/comments [get]
func (ctrl *CommentsController) Index(c *gin.Context) {
	request := requests.PaginationRequest{}
	if ok := requests.Validate(c, &request, requests.Pagination); !ok {
		return
	}

	// 从 Gin Context 创建请求 Context
	requestCtx := ctx.FromGinContext(c)

	listResponse, err := ctrl.commentService.List(requestCtx, c, 15)
	if err != nil {
		logger.LogErrorWithContext(c, err, "获取评论列表失败")
		response.ApiError(c, 500, err.Code, err.Message)
		return
	}
	response.JSON(c, gin.H{
		"data":  listResponse.Comments,
		"pager": listResponse.Paging,
	})
}

// Show 评论详情
// @Summary 获取评论详情
// @Description 根据ID获取评论详情
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param id path string true "评论ID"
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 404 {object} map[string]interface{} "评论不存在"
// @Router /api/v1/comments/{id} [get]
func (ctrl *CommentsController) Show(c *gin.Context) {
	// 从 Gin Context 创建请求 Context
	requestCtx := ctx.FromGinContext(c)

	commentModel, err := ctrl.commentService.GetByID(requestCtx, c.Param("id"))
	if err != nil {
		logger.LogErrorWithContext(c, err, "获取评论失败")
		if err.Code == 4001 {
			response.Abort404(c)
		} else {
			response.ApiError(c, 500, err.Code, err.Message)
		}
		return
	}
	response.Data(c, commentModel)
}

// Store 创建评论
// @Summary 创建评论
// @Description 创建新评论
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param comment body requests.CommentRequest true "评论信息"
// @Success 201 {object} map[string]interface{} "成功"
// @Failure 422 {object} map[string]interface{} "验证错误"
// @Router /api/v1/comments [post]
// @Security ApiKeyAuth
func (ctrl *CommentsController) Store(c *gin.Context) {
	request := requests.CommentRequest{}
	if ok := requests.Validate(c, &request, requests.CommentSave); !ok {
		return
	}

	// 获取当前用户ID
	userID := auth.CurrentUID(c)

	// 从 Gin Context 创建请求 Context
	requestCtx := ctx.FromGinContext(c)

	// 创建评论
	dto := &services.CommentCreateDTO{
		TopicID:  request.TopicID,
		UserID:   userID,
		Content:  request.Content,
		ParentID: request.ParentID,
	}

	commentModel, err := ctrl.commentService.Create(requestCtx, dto)
	if err != nil {
		logger.LogErrorWithContext(c, err, "创建评论失败")
		response.ApiError(c, 500, err.Code, err.Message)
		return
	}

	response.Created(c, commentModel)
}

// Update 更新评论
// @Summary 更新评论
// @Description 更新评论内容
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path string true "评论ID"
// @Param comment body requests.CommentRequest true "评论信息"
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 404 {object} map[string]interface{} "评论不存在"
// @Failure 422 {object} map[string]interface{} "验证错误"
// @Router /api/v1/comments/{id} [put]
// @Security ApiKeyAuth
func (ctrl *CommentsController) Update(c *gin.Context) {
	request := requests.CommentRequest{}
	if ok := requests.Validate(c, &request, requests.CommentUpdate); !ok {
		return
	}

	// 更新评论
	dto := &services.CommentUpdateDTO{
		Content: &request.Content,
	}

	// 从 Gin Context 创建请求 Context
	requestCtx := ctx.FromGinContext(c)

	commentModel, err := ctrl.commentService.Update(requestCtx, c.Param("id"), dto)
	if err != nil {
		logger.LogErrorWithContext(c, err, "更新评论失败")
		if err.Code == 4001 {
			response.Abort404(c)
		} else {
			response.ApiError(c, 500, err.Code, err.Message)
		}
		return
	}

	response.Data(c, commentModel)
}

// Delete 删除评论
// @Summary 删除评论
// @Description 删除指定评论
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path string true "评论ID"
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 404 {object} map[string]interface{} "评论不存在"
// @Router /api/v1/comments/{id} [delete]
// @Security ApiKeyAuth
func (ctrl *CommentsController) Delete(c *gin.Context) {
	// 从 Gin Context 创建请求 Context
	requestCtx := ctx.FromGinContext(c)

	err := ctrl.commentService.Delete(requestCtx, c.Param("id"))
	if err != nil {
		logger.LogErrorWithContext(c, err, "删除评论失败")
		if err.Code == 4001 {
			response.Abort404(c)
		} else {
			response.ApiError(c, 500, err.Code, err.Message)
		}
		return
	}

	response.Success(c)
}

// ListByTopicID 获取话题的评论列表
// @Summary 获取话题的评论列表
// @Description 获取指定话题的所有评论
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param id path string true "话题ID"
// @Param page query int false "页码" default(1)
// @Param per_page query int false "每页数量" default(15)
// @Success 200 {object} map[string]interface{} "成功"
// @Router /api/v1/topics/{id}/comments [get]
func (ctrl *CommentsController) ListByTopicID(c *gin.Context) {
	request := requests.PaginationRequest{}
	if ok := requests.Validate(c, &request, requests.Pagination); !ok {
		return
	}

	topicID := c.Param("id")

	// 从 Gin Context 创建请求 Context
	requestCtx := ctx.FromGinContext(c)

	listResponse, err := ctrl.commentService.ListByTopicID(requestCtx, c, topicID, 15)
	if err != nil {
		logger.LogErrorWithContext(c, err, "获取话题评论列表失败")
		response.ApiError(c, 500, err.Code, err.Message)
		return
	}

	response.JSON(c, gin.H{
		"data":  listResponse.Comments,
		"pager": listResponse.Paging,
	})
}

// ListByUserID 获取用户的评论列表
// @Summary 获取用户的评论列表
// @Description 获取指定用户的所有评论
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Param page query int false "页码" default(1)
// @Param per_page query int false "每页数量" default(15)
// @Success 200 {object} map[string]interface{} "成功"
// @Router /api/v1/users/{id}/comments [get]
func (ctrl *CommentsController) ListByUserID(c *gin.Context) {
	request := requests.PaginationRequest{}
	if ok := requests.Validate(c, &request, requests.Pagination); !ok {
		return
	}

	userID := c.Param("id")

	// 从 Gin Context 创建请求 Context
	requestCtx := ctx.FromGinContext(c)

	listResponse, err := ctrl.commentService.ListByUserID(requestCtx, c, userID, 15)
	if err != nil {
		logger.LogErrorWithContext(c, err, "获取用户评论列表失败")
		response.ApiError(c, 500, err.Code, err.Message)
		return
	}

	response.JSON(c, gin.H{
		"data":  listResponse.Comments,
		"pager": listResponse.Paging,
	})
}

// ListReplies 获取评论的回复列表
// @Summary 获取评论的回复列表
// @Description 获取指定评论的所有回复
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param id path string true "评论ID"
// @Param page query int false "页码" default(1)
// @Param per_page query int false "每页数量" default(15)
// @Success 200 {object} map[string]interface{} "成功"
// @Router /api/v1/comments/{id}/replies [get]
func (ctrl *CommentsController) ListReplies(c *gin.Context) {
	request := requests.PaginationRequest{}
	if ok := requests.Validate(c, &request, requests.Pagination); !ok {
		return
	}

	parentID := c.Param("id")

	// 从 Gin Context 创建请求 Context
	requestCtx := ctx.FromGinContext(c)

	listResponse, err := ctrl.commentService.ListReplies(requestCtx, c, parentID, 15)
	if err != nil {
		logger.LogErrorWithContext(c, err, "获取评论回复列表失败")
		response.ApiError(c, 500, err.Code, err.Message)
		return
	}

	response.JSON(c, gin.H{
		"data":  listResponse.Comments,
		"pager": listResponse.Paging,
	})
}

// Like 点赞评论
// @Summary 点赞评论
// @Description 对评论进行点赞
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path string true "评论ID"
// @Success 200 {object} map[string]interface{} "成功"
// @Router /api/v1/comments/{id}/like [post]
// @Security ApiKeyAuth
func (ctrl *CommentsController) Like(c *gin.Context) {
	// 从 Gin Context 创建请求 Context
	requestCtx := ctx.FromGinContext(c)

	err := ctrl.commentService.LikeComment(requestCtx, c.Param("id"))
	if err != nil {
		logger.LogErrorWithContext(c, err, "点赞评论失败")
		response.ApiError(c, 500, err.Code, err.Message)
		return
	}

	response.Success(c)
}

// Unlike 取消点赞
// @Summary 取消点赞
// @Description 取消对评论的点赞
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path string true "评论ID"
// @Success 200 {object} map[string]interface{} "成功"
// @Router /api/v1/comments/{id}/unlike [post]
// @Security ApiKeyAuth
func (ctrl *CommentsController) Unlike(c *gin.Context) {
	// 从 Gin Context 创建请求 Context
	requestCtx := ctx.FromGinContext(c)

	err := ctrl.commentService.UnlikeComment(requestCtx, c.Param("id"))
	if err != nil {
		logger.LogErrorWithContext(c, err, "取消点赞失败")
		response.ApiError(c, 500, err.Code, err.Message)
		return
	}

	response.Success(c)
}
