package security

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSignatureValidator_SignRequest(t *testing.T) {
	config := &SignatureConfig{
		Secret:         "test-secret-key-12345678",
		TimestampValid: 5 * time.Minute,
		NonceLength:    16,
	}
	validator := NewSignatureValidator(config)

	timestamp := int64(1609459200) // 2021-01-01 00:00:00
	nonce := "abcdef1234567890"
	body := `{"name":"test"}`

	t.Run("生成签名", func(t *testing.T) {
		signature := validator.SignRequest("POST", "/api/v1/users", timestamp, nonce, body)
		
		assert.NotEmpty(t, signature)
		assert.Len(t, signature, 64) // SHA256 hex = 64 字符
	})

	t.Run("相同参数生成相同签名", func(t *testing.T) {
		sig1 := validator.SignRequest("POST", "/api/v1/users", timestamp, nonce, body)
		sig2 := validator.SignRequest("POST", "/api/v1/users", timestamp, nonce, body)
		
		assert.Equal(t, sig1, sig2)
	})

	t.Run("不同参数生成不同签名", func(t *testing.T) {
		sig1 := validator.SignRequest("POST", "/api/v1/users", timestamp, nonce, body)
		sig2 := validator.SignRequest("PUT", "/api/v1/users", timestamp, nonce, body)
		
		assert.NotEqual(t, sig1, sig2)
	})
}

func TestSignatureValidator_VerifySignature(t *testing.T) {
	config := &SignatureConfig{
		Secret:         "test-secret-key-12345678",
		TimestampValid: 5 * time.Minute,
		NonceLength:    16,
	}
	validator := NewSignatureValidator(config)

	t.Run("验证有效签名", func(t *testing.T) {
		timestamp := time.Now().Unix()
		nonce := "abcdef1234567890"
		body := `{"name":"test"}`
		
		signature := validator.SignRequest("POST", "/api/v1/users", timestamp, nonce, body)
		result := validator.VerifySignature("POST", "/api/v1/users", timestamp, nonce, body, signature)
		
		assert.True(t, result.IsValid)
		assert.Empty(t, result.RiskType)
	})

	t.Run("签名不匹配", func(t *testing.T) {
		timestamp := time.Now().Unix()
		nonce := "abcdef1234567890"
		body := `{"name":"test"}`
		
		signature := "invalid_signature_1234567890abcdef1234567890abcdef1234567890abcd"
		result := validator.VerifySignature("POST", "/api/v1/users", timestamp, nonce, body, signature)
		
		assert.False(t, result.IsValid)
		assert.Equal(t, "signature_invalid", result.RiskType)
	})

	t.Run("时间戳过期", func(t *testing.T) {
		timestamp := time.Now().Add(-10 * time.Minute).Unix() // 10分钟前
		nonce := "abcdef1234567890"
		body := `{"name":"test"}`
		
		signature := validator.SignRequest("POST", "/api/v1/users", timestamp, nonce, body)
		result := validator.VerifySignature("POST", "/api/v1/users", timestamp, nonce, body, signature)
		
		assert.False(t, result.IsValid)
		assert.Equal(t, "timestamp_invalid", result.RiskType)
		assert.Contains(t, result.Reason, "expired")
	})

	t.Run("时间戳在未来", func(t *testing.T) {
		timestamp := time.Now().Add(10 * time.Minute).Unix() // 10分钟后
		nonce := "abcdef1234567890"
		body := `{"name":"test"}`
		
		signature := validator.SignRequest("POST", "/api/v1/users", timestamp, nonce, body)
		result := validator.VerifySignature("POST", "/api/v1/users", timestamp, nonce, body, signature)
		
		assert.False(t, result.IsValid)
		assert.Equal(t, "timestamp_invalid", result.RiskType)
		assert.Contains(t, result.Reason, "future")
	})

	t.Run("Nonce太短", func(t *testing.T) {
		timestamp := time.Now().Unix()
		nonce := "short" // 长度不足16
		body := `{"name":"test"}`
		
		signature := validator.SignRequest("POST", "/api/v1/users", timestamp, nonce, body)
		result := validator.VerifySignature("POST", "/api/v1/users", timestamp, nonce, body, signature)
		
		assert.False(t, result.IsValid)
		assert.Equal(t, "nonce_invalid", result.RiskType)
	})
}

func TestSignatureValidator_SignWithQuery(t *testing.T) {
	config := &SignatureConfig{
		Secret:         "test-secret-key-12345678",
		TimestampValid: 5 * time.Minute,
		NonceLength:    16,
	}
	validator := NewSignatureValidator(config)

	timestamp := time.Now().Unix()
	nonce := "abcdef1234567890"
	queryParams := map[string]string{
		"page":     "1",
		"per_page": "15",
		"sort":     "created_at",
	}

	t.Run("生成带查询参数的签名", func(t *testing.T) {
		signature := validator.SignWithQuery("GET", "/api/v1/topics", timestamp, nonce, queryParams)
		
		assert.NotEmpty(t, signature)
		assert.Len(t, signature, 64)
	})

	t.Run("验证带查询参数的签名", func(t *testing.T) {
		signature := validator.SignWithQuery("GET", "/api/v1/topics", timestamp, nonce, queryParams)
		result := validator.VerifySignatureWithQuery("GET", "/api/v1/topics", timestamp, nonce, queryParams, signature)
		
		assert.True(t, result.IsValid)
	})

	t.Run("查询参数顺序不影响签名", func(t *testing.T) {
		params1 := map[string]string{"a": "1", "b": "2", "c": "3"}
		params2 := map[string]string{"c": "3", "a": "1", "b": "2"}
		
		sig1 := validator.SignWithQuery("GET", "/api/v1/test", timestamp, nonce, params1)
		sig2 := validator.SignWithQuery("GET", "/api/v1/test", timestamp, nonce, params2)
		
		assert.Equal(t, sig1, sig2)
	})
}

func TestGenerateNonce(t *testing.T) {
	t.Run("生成指定长度的Nonce", func(t *testing.T) {
		nonce := GenerateNonce(16)
		assert.Len(t, nonce, 16)
	})

	t.Run("每次生成不同的Nonce", func(t *testing.T) {
		nonce1 := GenerateNonce(16)
		nonce2 := GenerateNonce(16)
		assert.NotEqual(t, nonce1, nonce2)
	})

	t.Run("生成32位Nonce", func(t *testing.T) {
		nonce := GenerateNonce(32)
		assert.Len(t, nonce, 32)
	})
}

func TestGetCurrentTimestamp(t *testing.T) {
	t.Run("获取当前时间戳", func(t *testing.T) {
		timestamp := GetCurrentTimestamp()
		now := time.Now().Unix()
		
		// 时间戳应该接近当前时间（误差在1秒内）
		assert.InDelta(t, now, timestamp, 1)
	})
}

// 基准测试
func BenchmarkSignatureValidator_SignRequest(b *testing.B) {
	config := &SignatureConfig{
		Secret:         "test-secret-key-12345678",
		TimestampValid: 5 * time.Minute,
		NonceLength:    16,
	}
	validator := NewSignatureValidator(config)
	
	timestamp := time.Now().Unix()
	nonce := "abcdef1234567890"
	body := `{"name":"test"}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.SignRequest("POST", "/api/v1/users", timestamp, nonce, body)
	}
}

func BenchmarkSignatureValidator_VerifySignature(b *testing.B) {
	config := &SignatureConfig{
		Secret:         "test-secret-key-12345678",
		TimestampValid: 5 * time.Minute,
		NonceLength:    16,
	}
	validator := NewSignatureValidator(config)
	
	timestamp := time.Now().Unix()
	nonce := "abcdef1234567890"
	body := `{"name":"test"}`
	signature := validator.SignRequest("POST", "/api/v1/users", timestamp, nonce, body)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.VerifySignature("POST", "/api/v1/users", timestamp, nonce, body, signature)
	}
}
