package admin

import (
	"GoHub-Service/app/models/follow"
	"GoHub-Service/app/repositories"
	"GoHub-Service/pkg/database"
	"GoHub-Service/pkg/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// FollowController 关注管理控制器
type FollowController struct{}

// AdminFollowResponse 管理员视图的关注响应结构
type AdminFollowResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	FollowID  string `json:"follow_id"`
	CreatedAt string `json:"created_at"`
}

// toAdminFollowResponse 转换为管理员响应格式
func toAdminFollowResponse(f *follow.Follow) AdminFollowResponse {
	return AdminFollowResponse{
		ID:        fmt.Sprintf("%d", f.ID),
		UserID:    f.UserID,
		FollowID:  f.FollowID,
		CreatedAt: f.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Index 关注列表
// @Summary 获取关注列表
// @Description 分页获取所有关注关系
// @Tags Follow
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param per_page query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/follows [get]
func (ctrl *FollowController) Index(c *gin.Context) {
	page := cast.ToInt(c.DefaultQuery("page", "1"))
	perPage := cast.ToInt(c.DefaultQuery("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	repo := repositories.NewFollowRepository()
	follows, total, err := repo.GetAll((page-1)*perPage, perPage)
	if err != nil {
		response.ApiError(c, http.StatusInternalServerError, response.CodeServerError, "获取关注列表失败")
		return
	}

	// 转换为响应格式
	adminFollows := make([]AdminFollowResponse, len(follows))
	for i, f := range follows {
		adminFollows[i] = toAdminFollowResponse(&f)
	}

	response.Data(c, gin.H{
		"follows": adminFollows,
		"paging": gin.H{
			"total":      total,
			"page":       page,
			"per_page":   perPage,
			"total_page": (total + int64(perPage) - 1) / int64(perPage),
		},
	})
}

// Delete 删除关注关系
// @Summary 删除关注关系
// @Description 删除指定的关注关系
// @Tags Follow
// @Accept json
// @Produce json
// @Param id path string true "关注ID"
// @Success 204
// @Router /api/v1/admin/follows/{id} [delete]
func (ctrl *FollowController) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, "关注ID无效")
		return
	}

	// 检查关注是否存在
	var f follow.Follow
	if err := database.DB.First(&f, "id = ?", id).Error; err != nil {
		response.Abort404(c, "关注不存在")
		return
	}

	// 删除关注
	if err := database.DB.Delete(&f).Error; err != nil {
		response.ApiError(c, http.StatusInternalServerError, response.CodeServerError, "删除关注失败")
		return
	}

	c.Status(http.StatusNoContent)
}

// Stats 关注统计
// @Summary 获取关注统计
// @Description 获取关注相关统计数据
// @Tags Follow
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/follows/stats [get]
func (ctrl *FollowController) Stats(c *gin.Context) {
	var totalFollows int64

	// 获取总关注数
	database.DB.Model(&follow.Follow{}).Count(&totalFollows)

	response.Data(c, gin.H{
		"total_follows": totalFollows,
	})
}

// GetFollowers 用户粉丝列表
// @Summary 获取用户粉丝列表
// @Description 获取指定用户的粉丝列表
// @Tags Follow
// @Accept json
// @Produce json
// @Param user_id path string true "用户ID"
// @Param page query int false "页码" default(1)
// @Param per_page query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/users/{user_id}/followers [get]
func (ctrl *FollowController) GetFollowers(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, "用户ID无效")
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

	repo := repositories.NewFollowRepository()
	followers, total, err := repo.GetFollowers(userID, (page-1)*perPage, perPage)
	if err != nil {
		response.ApiError(c, http.StatusInternalServerError, response.CodeServerError, "获取粉丝列表失败")
		return
	}

	// 转换为响应格式
	adminFollows := make([]AdminFollowResponse, len(followers))
	for i, f := range followers {
		adminFollows[i] = toAdminFollowResponse(&f)
	}

	response.Data(c, gin.H{
		"followers": adminFollows,
		"paging": gin.H{
			"total":      total,
			"page":       page,
			"per_page":   perPage,
			"total_page": (total + int64(perPage) - 1) / int64(perPage),
		},
	})
}

// GetFollowing 用户关注列表
// @Summary 获取用户关注列表
// @Description 获取指定用户的关注列表
// @Tags Follow
// @Accept json
// @Produce json
// @Param user_id path string true "用户ID"
// @Param page query int false "页码" default(1)
// @Param per_page query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/users/{user_id}/following [get]
func (ctrl *FollowController) GetFollowing(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		response.ApiError(c, http.StatusBadRequest, response.CodeValidationError, "用户ID无效")
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

	repo := repositories.NewFollowRepository()
	following, total, err := repo.GetFollowing(userID, (page-1)*perPage, perPage)
	if err != nil {
		response.ApiError(c, http.StatusInternalServerError, response.CodeServerError, "获取关注列表失败")
		return
	}

	// 转换为响应格式
	adminFollows := make([]AdminFollowResponse, len(following))
	for i, f := range following {
		adminFollows[i] = toAdminFollowResponse(&f)
	}

	response.Data(c, gin.H{
		"following": adminFollows,
		"paging": gin.H{
			"total":      total,
			"page":       page,
			"per_page":   perPage,
			"total_page": (total + int64(perPage) - 1) / int64(perPage),
		},
	})
}
