// Package response 统一 API 响应格式
package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// StandardResponse 标准 API 响应格式
type StandardResponse struct {
	Success   bool        `json:"success"`            // 请求是否成功
	Code      int         `json:"code"`               // 业务状态码
	Message   string      `json:"message"`            // 响应消息
	Data      interface{} `json:"data,omitempty"`     // 响应数据
	Error     *ErrorInfo  `json:"error,omitempty"`    // 错误信息（仅失败时）
	Meta      *MetaInfo   `json:"meta,omitempty"`     // 元数据（分页等）
	Timestamp int64       `json:"timestamp"`          // 时间戳
	RequestID string      `json:"request_id"`         // 请求ID
	Version   string      `json:"version,omitempty"`  // API版本
}

// ErrorInfo 错误详情
type ErrorInfo struct {
	Type    string            `json:"type"`              // 错误类型
	Details string            `json:"details,omitempty"` // 错误详情
	Fields  map[string]string `json:"fields,omitempty"`  // 字段错误
}

// MetaInfo 元数据信息
type MetaInfo struct {
	CurrentPage int   `json:"current_page,omitempty"` // 当前页
	PerPage     int   `json:"per_page,omitempty"`     // 每页数量
	Total       int64 `json:"total,omitempty"`        // 总数
	TotalPages  int   `json:"total_pages,omitempty"`  // 总页数
}

// NewStandardResponse 创建标准响应
func NewStandardResponse(c *gin.Context, success bool, code int, message string, data interface{}) *StandardResponse {
	requestID, _ := c.Get("request_id")
	version, _ := c.Get("api_version")
	
	resp := &StandardResponse{
		Success:   success,
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
		RequestID: requestID.(string),
	}
	
	if v, ok := version.(string); ok {
		resp.Version = v
	}
	
	return resp
}

// StandardSuccess 标准成功响应
func StandardSuccess(c *gin.Context, data interface{}) {
	resp := NewStandardResponse(c, true, CodeSuccess, "success", data)
	c.JSON(http.StatusOK, resp)
}

// StandardSuccessWithMessage 标准成功响应（自定义消息）
func StandardSuccessWithMessage(c *gin.Context, message string, data interface{}) {
	resp := NewStandardResponse(c, true, CodeSuccess, message, data)
	c.JSON(http.StatusOK, resp)
}

// StandardSuccessWithMeta 标准成功响应（带分页元数据）
func StandardSuccessWithMeta(c *gin.Context, data interface{}, meta *MetaInfo) {
	resp := NewStandardResponse(c, true, CodeSuccess, "success", data)
	resp.Meta = meta
	c.JSON(http.StatusOK, resp)
}

// StandardError 标准错误响应
func StandardError(c *gin.Context, httpCode, bizCode int, message string) {
	resp := NewStandardResponse(c, false, bizCode, message, nil)
	resp.Error = &ErrorInfo{
		Type:    GetErrorType(bizCode),
		Details: message,
	}
	c.JSON(httpCode, resp)
	c.Abort()
}

// StandardValidationError 标准验证错误响应
func StandardValidationError(c *gin.Context, fields map[string]string) {
	resp := NewStandardResponse(c, false, CodeValidationError, "Validation failed", nil)
	resp.Error = &ErrorInfo{
		Type:   "validation_error",
		Fields: fields,
	}
	c.JSON(http.StatusUnprocessableEntity, resp)
	c.Abort()
}

// GetErrorType 根据错误码获取错误类型
func GetErrorType(code int) string {
	switch {
	case code >= 40000 && code < 41000:
		return "bad_request"
	case code >= 41000 && code < 42000:
		return "unauthorized"
	case code >= 42000 && code < 43000:
		return "forbidden"
	case code >= 43000 && code < 44000:
		return "not_found"
	case code >= 44000 && code < 45000:
		return "conflict"
	case code >= 45000 && code < 50000:
		return "validation_error"
	case code >= 50000:
		return "internal_error"
	default:
		return "unknown"
	}
}
