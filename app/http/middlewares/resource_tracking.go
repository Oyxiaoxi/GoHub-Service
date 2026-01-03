package middlewares

import (
	"GoHub-Service/bootstrap"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ResourceTracking 资源追踪中间件
// 追踪HTTP请求的生命周期，检测响应是否及时返回
func ResourceTracking() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成资源ID
		requestID := GetRequestID(c)
		resourceID := fmt.Sprintf("http_request_%s", requestID)

		// 开始追踪
		bootstrap.Tracker.Track(resourceID, "http_request")

		// 记录开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 停止追踪
		bootstrap.Tracker.Untrack(resourceID)

		// 记录处理时间
		duration := time.Since(startTime)

		// 如果处理时间过长，记录警告
		if duration > 30*time.Second {
			logger := bootstrap.Logger
			if logger != nil {
				logger.Warn("HTTP请求处理时间过长",
					zap.String("request_id", requestID),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.Duration("duration", duration),
				)
			}
		}
	}
}

// DBConnectionTracking DB连接追踪包装器
// 用于追踪数据库连接的生命周期
func TrackDBConnection(connectionID string) func() {
	resourceID := fmt.Sprintf("db_connection_%s", connectionID)
	bootstrap.Tracker.Track(resourceID, "db_connection")

	return func() {
		bootstrap.Tracker.Untrack(resourceID)
	}
}

// GoroutineTracking Goroutine追踪包装器
// 用于追踪后台goroutine的生命周期
func TrackGoroutine(goroutineID string) func() {
	resourceID := fmt.Sprintf("goroutine_%s", goroutineID)
	bootstrap.Tracker.Track(resourceID, "goroutine")

	return func() {
		bootstrap.Tracker.Untrack(resourceID)
	}
}
