// Package routes 数据库监控路由
package routes

import (
	controllers "GoHub-Service/app/http/controllers/api/v1"
	"github.com/gin-gonic/gin"
)

// RegisterDatabaseStatsRoutes 注册数据库监控相关路由
func RegisterDatabaseStatsRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		dbCtrl := new(controllers.DatabaseStatsController)
		
		monitor := v1.Group("/monitor/database")
		{
			// 连接池统计信息
			monitor.GET("/stats", dbCtrl.Stats)
			
			// 健康检查
			monitor.GET("/health", dbCtrl.Health)
			
			// Prometheus 指标
			monitor.GET("/metrics", dbCtrl.Metrics)
			
			// 配置优化建议
			monitor.GET("/recommend", dbCtrl.Recommend)
		}
	}
}
