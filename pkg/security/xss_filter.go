// Package security XSS 防护
package security

import (
	"html"
	"regexp"
	"strings"
)

// XSSFilter XSS 过滤器
type XSSFilter struct {
	allowedTags       map[string]bool
	allowedAttributes map[string]bool
}

// NewXSSFilter 创建 XSS 过滤器
func NewXSSFilter(allowedTags, allowedAttributes []string) *XSSFilter {
	filter := &XSSFilter{
		allowedTags:       make(map[string]bool),
		allowedAttributes: make(map[string]bool),
	}

	for _, tag := range allowedTags {
		filter.allowedTags[strings.ToLower(tag)] = true
	}

	for _, attr := range allowedAttributes {
		filter.allowedAttributes[strings.ToLower(attr)] = true
	}

	return filter
}

// Sanitize 清理 HTML 内容，移除危险标签和属性
func (f *XSSFilter) Sanitize(input string) string {
	if input == "" {
		return input
	}

	// 移除所有 script 标签
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	input = scriptRegex.ReplaceAllString(input, "")

	// 移除所有 iframe 标签
	iframeRegex := regexp.MustCompile(`(?i)<iframe[^>]*>.*?</iframe>`)
	input = iframeRegex.ReplaceAllString(input, "")

	// 移除所有 style 标签
	styleRegex := regexp.MustCompile(`(?i)<style[^>]*>.*?</style>`)
	input = styleRegex.ReplaceAllString(input, "")

	// 移除 javascript: 协议
	jsProtocolRegex := regexp.MustCompile(`(?i)javascript:`)
	input = jsProtocolRegex.ReplaceAllString(input, "")

	// 移除 vbscript: 协议
	vbProtocolRegex := regexp.MustCompile(`(?i)vbscript:`)
	input = vbProtocolRegex.ReplaceAllString(input, "")

	// 移除 data: 协议（可能包含 base64 编码的恶意代码）
	dataProtocolRegex := regexp.MustCompile(`(?i)data:text/html`)
	input = dataProtocolRegex.ReplaceAllString(input, "")

	// 移除 on* 事件处理器
	eventRegex := regexp.MustCompile(`(?i)\s*on\w+\s*=\s*['"][^'"]*['"]`)
	input = eventRegex.ReplaceAllString(input, "")

	return input
}

// Escape HTML 转义
func (f *XSSFilter) Escape(input string) string {
	return html.EscapeString(input)
}

// StripTags 移除所有 HTML 标签
func (f *XSSFilter) StripTags(input string) string {
	// 移除所有 HTML 标签
	tagRegex := regexp.MustCompile(`<[^>]*>`)
	return tagRegex.ReplaceAllString(input, "")
}

// StripTagsExcept 移除除允许列表外的所有 HTML 标签
func (f *XSSFilter) StripTagsExcept(input string, allowedTags []string) string {
	if len(allowedTags) == 0 {
		return f.StripTags(input)
	}

	// 构建允许标签的正则表达式
	allowed := make(map[string]bool)
	for _, tag := range allowedTags {
		allowed[strings.ToLower(tag)] = true
	}

	// 查找所有标签
	tagRegex := regexp.MustCompile(`</?([a-zA-Z][a-zA-Z0-9]*)[^>]*>`)
	result := tagRegex.ReplaceAllStringFunc(input, func(match string) string {
		// 提取标签名
		tagMatch := regexp.MustCompile(`</?([a-zA-Z][a-zA-Z0-9]*)`).FindStringSubmatch(match)
		if len(tagMatch) > 1 {
			tagName := strings.ToLower(tagMatch[1])
			if allowed[tagName] {
				return match // 保留允许的标签
			}
		}
		return "" // 移除不允许的标签
	})

	return result
}

// RemoveDangerousAttributes 移除危险的 HTML 属性
func (f *XSSFilter) RemoveDangerousAttributes(input string) string {
	// 移除 on* 事件处理器
	eventRegex := regexp.MustCompile(`(?i)\s*on\w+\s*=\s*(['"])[^'"]*\1`)
	input = eventRegex.ReplaceAllString(input, "")

	// 移除 style 属性中的 expression
	styleExprRegex := regexp.MustCompile(`(?i)style\s*=\s*(['"])[^'"]*expression[^'"]*\1`)
	input = styleExprRegex.ReplaceAllString(input, "")

	return input
}

// CleanURL 清理 URL，防止 XSS
func (f *XSSFilter) CleanURL(url string) string {
	url = strings.TrimSpace(url)

	// 检查是否包含危险协议
	dangerousProtocols := []string{
		"javascript:",
		"vbscript:",
		"data:text/html",
		"data:text/javascript",
	}

	lowerURL := strings.ToLower(url)
	for _, protocol := range dangerousProtocols {
		if strings.HasPrefix(lowerURL, protocol) {
			return ""
		}
	}

	return url
}

// ValidateInput 验证输入是否包含潜在的 XSS 攻击
func (f *XSSFilter) ValidateInput(input string) bool {
	if input == "" {
		return true
	}

	// 检查危险模式
	dangerousPatterns := []string{
		`<script`,
		`javascript:`,
		`vbscript:`,
		`onload=`,
		`onerror=`,
		`onclick=`,
		`<iframe`,
		`<embed`,
		`<object`,
	}

	lowerInput := strings.ToLower(input)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerInput, pattern) {
			return false
		}
	}

	return true
}
