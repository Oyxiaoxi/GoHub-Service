// Package middlewares 内容安全中间件
package middlewares

import (
	"GoHub-Service/pkg/response"
	"GoHub-Service/pkg/security"

	"github.com/gin-gonic/gin"
)

// ContentSecurity 内容安全检查中间件
func ContentSecurity() gin.HandlerFunc {
	checker := security.GetContentChecker()

	return func(c *gin.Context) {
		// 对于某些方法，跳过检查
		if c.Request.Method == "GET" || c.Request.Method == "DELETE" {
			c.Next()
			return
		}

		// 检查常见的文本字段
		fieldsToCheck := []string{
			"title", "content", "body", "description",
			"name", "nickname", "introduction", "signature",
		}

		for _, field := range fieldsToCheck {
			if value, exists := c.GetPostForm(field); exists && value != "" {
				// 根据字段类型选择不同的检查方法
				var result *security.CheckResult
				if field == "title" || field == "name" {
					result = checker.CheckTitle(value)
				} else {
					result = checker.CheckContent(value)
				}

				if !result.IsValid {
					response.Abort403(c, result.Message)
					return
				}

				// 如果内容被过滤，更新表单值
				if result.FilteredText != value {
					c.Request.Form.Set(field, result.FilteredText)
				}

				// 记录敏感词
				if len(result.FoundWords) > 0 {
					c.Set("sensitive_words_found", result.FoundWords)
				}
			}
		}

		c.Next()
	}
}

// SensitiveWordFilter 敏感词过滤中间件（仅过滤，不拦截）
func SensitiveWordFilter() gin.HandlerFunc {
	checker := security.GetContentChecker()

	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "DELETE" {
			c.Next()
			return
		}

		// 过滤文本字段
		fieldsToFilter := []string{
			"title", "content", "body", "description",
			"name", "nickname", "introduction", "signature",
		}

		foundWords := make([]string, 0)
		for _, field := range fieldsToFilter {
			if value, exists := c.GetPostForm(field); exists && value != "" {
				// 查找敏感词
				words := checker.FindSensitiveWords(value)
				if len(words) > 0 {
					foundWords = append(foundWords, words...)
					// 过滤敏感词
					filtered := checker.FilterSensitiveWords(value)
					c.Request.Form.Set(field, filtered)
				}
			}
		}

		if len(foundWords) > 0 {
			c.Set("sensitive_words_filtered", foundWords)
		}

		c.Next()
	}
}

// ContentXSSProtection 内容 XSS 防护中间件（仅清理，不验证）
func ContentXSSProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "DELETE" {
			c.Next()
			return
		}

		// 检查所有表单字段
		for key, values := range c.Request.Form {
			for i, value := range values {
				// 清理内容
				checker := security.GetContentChecker()
				cleaned := checker.CleanHTML(value)
				if cleaned != value {
					c.Request.Form[key][i] = cleaned
				}
			}
		}

		c.Next()
	}
}

// ImageUploadSecurity 图片上传安全检查中间件
func ImageUploadSecurity() gin.HandlerFunc {
	checker := security.GetContentChecker()

	return func(c *gin.Context) {
		// 只检查包含文件上传的请求
		form, err := c.MultipartForm()
		if err != nil {
			c.Next()
			return
		}

		// 检查所有上传的文件
		for _, files := range form.File {
			for _, file := range files {
				result := checker.CheckImage(file)
				if !result.IsValid {
					response.Abort403(c, result.Message)
					return
				}
			}
		}

		c.Next()
	}
}
