package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigEncryptor_EncryptDecrypt(t *testing.T) {
	// 使用 32 字节密钥 (AES-256)
	key := "12345678901234567890123456789012"
	encryptor, err := NewConfigEncryptor(key)
	assert.NoError(t, err)

	t.Run("加密解密成功", func(t *testing.T) {
		plaintext := "my-secret-password"
		
		// 加密
		ciphertext, err := encryptor.Encrypt(plaintext)
		assert.NoError(t, err)
		assert.NotEmpty(t, ciphertext)
		assert.NotEqual(t, plaintext, ciphertext)
		
		// 解密
		decrypted, err := encryptor.Decrypt(ciphertext)
		assert.NoError(t, err)
		assert.Equal(t, plaintext, decrypted)
	})

	t.Run("加密不同文本产生不同密文", func(t *testing.T) {
		text1 := "password1"
		text2 := "password2"
		
		cipher1, _ := encryptor.Encrypt(text1)
		cipher2, _ := encryptor.Encrypt(text2)
		
		assert.NotEqual(t, cipher1, cipher2)
	})

	t.Run("相同文本多次加密产生不同密文", func(t *testing.T) {
		text := "same-password"
		
		cipher1, _ := encryptor.Encrypt(text)
		cipher2, _ := encryptor.Encrypt(text)
		
		// 由于使用随机 nonce，每次加密结果不同
		assert.NotEqual(t, cipher1, cipher2)
		
		// 但都能正确解密
		dec1, _ := encryptor.Decrypt(cipher1)
		dec2, _ := encryptor.Decrypt(cipher2)
		assert.Equal(t, text, dec1)
		assert.Equal(t, text, dec2)
	})
}

func TestNewConfigEncryptor_InvalidKey(t *testing.T) {
	t.Run("密钥长度无效", func(t *testing.T) {
		invalidKey := "short"
		_, err := NewConfigEncryptor(invalidKey)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidKey, err)
	})

	t.Run("密钥长度有效-16字节", func(t *testing.T) {
		validKey := "1234567890123456"
		encryptor, err := NewConfigEncryptor(validKey)
		assert.NoError(t, err)
		assert.NotNil(t, encryptor)
	})

	t.Run("密钥长度有效-24字节", func(t *testing.T) {
		validKey := "123456789012345678901234"
		encryptor, err := NewConfigEncryptor(validKey)
		assert.NoError(t, err)
		assert.NotNil(t, encryptor)
	})

	t.Run("密钥长度有效-32字节", func(t *testing.T) {
		validKey := "12345678901234567890123456789012"
		encryptor, err := NewConfigEncryptor(validKey)
		assert.NoError(t, err)
		assert.NotNil(t, encryptor)
	})
}

func TestConfigEncryptor_Decrypt_InvalidCiphertext(t *testing.T) {
	key := "12345678901234567890123456789012"
	encryptor, _ := NewConfigEncryptor(key)

	t.Run("解密无效密文", func(t *testing.T) {
		_, err := encryptor.Decrypt("invalid-ciphertext")
		assert.Error(t, err)
	})

	t.Run("解密空字符串", func(t *testing.T) {
		_, err := encryptor.Decrypt("")
		assert.Error(t, err)
	})
}

func TestEncryptSensitiveConfig(t *testing.T) {
	key := "12345678901234567890123456789012"
	encryptor, _ := NewConfigEncryptor(key)

	t.Run("加密敏感配置", func(t *testing.T) {
		config := &EncryptedConfig{
			DatabasePassword: "db-password-123",
			JWTSecret:        "jwt-secret-xyz",
			RedisPassword:    "redis-pass",
			SMSAPIKey:        "sms-api-key",
			MailPassword:     "mail-pass",
		}

		encrypted, err := EncryptSensitiveConfig(config, encryptor)
		assert.NoError(t, err)
		assert.NotNil(t, encrypted)
		
		// 验证已加密
		assert.NotEqual(t, config.DatabasePassword, encrypted.DatabasePassword)
		assert.NotEqual(t, config.JWTSecret, encrypted.JWTSecret)
		assert.NotEqual(t, config.RedisPassword, encrypted.RedisPassword)

		// 解密验证
		decrypted, err := DecryptSensitiveConfig(encrypted, encryptor)
		assert.NoError(t, err)
		assert.Equal(t, config.DatabasePassword, decrypted.DatabasePassword)
		assert.Equal(t, config.JWTSecret, decrypted.JWTSecret)
		assert.Equal(t, config.RedisPassword, decrypted.RedisPassword)
	})

	t.Run("部分字段加密", func(t *testing.T) {
		config := &EncryptedConfig{
			DatabasePassword: "db-password",
			JWTSecret:        "",
		}

		encrypted, err := EncryptSensitiveConfig(config, encryptor)
		assert.NoError(t, err)
		assert.NotEmpty(t, encrypted.DatabasePassword)
		assert.Empty(t, encrypted.JWTSecret)
	})
}
