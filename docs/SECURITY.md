# 🛡️ API 安全指南

API 安全最佳实践和防御措施。

## 1. 认证与授权

### JWT 令牌验证

```go
// 在请求头中传递 JWT
Authorization: Bearer <token>

// 令牌结构
{
    "user_id": 123,
    "username": "john",
    "roles": ["user"],
    "exp": 1735689600
}
```

### 权限检查

所有受保护的端点都必须通过中间件检查权限：

```go
// 中间件链
r.POST("/topics", 
    middlewares.Authenticate(),      // 验证认证
    middlewares.RequirePermission("topics.create"), // 检查权限
    controllers.TopicStore)
```

## 2. 输入验证

### 使用请求验证器

文件位置: `app/requests/`

```go
type CreateTopicRequest struct {
    Title       string `binding:"required,min=3,max=100"`
    Content     string `binding:"required,min=10,max=5000"`
    CategoryID  uint   `binding:"required,min=1"`
}

func (r *CreateTopicRequest) Validate() error {
    if len(r.Title) < 3 {
        return errors.New("标题长度至少3个字符")
    }
    return nil
}
```

### 防止 SQL 注入

✅ **正确做法**：使用参数化查询
```go
var user models.User
db.Where("username = ?", username).First(&user)  // 参数化

var user models.User
db.Where("username = ?", username).Where("status = ?", "active").First(&user)
```

❌ **错误做法**：字符串拼接
```go
var user models.User
db.Where("username = " + username).First(&user)  // 危险！
```

### 防止 XSS 攻击

在返回 HTML 响应时始终转义用户输入：

```go
// 在模板中自动转义
{{ .UserContent }}  // 自动转义

// 手动转义
import "html"
safeHTML := html.EscapeString(userInput)
```

### 防止 CSRF 攻击

✅ 使用 HTTPS
✅ 验证 Referer 或 Origin 头
✅ 为状态改变的请求使用令牌

```go
// 验证 Origin 头
origin := c.GetHeader("Origin")
if !isAllowedOrigin(origin) {
    c.AbortWithStatus(http.StatusForbidden)
    return
}
```

## 3. 数据安全

### 密码存储

✅ 使用强哈希算法（bcrypt）
```go
import "golang.org/x/crypto/bcrypt"

// 哈希密码
hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// 验证密码
bcrypt.CompareHashAndPassword(hash, []byte(password))
```

### 敏感数据隐藏

在 API 响应中隐藏敏感字段：

```go
type UserResponse struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    // ❌ 不返回
    // Password string `json:"password"`
}
```

### HTTPS 强制

在生产环境中必须使用 HTTPS：

```go
// 配置中启用 HTTPS
config.TLS.Enabled = true
config.TLS.CertFile = "/path/to/cert.pem"
config.TLS.KeyFile = "/path/to/key.pem"
```

## 4. 速率限制

### 限流配置

文件位置: `config/limiter.go`

```go
type LimiterConfig struct {
    Enabled     bool
    RequestsPerSecond int  // 每秒请求数
    BurstSize   int       // 突发大小
}
```

### 使用限流中间件

```go
r.Use(middlewares.RateLimit())
```

### 响应格式

```
429 Too Many Requests

{
    "error": "请求过于频繁，请稍后重试",
    "retry_after": 60
}
```

## 5. 日志与审计

### 审计日志

记录所有敏感操作：

```go
// 用户登录
auditLog.Create(&models.AuditLog{
    UserID: user.ID,
    Action: "login",
    IP: c.ClientIP(),
    CreatedAt: time.Now(),
})

// 权限变更
auditLog.Create(&models.AuditLog{
    UserID: adminID,
    Action: "assign_role",
    Details: fmt.Sprintf("给用户 %d 分配角色 %d", userID, roleID),
})
```

### 日志安全

✅ 不记录密码或令牌
✅ 定期轮换日志文件
✅ 限制日志访问权限
✅ 加密敏感日志内容

```go
// ✅ 安全的日志
logger.Info("用户登录", zap.Uint("user_id", user.ID))

// ❌ 不安全的日志
logger.Info("用户登录", zap.String("password", password))
```

## 6. 错误处理

### 信息泄露防止

不要在错误消息中暴露敏感信息：

```go
// ❌ 坏的做法
if err != nil {
    c.JSON(500, gin.H{"error": err.Error()})  // 暴露数据库错误
}

// ✅ 好的做法
if err != nil {
    logger.Error("数据库错误", zap.Error(err))
    c.JSON(500, gin.H{"error": "发生错误，请稍后重试"})
}
```

### 标准错误响应

```json
{
    "error": "用户不存在",
    "code": "NOT_FOUND",
    "status": 404
}
```

## 7. API 版本控制

### 使用版本前缀

```
/api/v1/topics      ✅ 推荐
/api/v2/topics      ✅ 新版本
/topics             ❌ 避免
```

### 版本兼容性

```go
// routes/api.go
v1 := r.Group("/api/v1")
{
    v1.GET("/topics", controllers.TopicIndex)
}

v2 := r.Group("/api/v2")
{
    v2.GET("/topics", controllers.TopicIndexV2)
}
```

## 8. 依赖安全

### 定期更新依赖

```bash
# 检查漏洞
go list -json -m all | nancy sleuth

# 更新依赖
go get -u ./...

# 审计依赖
go mod audit
```

### 依赖版本锁定

```
go.mod 使用精确版本
go.sum 记录校验和
```

## 9. CORS 配置

### 安全的 CORS 设置

```go
config := cors.Config{
    AllowOrigins:     []string{"https://example.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Authorization", "Content-Type"},
    ExposeHeaders:    []string{"X-Total-Count"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}
```

❌ 不要使用通配符
```go
AllowOrigins: []string{"*"}  // 危险！
```

## 10. 安全检查清单

- [ ] 所有输入验证
- [ ] SQL 注入防护（参数化查询）
- [ ] XSS 防护（输入转义）
- [ ] CSRF 令牌验证
- [ ] 密码强度检查（最少8字符，混合字符）
- [ ] 密码安全存储（bcrypt）
- [ ] HTTPS 强制
- [ ] 速率限制
- [ ] JWT 过期时间设置（推荐1小时）
- [ ] 刷新令牌机制
- [ ] 审计日志
- [ ] 错误信息不泄露
- [ ] CORS 正确配置
- [ ] 依赖定期更新
- [ ] 敏感数据不记录
- [ ] 实施密钥轮换

## 常见漏洞

### 1. 暴露用户 ID 序列
不要假设 ID 是难以猜测的，使用权限检查防止越权

### 2. API 端点暴露
不要在 API 文档中暴露管理员端点，使用版本控制隐藏

### 3. 信息过多的错误消息
永远提供通用错误消息，不要暴露系统细节

### 4. 过期令牌仍然有效
始终检查 JWT 的 exp 声明

### 5. 权限检查不完整
每个端点都要检查用户权限，不要假设路由就足够了

---

更多信息请查看 [ARCHITECTURE.md](./ARCHITECTURE.md)
