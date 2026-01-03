# 安全加固指南 v5.0

## 概述

GoHub-Service v5.0 实施了全面的安全加固措施，覆盖输入验证、SQL注入防护、XSS防护、限流等多个安全层面。

## 1. 输入验证

### 输入验证器

位置：`pkg/security/validator.go`

#### 核心功能

##### SQL 注入检测
```go
validator := security.NewInputValidator()
result := validator.CheckSQLInjection("user input")
if !result.IsValid {
    // 处理非法输入
    log.Warn("SQL注入尝试", "reason", result.Reason)
}
```

检测模式：
- UNION SELECT 注入
- INSERT/UPDATE/DELETE 注入
- 注释符号（--, #, /* */）
- OR/AND 布尔注入

##### XSS 攻击检测
```go
result := validator.CheckXSS("user input")
if !result.IsValid {
    // 处理XSS攻击
}
```

检测模式：
- `<script>` 标签
- `javascript:` 协议
- 事件处理器（onclick, onerror等）
- iframe, embed, object 标签

##### 路径遍历检测
```go
result := validator.CheckPathTraversal("user/path")
if !result.IsValid {
    // 处理路径遍历攻击
}
```

检测模式：
- `../` 和 `..\`
- URL编码的路径遍历
- 各种变体

##### 综合验证
```go
result := validator.Validate("user input")
if !result.IsValid {
    return fmt.Errorf("输入验证失败: %s (%s)", result.Reason, result.RiskType)
}
```

### 数据格式验证

#### 邮箱验证
```go
if !security.IsValidEmail("user@example.com") {
    return errors.New("邮箱格式不正确")
}
```

#### 手机号验证（中国）
```go
if !security.IsValidPhone("13800138000") {
    return errors.New("手机号格式不正确")
}
```

#### URL 验证
```go
if !security.IsValidURL("https://example.com") {
    return errors.New("URL格式不正确")
}
```

### 密码强度验证

```go
result := security.ValidatePasswordStrength(password)
if !result.IsValid {
    return fmt.Errorf("密码强度不足: %v", result.Issues)
}

// 密码得分：0-100
fmt.Printf("密码强度得分: %d/100\n", result.Score)

// 详细信息
fmt.Printf("包含大写: %v\n", result.HasUppercase)
fmt.Printf("包含小写: %v\n", result.HasLowercase)
fmt.Printf("包含数字: %v\n", result.HasDigit)
fmt.Printf("包含特殊字符: %v\n", result.HasSpecial)
```

密码要求：
- 最小长度：8 个字符
- 最大长度：128 个字符
- 至少包含 3 种字符类型（大写/小写/数字/特殊字符）

得分规则：
- 长度 8-11: 20分
- 长度 12-15: +10分
- 长度 16+: +10分
- 包含大写: +15分
- 包含小写: +15分
- 包含数字: +15分
- 包含特殊字符: +15分

## 2. 中间件安全

### 增强安全验证中间件

位置：`app/http/middlewares/security_enhanced.go`

#### EnhancedSecurityValidation
自动验证所有请求的查询参数和路径参数。

```go
// 在路由中启用
router.Use(middlewares.EnhancedSecurityValidation())
```

功能：
- 自动检测 SQL 注入
- 自动检测 XSS 攻击
- 自动检测路径遍历
- 长度限制验证

#### SQLInjectionProtection
专门的 SQL 注入防护中间件。

```go
router.Use(middlewares.SQLInjectionProtection())
```

#### EnhancedXSSProtection
增强的 XSS 防护中间件。

```go
router.Use(middlewares.EnhancedXSSProtection())
```

### 限流中间件

#### IP 级别限流
```go
// 每分钟最多 60 个请求
router.Use(middlewares.RateLimitMiddleware(60))
```

特性：
- 基于 IP 的限流
- 滑动时间窗口
- 自动封禁（超限后封禁 1 分钟）
- 自动清理过期访客

#### 路由级别限流
```go
// 针对特定路由的限流
router.POST("/api/login", 
    middlewares.LimitPerRoute("10-M"), // 每分钟10次
    loginHandler,
)
```

## 3. 现有安全措施

### 安全响应头

位置：`app/http/middlewares/security.go`

已启用的安全头：
```
X-Frame-Options: DENY
X-Content-Type-Options: nosniff
X-XSS-Protection: 1; mode=block
Content-Security-Policy: default-src 'self'
Referrer-Policy: strict-origin-when-cross-origin
Permissions-Policy: geolocation=(), microphone=(), camera=()
```

### GORM 参数化查询

GoHub-Service 使用 GORM ORM，所有数据库查询都是参数化的，天然防止 SQL 注入：

```go
// ✅ 安全 - 参数化查询
db.Where("email = ?", email).First(&user)

// ❌ 不安全 - 避免使用
db.Where("email = '" + email + "'").First(&user)
```

### 输入验证

所有用户输入都经过验证：
- 位置：`app/requests/*.go`
- 使用 validator 库进行数据验证
- 自定义验证规则

示例：
```go
type UserRequest struct {
    Email    string `json:"email" validate:"required,email,max=254"`
    Phone    string `json:"phone" validate:"required,len=11"`
    Password string `json:"password" validate:"required,min=8,max=128"`
}
```

## 4. 应用安全中间件

### 全局中间件

在 `bootstrap/route.go` 中启用：

```go
func registerGlobalMiddleWare(router *gin.Engine) {
    // 安全响应头
    router.Use(middlewares.SecureHeaders())
    
    // XSS 防护
    router.Use(middlewares.XSSProtection())
    
    // 增强安全验证（新增）
    router.Use(middlewares.EnhancedSecurityValidation())
    
    // SQL 注入防护（新增）
    router.Use(middlewares.SQLInjectionProtection())
    
    // 全局限流（新增）
    router.Use(middlewares.RateLimitMiddleware(200)) // 每分钟200请求
    
    // ... 其他中间件
}
```

### 路由级别限流

在 `routes/api.go` 中：

```go
// 认证相关路由 - 更严格的限流
authGroup := v1.Group("/auth")
authGroup.Use(middlewares.RateLimitMiddleware(20)) // 每分钟20次
{
    authGroup.POST("/login", loginCtrl.Login)
    authGroup.POST("/register", signupCtrl.Register)
}

// 敏感操作 - 极严格的限流
sensitiveGroup := v1.Group("/sensitive")
sensitiveGroup.Use(middlewares.RateLimitMiddleware(5)) // 每分钟5次
{
    sensitiveGroup.POST("/password/reset", passwordCtrl.Reset)
}
```

## 5. 最佳实践

### 密码处理

```go
// ❌ 不要明文存储密码
user.Password = request.Password

// ✅ 使用 bcrypt 哈希
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
user.Password = string(hashedPassword)

// ✅ 验证密码
err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
```

### 敏感数据

```go
// ✅ 从响应中排除敏感字段
type UserResponse struct {
    ID        uint64 `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    // Password 不包含在响应中
}

// ✅ 日志中过滤敏感信息
logger.Info("User logged in", 
    zap.String("user_id", userID),
    // 不记录密码、token 等敏感信息
)
```

### CSRF 防护

```go
// 对于状态改变的操作（POST/PUT/DELETE），验证 CSRF Token
func CSRFProtection() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 验证 CSRF token
        token := c.GetHeader("X-CSRF-Token")
        if !validateCSRFToken(token) {
            c.JSON(403, gin.H{"error": "Invalid CSRF token"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### HTTPS 配置

生产环境强制使用 HTTPS：

```go
// 在 main.go 中
if config.Get("app.env") == "production" {
    router.Use(func(c *gin.Context) {
        if c.Request.Header.Get("X-Forwarded-Proto") != "https" {
            c.Redirect(301, "https://"+c.Request.Host+c.Request.RequestURI)
            c.Abort()
            return
        }
        c.Next()
    })
}
```

## 6. 安全检查清单

### 代码审查

- [ ] 所有用户输入都经过验证
- [ ] 密码使用 bcrypt 哈希
- [ ] 敏感数据不出现在日志中
- [ ] 使用参数化查询（GORM）
- [ ] API 端点有适当的限流
- [ ] 认证端点使用 JWT
- [ ] 响应中不包含敏感信息

### 配置检查

- [ ] 生产环境启用 HTTPS
- [ ] 数据库密码足够复杂
- [ ] JWT secret 足够长且随机
- [ ] Redis 启用密码认证
- [ ] 限流阈值合理配置

### 测试验证

```bash
# 运行安全测试
go test ./pkg/security -v

# 测试限流
for i in {1..100}; do curl http://localhost:3000/api/test; done

# 测试 SQL 注入
curl "http://localhost:3000/api/users?id=1' OR '1'='1"

# 测试 XSS
curl "http://localhost:3000/api/search?q=<script>alert('xss')</script>"
```

## 7. 监控和告警

### 安全事件日志

```go
// 记录安全事件
logger.Warn("Security threat detected",
    zap.String("type", result.RiskType),
    zap.String("ip", c.ClientIP()),
    zap.String("path", c.Request.URL.Path),
    zap.String("input", suspiciousInput),
)
```

### 监控指标

- SQL 注入尝试次数
- XSS 攻击尝试次数
- 限流触发次数
- 异常登录尝试

## 8. 应急响应

### 发现攻击时

1. **立即封禁 IP**
   ```go
   // 在限流器中手动封禁
   limiter.BlockIP(attackerIP, 24*time.Hour)
   ```

2. **审查日志**
   ```bash
   # 查找攻击者的所有请求
   grep "攻击者IP" storage/logs/*.log
   ```

3. **评估损失**
   - 检查数据完整性
   - 审查数据库变更
   - 确认没有数据泄露

4. **加固措施**
   - 更新安全规则
   - 调整限流阈值
   - 更新依赖包

## 9. 依赖安全

### 定期更新

```bash
# 检查过期依赖
go list -u -m all

# 更新依赖
go get -u ./...

# 安全审计
go mod verify
```

### 漏洞扫描

```bash
# 使用 gosec 扫描
gosec ./...

# 使用 nancy 检查依赖漏洞
nancy sleuth
```

## 10. 合规性

### GDPR / 数据保护

- 用户数据加密存储
- 提供数据导出功能
- 提供数据删除功能
- 记录数据访问日志

### 审计日志

- 记录所有认证尝试
- 记录敏感操作
- 保留日志 90 天

## 总结

GoHub-Service v5.0 实施了多层安全防护：

1. **输入层**：全面的输入验证和清理
2. **传输层**：HTTPS 加密
3. **应用层**：认证、授权、限流
4. **数据层**：参数化查询、加密存储
5. **监控层**：安全事件日志和告警

定期审查和更新安全措施，保持系统安全性。
