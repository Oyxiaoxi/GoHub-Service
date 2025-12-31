package admin

import (
	"GoHub-Service/app/models/comment"
	"GoHub-Service/app/requests"
	"GoHub-Service/pkg/database"
	"GoHub-Service/pkg/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// CommentController 评论管理控制器
type CommentController struct{}

// AdminCommentResponse 管理员视图的评论响应结构
type AdminCommentResponse struct {
	ID        string `json:"id"`
	TopicID   string `json:"topic_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	ParentID  string `json:"parent_id,omitempty"`
	LikeCount int64  `json:"like_count"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// toAdminCommentResponse 转换为管理员响应格式
func toAdminCommentResponse(c *comment.Comment) AdminCommentResponse {
	return AdminCommentResponse{
		ID:        fmt.Sprintf("%d", c.ID),
		TopicID:   c.TopicID,
		UserID:    c.UserID,
		Content:   c.Content,
		ParentID:  c.ParentID,
		LikeCount: c.LikeCount,
		CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Index 评论列表
// @Summary 获取评论列表
// @Description 分页获取所有评论
// @Tags Comment
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param per_page query int false "每页数量" default(20)
// @Param keyword query string false "搜索关键词"
// @Param user_id query string false "用户ID过滤"
// @Param topic_id query string false "话题ID过滤"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/comments [get]
func (ctrl *CommentController) Index(c *gin.Context) {
	page := cast.ToInt(c.DefaultQuery("page", "1"))
	perPage := cast.ToInt(c.DefaultQuery("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	db := database.DB.Model(&comment.Comment{})

	// 搜索条件
	if keyword := c.Query("keyword"); keyword != "" {
		db = db.Where("content LIKE ?", "%"+keyword+"%")
	}

	// 用户ID过滤
	if userID := c.Query("user_id"); userID != "" {
		db = db.Where("user_id = ?", userID)
	}

	// 话题ID过滤
	if topicID := c.Query("topic_id"); topicID != "" {
		db = db.Where("topic_id = ?", topicID)
	}

	// 获取总数
	var total int64
	db.Count(&total)

	// 分页获取数据
	var comments []comment.Comment
	if err := db.Offset((page - 1) * perPage).
		Limit(perPage).
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		response.ApiError(c, http.StatusInternalServerError, response.CodeServerError, "获取评论列表失败")
		return
	}

	// 转换为响应格式
	adminComments := make([]AdminCommentResponse, len(comments))
	for i, comm := range comments {
		adminComments[i] = toAdminCommentResponse(&comm)
	}

	response.ApiSuccess(c, gin.H{
		"data": adminComments,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"per_page":  perPage,
			"last_page": (total + int64(perPage) - 1) / int64(perPage),
		},
	})
}

// Show 评论详情
// @Summary 获取评论详情
// @Description 根据ID获取评论详情
// @Tags Comment
// @Accept json
// @Produce json
// @Param id path string true "评论ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/comments/{id} [get]
func (ctrl *CommentController) Show(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, "评论ID无效")
		return
	}

	var comm comment.Comment
	if err := database.DB.First(&comm, "id = ?", id).Error; err != nil {
		response.ApiError(c, http.StatusNotFound, response.CodeNotFound, "评论不存在")
		return
	}

	response.ApiSuccess(c, toAdminCommentResponse(&comm))
}

// Delete 删除评论
// @Summary 删除评论
// @Description 删除指定的评论
// @Tags Comment
// @Accept json
// @Produce json
// @Param id path string true "评论ID"
// @Success 204
// @Router /api/v1/admin/comments/{id} [delete]
func (ctrl *CommentController) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, "评论ID无效")
		return
	}

	// 检查评论是否存在
	var comm comment.Comment
	if err := database.DB.First(&comm, "id = ?", id).Error; err != nil {
		response.ApiError(c, http.StatusNotFound, response.CodeNotFound, "评论不存在")
		return
	}

	// 删除评论
	if err := database.DB.Delete(&comm).Error; err != nil {
		response.ApiError(c, http.StatusInternalServerError, response.CodeServerError, "删除评论失败")
		return
	}

	c.Status(http.StatusNoContent)
}

// BatchDelete 批量删除评论
// @Summary 批量删除评论
// @Description 批量删除多条评论
// @Tags Comment
// @Accept json
// @Produce json
// @Param request body requests.BatchDeleteRequest true "ID列表"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/comments/batch-delete [post]
func (ctrl *CommentController) BatchDelete(c *gin.Context) {
	var req requests.BatchDeleteRequest
	if !requests.ValidateBatchDelete(c, &req) {
		return
	}

	if len(req.IDs) == 0 {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, "ID列表不能为空")
		return
	}

	// 删除评论
	if err := database.DB.Where("id IN ?", req.IDs).Delete(&comment.Comment{}).Error; err != nil {
		response.ApiError(c, http.StatusInternalServerError, response.CodeServerError, "批量删除评论失败")
		return
	}

	response.ApiSuccessWithMessage(c, "批量删除评论成功")
}

// Stats 评论统计
// @Summary 获取评论统计
// @Description 获取评论相关统计数据
// @Tags Comment
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/comments/stats [get]
func (ctrl *CommentController) Stats(c *gin.Context) {
	var totalComments int64
	var totalLikes int64

	// 获取总评论数
	database.DB.Model(&comment.Comment{}).Count(&totalComments)

	// 获取总点赞数
	database.DB.Model(&comment.Comment{}).Select("SUM(like_count)").Row().Scan(&totalLikes)

	// 获取今日评论数
	database.DB.Model(&comment.Comment{}).
		Where("DATE(created_at) = CURDATE()").
		Count(&totalComments)

	response.ApiSuccess(c, gin.H{
		"total_comments": totalComments,
		"total_likes":    totalLikes,
	})
}
