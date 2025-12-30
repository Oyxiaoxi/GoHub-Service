package admin

import (
	"GoHub-Service/app/models/category"
	"GoHub-Service/app/models/topic"
	"GoHub-Service/app/models/user"
	"GoHub-Service/pkg/database"
	"GoHub-Service/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

// DashboardController 管理后台仪表盘控制器
type DashboardController struct{}

// Statistics 系统统计信息
type Statistics struct {
	TotalUsers      int64 `json:"total_users"`
	TotalTopics     int64 `json:"total_topics"`
	TotalCategories int64 `json:"total_categories"`
	TodayUsers      int64 `json:"today_users"`
	TodayTopics     int64 `json:"today_topics"`
	ActiveUsers     int64 `json:"active_users"`     // 最近7天活跃用户
	PopularTopics   int64 `json:"popular_topics"`   // 热门话题数
}

// Overview 获取系统概览统计
func (ctrl *DashboardController) Overview(c *gin.Context) {
	db := database.DB

	var stats Statistics

	// 总用户数
	db.Model(&user.User{}).Count(&stats.TotalUsers)

	// 总话题数
	db.Model(&topic.Topic{}).Count(&stats.TotalTopics)

	// 总分类数
	db.Model(&category.Category{}).Count(&stats.TotalCategories)

	// 今日新增用户
	today := time.Now().Truncate(24 * time.Hour)
	db.Model(&user.User{}).Where("created_at >= ?", today).Count(&stats.TodayUsers)

	// 今日新增话题
	db.Model(&topic.Topic{}).Where("created_at >= ?", today).Count(&stats.TodayTopics)

	// 最近7天活跃用户（有发布话题或评论的）
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	db.Model(&user.User{}).
		Joins("LEFT JOIN topics ON users.id = topics.user_id").
		Where("topics.created_at >= ?", sevenDaysAgo).
		Group("users.id").
		Count(&stats.ActiveUsers)

	response.Data(c, gin.H{
		"statistics": stats,
		"timestamp":  time.Now(),
	})
}

// RecentUsers 获取最近注册用户
func (ctrl *DashboardController) RecentUsers(c *gin.Context) {
	db := database.DB

	var users []user.User
	db.Order("created_at DESC").Limit(10).Find(&users)

	response.Data(c, gin.H{
		"users": users,
	})
}

// RecentTopics 获取最近发布话题
func (ctrl *DashboardController) RecentTopics(c *gin.Context) {
	db := database.DB

	var topics []topic.Topic
	db.Preload("User").
		Preload("Category").
		Order("created_at DESC").
		Limit(10).
		Find(&topics)

	response.Data(c, gin.H{
		"topics": topics,
	})
}
