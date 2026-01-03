// Package routes 缓存监控路由
package routes

import (
	controllers "GoHub-Service/app/http/controllers/api/v1"
	"github.com/gin-gonic/gin"
)

// RegisterCacheMonitorRoutes 注册缓存监控路由
func RegisterCacheMonitorRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		cacheCtrl := new(controllers.CacheMonitorController)
		
		monitor := v1.Group("/monitor/cache")
		{
			monitor.GET("/stats", cacheCtrl.Stats)
			monitor.GET("/warmup", cacheCtrl.WarmupList)
			monitor.POST("/warmup", cacheCtrl.WarmupAll)
			monitor.POST("/warmup/:name", cacheCtrl.WarmupOne)
			monitor.POST("/degrade", cacheCtrl.Degrade)
			monitor.POST("/recover", cacheCtrl.Recover)
		}
	}
}
