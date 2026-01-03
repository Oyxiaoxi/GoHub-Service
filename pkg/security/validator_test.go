package security

import (
	"testing"
)

func TestInputValidator_CheckSQLInjection(t *testing.T) {
	validator := NewInputValidator()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"正常输入", "hello world", true},
		{"SQL UNION 注入", "1' UNION SELECT * FROM users--", false},
		{"SQL INSERT 注入", "'; INSERT INTO users VALUES('hacker')--", false},
		{"SQL UPDATE 注入", "admin' OR '1'='1", false},
		{"SQL DROP 注入", "'; DROP TABLE users--", false},
		{"正常 SQL 关键词", "select your favorite color", true}, // 非注入上下文
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.CheckSQLInjection(tt.input)
			if result.IsValid != tt.expected {
				t.Errorf("CheckSQLInjection(%q) = %v, want %v", tt.input, result.IsValid, tt.expected)
			}
		})
	}
}

func TestInputValidator_CheckXSS(t *testing.T) {
	validator := NewInputValidator()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"正常输入", "hello world", true},
		{"Script 标签", "<script>alert('xss')</script>", false},
		{"JavaScript 协议", "javascript:alert('xss')", false},
		{"Onclick 事件", "<img src=x onerror=alert('xss')>", false},
		{"Iframe 注入", "<iframe src='http://evil.com'></iframe>", false},
		{"正常 HTML", "<p>Hello</p>", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.CheckXSS(tt.input)
			if result.IsValid != tt.expected {
				t.Errorf("CheckXSS(%q) = %v, want %v (reason: %s)", 
					tt.input, result.IsValid, tt.expected, result.Reason)
			}
		})
	}
}

func TestInputValidator_CheckPathTraversal(t *testing.T) {
	validator := NewInputValidator()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"正常路径", "files/documents/report.pdf", true},
		{"路径遍历 Unix", "../../../etc/passwd", false},
		{"路径遍历 Windows", "..\\..\\windows\\system32", false},
		{"URL 编码路径遍历", "%2e%2e/etc/passwd", false},
		{"正常目录名", "dir..name/file.txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.CheckPathTraversal(tt.input)
			if result.IsValid != tt.expected {
				t.Errorf("CheckPathTraversal(%q) = %v, want %v", tt.input, result.IsValid, tt.expected)
			}
		})
	}
}

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name            string
		password        string
		expectedValid   bool
		expectedMinScore int
	}{
		{"弱密码 - 太短", "abc123", false, 0},
		{"弱密码 - 只有小写和数字", "abcd1234", false, 50},
		{"中等密码", "Abcd1234", true, 60},
		{"强密码", "Abcd1234!@#", true, 80},
		{"超强密码", "MyP@ssw0rd!2024", true, 90},
		{"太长", string(make([]byte, 129)), false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidatePasswordStrength(tt.password)
			if result.IsValid != tt.expectedValid {
				t.Errorf("ValidatePasswordStrength(%q).IsValid = %v, want %v", 
					tt.password, result.IsValid, tt.expectedValid)
			}
			if result.IsValid && result.Score < tt.expectedMinScore {
				t.Errorf("ValidatePasswordStrength(%q).Score = %d, want >= %d", 
					tt.password, result.Score, tt.expectedMinScore)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"有效邮箱", "user@example.com", true},
		{"有效邮箱带点", "user.name@example.com", true},
		{"有效邮箱带加号", "user+tag@example.com", true},
		{"无效邮箱 - 缺少@", "userexample.com", false},
		{"无效邮箱 - 缺少域名", "user@", false},
		{"无效邮箱 - 缺少后缀", "user@example", false},
		{"无效邮箱 - 特殊字符", "user<>@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidEmail(tt.email)
			if result != tt.expected {
				t.Errorf("IsValidEmail(%q) = %v, want %v", tt.email, result, tt.expected)
			}
		})
	}
}

func TestIsValidPhone(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected bool
	}{
		{"有效手机号 - 13x", "13800138000", true},
		{"有效手机号 - 15x", "15912345678", true},
		{"有效手机号 - 18x", "18612345678", true},
		{"无效手机号 - 太短", "1380013800", false},
		{"无效手机号 - 太长", "138001380000", false},
		{"无效手机号 - 非1开头", "28800138000", false},
		{"无效手机号 - 第二位无效", "12800138000", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidPhone(tt.phone)
			if result != tt.expected {
				t.Errorf("IsValidPhone(%q) = %v, want %v", tt.phone, result, tt.expected)
			}
		})
	}
}

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{"有效 HTTP URL", "http://example.com", true},
		{"有效 HTTPS URL", "https://example.com/path", true},
		{"有效 URL 带参数", "https://example.com/path?key=value", true},
		{"无效 URL - 无协议", "example.com", false},
		{"无效 URL - FTP 协议", "ftp://example.com", false},
		{"无效 URL - JavaScript", "javascript:alert(1)", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidURL(tt.url)
			if result != tt.expected {
				t.Errorf("IsValidURL(%q) = %v, want %v", tt.url, result, tt.expected)
			}
		})
	}
}

func BenchmarkInputValidator_Validate(b *testing.B) {
	validator := NewInputValidator()
	input := "This is a normal input string with some text"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.Validate(input)
	}
}

func BenchmarkValidatePasswordStrength(b *testing.B) {
	password := "MyP@ssw0rd!2024"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidatePasswordStrength(password)
	}
}
