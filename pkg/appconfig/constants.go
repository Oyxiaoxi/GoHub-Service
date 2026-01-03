// Package appconfig 应用配置辅助工具
// 提供便捷的配置访问方法，避免在代码中直接使用魔法字符串
package appconfig

import "GoHub-Service/pkg/config"

// ============= 分页配置 =============

// GetDefaultPerPage 获取默认每页条数
func GetDefaultPerPage() int {
	return config.GetInt("app_constants.pagination.default_per_page", 15)
}

// GetMinPerPage 获取最小每页条数
func GetMinPerPage() int {
	return config.GetInt("app_constants.pagination.min_per_page", 2)
}

// GetMaxPerPage 获取最大每页条数
func GetMaxPerPage() int {
	return config.GetInt("app_constants.pagination.max_per_page", 100)
}

// GetAdminPerPage 获取管理后台默认每页条数
func GetAdminPerPage() int {
	return config.GetInt("app_constants.pagination.admin_per_page", 20)
}

// ValidatePerPage 验证并修正每页条数
func ValidatePerPage(perPage int) int {
	min := GetMinPerPage()
	max := GetMaxPerPage()
	
	if perPage < min {
		return GetDefaultPerPage()
	}
	if perPage > max {
		return max
	}
	return perPage
}

// ============= 用户配置 =============

// GetUserNameMinLength 获取用户名最小长度
func GetUserNameMinLength() int {
	return config.GetInt("app_constants.user.name_min_length", 3)
}

// GetUserNameMaxLength 获取用户名最大长度
func GetUserNameMaxLength() int {
	return config.GetInt("app_constants.user.name_max_length", 20)
}

// GetCityMinLength 获取城市名最小长度
func GetCityMinLength() int {
	return config.GetInt("app_constants.user.city_min_length", 2)
}

// GetCityMaxLength 获取城市名最大长度
func GetCityMaxLength() int {
	return config.GetInt("app_constants.user.city_max_length", 20)
}

// GetIntroMaxLength 获取个人简介最大长度
func GetIntroMaxLength() int {
	return config.GetInt("app_constants.user.intro_max_length", 255)
}

// ============= 文件上传配置 =============

// GetMaxAvatarSizeMB 获取头像最大尺寸（MB）
func GetMaxAvatarSizeMB() int {
	return config.GetInt("app_constants.upload.max_avatar_size_mb", 2)
}

// GetMaxImageSizeMB 获取图片最大尺寸（MB）
func GetMaxImageSizeMB() int {
	return config.GetInt("app_constants.upload.max_image_size_mb", 20)
}

// GetMaxFileSizeMB 获取文件最大尺寸（MB）
func GetMaxFileSizeMB() int {
	return config.GetInt("app_constants.upload.max_file_size_mb", 50)
}

// GetMaxAvatarSizeBytes 获取头像最大尺寸（字节）
func GetMaxAvatarSizeBytes() int {
	return GetMaxAvatarSizeMB() * 1024 * 1024
}

// GetMaxImageSizeBytes 获取图片最大尺寸（字节）
func GetMaxImageSizeBytes() int {
	return GetMaxImageSizeMB() * 1024 * 1024
}

// GetMaxFileSizeBytes 获取文件最大尺寸（字节）
func GetMaxFileSizeBytes() int {
	return GetMaxFileSizeMB() * 1024 * 1024
}

// ============= 角色权限配置 =============

// GetRoleNameMinLength 获取角色名最小长度
func GetRoleNameMinLength() int {
	return config.GetInt("app_constants.role.name_min_length", 1)
}

// GetRoleNameMaxLength 获取角色名最大长度
func GetRoleNameMaxLength() int {
	return config.GetInt("app_constants.role.name_max_length", 50)
}

// GetRoleDisplayNameMinLength 获取角色显示名称最小长度
func GetRoleDisplayNameMinLength() int {
	return config.GetInt("app_constants.role.display_name_min_length", 1)
}

// GetRoleDisplayNameMaxLength 获取角色显示名称最大长度
func GetRoleDisplayNameMaxLength() int {
	return config.GetInt("app_constants.role.display_name_max_length", 100)
}

// GetRoleDescriptionMaxLength 获取角色描述最大长度
func GetRoleDescriptionMaxLength() int {
	return config.GetInt("app_constants.role.description_max_length", 255)
}

// ============= 评论配置 =============

// GetCommentMinLength 获取评论最小长度
func GetCommentMinLength() int {
	return config.GetInt("app_constants.comment.min_length", 2)
}

// GetCommentMaxLength 获取评论最大长度
func GetCommentMaxLength() int {
	return config.GetInt("app_constants.comment.max_length", 1000)
}

// ============= 话题配置 =============

// GetTopicTitleMinLength 获取话题标题最小长度
func GetTopicTitleMinLength() int {
	return config.GetInt("app_constants.topic.title_min_length", 3)
}

// GetTopicTitleMaxLength 获取话题标题最大长度
func GetTopicTitleMaxLength() int {
	return config.GetInt("app_constants.topic.title_max_length", 100)
}

// GetTopicContentMinLength 获取话题内容最小长度
func GetTopicContentMinLength() int {
	return config.GetInt("app_constants.topic.content_min_length", 10)
}

// GetTopicContentMaxLength 获取话题内容最大长度
func GetTopicContentMaxLength() int {
	return config.GetInt("app_constants.topic.content_max_length", 50000)
}

// ============= 时间格式配置 =============

// GetISO8601Format 获取ISO8601时间格式
func GetISO8601Format() string {
	return config.GetString("app_constants.datetime.format_iso8601", "2006-01-02T15:04:05Z07:00")
}

// GetDateFormat 获取日期格式
func GetDateFormat() string {
	return config.GetString("app_constants.datetime.format_date", "2006-01-02")
}

// GetTimeFormat 获取时间格式
func GetTimeFormat() string {
	return config.GetString("app_constants.datetime.format_time", "15:04:05")
}

// GetDateTimeFormat 获取日期时间格式
func GetDateTimeFormat() string {
	return config.GetString("app_constants.datetime.format_datetime", "2006-01-02 15:04:05")
}

// ============= 资源追踪配置 =============

// GetHTTPTimeoutWarningSeconds 获取HTTP请求超时警告阈值（秒）
func GetHTTPTimeoutWarningSeconds() int {
	return config.GetInt("app_constants.resource_tracking.http_timeout_warning_seconds", 30)
}

// GetResourceLeakThresholdMinutes 获取资源泄漏检测阈值（分钟）
func GetResourceLeakThresholdMinutes() int {
	return config.GetInt("app_constants.resource_tracking.resource_leak_threshold_minutes", 5)
}

// GetCheckIntervalMinutes 获取资源检查间隔（分钟）
func GetCheckIntervalMinutes() int {
	return config.GetInt("app_constants.resource_tracking.check_interval_minutes", 1)
}

// ============= 上下文超时配置 =============

// GetUserQueryTimeoutSeconds 获取用户查询超时（秒）
func GetUserQueryTimeoutSeconds() int {
	return config.GetInt("app_constants.context_timeout.user_query_seconds", 5)
}

// GetBatchNotifyTimeoutSeconds 获取批量通知超时（秒）
func GetBatchNotifyTimeoutSeconds() int {
	return config.GetInt("app_constants.context_timeout.batch_notify_seconds", 30)
}

// GetBatchOperationTimeoutSeconds 获取批量操作超时（秒）
func GetBatchOperationTimeoutSeconds() int {
	return config.GetInt("app_constants.context_timeout.batch_operation_seconds", 30)
}

// GetSearchAllTimeoutSeconds 获取搜索全部超时（秒）
func GetSearchAllTimeoutSeconds() int {
	return config.GetInt("app_constants.context_timeout.search_all_seconds", 10)
}

// GetDefaultOperationTimeoutSeconds 获取默认操作超时（秒）
func GetDefaultOperationTimeoutSeconds() int {
	return config.GetInt("app_constants.context_timeout.default_operation_seconds", 15)
}

// ============= 并发池配置 =============

// GetNotificationWorkers 获取通知服务工作协程数
func GetNotificationWorkers() int {
	return config.GetInt("app_constants.goroutine_pool.notification_workers", 20)
}

// GetInteractionWorkers 获取交互服务工作协程数
func GetInteractionWorkers() int {
	return config.GetInt("app_constants.goroutine_pool.interaction_workers", 30)
}

// GetDefaultWorkers 获取默认工作协程数
func GetDefaultWorkers() int {
	return config.GetInt("app_constants.goroutine_pool.default_workers", 10)
}
