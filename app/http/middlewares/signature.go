// Package middlewares API 签名验证中间件
package middlewares

import (
	"bytes"
	"io"
	"strconv"
	"sync"
	"time"

	"GoHub-Service/pkg/config"
	"GoHub-Service/pkg/redis"
	"GoHub-Service/pkg/response"
	"GoHub-Service/pkg/security"

	"github.com/gin-gonic/gin"
)

var (
	signatureValidator     *security.SignatureValidator
	signatureValidatorOnce sync.Once
)

// getSignatureValidator 获取签名验证器单例
func getSignatureValidator() *security.SignatureValidator {
	signatureValidatorOnce.Do(func() {
		secret := config.GetString("app.signature_secret")
		if secret == "" {
			// 使用 JWT secret 作为后备
			secret = config.GetString("app.key")
		}
		
		signatureValidator = security.NewSignatureValidator(
			&security.SignatureConfig{
				Secret:         secret,
				TimestampValid: 5 * time.Minute,
				NonceLength:    16,
			},
		)
	})
	return signatureValidator
}

// APISignatureVerification API 签名验证中间件
// 需要客户端在请求头中提供以下字段：
// - X-Timestamp: Unix 时间戳（秒）
// - X-Nonce: 随机字符串（至少 16 位）
// - X-Signature: HMAC-SHA256 签名
func APISignatureVerification() gin.HandlerFunc {
	validator := getSignatureValidator()

	return func(c *gin.Context) {
		// 1. 获取签名相关参数
		timestampStr := c.GetHeader("X-Timestamp")
		nonce := c.GetHeader("X-Nonce")
		signature := c.GetHeader("X-Signature")

		// 2. 检查必需参数
		if timestampStr == "" || nonce == "" || signature == "" {
			response.ApiError(c, 401, response.CodeUnauthorized,
				"缺少签名参数 (X-Timestamp, X-Nonce, X-Signature)")
			return
		}

		// 3. 解析时间戳
		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			response.ApiError(c, 401, response.CodeUnauthorized,
				"时间戳格式错误")
			return
		}

		// 4. 检查 Nonce 是否已使用（防重放）
		nonceKey := "api:nonce:" + nonce
		if redis.Redis.Has(c.Request.Context(), nonceKey) {
			response.ApiError(c, 401, response.CodeUnauthorized,
				"请求已被处理（重放攻击检测）")
			return
		}

		// 5. 读取请求体
		var body string
		if c.Request.Method != "GET" {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				response.ApiError(c, 400, response.CodeInvalidParams,
					"读取请求体失败")
				return
			}
			body = string(bodyBytes)
			// 重新设置请求体供后续使用
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 6. 验证签名
		method := c.Request.Method
		path := c.Request.URL.Path
		result := validator.VerifySignature(method, path, timestamp, nonce, body, signature)

		if !result.IsValid {
			response.ApiError(c, 401, response.CodeUnauthorized,
				"签名验证失败: "+result.Reason)
			return
		}

		// 7. 记录 Nonce（5分钟有效期）
		redis.Redis.Set(
			c.Request.Context(),
			nonceKey,
			"1",
			5*time.Minute,
		)

		c.Next()
	}
}

// APISignatureVerificationWithQuery GET 请求签名验证（包含查询参数）
func APISignatureVerificationWithQuery() gin.HandlerFunc {
	validator := getSignatureValidator()

	return func(c *gin.Context) {
		// 1. 获取签名相关参数
		timestampStr := c.GetHeader("X-Timestamp")
		nonce := c.GetHeader("X-Nonce")
		signature := c.GetHeader("X-Signature")

		if timestampStr == "" || nonce == "" || signature == "" {
			response.ApiError(c, 401, response.CodeUnauthorized,
				"缺少签名参数")
			return
		}

		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			response.ApiError(c, 401, response.CodeUnauthorized,
				"时间戳格式错误")
			return
		}

		// 2. 检查 Nonce 是否已使用
		nonceKey := "api:nonce:" + nonce
		if redis.Redis.Has(c.Request.Context(), nonceKey) {
			response.ApiError(c, 401, response.CodeUnauthorized,
				"请求已被处理（重放攻击检测）")
			return
		}

		// 3. 获取查询参数
		queryParams := make(map[string]string)
		for key, values := range c.Request.URL.Query() {
			if len(values) > 0 {
				queryParams[key] = values[0]
			}
		}

		// 4. 验证签名（包含查询参数）
		method := c.Request.Method
		path := c.Request.URL.Path
		result := validator.VerifySignatureWithQuery(method, path, timestamp, nonce, queryParams, signature)

		if !result.IsValid {
			response.ApiError(c, 401, response.CodeUnauthorized,
				"签名验证失败: "+result.Reason)
			return
		}

		// 5. 记录 Nonce
		redis.Redis.Set(
			c.Request.Context(),
			nonceKey,
			"1",
			5*time.Minute,
		)

		c.Next()
	}
}

// OptionalSignatureVerification 可选的签名验证
// 如果提供了签名参数则验证，否则跳过
func OptionalSignatureVerification() gin.HandlerFunc {
	validator := getSignatureValidator()

	return func(c *gin.Context) {
		signature := c.GetHeader("X-Signature")
		
		// 如果没有提供签名，直接放行
		if signature == "" {
			c.Next()
			return
		}

		// 有签名则进行验证
		timestampStr := c.GetHeader("X-Timestamp")
		nonce := c.GetHeader("X-Nonce")

		if timestampStr == "" || nonce == "" {
			response.ApiError(c, 401, response.CodeUnauthorized,
				"签名参数不完整")
			return
		}

		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			response.ApiError(c, 401, response.CodeUnauthorized,
				"时间戳格式错误")
			return
		}

		// 读取请求体
		var body string
		if c.Request.Method != "GET" {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				response.ApiError(c, 400, response.CodeInvalidParams,
					"读取请求体失败")
				return
			}
			body = string(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 验证签名
		method := c.Request.Method
		path := c.Request.URL.Path
		result := validator.VerifySignature(method, path, timestamp, nonce, body, signature)

		if !result.IsValid {
			response.ApiError(c, 401, response.CodeUnauthorized,
				"签名验证失败: "+result.Reason)
			return
		}

		c.Next()
	}
}
