// Package middlewares RequestID追踪中间件
package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// RequestIDKey 在gin.Context中存储RequestID的key
	RequestIDKey = "X-Request-ID"
)

// RequestID 为每个请求生成唯一的RequestID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头获取
		requestID := c.GetHeader(RequestIDKey)
		
		// 如果没有，生成新的UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		// 存储到context
		c.Set(RequestIDKey, requestID)
		
		// 设置响应头
		c.Header(RequestIDKey, requestID)
		
		c.Next()
	}
}

// GetRequestID 从gin.Context获取RequestID
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		return requestID.(string)
	}
	return ""
}
