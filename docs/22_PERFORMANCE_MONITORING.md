# 性能监控与敏感配置加密 v3.9

本文档说明了 v3.9 版本新增的性能监控和敏感配置加密功能。

## 1. 慢查询日志记录

### 功能说明
自动记录执行时间超过阈值的数据库查询，帮助识别性能瓶颈。

### 使用方法

#### 启用慢查询日志
```go
import (
    "github.com/Oyxiaoxi/GoHub-Service/pkg/database"
    "time"
)

// 在数据库初始化时启用（200ms 阈值）
db = database.EnableSlowQueryLog(db, 200*time.Millisecond)
```

#### 查看慢查询统计
```go
stats := database.GetSlowQueryStats()
fmt.Printf("慢查询总数: %d\n", stats.SlowCount)
fmt.Printf("平均执行时间: %v\n", stats.AverageTime)
fmt.Printf("最大执行时间: %v\n", stats.MaxTime)

for _, query := range stats.SlowQueries {
    fmt.Printf("SQL: %s, 耗时: %v\n", query.SQL, query.Duration)
}
```

#### 清除统计
```go
database.ClearSlowQueryStats()
```

### 配置说明
- 默认阈值：200ms
- 最大记录数：100条（超出后自动清理最旧记录）
- 日志级别：WARN

### 日志格式
```json
{
    "level": "warn",
    "msg": "slow_query",
    "sql": "SELECT * FROM users WHERE ...",
    "elapsed_ms": 350,
    "rows_affected": 10,
    "threshold_ms": 200
}
```

## 2. 接口响应时间统计

### 功能说明
自动记录每个 API 接口的响应时间，提供统计分析功能。

### 使用方法

#### 启用性能监控中间件
```go
import "github.com/Oyxiaoxi/GoHub-Service/app/http/middlewares"

// 方式1：仅记录日志
router.Use(middlewares.PerformanceMonitor())

// 方式2：记录日志 + 统计
router.Use(middlewares.PerformanceStats())
```

#### 查看接口统计
```go
stats := middlewares.GetAPIStats()

for endpoint, stat := range stats {
    fmt.Printf("接口: %s\n", endpoint)
    fmt.Printf("  总调用次数: %d\n", stat.Count)
    fmt.Printf("  成功次数: %d\n", stat.SuccessCount)
    fmt.Printf("  失败次数: %d\n", stat.ErrorCount)
    fmt.Printf("  平均响应时间: %v\n", stat.AverageTime)
    fmt.Printf("  最大响应时间: %v\n", stat.MaxTime)
    fmt.Printf("  最小响应时间: %v\n", stat.MinTime)
}
```

#### 清除统计
```go
middlewares.ClearAPIStats()
```

### 特性
- 自动记录每个接口的响应时间
- 区分成功/失败请求
- 计算平均/最大/最小响应时间
- 慢接口警告（超过 1 秒）
- 响应头包含执行时间：`X-Response-Time`

### 日志格式
```json
{
    "level": "info",
    "msg": "api_performance",
    "method": "GET",
    "path": "/api/v1/users",
    "status": 200,
    "duration_ms": 125,
    "client_ip": "127.0.0.1"
}
```

### 慢接口日志
```json
{
    "level": "warn",
    "msg": "slow_api",
    "method": "POST",
    "path": "/api/v1/topics",
    "duration_ms": 1500,
    "threshold_ms": 1000
}
```

## 3. 敏感配置加密存储

### 功能说明
使用 AES-GCM 加密算法保护敏感配置（数据库密码、API密钥等）。

### 使用方法

#### 设置加密密钥
```bash
# 环境变量方式（推荐）
export CONFIG_ENCRYPTION_KEY="your-32-byte-secret-key-here-123"

# 密钥长度要求：16/24/32 字节（AES-128/192/256）
```

#### 加密敏感配置
```go
import "github.com/Oyxiaoxi/GoHub-Service/pkg/security"

// 从环境变量创建加密器
encryptor, err := security.NewConfigEncryptorFromEnv()
if err != nil {
    panic(err)
}

// 加密单个值
ciphertext, err := encryptor.Encrypt("my-database-password")

// 加密配置结构
config := &security.EncryptedConfig{
    DatabasePassword: "db-pass-123",
    JWTSecret:        "jwt-secret-xyz",
    RedisPassword:    "redis-pass",
}

encrypted, err := security.EncryptSensitiveConfig(config, encryptor)
```

#### 解密敏感配置
```go
// 解密单个值
plaintext, err := encryptor.Decrypt(ciphertext)

// 解密配置结构
decrypted, err := security.DecryptSensitiveConfig(encrypted, encryptor)
```

### 安全特性
- **AES-GCM 加密**：认证加密模式，防篡改
- **随机 Nonce**：每次加密使用不同随机数
- **Base64 编码**：便于存储和传输
- **环境变量密钥**：密钥不硬编码在代码中

### 支持的敏感配置
- DatabasePassword - 数据库密码
- JWTSecret - JWT 签名密钥
- RedisPassword - Redis 密码
- SMSAPIKey - 短信服务 API 密钥
- MailPassword - 邮件服务密码

### 最佳实践
1. **密钥管理**
   - 使用 32 字节密钥（AES-256）
   - 通过环境变量传递密钥
   - 不要提交密钥到代码仓库
   - 定期轮换密钥

2. **配置存储**
   - 加密后的配置可以安全存储在配置文件
   - 可以提交到版本控制系统
   - 不同环境使用不同密钥

3. **应用启动**
   - 在应用启动时解密配置
   - 仅在内存中保存明文配置
   - 日志不记录敏感信息

## 4. 集成示例

### bootstrap/database.go 集成
```go
import (
    "github.com/Oyxiaoxi/GoHub-Service/pkg/database"
    "github.com/Oyxiaoxi/GoHub-Service/pkg/security"
)

func InitDatabase() {
    // 解密数据库密码
    encryptor, _ := security.NewConfigEncryptorFromEnv()
    dbPassword, _ := encryptor.Decrypt(config.Get("database.password"))
    
    // 连接数据库
    db := database.Connect(dbConfig)
    
    // 启用慢查询日志
    db = database.EnableSlowQueryLog(db, 200*time.Millisecond)
}
```

### bootstrap/route.go 集成
```go
import "github.com/Oyxiaoxi/GoHub-Service/app/http/middlewares"

func RegisterGlobalMiddleware(router *gin.Engine) {
    // 性能监控
    router.Use(middlewares.PerformanceStats())
    
    // 其他中间件...
}
```

### 监控接口示例
```go
// controllers/monitor_controller.go
func GetPerformanceStats(c *gin.Context) {
    // 数据库慢查询统计
    dbStats := database.GetSlowQueryStats()
    
    // API 接口统计
    apiStats := middlewares.GetAPIStats()
    
    response.JSON(c, gin.H{
        "database": dbStats,
        "api":      apiStats,
    })
}
```

## 5. 性能影响

- **慢查询日志**：< 1ms overhead
- **接口监控**：< 0.5ms overhead
- **配置加密**：仅在启动时执行，无运行时影响

## 6. 配置建议

### 慢查询阈值
- **开发环境**：100ms
- **测试环境**：200ms
- **生产环境**：300ms

### 慢接口阈值
- **查询接口**：500ms
- **写入接口**：1000ms
- **批量操作**：2000ms

## 7. 监控告警

### 慢查询告警
- 单次查询超过 1 秒：立即告警
- 1 分钟内慢查询超过 10 次：告警
- 慢查询数量持续增长：告警

### 慢接口告警
- 单次请求超过 3 秒：立即告警
- 接口平均响应时间超过 1 秒：告警
- 错误率超过 5%：告警

## 更新日志

### v3.9 (2026-01-03)
- ✅ 新增慢查询日志记录功能
- ✅ 新增接口响应时间统计功能
- ✅ 新增敏感配置加密存储功能
- ✅ 提供完整的测试用例
- ✅ 提供使用文档和示例
