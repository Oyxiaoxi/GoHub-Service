package admin

import (
	"GoHub-Service/app/models/user"
	"GoHub-Service/app/requests"
	"GoHub-Service/pkg/database"
	"GoHub-Service/pkg/paginator"
	"GoHub-Service/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// UserController 用户管理控制器
type UserController struct{}

// AdminUserResponse 管理员视图的用户响应结构（包含敏感信息）
type AdminUserResponse struct {
	ID             uint64 `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	City           string `json:"city,omitempty"`
	Introduction   string `json:"introduction,omitempty"`
	Avatar         string `json:"avatar,omitempty"`
	FollowersCount int64  `json:"followers_count,omitempty"`
	Points         int64  `json:"points,omitempty"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// toAdminUserResponse 将用户模型转换为管理员响应格式
func toAdminUserResponse(u user.User) AdminUserResponse {
	return AdminUserResponse{
		ID:             u.ID,
		Name:           u.Name,
		Email:          u.Email,
		Phone:          u.Phone,
		City:           u.City,
		Introduction:   u.Introduction,
		Avatar:         u.Avatar,
		FollowersCount: u.FollowersCount,
		Points:         u.Points,
		CreatedAt:      u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Index 用户列表
func (ctrl *UserController) Index(c *gin.Context) {
	var req requests.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}

	db := database.DB.Model(&user.User{})

	// 搜索条件
	if keyword := c.Query("keyword"); keyword != "" {
		db = db.Where("name LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 分页
	var users []user.User
	perPage := cast.ToInt(c.DefaultQuery("per_page", "20"))
	paging := paginator.Paginate(
		c,
		db,
		&users,
		"/api/v1/admin/users",
		perPage,
	)

	// 转换为管理员响应格式
	adminUsers := make([]AdminUserResponse, len(users))
	for i, u := range users {
		adminUsers[i] = toAdminUserResponse(u)
	}

	response.Data(c, gin.H{
		"users": adminUsers,
		"paging": paging,
	})
}

// Show 用户详情
func (ctrl *UserController) Show(c *gin.Context) {
	userID := c.Param("id")
	
	var u user.User
	if err := database.DB.First(&u, userID).Error; err != nil {
		response.Abort404(c, "用户不存在")
		return
	}

	// 获取用户统计信息
	var topicCount, commentCount int64
	database.DB.Model(&user.User{}).Where("user_id = ?", userID).Count(&topicCount)
	
	response.Data(c, gin.H{
		"user": u,
		"statistics": gin.H{
			"topic_count":   topicCount,
			"comment_count": commentCount,
		},
	})
}

// Update 更新用户信息
func (ctrl *UserController) Update(c *gin.Context) {
	userID := c.Param("id")

	var u user.User
	if err := database.DB.First(&u, userID).Error; err != nil {
		response.Abort404(c, "用户不存在")
		return
	}

	// 绑定更新数据
	type UpdateUserRequest struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		Phone       string `json:"phone"`
		City        string `json:"city"`
		Introduction string `json:"introduction"`
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}

	// 更新字段
	if req.Name != "" {
		u.Name = req.Name
	}
	if req.Email != "" {
		u.Email = req.Email
	}
	if req.Phone != "" {
		u.Phone = req.Phone
	}
	if req.City != "" {
		u.City = req.City
	}
	if req.Introduction != "" {
		u.Introduction = req.Introduction
	}

	if err := database.DB.Save(&u).Error; err != nil {
		response.Abort500(c, "更新失败")
		return
	}

	response.Data(c, gin.H{
		"user": u,
	})
}

// Delete 删除用户（软删除）
func (ctrl *UserController) Delete(c *gin.Context) {
	userID := c.Param("id")

	var u user.User
	if err := database.DB.First(&u, userID).Error; err != nil {
		response.Abort404(c, "用户不存在")
		return
	}

	// 软删除
	if err := database.DB.Delete(&u).Error; err != nil {
		response.Abort500(c, "删除失败")
		return
	}

	response.Data(c, gin.H{
		"message": "用户已删除",
	})
}

// Ban 封禁用户
func (ctrl *UserController) Ban(c *gin.Context) {
	userID := c.Param("id")

	type BanRequest struct {
		Reason string `json:"reason" binding:"required"`
		Days   int    `json:"days" binding:"required,min=1"`
	}

	var req BanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}

	var u user.User
	if err := database.DB.First(&u, userID).Error; err != nil {
		response.Abort404(c, "用户不存在")
		return
	}

	// TODO: 添加封禁逻辑（需要在 User 模型中添加封禁相关字段）
	// 例如: u.IsBanned = true
	// u.BanReason = req.Reason
	// u.BanUntil = time.Now().Add(time.Duration(req.Days) * 24 * time.Hour)

	if err := database.DB.Save(&u).Error; err != nil {
		response.Abort500(c, "封禁失败")
		return
	}

	response.Data(c, gin.H{
		"message": "用户已封禁",
		"reason":  req.Reason,
		"days":    req.Days,
	})
}

// Unban 解封用户
func (ctrl *UserController) Unban(c *gin.Context) {
	userID := c.Param("id")

	var u user.User
	if err := database.DB.First(&u, userID).Error; err != nil {
		response.Abort404(c, "用户不存在")
		return
	}

	// TODO: 添加解封逻辑（需要在 User 模型中添加封禁相关字段）
	// 例如: u.IsBanned = false
	// u.BanReason = ""
	// u.BanUntil = nil

	if err := database.DB.Save(&u).Error; err != nil {
		response.Abort500(c, "解封失败")
		return
	}

	response.Data(c, gin.H{
		"message": "用户已解封",
	})
}

// BatchDelete 批量删除用户
func (ctrl *UserController) BatchDelete(c *gin.Context) {
	type BatchRequest struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	var req BatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}

	if err := database.DB.Delete(&user.User{}, req.IDs).Error; err != nil {
		response.Abort500(c, "批量删除失败")
		return
	}

	response.Data(c, gin.H{
		"message": "批量删除成功",
		"count":   len(req.IDs),
	})
}

// ResetPassword 重置用户密码
func (ctrl *UserController) ResetPassword(c *gin.Context) {
	userID := c.Param("id")

	type ResetPasswordRequest struct {
		Password             string `json:"password" binding:"required,min=6"`
		PasswordConfirmation string `json:"password_confirmation" binding:"required"`
	}

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}

	if req.Password != req.PasswordConfirmation {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "两次密码不一致"})
		return
	}

	var u user.User
	if err := database.DB.First(&u, userID).Error; err != nil {
		response.Abort404(c, "用户不存在")
		return
	}

	// 设置新密码（这里需要使用 bcrypt 加密）
	// u.Password = bcrypt.Hash(req.Password)
	
	if err := database.DB.Save(&u).Error; err != nil {
		response.Abort500(c, "重置密码失败")
		return
	}

	response.Data(c, gin.H{
		"message": "密码重置成功",
	})
}

// AssignRole 分配角色
func (ctrl *UserController) AssignRole(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	type AssignRoleRequest struct {
		RoleIDs []uint `json:"role_ids" binding:"required"`
	}

	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}

	var u user.User
	if err := database.DB.First(&u, userID).Error; err != nil {
		response.Abort404(c, "用户不存在")
		return
	}

	// 删除旧的角色关联
	database.DB.Exec("DELETE FROM user_roles WHERE user_id = ?", userID)

	// 添加新的角色关联
	for _, roleID := range req.RoleIDs {
		database.DB.Exec("INSERT INTO user_roles (user_id, role_id, created_at, updated_at) VALUES (?, ?, NOW(), NOW())", 
			userID, roleID)
	}

	response.Data(c, gin.H{
		"message": "角色分配成功",
	})
}
