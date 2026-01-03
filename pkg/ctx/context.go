// Package ctx Context 助手工具包
package ctx

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// 常用的 Context key 类型
type contextKey string

const (
	// RequestIDKey 请求ID的 context key
	RequestIDKey contextKey = "request_id"
	// UserIDKey 用户ID的 context key
	UserIDKey contextKey = "user_id"
	// TraceIDKey 链路追踪ID的 context key
	TraceIDKey contextKey = "trace_id"
)

// DefaultTimeout 默认超时时间
const DefaultTimeout = 30 * time.Second

// WithRequestID 将请求ID添加到 context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetRequestID 从 context 获取请求ID
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// WithUserID 将用户ID添加到 context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserID 从 context 获取用户ID
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// WithTraceID 将链路追踪ID添加到 context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// GetTraceID 从 context 获取链路追踪ID
func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// FromGinContext 从 Gin Context 创建带超时的 context.Context
func FromGinContext(c *gin.Context) context.Context {
	return c.Request.Context()
}

// WithTimeout 创建带超时的 context
func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return context.WithTimeout(parent, timeout)
}

// WithDefaultTimeout 使用默认超时创建 context
func WithDefaultTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	return WithTimeout(parent, DefaultTimeout)
}

// Background 返回一个非空的 Context，仅用于初始化、测试和特殊场景
// 在大多数业务代码中，应该使用 FromGinContext 或带超时的 context
func Background() context.Context {
	return context.Background()
}

// TODO 返回一个非空的 Context，用于临时代码或TODO标记
// 在生产代码中应该替换为具体的 context
func TODO() context.Context {
	return context.TODO()
}
