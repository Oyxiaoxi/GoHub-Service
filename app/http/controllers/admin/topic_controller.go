package admin

import (
	"GoHub-Service/app/models/topic"
	"GoHub-Service/app/requests"
	"GoHub-Service/pkg/database"
	"GoHub-Service/pkg/paginator"
	"GoHub-Service/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// TopicController 话题管理控制器
type TopicController struct{}

// Index 话题列表
func (ctrl *TopicController) Index(c *gin.Context) {
	var req requests.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}

	db := database.DB.Model(&topic.Topic{})

	// 搜索条件
	if keyword := c.Query("keyword"); keyword != "" {
		db = db.Where("title LIKE ? OR body LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 分类筛选
	if categoryID := c.Query("category_id"); categoryID != "" {
		db = db.Where("category_id = ?", categoryID)
	}

	// 用户筛选
	if userID := c.Query("user_id"); userID != "" {
		db = db.Where("user_id = ?", userID)
	}

	// 预加载关联
	db = db.Preload("User").Preload("Category")

	// 分页
	var topics []topic.Topic
	perPage := 20
	if pp := c.Query("per_page"); pp != "" {
		if ppInt, err := strconv.Atoi(pp); err == nil {
			perPage = ppInt
		}
	}
	paging := paginator.Paginate(
		c,
		db,
		&topics,
		"/api/v1/admin/topics",
		perPage,
	)

	response.Data(c, gin.H{
		"topics": topics,
		"paging": paging,
	})
}

// Show 话题详情
func (ctrl *TopicController) Show(c *gin.Context) {
	topicID := c.Param("id")

	var t topic.Topic
	if err := database.DB.
		Preload("User").
		Preload("Category").
		First(&t, topicID).Error; err != nil {
		response.Abort404(c, "话题不存在")
		return
	}

	response.Data(c, gin.H{
		"topic": t,
	})
}

// Update 更新话题
func (ctrl *TopicController) Update(c *gin.Context) {
	topicID := c.Param("id")

	var t topic.Topic
	if err := database.DB.First(&t, topicID).Error; err != nil {
		response.Abort404(c, "话题不存在")
		return
	}

	type UpdateTopicRequest struct {
		Title      string `json:"title"`
		Body       string `json:"body"`
		CategoryID uint   `json:"category_id"`
		Status     int    `json:"status"`
	}

	var req UpdateTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}

	// 更新字段
	if req.Title != "" {
		t.Title = req.Title
	}
	if req.Body != "" {
		t.Body = req.Body
	}
	if req.CategoryID > 0 {
		t.CategoryID = cast.ToString(req.CategoryID)
	}

	if err := database.DB.Save(&t).Error; err != nil {
		response.Abort500(c, "更新失败")
		return
	}

	response.Data(c, gin.H{
		"topic": t,
	})
}

// Delete 删除话题
func (ctrl *TopicController) Delete(c *gin.Context) {
	topicID := c.Param("id")

	var t topic.Topic
	if err := database.DB.First(&t, topicID).Error; err != nil {
		response.Abort404(c, "话题不存在")
		return
	}

	if err := database.DB.Delete(&t).Error; err != nil {
		response.Abort500(c, "删除失败")
		return
	}

	response.Data(c, gin.H{
		"message": "话题已删除",
	})
}

// Pin 置顶话题
func (ctrl *TopicController) Pin(c *gin.Context) {
	topicID := c.Param("id")

	var t topic.Topic
	if err := database.DB.First(&t, topicID).Error; err != nil {
		response.Abort404(c, "话题不存在")
		return
	}

	// TODO: 添加置顶逻辑（需要在 Topic 模型中添加 IsPinned 或 PinnedAt 字段）
	// 例如: t.IsPinned = true
	// t.PinnedAt = time.Now()

	if err := database.DB.Save(&t).Error; err != nil {
		response.Abort500(c, "置顶失败")
		return
	}

	response.Data(c, gin.H{
		"message": "话题已置顶",
	})
}

// Unpin 取消置顶
func (ctrl *TopicController) Unpin(c *gin.Context) {
	topicID := c.Param("id")

	var t topic.Topic
	if err := database.DB.First(&t, topicID).Error; err != nil {
		response.Abort404(c, "话题不存在")
		return
	}

	// TODO: 添加取消置顶逻辑（需要在 Topic 模型中添加 IsPinned 或 PinnedAt 字段）
	// 例如: t.IsPinned = false
	// t.PinnedAt = nil

	if err := database.DB.Save(&t).Error; err != nil {
		response.Abort500(c, "取消置顶失败")
		return
	}

	response.Data(c, gin.H{
		"message": "已取消置顶",
	})
}

// BatchDelete 批量删除话题
func (ctrl *TopicController) BatchDelete(c *gin.Context) {
	type BatchRequest struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	var req BatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}

	if err := database.DB.Delete(&topic.Topic{}, req.IDs).Error; err != nil {
		response.Abort500(c, "批量删除失败")
		return
	}

	response.Data(c, gin.H{
		"message": "批量删除成功",
		"count":   len(req.IDs),
	})
}

// Approve 审核通过话题
func (ctrl *TopicController) Approve(c *gin.Context) {
	topicID := c.Param("id")

	var t topic.Topic
	if err := database.DB.First(&t, topicID).Error; err != nil {
		response.Abort404(c, "话题不存在")
		return
	}

	// 审核通过状态设置为 1
	// t.Status = 1

	if err := database.DB.Save(&t).Error; err != nil {
		response.Abort500(c, "审核失败")
		return
	}

	response.Data(c, gin.H{
		"message": "审核通过",
	})
}

// Reject 审核拒绝话题
func (ctrl *TopicController) Reject(c *gin.Context) {
	topicID := c.Param("id")

	type RejectRequest struct {
		Reason string `json:"reason" binding:"required"`
	}

	var req RejectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}

	var t topic.Topic
	if err := database.DB.First(&t, topicID).Error; err != nil {
		response.Abort404(c, "话题不存在")
		return
	}

	// 审核拒绝状态设置为 -1
	// t.Status = -1

	if err := database.DB.Save(&t).Error; err != nil {
		response.Abort500(c, "审核失败")
		return
	}

	response.Data(c, gin.H{
		"message": "已拒绝",
		"reason":  req.Reason,
	})
}
