// Package middlewares 性能监控中间件
package middlewares

import (
	"time"

	"GoHub-Service/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PerformanceMonitor 性能监控中间件
func PerformanceMonitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 计算响应时间
		duration := time.Since(startTime)
		
		// 记录响应时间
		logger.InfoContext(c, "api_performance",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Int64("duration_ms", duration.Milliseconds()),
			zap.String("client_ip", c.ClientIP()),
		)

		// 如果响应时间超过阈值，记录警告
		if duration > 1*time.Second {
			logger.WarnContext(c, "slow_api",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int64("duration_ms", duration.Milliseconds()),
				zap.Int64("threshold_ms", 1000),
			)
		}

		// 添加响应时间到响应头
		c.Header("X-Response-Time", duration.String())
	}
}

// APIStats 接口统计
type APIStats struct {
	Path         string        `json:"path"`
	Method       string        `json:"method"`
	Count        int64         `json:"count"`
	TotalTime    time.Duration `json:"total_time"`
	AverageTime  time.Duration `json:"average_time"`
	MaxTime      time.Duration `json:"max_time"`
	MinTime      time.Duration `json:"min_time"`
	ErrorCount   int64         `json:"error_count"`
	SuccessCount int64         `json:"success_count"`
}

// statsRecorder 全局统计记录器
var statsRecorder = make(map[string]*APIStats)

// RecordAPIStats 记录接口统计
func RecordAPIStats(path, method string, duration time.Duration, isError bool) {
	key := method + " " + path
	
	stats, exists := statsRecorder[key]
	if !exists {
		stats = &APIStats{
			Path:    path,
			Method:  method,
			MinTime: duration,
			MaxTime: duration,
		}
		statsRecorder[key] = stats
	}
	
	stats.Count++
	stats.TotalTime += duration
	stats.AverageTime = stats.TotalTime / time.Duration(stats.Count)
	
	if duration > stats.MaxTime {
		stats.MaxTime = duration
	}
	if duration < stats.MinTime {
		stats.MinTime = duration
	}
	
	if isError {
		stats.ErrorCount++
	} else {
		stats.SuccessCount++
	}
}

// GetAPIStats 获取接口统计
func GetAPIStats() map[string]*APIStats {
	return statsRecorder
}

// ClearAPIStats 清除接口统计
func ClearAPIStats() {
	statsRecorder = make(map[string]*APIStats)
}

// PerformanceStats 性能统计中间件（带统计功能）
func PerformanceStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 计算响应时间
		duration := time.Since(startTime)
		
		// 记录统计
		isError := c.Writer.Status() >= 400
		RecordAPIStats(c.Request.URL.Path, c.Request.Method, duration, isError)
		
		// 记录日志
		logger.InfoContext(c, "api_stats",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Int64("duration_ms", duration.Milliseconds()),
		)
	}
}
