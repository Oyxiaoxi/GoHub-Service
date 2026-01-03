// Package middlewares 安全增强中间件
package middlewares

import (
	"GoHub-Service/pkg/response"
	"GoHub-Service/pkg/security"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	inputValidator     *security.InputValidator
	inputValidatorOnce sync.Once
)

// getInputValidator 获取输入验证器单例
func getInputValidator() *security.InputValidator {
	inputValidatorOnce.Do(func() {
		inputValidator = security.NewInputValidator()
	})
	return inputValidator
}

// EnhancedSecurityValidation 增强安全验证中间件
func EnhancedSecurityValidation() gin.HandlerFunc {
	validator := getInputValidator()

	return func(c *gin.Context) {
		// 1. 验证查询参数
		for _, values := range c.Request.URL.Query() {
			for _, value := range values {
				result := validator.Validate(value)
				if !result.IsValid {
					response.ApiError(c, 400, response.CodeInvalidParams, 
						"请求参数包含非法内容: "+result.Reason)
					return
				}
			}
		}

		// 2. 验证路径参数
		for _, param := range c.Params {
			result := validator.Validate(param.Value)
			if !result.IsValid {
				response.ApiError(c, 400, response.CodeInvalidParams, 
					"路径参数包含非法内容: "+result.Reason)
				return
			}
		}

		c.Next()
	}
}

// EnhancedSQLInjectionProtection 增强型 SQL 注入防护
func EnhancedSQLInjectionProtection() gin.HandlerFunc {
	validator := getInputValidator()

	return func(c *gin.Context) {
		// 检查所有查询参数和路径参数
		allParams := make(map[string]string)
		
		// 添加查询参数
		for _, values := range c.Request.URL.Query() {
			for _, value := range values {
				allParams[value] = value
			}
		}
		
		// 添加路径参数
		for _, param := range c.Params {
			allParams[param.Key] = param.Value
		}

		// 验证每个参数
		for _, value := range allParams {
			result := validator.CheckSQLInjection(value)
			if !result.IsValid {
				response.ApiError(c, 400, response.CodeInvalidParams, 
					"检测到潜在的SQL注入攻击")
				return
			}
		}

		c.Next()
	}
}

// XSSProtection 增强型 XSS 防护
func EnhancedXSSProtection() gin.HandlerFunc {
	validator := getInputValidator()

	return func(c *gin.Context) {
		// 检查查询参数
		for _, values := range c.Request.URL.Query() {
			for _, value := range values {
				result := validator.CheckXSS(value)
				if !result.IsValid {
					response.ApiError(c, 400, response.CodeInvalidParams, 
						"检测到潜在的XSS攻击")
					return
				}
			}
		}

		c.Next()
	}
}

// RateLimitByIP IP 级别的限流（增强版）
type IPRateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int           // 每分钟允许的请求数
	window   time.Duration // 时间窗口
}

type visitor struct {
	count      int
	lastSeen   time.Time
	blockedUntil time.Time
}

// NewIPRateLimiter 创建 IP 限流器
func NewIPRateLimiter(requestsPerMinute int) *IPRateLimiter {
	limiter := &IPRateLimiter{
		visitors: make(map[string]*visitor),
		rate:     requestsPerMinute,
		window:   time.Minute,
	}

	// 定期清理过期访客
	go limiter.cleanupVisitors()

	return limiter
}

func (l *IPRateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		now := time.Now()
		for ip, v := range l.visitors {
			if now.Sub(v.lastSeen) > 10*time.Minute {
				delete(l.visitors, ip)
			}
		}
		l.mu.Unlock()
	}
}

// Allow 检查是否允许请求
func (l *IPRateLimiter) Allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	v, exists := l.visitors[ip]

	if !exists {
		l.visitors[ip] = &visitor{
			count:    1,
			lastSeen: now,
		}
		return true
	}

	// 检查是否在封禁期
	if now.Before(v.blockedUntil) {
		return false
	}

	// 重置计数器（新的时间窗口）
	if now.Sub(v.lastSeen) > l.window {
		v.count = 1
		v.lastSeen = now
		return true
	}

	// 增加计数
	v.count++
	v.lastSeen = now

	// 检查是否超过限制
	if v.count > l.rate {
		// 封禁 1 分钟
		v.blockedUntil = now.Add(time.Minute)
		return false
	}

	return true
}

// RateLimitMiddleware 创建限流中间件
func RateLimitMiddleware(requestsPerMinute int) gin.HandlerFunc {
	limiter := NewIPRateLimiter(requestsPerMinute)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.Allow(ip) {
			response.ApiError(c, 429, response.CodeTooManyRequests, 
				"请求过于频繁，请稍后再试")
			return
		}

		c.Next()
	}
}
