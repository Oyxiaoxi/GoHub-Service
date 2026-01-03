// Package security 敏感配置加密存储
package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

var (
	// ErrInvalidKey 无效的密钥
	ErrInvalidKey = errors.New("invalid encryption key")
	// ErrInvalidCiphertext 无效的密文
	ErrInvalidCiphertext = errors.New("invalid ciphertext")
)

// ConfigEncryptor 配置加密器
type ConfigEncryptor struct {
	key []byte
}

// NewConfigEncryptor 创建配置加密器
// key 必须是 16, 24 或 32 字节（对应 AES-128, AES-192 或 AES-256）
func NewConfigEncryptor(key string) (*ConfigEncryptor, error) {
	keyBytes := []byte(key)
	
	// 验证密钥长度
	if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
		return nil, ErrInvalidKey
	}
	
	return &ConfigEncryptor{
		key: keyBytes,
	}, nil
}

// NewConfigEncryptorFromEnv 从环境变量创建加密器
func NewConfigEncryptorFromEnv() (*ConfigEncryptor, error) {
	key := os.Getenv("CONFIG_ENCRYPTION_KEY")
	if key == "" {
		return nil, errors.New("CONFIG_ENCRYPTION_KEY environment variable not set")
	}
	return NewConfigEncryptor(key)
}

// Encrypt 加密配置值
func (e *ConfigEncryptor) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	// 创建 GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 生成随机 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	
	// Base64 编码
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密配置值
func (e *ConfigEncryptor) Decrypt(ciphertext string) (string, error) {
	// Base64 解码
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", ErrInvalidCiphertext
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// EncryptConfig 加密敏感配置
type EncryptedConfig struct {
	DatabasePassword string `json:"database_password"`
	JWTSecret        string `json:"jwt_secret"`
	RedisPassword    string `json:"redis_password"`
	SMSAPIKey        string `json:"sms_api_key"`
	MailPassword     string `json:"mail_password"`
}

// EncryptSensitiveConfig 加密所有敏感配置
func EncryptSensitiveConfig(config *EncryptedConfig, encryptor *ConfigEncryptor) (*EncryptedConfig, error) {
	encrypted := &EncryptedConfig{}
	
	if config.DatabasePassword != "" {
		enc, err := encryptor.Encrypt(config.DatabasePassword)
		if err != nil {
			return nil, err
		}
		encrypted.DatabasePassword = enc
	}
	
	if config.JWTSecret != "" {
		enc, err := encryptor.Encrypt(config.JWTSecret)
		if err != nil {
			return nil, err
		}
		encrypted.JWTSecret = enc
	}
	
	if config.RedisPassword != "" {
		enc, err := encryptor.Encrypt(config.RedisPassword)
		if err != nil {
			return nil, err
		}
		encrypted.RedisPassword = enc
	}
	
	if config.SMSAPIKey != "" {
		enc, err := encryptor.Encrypt(config.SMSAPIKey)
		if err != nil {
			return nil, err
		}
		encrypted.SMSAPIKey = enc
	}
	
	if config.MailPassword != "" {
		enc, err := encryptor.Encrypt(config.MailPassword)
		if err != nil {
			return nil, err
		}
		encrypted.MailPassword = enc
	}
	
	return encrypted, nil
}

// DecryptSensitiveConfig 解密所有敏感配置
func DecryptSensitiveConfig(encrypted *EncryptedConfig, encryptor *ConfigEncryptor) (*EncryptedConfig, error) {
	config := &EncryptedConfig{}
	
	if encrypted.DatabasePassword != "" {
		dec, err := encryptor.Decrypt(encrypted.DatabasePassword)
		if err != nil {
			return nil, err
		}
		config.DatabasePassword = dec
	}
	
	if encrypted.JWTSecret != "" {
		dec, err := encryptor.Decrypt(encrypted.JWTSecret)
		if err != nil {
			return nil, err
		}
		config.JWTSecret = dec
	}
	
	if encrypted.RedisPassword != "" {
		dec, err := encryptor.Decrypt(encrypted.RedisPassword)
		if err != nil {
			return nil, err
		}
		config.RedisPassword = dec
	}
	
	if encrypted.SMSAPIKey != "" {
		dec, err := encryptor.Decrypt(encrypted.SMSAPIKey)
		if err != nil {
			return nil, err
		}
		config.SMSAPIKey = dec
	}
	
	if encrypted.MailPassword != "" {
		dec, err := encryptor.Decrypt(encrypted.MailPassword)
		if err != nil {
			return nil, err
		}
		config.MailPassword = dec
	}
	
	return config, nil
}
