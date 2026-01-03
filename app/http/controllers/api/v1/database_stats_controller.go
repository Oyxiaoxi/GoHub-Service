// Package controllers 数据库监控控制器
package controllers

import (
	"GoHub-Service/pkg/database"
	"GoHub-Service/pkg/response"
	"github.com/gin-gonic/gin"
)

// DatabaseStatsController 数据库统计控制器
type DatabaseStatsController struct {
	BaseAPIController
}

// Stats 获取数据库连接池统计信息
// @Summary 数据库连接池统计
// @Description 获取当前数据库连接池的详细统计信息
// @Tags 监控
// @Produce json
// @Success 200 {object} database.Stats
// @Router /api/v1/monitor/database/stats [get]
func (ctrl *DatabaseStatsController) Stats(c *gin.Context) {
	stats := database.GetStats()
	if stats == nil {
		response.Abort500(c, "数据库未初始化")
		return
	}
	response.JSON(c, stats)
}

// Health 检查数据库连接池健康状态
// @Summary 数据库健康检查
// @Description 检查数据库连接池是否健康，返回警告信息
// @Tags 监控
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/monitor/database/health [get]
func (ctrl *DatabaseStatsController) Health(c *gin.Context) {
	healthy, warnings := database.CheckHealth()
	response.JSON(c, gin.H{
		"healthy":  healthy,
		"warnings": warnings,
	})
}

// Metrics 获取 Prometheus 格式的指标
// @Summary 数据库指标
// @Description 获取适合 Prometheus 采集的数据库指标
// @Tags 监控
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/monitor/database/metrics [get]
func (ctrl *DatabaseStatsController) Metrics(c *gin.Context) {
	metrics := database.GetMetrics()
	response.JSON(c, metrics)
}

// Recommend 获取配置优化建议
// @Summary 配置优化建议
// @Description 根据当前统计数据推荐连接池配置调整
// @Tags 监控
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/monitor/database/recommend [get]
func (ctrl *DatabaseStatsController) Recommend(c *gin.Context) {
	recommendations := database.RecommendConfig()
	response.JSON(c, recommendations)
}
