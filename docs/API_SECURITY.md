# API 安全加固文档

> 创建时间：2025年12月29日  
> 最后更新：2025年12月29日 v1.0  
> 状态：已完成基础安全加固

---

## 📋 目录

- [概述](#概述)
- [已实现的安全措施](#已实现的安全措施)
- [CORS 跨域配置](#cors-跨域配置)
- [安全响应头](#安全响应头)
- [XSS 防护](#xss-防护)
- [SQL 注入防护](#sql-注入防护)
- [限流增强](#限流增强)
- [使用指南](#使用指南)
- [生产环境配置建议](#生产环境配置建议)
- [安全检查清单](#安全检查清单)

---

## 概述

本文档详细说明 GoHub-Service 项目的 API 安全加固措施，包括跨域资源共享(CORS)、XSS 防护、SQL 注入防护、安全响应头和增强的限流机制。

### 安全架构图

```
请求流程：
客户端请求
    ↓
CORS 验证 (middlewares.CORS)
    ↓
安全响应头 (middlewares.SecureHeaders)
    ↓
XSS 防护 (middlewares.XSSProtection)
    ↓
限流检查 (middlewares.LimitIP/LimitPerRoute)
    ↓
认证授权 (middlewares.AuthJWT)
    ↓
业务逻辑处理
    ↓
响应（带安全头）
```

---

## 已实现的安全措施

### ✅ 已完成

1. **CORS 跨域配置**
   - 细粒度源控制
   - 方法白名单
   - 请求头/响应头控制
   - 预检请求缓存

2. **安全响应头**
   - X-Frame-Options (防点击劫持)
   - X-Content-Type-Options (防 MIME 嗅探)
   - X-XSS-Protection (XSS 防护)
   - Content-Security-Policy (CSP)
   - Referrer-Policy
   - Permissions-Policy

3. **XSS 防护**
   - HTML 实体转义
   - 脚本标签过滤
   - 事件处理器清理
   - JavaScript 协议过滤

4. **SQL 注入防护**
   - GORM 参数化查询（内置）
   - 关键词模式检测（额外保护）
   - 查询参数验证

5. **限流增强**
   - IP 限流
   - 路由限流
   - 用户限流
   - 可配置限流策略
   - 速率限制响应头

6. **Content-Type 验证**
   - 请求类型白名单
   - 防止 MIME 混淆攻击

---

## CORS 跨域配置

### 中间件文件

**位置**: `app/http/middlewares/cors.go`

### 三种 CORS 策略

#### 1. 标准 CORS 配置

```go
router.Use(middlewares.CORS())
```

**特性**：
- 允许指定源列表（开发环境默认 localhost）
- 支持常用 HTTP 方法
- 允许携带 Cookie (AllowCredentials: true)
- 预检请求缓存 12 小时

**配置详情**：
```go
AllowOrigins: []string{
    "http://localhost:3000",
    "http://localhost:3001",
    "http://localhost:8080",
    // 生产环境需配置具体域名
}
AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
AllowHeaders: []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"}
AllowCredentials: true
MaxAge: 12 * time.Hour
```

#### 2. 公开 API CORS 配置

```go
publicRouter.Use(middlewares.CORSPublic())
```

**特性**：
- 允许所有源 (AllowAllOrigins: true)
- 仅允许只读操作 (GET, OPTIONS)
- 适用于完全公开的只读 API

#### 3. 自定义源 CORS 配置

```go
router.Use(middlewares.CORSWithOrigins([]string{
    "https://app.example.com",
    "https://admin.example.com",
}))
```

**特性**：
- 灵活指定允许的源
- 适用于需要特定源配置的场景

### 生产环境配置建议

```go
// 生产环境 CORS 配置示例
AllowOrigins: []string{
    "https://yourdomain.com",
    "https://www.yourdomain.com",
    "https://app.yourdomain.com",
}
```

⚠️ **重要**：生产环境必须移除 localhost 和 127.0.0.1

---

## 安全响应头

### 中间件文件

**位置**: `app/http/middlewares/security.go`

### 启用方式

```go
router.Use(middlewares.SecureHeaders())
```

### 响应头详解

| 响应头 | 值 | 作用 |
|-------|-----|------|
| X-Frame-Options | DENY | 防止页面被嵌入 iframe，防止点击劫持 |
| X-Content-Type-Options | nosniff | 防止浏览器 MIME 类型嗅探 |
| X-XSS-Protection | 1; mode=block | 启用浏览器 XSS 过滤器 |
| Content-Security-Policy | default-src 'self' | 限制资源加载源，防止 XSS |
| Referrer-Policy | strict-origin-when-cross-origin | 控制 Referrer 信息泄露 |
| Permissions-Policy | geolocation=(), microphone=(), camera=() | 禁用敏感浏览器 API |

### HTTPS 强制（生产环境）

```go
// 取消注释以启用 HSTS（仅 HTTPS 环境）
c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
```

---

## XSS 防护

### 中间件文件

**位置**: `app/http/middlewares/security.go`

### 启用方式

```go
router.Use(middlewares.XSSProtection())
```

### 防护策略

#### 1. HTML 实体转义

```go
input = html.EscapeString(input)
// "<script>" → "&lt;script&gt;"
```

#### 2. 脚本标签过滤

```go
scriptPattern := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
input = scriptPattern.ReplaceAllString(input, "")
```

#### 3. 事件处理器清理

```go
eventPattern := regexp.MustCompile(`(?i)on\w+\s*=`)
input = eventPattern.ReplaceAllString(input, "")
// onclick= → 移除
```

#### 4. JavaScript 协议过滤

```go
input = strings.ReplaceAll(input, "javascript:", "")
// href="javascript:alert(1)" → href="alert(1)"
```

### 作用范围

- ✅ URL 查询参数自动清理
- ⚠️ POST Body 需要在业务层额外处理

### 业务层 XSS 防护示例

```go
import "html"

// 在 Service 层对用户输入进行转义
func (s *TopicService) CreateTopic(dto TopicCreateDTO) error {
    // 转义 HTML 内容
    dto.Title = html.EscapeString(dto.Title)
    dto.Body = html.EscapeString(dto.Body)
    
    // ... 业务逻辑
}
```

---

## SQL 注入防护

### 中间件文件

**位置**: `app/http/middlewares/security.go`

### 多层防护策略

#### 1. GORM 参数化查询（主要防护）

GORM 默认使用参数化查询，有效防止 SQL 注入：

```go
// ✅ 安全：参数化查询
db.Where("name = ?", userInput).First(&user)

// ❌ 危险：字符串拼接（避免使用）
db.Where(fmt.Sprintf("name = '%s'", userInput)).First(&user)
```

#### 2. 中间件关键词检测（额外保护）

```go
router.Use(middlewares.SQLInjectionProtection())
```

**检测模式**：
```go
sqlPattern := regexp.MustCompile(`(?i)(union|select|insert|update|delete|drop|create|alter|exec|script|javascript|<script|</script>)`)
```

**拦截示例**：
```
❌ /api/users?name=admin' OR '1'='1
❌ /api/topics?search=<script>alert(1)</script>
❌ /api/categories?id=1; DROP TABLE users;
```

### 最佳实践

1. **始终使用 GORM 占位符**
   ```go
   db.Where("email = ?", email).First(&user)
   ```

2. **避免原生 SQL**
   ```go
   // 如必须使用，确保参数化
   db.Raw("SELECT * FROM users WHERE id = ?", id).Scan(&user)
   ```

3. **验证输入类型**
   ```go
   id, err := strconv.ParseUint(c.Param("id"), 10, 64)
   if err != nil {
       return errors.New("Invalid ID")
   }
   ```

---

## 限流增强

### 中间件文件

**位置**: `app/http/middlewares/limit.go`

### 三种限流策略

#### 1. IP 全局限流

```go
// 基础使用
router.Use(middlewares.LimitIP("200-H")) // 每小时 200 次

// 带配置使用
router.Use(middlewares.LimitIPWithConfig(middlewares.LimitConfig{
    Rate:          "100-M",
    Message:       "全局访问过于频繁，请稍后再试",
    ShowRemaining: true,
}))
```

#### 2. 路由限流

```go
// 限制特定路由
router.POST("/api/v1/topics", 
    middlewares.LimitPerRoute("10-M"), // 每分钟 10 次
    topicsController.Store,
)

// 带配置使用
router.POST("/api/v1/auth/login",
    middlewares.LimitPerRouteWithConfig(middlewares.LimitConfig{
        Rate:          "5-M",
        Message:       "登录尝试过多，请 5 分钟后再试",
        ShowRemaining: false, // 不显示剩余次数
    }),
    authController.Login,
)
```

#### 3. 用户限流（新增）

```go
// 需要在认证中间件之后使用
router.POST("/api/v1/posts",
    middlewares.AuthJWT(),
    middlewares.LimitByUser("50-H"), // 每用户每小时 50 次
    postsController.Create,
)

// 带配置使用
router.POST("/api/v1/comments",
    middlewares.AuthJWT(),
    middlewares.LimitByUserWithConfig(middlewares.LimitConfig{
        Rate:          "30-H",
        Message:       "您的评论过于频繁，请稍后再试",
        ShowRemaining: true,
    }),
    commentsController.Create,
)
```

### 限流响应头

所有限流中间件自动添加以下响应头：

```http
X-RateLimit-Limit: 100        # 最大请求次数
X-RateLimit-Remaining: 95     # 剩余请求次数
X-RateLimit-Reset: 1735459200 # 重置时间戳
```

### 超限响应

```json
{
    "code": 429,
    "message": "请求过于频繁，请稍后再试",
    "retry_after": 1735459200
}
```

### 限流格式说明

| 格式 | 说明 | 示例 |
|-----|------|------|
| N-S | N 次/秒 | "5-S" = 每秒 5 次 |
| N-M | N 次/分钟 | "10-M" = 每分钟 10 次 |
| N-H | N 次/小时 | "1000-H" = 每小时 1000 次 |
| N-D | N 次/天 | "2000-D" = 每天 2000 次 |

---

## 使用指南

### 完整中间件链配置

#### 全局中间件（bootstrap/route.go）

```go
func registerGlobalMiddleWare(router *gin.Engine) {
    router.Use(
        middlewares.Logger(),               // 日志记录
        middlewares.Recovery(),             // 恢复 panic
        middlewares.CORS(),                 // CORS 跨域
        middlewares.SecureHeaders(),        // 安全响应头
        middlewares.XSSProtection(),        // XSS 防护
        gzip.Gzip(gzip.DefaultCompression), // Gzip 压缩
    )
}
```

#### 路由级别配置示例

```go
// 公开 API（只读，宽松 CORS）
publicAPI := router.Group("/api/v1/public")
publicAPI.Use(middlewares.CORSPublic())
{
    publicAPI.GET("/posts", postsController.Index)
    publicAPI.GET("/posts/:id", postsController.Show)
}

// 认证 API（严格限流）
authAPI := router.Group("/api/v1")
authAPI.Use(
    middlewares.LimitIP("200-H"),          // 全局 IP 限流
    middlewares.AuthJWT(),                  // JWT 认证
)
{
    // 创建操作 - 用户限流
    authAPI.POST("/topics",
        middlewares.LimitByUser("20-H"),   // 每用户每小时 20 次
        topicsController.Store,
    )
    
    // 敏感操作 - 严格限流
    authAPI.POST("/admin/users",
        middlewares.LimitPerRoute("5-M"),  // 每 IP 每分钟 5 次
        adminController.CreateUser,
    )
}
```

### Content-Type 验证（可选）

```go
// 对需要 body 的路由验证 Content-Type
router.POST("/api/v1/topics",
    middlewares.ContentTypeValidation(),
    topicsController.Store,
)
```

---

## 生产环境配置建议

### 1. CORS 配置

```go
// ✅ 生产环境配置
AllowOrigins: []string{
    "https://yourdomain.com",
    "https://www.yourdomain.com",
    "https://app.yourdomain.com",
}

// ❌ 禁止使用
AllowOrigins: []string{"http://localhost:*"}  // 移除所有 localhost
AllowAllOrigins: true                          // 除非完全公开 API
```

### 2. 安全响应头

```go
// 启用 HSTS（强制 HTTPS）
c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

// 更严格的 CSP
c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; font-src 'self'; connect-src 'self'; frame-ancestors 'none'")
```

### 3. 限流配置

```go
// 全局限流（根据实际业务调整）
middlewares.LimitIP("1000-H")     // 每 IP 每小时 1000 次

// 登录接口（防暴力破解）
middlewares.LimitPerRoute("5-M")  // 每 IP 每分钟 5 次

// 创建操作（防滥用）
middlewares.LimitByUser("50-H")   // 每用户每小时 50 次
```

### 4. HTTPS 强制

```go
// Nginx 配置
server {
    listen 80;
    server_name yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    # SSL 配置...
}
```

---

## 安全检查清单

### 部署前检查

- [ ] CORS 配置已更新为生产域名
- [ ] 已移除所有 localhost 源
- [ ] HSTS 头已启用（HTTPS 环境）
- [ ] CSP 策略已配置
- [ ] 限流阈值已根据业务调整
- [ ] XSS 防护已全局启用
- [ ] SQL 注入防护已启用
- [ ] 所有数据库查询使用参数化
- [ ] 敏感操作已添加额外限流
- [ ] Content-Type 验证已启用

### 定期审查

- [ ] 审查 CORS 允许的源列表
- [ ] 检查限流日志，调整阈值
- [ ] 审查安全日志，检测异常
- [ ] 更新依赖包，修复安全漏洞
- [ ] 进行渗透测试
- [ ] 审查用户输入处理逻辑

---

## 测试建议

### 1. CORS 测试

```bash
# 测试预检请求
curl -X OPTIONS http://localhost:3000/api/v1/topics \
  -H "Origin: http://localhost:3001" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" \
  -v

# 验证响应头
Access-Control-Allow-Origin: http://localhost:3001
Access-Control-Allow-Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
Access-Control-Allow-Credentials: true
```

### 2. XSS 防护测试

```bash
# 测试脚本注入
curl "http://localhost:3000/api/v1/topics?search=<script>alert(1)</script>"

# 预期：脚本标签被转义或移除
```

### 3. SQL 注入测试

```bash
# 测试 SQL 注入
curl "http://localhost:3000/api/v1/users?name=admin' OR '1'='1"

# 预期：400 响应，拒绝请求
```

### 4. 限流测试

```bash
# 快速发送多个请求
for i in {1..10}; do
  curl http://localhost:3000/api/v1/topics \
    -H "Authorization: Bearer $TOKEN"
done

# 检查响应头
X-RateLimit-Limit: 200
X-RateLimit-Remaining: 190
X-RateLimit-Reset: 1735459200

# 超限时响应 429
```

---

## 相关文档

- [PERFORMANCE_OPTIMIZATION.md](PERFORMANCE_OPTIMIZATION.md) - 性能优化文档
- [DEVELOPMENT_GUIDE.md](DEVELOPMENT_GUIDE.md) - 开发与配置指南

---

## 版本历史

| 版本 | 日期 | 变更内容 |
|-----|------|---------|
| v1.0 | 2025-12-29 | 初始版本，完成基础安全加固 |

---

**维护者**: GoHub-Service Team  
**最后审核**: 2025-12-29
