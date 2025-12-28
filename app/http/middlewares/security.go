package middlewares

import (
	"html"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// SecureHeaders 添加安全响应头中间件
func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止点击劫持攻击
		c.Header("X-Frame-Options", "DENY")

		// 防止 MIME 类型嗅探
		c.Header("X-Content-Type-Options", "nosniff")

		// 启用浏览器 XSS 防护
		c.Header("X-XSS-Protection", "1; mode=block")

		// 内容安全策略 (CSP)
		c.Header("Content-Security-Policy", "default-src 'self'")

		// 强制使用 HTTPS (仅生产环境启用)
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// 控制 Referrer 信息
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// 权限策略
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}

// XSSProtection XSS 防护中间件 - 对输入数据进行清理
func XSSProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 处理查询参数
		query := c.Request.URL.Query()
		for key, values := range query {
			for i, value := range values {
				query[key][i] = sanitizeInput(value)
			}
		}
		c.Request.URL.RawQuery = query.Encode()

		// 注意：对于 POST body，应该在 binding 后在业务层处理
		// 这里只处理 URL 参数，避免影响 JSON 解析

		c.Next()
	}
}

// sanitizeInput 清理输入数据，防止 XSS 攻击
func sanitizeInput(input string) string {
	// HTML 实体转义
	input = html.EscapeString(input)

	// 移除潜在的脚本标签
	scriptPattern := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	input = scriptPattern.ReplaceAllString(input, "")

	// 移除事件处理器
	eventPattern := regexp.MustCompile(`(?i)on\w+\s*=`)
	input = eventPattern.ReplaceAllString(input, "")

	// 移除 javascript: 协议
	input = strings.ReplaceAll(input, "javascript:", "")
	input = strings.ReplaceAll(input, "Javascript:", "")
	input = strings.ReplaceAll(input, "JAVASCRIPT:", "")

	return input
}

// SQLInjectionProtection SQL 注入防护检查中间件
// 注意：GORM 已有基础防护，这只是额外检查
func SQLInjectionProtection() gin.HandlerFunc {
	// SQL 注入常见关键词模式
	sqlPattern := regexp.MustCompile(`(?i)(union|select|insert|update|delete|drop|create|alter|exec|script|javascript|<script|</script>)`)

	return func(c *gin.Context) {
		// 检查查询参数
		for _, values := range c.Request.URL.Query() {
			for _, value := range values {
				if sqlPattern.MatchString(value) {
					c.JSON(400, gin.H{
						"code":    400,
						"message": "非法请求：检测到潜在的安全威胁",
					})
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}

// ContentTypeValidation 验证 Content-Type 中间件
func ContentTypeValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 对于需要 body 的请求方法，验证 Content-Type
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")

			// 允许的 Content-Type 列表
			allowedTypes := []string{
				"application/json",
				"application/x-www-form-urlencoded",
				"multipart/form-data",
			}

			valid := false
			for _, allowedType := range allowedTypes {
				if strings.Contains(contentType, allowedType) {
					valid = true
					break
				}
			}

			if !valid && contentType != "" {
				c.JSON(415, gin.H{
					"code":    415,
					"message": "不支持的 Content-Type",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
