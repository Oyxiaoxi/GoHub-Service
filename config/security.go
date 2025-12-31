// Package config 安全配置
package config

import "GoHub-Service/pkg/config"

func init() {
	config.Add("security", func() map[string]interface{} {
		return map[string]interface{}{
			// 是否启用内容安全检查
			"content_check_enabled": config.Env("CONTENT_CHECK_ENABLED", true),

			// 是否启用敏感词过滤
			"sensitive_word_filter_enabled": config.Env("SENSITIVE_WORD_FILTER_ENABLED", true),

			// 敏感词替换字符
			"sensitive_word_replacement": config.Env("SENSITIVE_WORD_REPLACEMENT", "***"),

			// 是否启用 XSS 防护
			"xss_protection_enabled": config.Env("XSS_PROTECTION_ENABLED", true),

			// 内容最大长度限制
			"max_content_length": config.Env("MAX_CONTENT_LENGTH", 10000),

			// 标题最大长度限制
			"max_title_length": config.Env("MAX_TITLE_LENGTH", 200),

			// 是否允许 HTML 标签
			"allow_html_tags": config.Env("ALLOW_HTML_TAGS", false),

			// 允许的 HTML 标签白名单（如果启用 HTML）
			"allowed_html_tags": []string{
				"p", "br", "strong", "em", "u", "a", "img",
				"ul", "ol", "li", "blockquote", "code", "pre",
			},

			// 允许的 HTML 属性白名单
			"allowed_html_attributes": []string{
				"href", "src", "alt", "title", "class",
			},

			// 图片上传安全检查
			"image_check_enabled": config.Env("IMAGE_CHECK_ENABLED", true),

			// 允许的图片格式
			"allowed_image_types": []string{"jpg", "jpeg", "png", "gif", "webp"},

			// 图片最大尺寸（字节）
			"max_image_size": config.Env("MAX_IMAGE_SIZE", 2*1024*1024), // 2MB

			// 是否启用内容审核日志
			"audit_log_enabled": config.Env("AUDIT_LOG_ENABLED", true),
		}
	})
}
