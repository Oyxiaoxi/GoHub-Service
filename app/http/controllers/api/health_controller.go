package api

import (
	"time"

	"GoHub-Service/pkg/database"
	"GoHub-Service/pkg/redis"

	"github.com/gin-gonic/gin"
)

// HealthController 健康检查控制器
type HealthController struct{}

// Health 基础健康检查
// @Summary 健康检查
// @Description 返回服务健康状态
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (ctrl *HealthController) Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0",
		"service":   "GoHub-Service",
	})
}

// Readiness 就绪探针
// @Summary 就绪探针
// @Description 检查服务是否就绪（数据库、Redis等）
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /readiness [get]
func (ctrl *HealthController) Readiness(c *gin.Context) {
	checks := make(map[string]bool)
	allReady := true

	// 检查数据库连接
	sqlDB, err := database.DB.DB()
	if err != nil || sqlDB.Ping() != nil {
		checks["database"] = false
		allReady = false
	} else {
		checks["database"] = true
	}

	// 检查 Redis 连接
	if redis.Redis != nil && redis.Redis.Client != nil {
		err := redis.Redis.Client.Ping(c.Request.Context()).Err()
		if err != nil {
			checks["redis"] = false
			allReady = false
		} else {
			checks["redis"] = true
		}
	} else {
		checks["redis"] = false
		allReady = false
	}

	status := 200
	message := "ready"
	if !allReady {
		status = 503
		message = "not ready"
	}

	c.JSON(status, gin.H{
		"status":    message,
		"checks":    checks,
		"timestamp": time.Now().Unix(),
	})
}

// Liveness 存活探针
// @Summary 存活探针
// @Description 检查服务是否存活
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /liveness [get]
func (ctrl *HealthController) Liveness(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":    "alive",
		"timestamp": time.Now().Unix(),
	})
}
