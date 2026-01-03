// Package security 安全验证工具
package security

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// InputValidator 输入验证器
type InputValidator struct {
	// SQL 注入模式
	sqlInjectionPatterns []*regexp.Regexp
	// XSS 攻击模式
	xssPatterns []*regexp.Regexp
	// 路径遍历模式
	pathTraversalPatterns []*regexp.Regexp
}

// NewInputValidator 创建输入验证器
func NewInputValidator() *InputValidator {
	return &InputValidator{
		sqlInjectionPatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(union\s+select|select\s+.*\s+from|insert\s+into|update\s+.*\s+set|delete\s+from|drop\s+table|truncate\s+table)`),
			regexp.MustCompile(`(?i)(exec|execute|sp_executesql|xp_cmdshell)`),
			regexp.MustCompile(`(?i)(--|\#|\/\*|\*\/|;)\s*$`),
			regexp.MustCompile(`(?i)('|")\s*(or|and)\s*('|")?\d+('|")?(\s*=\s*\d+)?`),
		},
		xssPatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`),
			regexp.MustCompile(`(?i)javascript:`),
			regexp.MustCompile(`(?i)on(load|error|click|mouse\w+|key\w+|focus|blur)\s*=`),
			regexp.MustCompile(`(?i)<iframe[^>]*>.*?</iframe>`),
			regexp.MustCompile(`(?i)<embed[^>]*>.*?</embed>`),
			regexp.MustCompile(`(?i)<object[^>]*>.*?</object>`),
		},
		pathTraversalPatterns: []*regexp.Regexp{
			regexp.MustCompile(`\.\./|\.\.\\`),
			regexp.MustCompile(`%2e%2e/|%2e%2e\\`),
			regexp.MustCompile(`\.\.%2f|\.\.%5c`),
		},
	}
}

// ValidateInput 综合输入验证
type ValidationResult struct {
	IsValid  bool
	Reason   string
	RiskType string // "sql_injection", "xss", "path_traversal", "invalid_chars"
}

// Validate 验证输入
func (v *InputValidator) Validate(input string) ValidationResult {
	// 检查长度
	if utf8.RuneCountInString(input) > 10000 {
		return ValidationResult{
			IsValid:  false,
			Reason:   "输入长度超过限制",
			RiskType: "invalid_length",
		}
	}

	// 检查 SQL 注入
	if result := v.CheckSQLInjection(input); !result.IsValid {
		return result
	}

	// 检查 XSS
	if result := v.CheckXSS(input); !result.IsValid {
		return result
	}

	// 检查路径遍历
	if result := v.CheckPathTraversal(input); !result.IsValid {
		return result
	}

	return ValidationResult{IsValid: true}
}

// CheckSQLInjection 检查 SQL 注入
func (v *InputValidator) CheckSQLInjection(input string) ValidationResult {
	for _, pattern := range v.sqlInjectionPatterns {
		if pattern.MatchString(input) {
			return ValidationResult{
				IsValid:  false,
				Reason:   "检测到潜在的 SQL 注入攻击",
				RiskType: "sql_injection",
			}
		}
	}
	return ValidationResult{IsValid: true}
}

// CheckXSS 检查 XSS 攻击
func (v *InputValidator) CheckXSS(input string) ValidationResult {
	for _, pattern := range v.xssPatterns {
		if pattern.MatchString(input) {
			return ValidationResult{
				IsValid:  false,
				Reason:   "检测到潜在的 XSS 攻击",
				RiskType: "xss",
			}
		}
	}
	return ValidationResult{IsValid: true}
}

// CheckPathTraversal 检查路径遍历攻击
func (v *InputValidator) CheckPathTraversal(input string) ValidationResult {
	for _, pattern := range v.pathTraversalPatterns {
		if pattern.MatchString(input) {
			return ValidationResult{
				IsValid:  false,
				Reason:   "检测到潜在的路径遍历攻击",
				RiskType: "path_traversal",
			}
		}
	}
	return ValidationResult{IsValid: true}
}

// SanitizeInput 清理输入（保留用于向后兼容）
func (v *InputValidator) SanitizeInput(input string) string {
	// HTML 实体转义
	replacements := map[string]string{
		"<":  "&lt;",
		">":  "&gt;",
		"&":  "&amp;",
		"\"": "&quot;",
		"'":  "&#x27;",
		"/":  "&#x2F;",
	}

	result := input
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}

	return result
}

// IsValidEmail 验证邮箱格式
func IsValidEmail(email string) bool {
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailPattern.MatchString(email) && len(email) <= 254
}

// IsValidPhone 验证手机号格式（中国）
func IsValidPhone(phone string) bool {
	phonePattern := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return phonePattern.MatchString(phone)
}

// IsValidURL 验证 URL 格式
func IsValidURL(url string) bool {
	urlPattern := regexp.MustCompile(`^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$`)
	return urlPattern.MatchString(url) && len(url) <= 2048
}

// IsAlphanumeric 验证是否只包含字母和数字
func IsAlphanumeric(s string) bool {
	alphanumericPattern := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return alphanumericPattern.MatchString(s)
}

// HasSpecialChars 检查是否包含特殊字符
func HasSpecialChars(s string) bool {
	specialCharsPattern := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]`)
	return specialCharsPattern.MatchString(s)
}

// ValidatePassword 验证密码强度
type PasswordStrength struct {
	IsValid      bool
	Score        int    // 0-100
	Issues       []string
	HasUppercase bool
	HasLowercase bool
	HasDigit     bool
	HasSpecial   bool
	Length       int
}

// ValidatePasswordStrength 验证密码强度
func ValidatePasswordStrength(password string) PasswordStrength {
	result := PasswordStrength{
		Length: utf8.RuneCountInString(password),
		Issues: []string{},
	}

	// 检查长度
	if result.Length < 8 {
		result.Issues = append(result.Issues, "密码长度至少 8 个字符")
	}
	if result.Length > 128 {
		result.Issues = append(result.Issues, "密码长度不能超过 128 个字符")
		result.IsValid = false
		return result
	}

	// 检查字符类型
	result.HasUppercase = regexp.MustCompile(`[A-Z]`).MatchString(password)
	result.HasLowercase = regexp.MustCompile(`[a-z]`).MatchString(password)
	result.HasDigit = regexp.MustCompile(`\d`).MatchString(password)
	result.HasSpecial = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]`).MatchString(password)

	// 计算得分
	score := 0
	if result.Length >= 8 {
		score += 20
	}
	if result.Length >= 12 {
		score += 10
	}
	if result.Length >= 16 {
		score += 10
	}
	if result.HasUppercase {
		score += 15
	}
	if result.HasLowercase {
		score += 15
	}
	if result.HasDigit {
		score += 15
	}
	if result.HasSpecial {
		score += 15
	}

	result.Score = score

	// 建议
	if !result.HasUppercase {
		result.Issues = append(result.Issues, "建议包含大写字母")
	}
	if !result.HasLowercase {
		result.Issues = append(result.Issues, "建议包含小写字母")
	}
	if !result.HasDigit {
		result.Issues = append(result.Issues, "建议包含数字")
	}
	if !result.HasSpecial {
		result.Issues = append(result.Issues, "建议包含特殊字符")
	}

	// 最低要求：长度 >= 8 且至少包含 3 种字符类型
	charTypes := 0
	if result.HasUppercase {
		charTypes++
	}
	if result.HasLowercase {
		charTypes++
	}
	if result.HasDigit {
		charTypes++
	}
	if result.HasSpecial {
		charTypes++
	}

	result.IsValid = result.Length >= 8 && charTypes >= 3

	return result
}
