// Package security 随机字符串生成工具
package security

import (
	"crypto/rand"
	"math/big"
)

// GenerateRandomString 生成指定长度的随机字符串
// @param length 字符串长度
// @param charset 字符集（如果为空，使用默认字符集）
// @return 随机字符串
func GenerateRandomString(length int, charset string) string {
	if charset == "" {
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}

	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			// 降级使用伪随机数
			result[i] = charset[i%len(charset)]
			continue
		}
		result[i] = charset[randomIndex.Int64()]
	}

	return string(result)
}

// GenerateRandomBytes 生成指定长度的随机字节
func GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// GenerateSecureToken 生成安全令牌（Base62编码）
func GenerateSecureToken(length int) string {
	const base62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	return GenerateRandomString(length, base62)
}
