// Package logger 敏感信息过滤
package logger

import (
	"regexp"
	"strings"
)

// SensitiveFilter 敏感信息过滤器
type SensitiveFilter struct {
	patterns map[string]*regexp.Regexp
}

// NewSensitiveFilter 创建过滤器
func NewSensitiveFilter() *SensitiveFilter {
	return &SensitiveFilter{
		patterns: map[string]*regexp.Regexp{
			"password":     regexp.MustCompile(`(?i)(password|passwd|pwd)["']?\s*[:=]\s*["']?([^"'\s,}]+)`),
			"token":        regexp.MustCompile(`(?i)(token|access_token|refresh_token)["']?\s*[:=]\s*["']?([^"'\s,}]+)`),
			"secret":       regexp.MustCompile(`(?i)(secret|api_key|apikey)["']?\s*[:=]\s*["']?([^"'\s,}]+)`),
			"card":         regexp.MustCompile(`\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b`),
			"id_card":      regexp.MustCompile(`\b\d{15}|\d{18}\b`),
			"phone":        regexp.MustCompile(`\b1[3-9]\d{9}\b`),
			"email":        regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`),
			"authorization": regexp.MustCompile(`(?i)Authorization["']?\s*[:=]\s*["']?([^"'\s,}]+)`),
		},
	}
}

// Filter 过滤敏感信息
func (f *SensitiveFilter) Filter(message string) string {
	result := message
	
	// 密码类
	if f.patterns["password"].MatchString(result) {
		result = f.patterns["password"].ReplaceAllString(result, `$1:"***"`)
	}
	
	// Token类
	if f.patterns["token"].MatchString(result) {
		result = f.patterns["token"].ReplaceAllString(result, `$1:"***"`)
	}
	
	// Secret类
	if f.patterns["secret"].MatchString(result) {
		result = f.patterns["secret"].ReplaceAllString(result, `$1:"***"`)
	}
	
	// Authorization
	if f.patterns["authorization"].MatchString(result) {
		result = f.patterns["authorization"].ReplaceAllString(result, `Authorization:"Bearer ***"`)
	}
	
	// 银行卡号 (保留前4位后4位)
	if f.patterns["card"].MatchString(result) {
		result = f.patterns["card"].ReplaceAllStringFunc(result, func(s string) string {
			s = strings.ReplaceAll(s, " ", "")
			s = strings.ReplaceAll(s, "-", "")
			if len(s) == 16 {
				return s[:4] + "****" + "****" + s[12:]
			}
			return "****-****-****-****"
		})
	}
	
	// 身份证号 (保留前6位后4位)
	if f.patterns["id_card"].MatchString(result) {
		result = f.patterns["id_card"].ReplaceAllStringFunc(result, func(s string) string {
			if len(s) == 18 {
				return s[:6] + "********" + s[14:]
			} else if len(s) == 15 {
				return s[:6] + "*****" + s[11:]
			}
			return "******************"
		})
	}
	
	// 手机号 (保留前3位后4位)
	if f.patterns["phone"].MatchString(result) {
		result = f.patterns["phone"].ReplaceAllStringFunc(result, func(s string) string {
			if len(s) == 11 {
				return s[:3] + "****" + s[7:]
			}
			return "***-****-****"
		})
	}
	
	// 邮箱 (保留前2位和域名)
	if f.patterns["email"].MatchString(result) {
		result = f.patterns["email"].ReplaceAllStringFunc(result, func(s string) string {
			parts := strings.Split(s, "@")
			if len(parts) == 2 && len(parts[0]) > 2 {
				return parts[0][:2] + "***@" + parts[1]
			}
			return "***@***"
		})
	}
	
	return result
}

// FilterMap 过滤 map 中的敏感信息
func (f *SensitiveFilter) FilterMap(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	
	sensitiveKeys := []string{
		"password", "passwd", "pwd",
		"token", "access_token", "refresh_token",
		"secret", "api_key", "apikey",
		"authorization",
	}
	
	for k, v := range data {
		lowerKey := strings.ToLower(k)
		
		// 检查是否为敏感字段
		isSensitive := false
		for _, sk := range sensitiveKeys {
			if strings.Contains(lowerKey, sk) {
				isSensitive = true
				break
			}
		}
		
		if isSensitive {
			result[k] = "***"
		} else {
			// 递归处理嵌套 map
			if nestedMap, ok := v.(map[string]interface{}); ok {
				result[k] = f.FilterMap(nestedMap)
			} else if str, ok := v.(string); ok {
				result[k] = f.Filter(str)
			} else {
				result[k] = v
			}
		}
	}
	
	return result
}

// 全局过滤器实例
var globalFilter = NewSensitiveFilter()

// FilterSensitive 过滤敏感信息的快捷方法
func FilterSensitive(message string) string {
	return globalFilter.Filter(message)
}

// FilterSensitiveMap 过滤 map 的快捷方法
func FilterSensitiveMap(data map[string]interface{}) map[string]interface{} {
	return globalFilter.FilterMap(data)
}
