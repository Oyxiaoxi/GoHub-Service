// Package examples API 签名验证应用示例
package main

import (
	"fmt"
	"time"

	"GoHub-Service/pkg/security"
)

// 示例1：POST 请求签名
func Example_POSTRequest() {
	config := &security.SignatureConfig{
		Secret:         "your-secret-key-12345678",
		TimestampValid: 5 * time.Minute,
		NonceLength:    16,
	}
	validator := security.NewSignatureValidator(config)

	// 客户端生成签名
	timestamp := security.GetCurrentTimestamp()
	nonce := security.GenerateNonce(16)
	method := "POST"
	path := "/api/v1/users"
	body := `{"username":"test","email":"test@example.com"}`

	signature := validator.SignRequest(method, path, timestamp, nonce, body)

	fmt.Printf("请求头:\n")
	fmt.Printf("X-Timestamp: %d\n", timestamp)
	fmt.Printf("X-Nonce: %s\n", nonce)
	fmt.Printf("X-Signature: %s\n", signature)

	// 服务端验证签名
	result := validator.VerifySignature(method, path, timestamp, nonce, body, signature)
	fmt.Printf("\n验证结果: %v\n", result.IsValid)
	// Output:
	// 验证结果: true
}

// 示例2：GET 请求带查询参数
func Example_GETRequestWithQuery() {
	config := security.DefaultSignatureConfig("your-secret-key-12345678")
	validator := security.NewSignatureValidator(config)

	timestamp := security.GetCurrentTimestamp()
	nonce := security.GenerateNonce(16)
	method := "GET"
	path := "/api/v1/topics"

	// 查询参数
	queryParams := map[string]string{
		"page":     "1",
		"per_page": "15",
		"sort":     "created_at",
	}

	// 生成签名
	signature := validator.SignWithQuery(method, path, timestamp, nonce, queryParams)

	fmt.Printf("请求 URL: %s?page=1&per_page=15&sort=created_at\n", path)
	fmt.Printf("X-Signature: %s\n", signature)

	// 验证签名
	result := validator.VerifySignatureWithQuery(method, path, timestamp, nonce, queryParams, signature)
	fmt.Printf("验证结果: %v\n", result.IsValid)
	// Output:
	// 验证结果: true
}

// 示例3：签名失败场景
func Example_SignatureFailure() {
	config := security.DefaultSignatureConfig("your-secret-key-12345678")
	validator := security.NewSignatureValidator(config)

	timestamp := security.GetCurrentTimestamp()
	nonce := security.GenerateNonce(16)

	// 1. 签名不匹配
	fmt.Println("场景1: 签名不匹配")
	result := validator.VerifySignature(
		"POST", "/api/v1/users", timestamp, nonce, `{"test":"data"}`,
		"invalid_signature",
	)
	fmt.Printf("验证结果: %v, 原因: %s\n\n", result.IsValid, result.Reason)

	// 2. 时间戳过期（10分钟前）
	fmt.Println("场景2: 时间戳过期")
	expiredTimestamp := time.Now().Add(-10 * time.Minute).Unix()
	expiredSignature := validator.SignRequest("POST", "/api/v1/users", expiredTimestamp, nonce, `{"test":"data"}`)
	result = validator.VerifySignature(
		"POST", "/api/v1/users", expiredTimestamp, nonce, `{"test":"data"}`,
		expiredSignature,
	)
	fmt.Printf("验证结果: %v, 原因: %s\n\n", result.IsValid, result.Reason)

	// 3. Nonce 太短
	fmt.Println("场景3: Nonce 长度不足")
	shortNonce := "short"
	shortSignature := validator.SignRequest("POST", "/api/v1/users", timestamp, shortNonce, `{"test":"data"}`)
	result = validator.VerifySignature(
		"POST", "/api/v1/users", timestamp, shortNonce, `{"test":"data"}`,
		shortSignature,
	)
	fmt.Printf("验证结果: %v, 原因: %s\n", result.IsValid, result.Reason)
	// Output:
	// 验证结果: false, 原因: signature mismatch
}

// 示例4：完整的客户端请求流程
func Example_CompleteFlow() {
	// 1. 初始化客户端配置
	secret := "your-secret-key-12345678"
	config := security.DefaultSignatureConfig(secret)
	validator := security.NewSignatureValidator(config)

	// 2. 准备请求参数
	method := "POST"
	path := "/api/v1/auth/login"
	body := `{"username":"admin","password":"123456"}`

	// 3. 生成时间戳和 Nonce
	timestamp := security.GetCurrentTimestamp()
	nonce := security.GenerateNonce(16)

	// 4. 生成签名
	signature := validator.SignRequest(method, path, timestamp, nonce, body)

	// 5. 构建请求头
	fmt.Println("=== 客户端请求 ===")
	fmt.Printf("Method: %s\n", method)
	fmt.Printf("Path: %s\n", path)
	fmt.Printf("Body: %s\n", body)
	fmt.Println("\n请求头:")
	fmt.Printf("  X-Timestamp: %d\n", timestamp)
	fmt.Printf("  X-Nonce: %s\n", nonce)
	fmt.Printf("  X-Signature: %s\n", signature)

	// 6. 服务端验证
	fmt.Println("\n=== 服务端验证 ===")
	result := validator.VerifySignature(method, path, timestamp, nonce, body, signature)
	fmt.Printf("验证结果: %v\n", result.IsValid)
	if result.IsValid {
		fmt.Println("✅ 签名验证通过，允许请求")
	} else {
		fmt.Printf("❌ 签名验证失败: %s\n", result.Reason)
	}
	// Output:
	// 验证结果: true
}

// 示例5：并发安全测试
func Example_ConcurrentSafety() {
	config := security.DefaultSignatureConfig("your-secret-key-12345678")
	validator := security.NewSignatureValidator(config)

	// 模拟并发请求
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(index int) {
			timestamp := security.GetCurrentTimestamp()
			nonce := security.GenerateNonce(16)
			body := fmt.Sprintf(`{"request_id":%d}`, index)

			signature := validator.SignRequest("POST", "/api/v1/test", timestamp, nonce, body)
			result := validator.VerifySignature("POST", "/api/v1/test", timestamp, nonce, body, signature)

			fmt.Printf("Request %d: %v\n", index, result.IsValid)
			done <- true
		}(i)
	}

	// 等待所有请求完成
	for i := 0; i < 10; i++ {
		<-done
	}
	fmt.Println("所有并发请求验证完成")
}

func main() {
	fmt.Println("=== API 签名验证示例 ===\n")

	fmt.Println("示例1: POST 请求签名")
	Example_POSTRequest()
	fmt.Println()

	fmt.Println("示例2: GET 请求带查询参数")
	Example_GETRequestWithQuery()
	fmt.Println()

	fmt.Println("示例3: 签名失败场景")
	Example_SignatureFailure()
	fmt.Println()

	fmt.Println("示例4: 完整的客户端请求流程")
	Example_CompleteFlow()
	fmt.Println()

	fmt.Println("示例5: 并发安全测试")
	Example_ConcurrentSafety()
}
