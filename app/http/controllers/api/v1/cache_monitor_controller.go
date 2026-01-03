// Package v1 缓存监控控制器
package v1

import (
	"GoHub-Service/pkg/cache"
	"GoHub-Service/pkg/response"
	"github.com/gin-gonic/gin"
)

// CacheMonitorController 缓存监控控制器
type CacheMonitorController struct {
	BaseAPIController
}

// Stats 获取缓存统计信息
func (ctrl *CacheMonitorController) Stats(c *gin.Context) {
	dm := cache.GetDegradationManager()
	stats := dm.GetStats()
	response.JSON(c, stats)
}

// WarmupList 列出预热任务
func (ctrl *CacheMonitorController) WarmupList(c *gin.Context) {
	tasks := cache.GetScheduler().ListTasks()
	response.JSON(c, gin.H{
		"tasks": tasks,
		"count": len(tasks),
	})
}

// WarmupAll 执行所有预热任务
func (ctrl *CacheMonitorController) WarmupAll(c *gin.Context) {
	errors := cache.WarmupAllTasks(c.Request.Context())
	if len(errors) > 0 {
		response.Abort500(c, "预热失败")
		return
	}
	response.Success(c)
}

// WarmupOne 执行单个预热任务
func (ctrl *CacheMonitorController) WarmupOne(c *gin.Context) {
	name := c.Param("name")
	err := cache.GetScheduler().ExecuteOne(c.Request.Context(), name)
	if err != nil {
		response.Abort500(c, err.Error())
		return
	}
	response.Success(c)
}

// Degrade 手动触发降级
func (ctrl *CacheMonitorController) Degrade(c *gin.Context) {
	cache.GetDegradationManager().Degrade()
	response.Success(c)
}

// Recover 恢复正常
func (ctrl *CacheMonitorController) Recover(c *gin.Context) {
	cache.GetDegradationManager().Recover()
	response.Success(c)
}
