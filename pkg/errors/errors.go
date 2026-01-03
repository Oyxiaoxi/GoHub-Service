// Package errors 自定义错误类型和错误码管理
package errors

import (
	"errors"
	"fmt"
	"runtime"
)

// ErrorType 错误类型
type ErrorType string

const (
	ErrorTypeBusiness      ErrorType = "BUSINESS_ERROR"      // 业务错误
	ErrorTypeValidation    ErrorType = "VALIDATION_ERROR"    // 验证错误
	ErrorTypeAuthorization ErrorType = "AUTHORIZATION_ERROR" // 授权错误
	ErrorTypeNotFound      ErrorType = "NOT_FOUND_ERROR"     // 资源不存在
	ErrorTypeDatabase      ErrorType = "DATABASE_ERROR"      // 数据库错误
	ErrorTypeExternal      ErrorType = "EXTERNAL_ERROR"      // 外部服务错误
	ErrorTypeInternal      ErrorType = "INTERNAL_ERROR"      // 内部错误
	ErrorTypeNetwork       ErrorType = "NETWORK_ERROR"       // 网络错误
	ErrorTypeTimeout       ErrorType = "TIMEOUT_ERROR"       // 超时错误
	ErrorTypeConflict      ErrorType = "CONFLICT_ERROR"      // 冲突错误
)

// 错误码定义 (1000-1999: 业务错误, 2000-2999: 系统错误, 3000-3999: 第三方错误)
const (
	// 通用错误码 1000-1099
	CodeSuccess           = 0    // 成功
	CodeUnknownError      = 1000 // 未知错误
	CodeInvalidParameter  = 1001 // 参数错误
	CodeUnauthorized      = 1002 // 未授权
	CodeForbidden         = 1003 // 禁止访问
	CodeNotFound          = 1004 // 资源不存在
	CodeInternalError     = 1005 // 内部错误
	CodeValidationError   = 1006 // 验证错误
	CodeConflict          = 1009 // 冲突（资源已存在）
	CodeTooManyRequests   = 1029 // 请求过于频繁

	// 数据库错误码 2000-2099
	CodeDatabaseError      = 2001 // 数据库错误
	CodeDatabaseConnection = 2002 // 数据库连接失败
	CodeDatabaseQuery      = 2003 // 查询失败
	CodeDatabaseCreate     = 2004 // 创建失败
	CodeDatabaseUpdate     = 2005 // 更新失败
	CodeDatabaseDelete     = 2006 // 删除失败
	CodeDatabaseDuplicate  = 2007 // 重复记录
	CodeDatabaseTransaction = 2008 // 事务失败

	// 缓存错误码 2100-2199
	CodeCacheError     = 2101 // 缓存错误
	CodeCacheMiss      = 2102 // 缓存未命中
	CodeCacheSet       = 2103 // 缓存设置失败
	CodeCacheDelete    = 2104 // 缓存删除失败

	// 网络错误码 2200-2299
	CodeNetworkError   = 2201 // 网络错误
	CodeTimeout        = 2202 // 超时
	CodeServiceUnavailable = 2203 // 服务不可用

	// 第三方服务错误码 3000-3999
	CodeExternalService = 3001 // 外部服务错误
	CodeSMSError        = 3101 // 短信服务错误
	CodeEmailError      = 3102 // 邮件服务错误
	CodeStorageError    = 3103 // 存储服务错误
	CodePaymentError    = 3104 // 支付服务错误
	
	// 业务错误码 4000-4999 (按模块分配)
	// 用户模块 4000-4099
	CodeUserNotFound       = 4001 // 用户不存在
	CodeUserAlreadyExists  = 4002 // 用户已存在
	CodeInvalidCredentials = 4003 // 凭证无效
	CodeUserDisabled       = 4004 // 用户已禁用
	CodeInvalidToken       = 4005 // Token无效
	CodeTokenExpired       = 4006 // Token已过期

	// 话题模块 4100-4199
	CodeTopicNotFound      = 4101 // 话题不存在
	CodeTopicAlreadyExists = 4102 // 话题已存在
	CodeTopicDisabled      = 4103 // 话题已禁用

	// 评论模块 4200-4299
	CodeCommentNotFound      = 4201 // 评论不存在
	CodeCommentAlreadyDeleted = 4202 // 评论已删除

	// 分类模块 4300-4399
	CodeCategoryNotFound      = 4301 // 分类不存在
	CodeCategoryAlreadyExists = 4302 // 分类已存在

	// 权限模块 4400-4499
	CodeRoleNotFound      = 4401 // 角色不存在
	CodeRoleAlreadyExists = 4402 // 角色已存在
	CodePermissionDenied  = 4403 // 权限不足
	CodePermissionNotFound = 4404 // 权限不存在
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

// AuthorizationError 授权错误
func AuthorizationError(message string) *AppError {
	return NewAppError(ErrorTypeAuthorization, CodeForbidden, message, nil)
}

// UnauthorizedError 未授权错误
func UnauthorizedError(message string) *AppError {
	return NewAppError(ErrorTypeAuthorization, CodeUnauthorized, message, nil)
}

// ValidationError 验证错误
func ValidationError(message string, details map[string]interface{}) *AppError {
	err := NewAppError(ErrorTypeValidation, CodeValidationError, message, nil)
	err.Details = details
	return err
}

// NotFoundError 资源不存在错误
func NotFoundError(resource string) *AppError {
	return NewAppError(ErrorTypeNotFound, CodeNotFound, fmt.Sprintf("%s不存在", resource), nil)
}

// NotFoundErrorWithCode 带自定义错误码的资源不存在错误
func NotFoundErrorWithCode(code int, resource string) *AppError {
	return NewAppError(ErrorTypeNotFound, code, fmt.Sprintf("%s不存在", resource), nil)
}

// DatabaseError 数据库错误
func DatabaseError(operation string, err error) *AppError {
	return NewAppError(ErrorTypeDatabase, CodeDatabaseError, fmt.Sprintf("数据库%s失败", operation), err)
}

// DatabaseCreateError 创建失败错误
func DatabaseCreateError(resource string, err error) *AppError {
	return NewAppError(ErrorTypeDatabase, CodeDatabaseCreate, fmt.Sprintf("创建%s失败", resource), err)
}

// DatabaseUpdateError 更新失败错误
func DatabaseUpdateError(resource string, err error) *AppError {
	return NewAppError(ErrorTypeDatabase, CodeDatabaseUpdate, fmt.Sprintf("更新%s失败", resource), err)
}

// DatabaseDeleteError 删除失败错误
func DatabaseDeleteError(resource string, err error) *AppError {
	return NewAppError(ErrorTypeDatabase, CodeDatabaseDelete, fmt.Sprintf("删除%s失败", resource), err)
}

// DatabaseDuplicateError 重复记录错误
func DatabaseDuplicateError(resource string) *AppError {
	return NewAppError(ErrorTypeConflict, CodeDatabaseDuplicate, fmt.Sprintf("%s已存在", resource), nil)
}

// InternalError 内部错误
func InternalError(message string, err error) *AppError {
	return NewAppError(ErrorTypeInternal, CodeInternalError, message, err)
}

// ExternalError 外部服务错误
func ExternalError(service string, err error) *AppError {
	return NewAppError(ErrorTypeExternal, CodeExternalService, fmt.Sprintf("%s服务错误", service), err)
}

// ConflictError 冲突错误（资源已存在）
func ConflictError(resource string) *AppError {
	return NewAppError(ErrorTypeConflict, CodeConflict, fmt.Sprintf("%s已存在", resource), nil)
}

// TimeoutError 超时错误
func TimeoutError(operation string) *AppError {
	return NewAppError(ErrorTypeTimeout, CodeTimeout, fmt.Sprintf("%s超时", operation), nil)
}

// NetworkError 网络错误
func NetworkError(message string, err error) *AppError {
	return NewAppError(ErrorTypeNetwork, CodeNetworkError, message, err)
}

// CacheError 缓存错误
func CacheError(operation string, err error) *AppError {
	return NewAppError(ErrorTypeInternal, CodeCacheError, fmt.Sprintf("缓存%s失败", operation), err)
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
	
	return NewAppError(ErrorTypeInternal, CodeInternalError, message, err)
}

// IsAppError 判断是否为AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError 获取AppError
func GetAppError(err error) (*AppError, bool) {
	if err == nil {
		return nil, false
	}
	
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// Is 判断错误类型（支持errors.Is）
func (e *AppError) Is(target error) bool {
	if target == nil {
		return false
	}
	
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	
	// 比较错误码
	return e.Code == t.Code
}

// As 支持errors.As转换
func (e *AppError) As(target interface{}) bool {
	if t, ok := target.(**AppError); ok {
		*t = e
		return true
	}
	return false
}

// WithError 添加原始错误
func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

// GetCode 获取错误码
func (e *AppError) GetCode() int {
	return e.Code
}

// GetType 获取错误类型
func (e *AppError) GetType() ErrorType {
	return e.Type
}

// IsType 判断错误类型
func (e *AppError) IsType(errorType ErrorType) bool {
	return e.Type == errorType
}

// 预定义的常见错误（用于errors.Is比较）
var (
	ErrNotFound          = &AppError{Type: ErrorTypeNotFound, Code: CodeNotFound, Message: "资源不存在"}
	ErrUnauthorized      = &AppError{Type: ErrorTypeAuthorization, Code: CodeUnauthorized, Message: "未授权"}
	ErrForbidden         = &AppError{Type: ErrorTypeAuthorization, Code: CodeForbidden, Message: "禁止访问"}
	ErrInvalidParameter  = &AppError{Type: ErrorTypeValidation, Code: CodeInvalidParameter, Message: "参数错误"}
	ErrConflict          = &AppError{Type: ErrorTypeConflict, Code: CodeConflict, Message: "资源已存在"}
	ErrDatabaseError     = &AppError{Type: ErrorTypeDatabase, Code: CodeDatabaseError, Message: "数据库错误"}
	ErrInternalError     = &AppError{Type: ErrorTypeInternal, Code: CodeInternalError, Message: "内部错误"}
	ErrTimeout           = &AppError{Type: ErrorTypeTimeout, Code: CodeTimeout, Message: "操作超时"}
)

// New 创建简单错误（兼容标准库errors.New）
func New(message string) *AppError {
	return &AppError{
		Type:       ErrorTypeInternal,
		Code:       CodeInternalError,
		Message:    message,
		StackTrace: captureStackTrace(),
		Details:    make(map[string]interface{}),
	}
}

// Errorf 创建格式化错误（兼容fmt.Errorf）
func Errorf(format string, args ...interface{}) *AppError {
	return &AppError{
		Type:       ErrorTypeInternal,
		Code:       CodeInternalError,
		Message:    fmt.Sprintf(format, args...),
		StackTrace: captureStackTrace(),
		Details:    make(map[string]interface{}),
	}
}
