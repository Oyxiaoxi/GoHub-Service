// Package logger 日志追踪辅助函数
package logger

import (
	apperrors "GoHub-Service/pkg/errors"
	
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LogErrorWithContext 记录带上下文的错误日志
func LogErrorWithContext(c *gin.Context, err error, message string, fields ...zap.Field) {
	// 获取RequestID
	requestID := ""
	if c != nil {
		if id, exists := c.Get("X-Request-ID"); exists {
			requestID = id.(string)
		}
	}
	
	// 构建日志字段
	logFields := []zap.Field{
		zap.String("request_id", requestID),
	}
	
	// 如果是AppError，添加额外信息
	if appErr := apperrors.GetAppError(err); appErr != nil {
		logFields = append(logFields,
			zap.String("error_type", string(appErr.Type)),
			zap.Int("error_code", appErr.Code),
			zap.Any("error_details", appErr.Details),
			zap.String("stack_trace", appErr.StackTrace),
		)
		if appErr.Err != nil {
			logFields = append(logFields, zap.Error(appErr.Err))
		}
	} else if err != nil {
		logFields = append(logFields, zap.Error(err))
	}
	
	// 添加自定义字段
	logFields = append(logFields, fields...)
	
	// 记录日志
	Error(message, logFields...)
}

// LogWithRequestID 记录带RequestID的日志
func LogWithRequestID(c *gin.Context, level string, message string, fields ...zap.Field) {
	requestID := ""
	if c != nil {
		if id, exists := c.Get("X-Request-ID"); exists {
			requestID = id.(string)
		}
	}
	
	logFields := append([]zap.Field{zap.String("request_id", requestID)}, fields...)
	
	switch level {
	case "debug":
		Debug(message, logFields...)
	case "info":
		Info(message, logFields...)
	case "warn":
		Warn(message, logFields...)
	case "error":
		Error(message, logFields...)
	default:
		Info(message, logFields...)
	}
}
