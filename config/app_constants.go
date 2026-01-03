package config

import "GoHub-Service/pkg/config"

// init 应用常量配置
// 集中管理应用中的硬编码常量，避免魔法数字分散在代码各处
func init() {
	config.Add("app_constants", func() map[string]interface{} {
		return map[string]interface{}{
			// ============= 分页配置 =============
			"pagination": map[string]interface{}{
				"default_per_page": 15,   // 默认每页条数
				"min_per_page":     2,    // 最小每页条数
				"max_per_page":     100,  // 最大每页条数
				"admin_per_page":   20,   // 管理后台默认每页条数
			},

			// ============= 用户配置 =============
			"user": map[string]interface{}{
				"name_min_length":   3,   // 用户名最小长度
				"name_max_length":   20,  // 用户名最大长度
				"city_min_length":   2,   // 城市名最小长度（中文字符）
				"city_max_length":   20,  // 城市名最大长度（中文字符）
				"intro_max_length":  255, // 个人简介最大长度
			},

			// ============= 文件上传配置 =============
			"upload": map[string]interface{}{
				"max_avatar_size_mb": 2, // 头像最大尺寸（MB）
				"max_image_size_mb":  2, // 图片最大尺寸（MB）
				"max_file_size_mb":   5, // 文件最大尺寸（MB）
				"allowed_image_exts": []string{"jpg", "jpeg", "png", "gif", "webp"},
				"allowed_file_exts":  []string{"pdf", "doc", "docx", "xls", "xlsx", "zip", "rar"},
			},

			// ============= 角色权限配置 =============
			"role": map[string]interface{}{
				"name_min_length":         1,   // 角色名最小长度
				"name_max_length":         50,  // 角色名最大长度
				"display_name_min_length": 1,   // 显示名称最小长度
				"display_name_max_length": 100, // 显示名称最大长度
				"description_max_length":  255, // 描述最大长度
			},

			// ============= 评论配置 =============
			"comment": map[string]interface{}{
				"min_length": 2,    // 评论最小长度
				"max_length": 1000, // 评论最大长度
			},

			// ============= 话题配置 =============
			"topic": map[string]interface{}{
				"title_min_length":   3,     // 标题最小长度
				"title_max_length":   100,   // 标题最大长度
				"content_min_length": 10,    // 内容最小长度
				"content_max_length": 50000, // 内容最大长度
			},

			// ============= 时间格式配置 =============
			"datetime": map[string]interface{}{
				"format_iso8601":  "2006-01-02T15:04:05Z07:00", // ISO8601格式
				"format_date":     "2006-01-02",                 // 日期格式
				"format_time":     "15:04:05",                   // 时间格式
				"format_datetime": "2006-01-02 15:04:05",        // 日期时间格式
			},

			// ============= 资源追踪配置 =============
			"resource_tracking": map[string]interface{}{
				"http_timeout_warning_seconds":  30,  // HTTP请求超时警告阈值（秒）
				"resource_leak_threshold_minutes": 5, // 资源泄漏检测阈值（分钟）
				"check_interval_minutes":          1, // 资源检查间隔（分钟）
			},

			// ============= 上下文超时配置 =============
			"context_timeout": map[string]interface{}{
				"user_query_seconds":        5,  // 用户查询超时
				"batch_notify_seconds":      30, // 批量通知超时
				"batch_operation_seconds":   30, // 批量操作超时
				"search_all_seconds":        10, // 搜索全部超时
				"default_operation_seconds": 15, // 默认操作超时
			},

			// ============= 并发池配置 =============
			"goroutine_pool": map[string]interface{}{
				"notification_workers": 20, // 通知服务工作协程数
				"interaction_workers":  30, // 交互服务工作协程数
				"default_workers":      10, // 默认工作协程数
			},
		}
	})
}
