// Package security API 签名验证
package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

// SignatureConfig 签名配置
type SignatureConfig struct {
	Secret         string        // 签名密钥
	TimestampValid time.Duration // 时间戳有效期（默认 5 分钟）
	NonceLength    int           // Nonce 最小长度（默认 16）
}

// DefaultSignatureConfig 默认签名配置
func DefaultSignatureConfig(secret string) *SignatureConfig {
	return &SignatureConfig{
		Secret:         secret,
		TimestampValid: 5 * time.Minute,
		NonceLength:    16,
	}
}

// SignatureValidator 签名验证器
type SignatureValidator struct {
	config *SignatureConfig
}

// NewSignatureValidator 创建签名验证器
func NewSignatureValidator(config *SignatureConfig) *SignatureValidator {
	if config == nil {
		config = DefaultSignatureConfig("")
	}
	return &SignatureValidator{
		config: config,
	}
}

// SignRequest 生成请求签名
// @param method HTTP 方法（GET/POST/PUT/DELETE）
// @param path 请求路径（不含域名）
// @param timestamp 时间戳（Unix 秒）
// @param nonce 随机字符串
// @param body 请求体（可选，GET 请求传空字符串）
// @return 签名字符串
func (sv *SignatureValidator) SignRequest(method, path string, timestamp int64, nonce, body string) string {
	// 构建待签名字符串
	signString := sv.buildSignString(method, path, timestamp, nonce, body)
	
	// HMAC-SHA256 签名
	h := hmac.New(sha256.New, []byte(sv.config.Secret))
	h.Write([]byte(signString))
	signature := hex.EncodeToString(h.Sum(nil))
	
	return signature
}

// VerifySignature 验证签名
// @param method HTTP 方法
// @param path 请求路径
// @param timestamp 时间戳
// @param nonce 随机字符串
// @param body 请求体
// @param signature 客户端提供的签名
// @return 验证结果
func (sv *SignatureValidator) VerifySignature(method, path string, timestamp int64, nonce, body, signature string) *ValidationResult {
	// 1. 验证时间戳
	if err := sv.validateTimestamp(timestamp); err != nil {
		return &ValidationResult{
			IsValid:  false,
			Reason:   err.Error(),
			RiskType: "timestamp_invalid",
		}
	}
	
	// 2. 验证 Nonce 长度
	if len(nonce) < sv.config.NonceLength {
		return &ValidationResult{
			IsValid:  false,
			Reason:   fmt.Sprintf("nonce length must be at least %d characters", sv.config.NonceLength),
			RiskType: "nonce_invalid",
		}
	}
	
	// 3. 计算期望的签名
	expectedSignature := sv.SignRequest(method, path, timestamp, nonce, body)
	
	// 4. 对比签名（防止时序攻击）
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return &ValidationResult{
			IsValid:  false,
			Reason:   "signature mismatch",
			RiskType: "signature_invalid",
		}
	}
	
	return &ValidationResult{
		IsValid:  true,
		Reason:   "signature valid",
		RiskType: "",
	}
}

// validateTimestamp 验证时间戳
func (sv *SignatureValidator) validateTimestamp(timestamp int64) error {
	now := time.Now().Unix()
	diff := now - timestamp
	
	// 时间戳不能是未来时间
	if diff < 0 {
		return fmt.Errorf("timestamp is in the future")
	}
	
	// 时间戳不能过期
	if diff > int64(sv.config.TimestampValid.Seconds()) {
		return fmt.Errorf("timestamp expired, max age: %v", sv.config.TimestampValid)
	}
	
	return nil
}

// buildSignString 构建待签名字符串
// 格式：METHOD\nPATH\nTIMESTAMP\nNONCE\nBODY
func (sv *SignatureValidator) buildSignString(method, path string, timestamp int64, nonce, body string) string {
	return fmt.Sprintf("%s\n%s\n%d\n%s\n%s",
		strings.ToUpper(method),
		path,
		timestamp,
		nonce,
		body,
	)
}

// SignWithQuery 生成带查询参数的签名
// 用于 GET 请求，将查询参数排序后参与签名
func (sv *SignatureValidator) SignWithQuery(method, path string, timestamp int64, nonce string, queryParams map[string]string) string {
	// 查询参数排序
	keys := make([]string, 0, len(queryParams))
	for k := range queryParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	
	// 构建查询字符串
	var queryParts []string
	for _, k := range keys {
		queryParts = append(queryParts, fmt.Sprintf("%s=%s", k, queryParams[k]))
	}
	queryString := strings.Join(queryParts, "&")
	
	// 使用查询字符串作为 body
	return sv.SignRequest(method, path, timestamp, nonce, queryString)
}

// VerifySignatureWithQuery 验证带查询参数的签名
func (sv *SignatureValidator) VerifySignatureWithQuery(method, path string, timestamp int64, nonce string, queryParams map[string]string, signature string) *ValidationResult {
	// 查询参数排序
	keys := make([]string, 0, len(queryParams))
	for k := range queryParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	
	// 构建查询字符串
	var queryParts []string
	for _, k := range keys {
		queryParts = append(queryParts, fmt.Sprintf("%s=%s", k, queryParams[k]))
	}
	queryString := strings.Join(queryParts, "&")
	
	return sv.VerifySignature(method, path, timestamp, nonce, queryString, signature)
}

// GenerateNonce 生成随机 Nonce
func GenerateNonce(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return GenerateRandomString(length, charset)
}

// GetCurrentTimestamp 获取当前时间戳
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}
