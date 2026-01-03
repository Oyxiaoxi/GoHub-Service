// Package logger 日志追踪辅助函数
package logger

import (
	"context"
	apperrors "GoHub-Service/pkg/errors"
	pkgctx "GoHub-Service/pkg/ctx"
	
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	if appErr, ok := apperrors.GetAppError(err); ok && appErr != nil {
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

// ContextLogger 带 context 的日志记录器
type ContextLogger struct {
	logger *zap.Logger
	ctx    context.Context
}

// WithContext 创建带 context 的 logger
func WithContext(c context.Context) *ContextLogger {
	return &ContextLogger{
		logger: Logger,
		ctx:    c,
	}
}

// fields 获取 context 中的结构化字段
func (l *ContextLogger) fields() []zapcore.Field {
	fields := []zapcore.Field{}
	
	// 添加 TraceID
	if traceID := pkgctx.GetTraceID(l.ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}
	
	// 添加 RequestID
	if requestID := pkgctx.GetRequestID(l.ctx); requestID != "" {
		fields = append(fields, zap.String("request_id", requestID))
	}
	
	// 添加 UserID
	if userID := pkgctx.GetUserID(l.ctx); userID != "" {
		fields = append(fields, zap.String("user_id", userID))
	}
	
	return fields
}

// Debug 日志，自动过滤敏感信息
func (l *ContextLogger) Debug(msg string, fields ...zapcore.Field) {
	l.logger.Debug(FilterSensitive(msg), append(l.fields(), fields...)...)
}

// Info 日志，自动过滤敏感信息
func (l *ContextLogger) Info(msg string, fields ...zapcore.Field) {
	l.logger.Info(FilterSensitive(msg), append(l.fields(), fields...)...)
}

// Warn 日志，自动过滤敏感信息
func (l *ContextLogger) Warn(msg string, fields ...zapcore.Field) {
	l.logger.Warn(FilterSensitive(msg), append(l.fields(), fields...)...)
}

// Error 日志，自动过滤敏感信息
func (l *ContextLogger) Error(msg string, fields ...zapcore.Field) {
	l.logger.Error(FilterSensitive(msg), append(l.fields(), fields...)...)
}

// Fatal 日志，自动过滤敏感信息
func (l *ContextLogger) Fatal(msg string, fields ...zapcore.Field) {
	l.logger.Fatal(FilterSensitive(msg), append(l.fields(), fields...)...)
}

// Panic 日志，自动过滤敏感信息
func (l *ContextLogger) Panic(msg string, fields ...zapcore.Field) {
	l.logger.Panic(FilterSensitive(msg), append(l.fields(), fields...)...)
}

// WithFields 添加自定义字段
func (l *ContextLogger) WithFields(fields ...zapcore.Field) *ContextLogger {
	return &ContextLogger{
		logger: l.logger.With(fields...),
		ctx:    l.ctx,
	}
}

// 便捷方法：直接使用全局 Logger + Context
func DebugContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	WithContext(ctx).Debug(msg, fields...)
}

func InfoContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	WithContext(ctx).Info(msg, fields...)
}

func WarnContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	WithContext(ctx).Warn(msg, fields...)
}

func ErrorContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	WithContext(ctx).Error(msg, fields...)
}

func FatalContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	WithContext(ctx).Fatal(msg, fields...)
}

func PanicContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	WithContext(ctx).Panic(msg, fields...)
}
