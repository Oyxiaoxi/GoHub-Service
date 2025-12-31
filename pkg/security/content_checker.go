// Package security 内容安全检查器
package security

import (
	"GoHub-Service/pkg/config"
	"mime/multipart"
	"path/filepath"
	"strings"
)

// ContentChecker 内容安全检查器
type ContentChecker struct {
	sensitiveFilter *SensitiveWordFilter
	xssFilter       *XSSFilter
	config          ContentCheckConfig
}

// ContentCheckConfig 内容检查配置
type ContentCheckConfig struct {
	Enabled              bool
	SensitiveWordEnabled bool
	XSSProtectionEnabled bool
	MaxContentLength     int
	MaxTitleLength       int
	AllowHTMLTags        bool
	AllowedHTMLTags      []string
	AllowedHTMLAttrs     []string
	AllowedImageTypes    []string
	MaxImageSize         int64
}

var contentChecker *ContentChecker

// InitContentChecker 初始化内容检查器
func InitContentChecker() *ContentChecker {
	if contentChecker != nil {
		return contentChecker
	}

	// 默认允许的 HTML 标签
	allowedTags := []string{
		"p", "br", "strong", "em", "u", "a", "img",
		"ul", "ol", "li", "blockquote", "code", "pre",
	}

	// 默认允许的 HTML 属性
	allowedAttrs := []string{
		"href", "src", "alt", "title", "class",
	}

	// 默认允许的图片类型
	allowedImages := []string{"jpg", "jpeg", "png", "gif", "webp"}

	cfg := ContentCheckConfig{
		Enabled:              config.GetBool("security.content_check_enabled"),
		SensitiveWordEnabled: config.GetBool("security.sensitive_word_filter_enabled"),
		XSSProtectionEnabled: config.GetBool("security.xss_protection_enabled"),
		MaxContentLength:     config.GetInt("security.max_content_length"),
		MaxTitleLength:       config.GetInt("security.max_title_length"),
		AllowHTMLTags:        config.GetBool("security.allow_html_tags"),
		AllowedHTMLTags:      allowedTags,
		AllowedHTMLAttrs:     allowedAttrs,
		AllowedImageTypes:    allowedImages,
		MaxImageSize:         int64(config.GetInt("security.max_image_size")),
	}

	replacement := config.GetString("security.sensitive_word_replacement")
	contentChecker = &ContentChecker{
		sensitiveFilter: GetFilter(),
		xssFilter:       NewXSSFilter(cfg.AllowedHTMLTags, cfg.AllowedHTMLAttrs),
		config:          cfg,
	}

	contentChecker.sensitiveFilter.SetReplacement(replacement)

	return contentChecker
}

// GetContentChecker 获取内容检查器实例
func GetContentChecker() *ContentChecker {
	if contentChecker == nil {
		return InitContentChecker()
	}
	return contentChecker
}

// CheckResult 检查结果
type CheckResult struct {
	IsValid      bool
	Message      string
	FilteredText string
	FoundWords   []string
}

// CheckText 检查文本内容
func (c *ContentChecker) CheckText(text string, maxLength int) *CheckResult {
	result := &CheckResult{
		IsValid:      true,
		FilteredText: text,
	}

	if !c.config.Enabled {
		return result
	}

	// 检查长度
	if maxLength > 0 && len(text) > maxLength {
		result.IsValid = false
		result.Message = "内容长度超出限制"
		return result
	}

	// XSS 防护
	if c.config.XSSProtectionEnabled {
		if !c.xssFilter.ValidateInput(text) {
			result.IsValid = false
			result.Message = "内容包含不安全的代码"
			return result
		}

		// 清理 HTML
		if !c.config.AllowHTMLTags {
			result.FilteredText = c.xssFilter.StripTags(result.FilteredText)
		} else {
			result.FilteredText = c.xssFilter.Sanitize(result.FilteredText)
			result.FilteredText = c.xssFilter.RemoveDangerousAttributes(result.FilteredText)
		}
	}

	// 敏感词过滤
	if c.config.SensitiveWordEnabled {
		foundWords := c.sensitiveFilter.FindAll(result.FilteredText)
		if len(foundWords) > 0 {
			result.FoundWords = foundWords
			result.FilteredText = c.sensitiveFilter.Filter(result.FilteredText)
			// 注意：这里不设置 IsValid = false，而是自动过滤
			// 如果需要拒绝包含敏感词的内容，可以取消下面的注释
			// result.IsValid = false
			// result.Message = "内容包含敏感词"
			// return result
		}
	}

	return result
}

// CheckTitle 检查标题
func (c *ContentChecker) CheckTitle(title string) *CheckResult {
	return c.CheckText(title, c.config.MaxTitleLength)
}

// CheckContent 检查内容
func (c *ContentChecker) CheckContent(content string) *CheckResult {
	return c.CheckText(content, c.config.MaxContentLength)
}

// CheckImage 检查图片上传
func (c *ContentChecker) CheckImage(file *multipart.FileHeader) *CheckResult {
	result := &CheckResult{
		IsValid: true,
	}

	if !c.config.Enabled {
		return result
	}

	// 检查文件大小
	if file.Size > c.config.MaxImageSize {
		result.IsValid = false
		result.Message = "图片文件过大"
		return result
	}

	// 检查文件类型
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file.Filename), "."))
	isAllowed := false
	for _, allowedExt := range c.config.AllowedImageTypes {
		if ext == allowedExt {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		result.IsValid = false
		result.Message = "不支持的图片格式"
		return result
	}

	return result
}

// CleanHTML 清理 HTML 内容
func (c *ContentChecker) CleanHTML(html string) string {
	if !c.config.AllowHTMLTags {
		return c.xssFilter.StripTags(html)
	}
	return c.xssFilter.Sanitize(html)
}

// EscapeHTML HTML 转义
func (c *ContentChecker) EscapeHTML(text string) string {
	return c.xssFilter.Escape(text)
}

// ContainsSensitiveWord 检查是否包含敏感词
func (c *ContentChecker) ContainsSensitiveWord(text string) bool {
	if !c.config.SensitiveWordEnabled {
		return false
	}
	return c.sensitiveFilter.Contains(text)
}

// FilterSensitiveWords 过滤敏感词
func (c *ContentChecker) FilterSensitiveWords(text string) string {
	if !c.config.SensitiveWordEnabled {
		return text
	}
	return c.sensitiveFilter.Filter(text)
}

// FindSensitiveWords 查找所有敏感词
func (c *ContentChecker) FindSensitiveWords(text string) []string {
	if !c.config.SensitiveWordEnabled {
		return []string{}
	}
	return c.sensitiveFilter.FindAll(text)
}
