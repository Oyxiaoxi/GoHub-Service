// Package errors 自定义错误类型
package errors

import (
	"fmt"
	"runtime"
)

// ErrorType 错误类型
type ErrorType string

const (
	ErrorTypeBusiness       ErrorType = "BUSINESS_ERROR"       // 业务错误
	ErrorTypeValidation     ErrorType = "VALIDATION_ERROR"     // 验证错误
	ErrorTypeAuthorization  ErrorType = "AUTHORIZATION_ERROR"  // 授权错误
	ErrorTypeNotFound       ErrorType = "NOT_FOUND_ERROR"      // 资源不存在
	ErrorTypeDatabase       ErrorType = "DATABASE_ERROR"       // 数据库错误
	ErrorTypeExternal       ErrorType = "EXTERNAL_ERROR"       // 外部服务错误
	ErrorTypeInternal       ErrorType = "INTERNAL_ERROR"       // 内部错误
)

// AppError 应用错误基础结构
type AppError struct {
	Type       ErrorType              // 错误类型
	Code       int                    // 业务错误码
	Message    string                 // 错误消息
	Details    map[string]interface{} // 错误详情
	Err        error                  // 原始错误
	StackTrace string                 // 堆栈信息
	RequestID  string                 // 请求ID
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// Unwrap 支持errors.Unwrap
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithRequestID 添加请求ID
func (e *AppError) WithRequestID(requestID string) *AppError {
	e.RequestID = requestID
	return e
}

// WithDetails 添加错误详情
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

// captureStackTrace 捕获堆栈信息
func captureStackTrace() string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	
	var trace string
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		trace += fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	return trace
}

// NewAppError 创建应用错误
func NewAppError(errorType ErrorType, code int, message string, err error) *AppError {
	return &AppError{
		Type:       errorType,
		Code:       code,
		Message:    message,
		Err:        err,
		StackTrace: captureStackTrace(),
		Details:    make(map[string]interface{}),
	}
}

// BusinessError 业务错误
func BusinessError(code int, message string) *AppError {
	return NewAppError(ErrorTypeBusiness, code, message, nil)
}

// BusinessErrorf 格式化业务错误
func BusinessErrorf(code int, format string, args ...interface{}) *AppError {
	return NewAppError(ErrorTypeBusiness, code, fmt.Sprintf(format, args...), nil)
}

// ValidationError 验证错误
func ValidationError(message string, details map[string]interface{}) *AppError {
	err := NewAppError(ErrorTypeValidation, 1006, message, nil)
	err.Details = details
	return err
}

// AuthorizationError 授权错误
func AuthorizationError(message string) *AppError {
	return NewAppError(ErrorTypeAuthorization, 1003, message, nil)
}

// NotFoundError 资源不存在错误
func NotFoundError(resource string) *AppError {
	return NewAppError(ErrorTypeNotFound, 1004, fmt.Sprintf("%s不存在", resource), nil)
}

// DatabaseError 数据库错误
func DatabaseError(operation string, err error) *AppError {
	return NewAppError(ErrorTypeDatabase, 5001, fmt.Sprintf("数据库%s失败", operation), err)
}

// InternalError 内部错误
func InternalError(message string, err error) *AppError {
	return NewAppError(ErrorTypeInternal, 1005, message, err)
}

// ExternalError 外部服务错误
func ExternalError(service string, err error) *AppError {
	return NewAppError(ErrorTypeExternal, 1000, fmt.Sprintf("%s服务错误", service), err)
}

// WrapError 包装错误
func WrapError(err error, message string) *AppError {
	if err == nil {
		return nil
	}
	
	// 如果已经是AppError，直接返回
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	
	return NewAppError(ErrorTypeInternal, 1005, message, err)
}

// IsAppError 判断是否为AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError 获取AppError
func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return nil
}
