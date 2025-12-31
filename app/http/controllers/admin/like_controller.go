package admin

import (
	"GoHub-Service/app/models/like"
	"GoHub-Service/app/repositories"
	"GoHub-Service/pkg/database"
	"GoHub-Service/pkg/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// LikeController 点赞管理控制器
type LikeController struct{}

// AdminLikeResponse 管理员视图的点赞响应结构
type AdminLikeResponse struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	TargetType string `json:"target_type"`
	TargetID   string `json:"target_id"`
	CreatedAt  string `json:"created_at"`
}

// toAdminLikeResponse 转换为管理员响应格式
func toAdminLikeResponse(l *like.Like) AdminLikeResponse {
	return AdminLikeResponse{
		ID:         fmt.Sprintf("%d", l.ID),
		UserID:     l.UserID,
		TargetType: l.TargetType,
		TargetID:   l.TargetID,
		CreatedAt:  l.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Index 点赞列表
// @Summary 获取点赞列表
// @Description 分页获取所有点赞
// @Tags Like
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param per_page query int false "每页数量" default(20)
// @Param target_type query string false "目标类型过滤"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/likes [get]
func (ctrl *LikeController) Index(c *gin.Context) {
	page := cast.ToInt(c.DefaultQuery("page", "1"))
	perPage := cast.ToInt(c.DefaultQuery("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	db := database.DB.Model(&like.Like{})

	// 目标类型过滤
	if targetType := c.Query("target_type"); targetType != "" {
		db = db.Where("target_type = ?", targetType)
	}

	// 获取总数
	var total int64
	db.Count(&total)

	// 分页获取数据
	var likes []like.Like
	if err := db.Offset((page - 1) * perPage).
		Limit(perPage).
		Order("created_at DESC").
		Find(&likes).Error; err != nil {
		response.ApiError(c, http.StatusInternalServerError, response.CodeServerError, "获取点赞列表失败")
		return
	}

	// 转换为响应格式
	adminLikes := make([]AdminLikeResponse, len(likes))
	for i, l := range likes {
		adminLikes[i] = toAdminLikeResponse(&l)
	}

	response.ApiSuccess(c, gin.H{
		"data": adminLikes,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"per_page":  perPage,
			"last_page": (total + int64(perPage) - 1) / int64(perPage),
		},
	})
}

// Delete 删除点赞
// @Summary 删除点赞
// @Description 删除指定的点赞
// @Tags Like
// @Accept json
// @Produce json
// @Param id path string true "点赞ID"
// @Success 204
// @Router /api/v1/admin/likes/{id} [delete]
func (ctrl *LikeController) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, "点赞ID无效")
		return
	}

	// 检查点赞是否存在
	var l like.Like
	if err := database.DB.First(&l, "id = ?", id).Error; err != nil {
		response.ApiError(c, http.StatusNotFound, response.CodeNotFound, "点赞不存在")
		return
	}

	// 删除点赞
	if err := database.DB.Delete(&l).Error; err != nil {
		response.ApiError(c, http.StatusInternalServerError, response.CodeServerError, "删除点赞失败")
		return
	}

	c.Status(http.StatusNoContent)
}

// Stats 点赞统计
// @Summary 获取点赞统计
// @Description 获取点赞相关统计数据
// @Tags Like
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/likes/stats [get]
func (ctrl *LikeController) Stats(c *gin.Context) {
	var totalLikes int64
	var topicLikes int64
	var commentLikes int64

	// 获取总点赞数
	database.DB.Model(&like.Like{}).Count(&totalLikes)

	// 获取话题点赞数
	database.DB.Model(&like.Like{}).Where("target_type = ?", "topic").Count(&topicLikes)

	// 获取评论点赞数
	database.DB.Model(&like.Like{}).Where("target_type = ?", "comment").Count(&commentLikes)

	response.ApiSuccess(c, gin.H{
		"total_likes":   totalLikes,
		"topic_likes":   topicLikes,
		"comment_likes": commentLikes,
	})
}

// GetTargetLikes 获取目标的点赞列表
// @Summary 获取目标的点赞列表
// @Description 获取指定目标（话题或评论）的点赞列表
// @Tags Like
// @Accept json
// @Produce json
// @Param target_type path string true "目标类型 (topic/comment)"
// @Param target_id path string true "目标ID"
// @Param page query int false "页码" default(1)
// @Param per_page query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/{target_type}/{target_id}/likes [get]
func (ctrl *LikeController) GetTargetLikes(c *gin.Context) {
	targetType := c.Param("target_type")
	targetID := c.Param("target_id")

	if targetType == "" || targetID == "" {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, "目标类型或ID无效")
		return
	}

	if targetType != "topic" && targetType != "comment" {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, "目标类型只能是 topic 或 comment")
		return
	}

	page := cast.ToInt(c.DefaultQuery("page", "1"))
	perPage := cast.ToInt(c.DefaultQuery("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	repo := repositories.NewLikeRepository()
	likes, total, err := repo.GetByTarget(targetType, targetID, (page-1)*perPage, perPage)
	if err != nil {
		response.ApiError(c, http.StatusInternalServerError, response.CodeServerError, "获取点赞列表失败")
		return
	}

	// 转换为响应格式
	adminLikes := make([]AdminLikeResponse, len(likes))
	for i, l := range likes {
		adminLikes[i] = toAdminLikeResponse(&l)
	}

	response.ApiSuccess(c, gin.H{
		"data": adminLikes,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"per_page":  perPage,
			"last_page": (total + int64(perPage) - 1) / int64(perPage),
		},
	})
}
