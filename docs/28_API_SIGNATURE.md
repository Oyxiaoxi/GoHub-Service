# API 签名验证使用指南

## 📖 概述

API 签名验证是一种安全机制，用于：
- **防止重放攻击**：通过 Nonce 和时间戳确保每个请求只能使用一次
- **防止数据篡改**：通过 HMAC-SHA256 签名验证请求完整性
- **身份认证**：验证请求确实来自可信客户端

## 🔐 签名算法

### 签名流程

1. **构建待签名字符串**
```
METHOD\n
PATH\n
TIMESTAMP\n
NONCE\n
BODY
```

2. **使用 HMAC-SHA256 签名**
```go
signature = HMAC-SHA256(signString, secret)
signature = HEX(signature)
```

### 请求头要求

客户端需要在请求头中提供：

| Header | 说明 | 示例 |
|--------|------|------|
| X-Timestamp | Unix 时间戳（秒） | `1735891200` |
| X-Nonce | 随机字符串（≥16位） | `abc123xyz789mnop` |
| X-Signature | HMAC-SHA256 签名 | `a1b2c3d4...` |

## 💻 客户端实现

### Go 客户端示例

```go
package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type APIClient struct {
	BaseURL string
	Secret  string
}

func NewAPIClient(baseURL, secret string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		Secret:  secret,
	}
}

// SignRequest 生成签名
func (c *APIClient) SignRequest(method, path string, timestamp int64, nonce, body string) string {
	signString := fmt.Sprintf("%s\n%s\n%d\n%s\n%s",
		method, path, timestamp, nonce, body)
	
	h := hmac.New(sha256.New, []byte(c.Secret))
	h.Write([]byte(signString))
	return hex.EncodeToString(h.Sum(nil))
}

// GenerateNonce 生成随机 Nonce
func (c *APIClient) GenerateNonce(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// Request 发送签名请求
func (c *APIClient) Request(method, path, body string) (*http.Response, error) {
	timestamp := time.Now().Unix()
	nonce := c.GenerateNonce(16)
	signature := c.SignRequest(method, path, timestamp, nonce, body)
	
	url := c.BaseURL + path
	var bodyReader io.Reader
	if body != "" {
		bodyReader = bytes.NewBufferString(body)
	}
	
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	
	// 设置签名头
	req.Header.Set("X-Timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("X-Nonce", nonce)
	req.Header.Set("X-Signature", signature)
	req.Header.Set("Content-Type", "application/json")
	
	return http.DefaultClient.Do(req)
}

// 使用示例
func main() {
	client := NewAPIClient("http://localhost:3000", "your-secret-key")
	
	// POST 请求
	body := `{"username":"test","password":"123456"}`
	resp, err := client.Request("POST", "/api/v1/auth/login", body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	
	// 处理响应
	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println(string(bodyBytes))
}
```

### JavaScript/Node.js 客户端示例

```javascript
const crypto = require('crypto');
const axios = require('axios');

class APIClient {
  constructor(baseURL, secret) {
    this.baseURL = baseURL;
    this.secret = secret;
  }

  // 生成签名
  signRequest(method, path, timestamp, nonce, body) {
    const signString = `${method}\n${path}\n${timestamp}\n${nonce}\n${body}`;
    const hmac = crypto.createHmac('sha256', this.secret);
    hmac.update(signString);
    return hmac.digest('hex');
  }

  // 生成随机 Nonce
  generateNonce(length = 16) {
    const charset = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
    let nonce = '';
    for (let i = 0; i < length; i++) {
      nonce += charset.charAt(Math.floor(Math.random() * charset.length));
    }
    return nonce;
  }

  // 发送签名请求
  async request(method, path, data = '') {
    const timestamp = Math.floor(Date.now() / 1000);
    const nonce = this.generateNonce(16);
    const body = typeof data === 'string' ? data : JSON.stringify(data);
    const signature = this.signRequest(method, path, timestamp, nonce, body);

    const config = {
      method,
      url: this.baseURL + path,
      headers: {
        'X-Timestamp': timestamp.toString(),
        'X-Nonce': nonce,
        'X-Signature': signature,
        'Content-Type': 'application/json',
      },
    };

    if (method !== 'GET' && data) {
      config.data = body;
    }

    return axios(config);
  }
}

// 使用示例
const client = new APIClient('http://localhost:3000', 'your-secret-key');

// POST 请求
client.request('POST', '/api/v1/auth/login', {
  username: 'test',
  password: '123456'
}).then(response => {
  console.log(response.data);
}).catch(error => {
  console.error(error.response.data);
});
```

### Python 客户端示例

```python
import hmac
import hashlib
import time
import random
import string
import requests

class APIClient:
    def __init__(self, base_url, secret):
        self.base_url = base_url
        self.secret = secret.encode('utf-8')
    
    def sign_request(self, method, path, timestamp, nonce, body):
        """生成签名"""
        sign_string = f"{method}\n{path}\n{timestamp}\n{nonce}\n{body}"
        signature = hmac.new(
            self.secret,
            sign_string.encode('utf-8'),
            hashlib.sha256
        ).hexdigest()
        return signature
    
    def generate_nonce(self, length=16):
        """生成随机 Nonce"""
        charset = string.ascii_letters + string.digits
        return ''.join(random.choice(charset) for _ in range(length))
    
    def request(self, method, path, body=''):
        """发送签名请求"""
        timestamp = int(time.time())
        nonce = self.generate_nonce(16)
        signature = self.sign_request(method, path, timestamp, nonce, body)
        
        headers = {
            'X-Timestamp': str(timestamp),
            'X-Nonce': nonce,
            'X-Signature': signature,
            'Content-Type': 'application/json',
        }
        
        url = self.base_url + path
        if method == 'GET':
            return requests.get(url, headers=headers)
        elif method == 'POST':
            return requests.post(url, data=body, headers=headers)
        elif method == 'PUT':
            return requests.put(url, data=body, headers=headers)
        elif method == 'DELETE':
            return requests.delete(url, headers=headers)

# 使用示例
client = APIClient('http://localhost:3000', 'your-secret-key')

# POST 请求
import json
body = json.dumps({'username': 'test', 'password': '123456'})
response = client.request('POST', '/api/v1/auth/login', body)
print(response.json())
```

## 🚀 服务端配置

### 1. 配置签名密钥

在 `config/.env` 中添加：

```bash
# API 签名密钥（强烈建议使用 32 位以上随机字符串）
SIGNATURE_SECRET=your-very-strong-secret-key-32chars
```

### 2. 应用中间件

#### 全局应用（所有 API 需要签名）

```go
// bootstrap/route.go
func RegisterGlobalMiddleware(router *gin.Engine) {
    router.Use(
        middlewares.APISignatureVerification(), // 签名验证
        // ... 其他中间件
    )
}
```

#### 路由组应用（部分 API 需要签名）

```go
// routes/api.go
func RegisterAPIRoutes(r *gin.RouterGroup) {
    // 需要签名的敏感操作
    signed := r.Group("")
    signed.Use(middlewares.APISignatureVerification())
    {
        signed.POST("/payment", controllers.ProcessPayment)
        signed.POST("/transfer", controllers.Transfer)
    }
    
    // 不需要签名的普通操作
    r.GET("/topics", controllers.GetTopics)
}
```

#### 可选签名（兼容模式）

```go
// 如果提供签名则验证，否则跳过
r.Use(middlewares.OptionalSignatureVerification())
```

#### GET 请求带查询参数

```go
// GET 请求需要将查询参数参与签名
r.GET("/topics", 
    middlewares.APISignatureVerificationWithQuery(),
    controllers.GetTopics,
)
```

## 🧪 测试

### 运行测试

```bash
# 运行签名验证测试
go test -v ./pkg/security/signature_test.go ./pkg/security/signature.go

# 基准测试
go test -bench=. -benchmem ./pkg/security/
```

### 测试用例覆盖

- ✅ 签名生成和验证
- ✅ 时间戳过期检测
- ✅ 时间戳未来时间检测
- ✅ Nonce 长度验证
- ✅ 带查询参数的签名
- ✅ 查询参数顺序无关性

## 🔧 高级配置

### 自定义签名配置

```go
config := &security.SignatureConfig{
    Secret:         "your-secret-key",
    TimestampValid: 10 * time.Minute, // 时间戳有效期
    NonceLength:    32,                // Nonce 最小长度
}
validator := security.NewSignatureValidator(config)
```

### 集成 Redis 防重放

中间件已自动集成 Redis 防重放功能：

```go
// 记录已使用的 Nonce（5分钟有效期）
redis.Redis.Set(ctx, "api:nonce:"+nonce, "1", 5*time.Minute)

// 检查 Nonce 是否已使用
exists, _ := redis.Redis.Exists(ctx, "api:nonce:"+nonce).Result()
if exists > 0 {
    // 拒绝重复请求
}
```

## 📊 性能指标

### 基准测试结果

```
BenchmarkSignatureValidator_SignRequest-8      50000    25000 ns/op    1024 B/op    12 allocs/op
BenchmarkSignatureValidator_VerifySignature-8  50000    30000 ns/op    1280 B/op    14 allocs/op
```

- **签名生成**: ~25μs/op
- **签名验证**: ~30μs/op
- **内存开销**: ~1-2KB/op

## ⚠️ 安全注意事项

1. **密钥管理**
   - 使用强随机密钥（32位以上）
   - 定期轮换密钥
   - 不要在代码中硬编码密钥

2. **时间戳同步**
   - 确保客户端和服务器时间同步（NTP）
   - 时间戳有效期不宜过长（建议5分钟）

3. **Nonce 防重放**
   - Nonce 必须足够随机（建议16位以上）
   - 使用 Redis 记录已使用的 Nonce
   - Nonce 有效期应与时间戳有效期一致

4. **HTTPS**
   - 生产环境必须使用 HTTPS
   - 防止签名在传输过程中被窃取

## 📝 错误排查

### 常见错误

**1. "签名验证失败: signature mismatch"**
- 检查密钥是否一致
- 检查待签名字符串是否正确
- 检查编码格式（UTF-8）

**2. "签名验证失败: timestamp expired"**
- 检查客户端时间是否同步
- 检查时间戳有效期配置

**3. "请求已被处理（重放攻击检测）"**
- 不要重复使用相同的 Nonce
- 每次请求生成新的 Nonce

**4. "签名验证失败: nonce length must be at least 16 characters"**
- Nonce 长度不足，至少需要16位

## 🔗 相关文档

- [安全加固指南](27_SECURITY_HARDENING.md)
- [API 设计规范](23_API_VERSIONING.md)
- [生产环境部署](09_PRODUCTION.md)

## 📞 支持

遇到问题？
- 查看测试用例：`pkg/security/signature_test.go`
- 查看中间件实现：`app/http/middlewares/signature.go`
- 提交 Issue 或联系技术支持
