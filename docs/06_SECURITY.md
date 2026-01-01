# 🔐 安全防护指南

**最后更新**: 2026年1月1日 | **版本**: v2.0 | **安全级别**: 高

---

## 📖 目录

1. [安全概述](#安全概述)
2. [认证与授权](#认证与授权)
3. [数据保护](#数据保护)
4. [内容安全](#内容安全)
5. [API安全](#api安全)
6. [基础设施安全](#基础设施安全)
7. [审计与监控](#审计与监控)
8. [应急响应](#应急响应)

---

## 🎯 安全概述

### 安全原则

本项目遵循以下核心安全原则：

| 原则 | 说明 | 实现 |
|------|------|------|
| **最小权限** | 只授予必要权限 | RBAC系统严格控制 |
| **纵深防御** | 多层防护机制 | 认证→授权→审计 |
| **加密优先** | 敏感数据必须加密 | AES-256, bcrypt等 |
| **审计日志** | 记录所有重要操作 | 完整的操作日志 |
| **快速响应** | 快速处理安全事件 | 告警+应急预案 |

### 威胁模型

```
┌──────────────────────────────────────┐
│        面临的主要威胁                  │
├──────────────────────────────────────┤
│ 1. 身份认证攻击                       │
│    - 暴力破解                        │
│    - Token伪造                       │
│    → 防护: JWT + 限流 + 日志审计      │
│                                      │
│ 2. 权限提升                          │
│    - 越权访问                        │
│    - 角色提升                        │
│    → 防护: RBAC + 权限检查 + 审计     │
│                                      │
│ 3. 数据泄露                          │
│    - SQL注入                         │
│    - 敏感信息泄露                    │
│    → 防护: 参数化查询 + 加密 + 日志  │
│                                      │
│ 4. 内容安全                          │
│    - XSS攻击                         │
│    - 敏感词违规                      │
│    → 防护: 内容过滤 + 敏感词库      │
│                                      │
│ 5. 拒绝服务                          │
│    - 速率限制                        │
│    - 资源枯竭                        │
│    → 防护: 限流 + 缓存 + 监控       │
└──────────────────────────────────────┘
```

---

## 🔑 认证与授权

### JWT认证机制

#### Token结构

```go
// Token包含以下声明
type Claims struct {
    UserID   int64    `json:"user_id"`
    Email    string   `json:"email"`
    Name     string   `json:"name"`
    Roles    []string `json:"roles"`      // 用户角色
    Standard jwt.StandardClaims
}

// Token格式: header.payload.signature
```

#### Token生成

```go
// 配置
const (
    TokenExpiration = 24 * time.Hour      // 24小时有效期
    RefreshWindow   = 1 * time.Hour       // 剩余1小时时刷新
)

// 生成Token
func GenerateToken(user *User) (string, error) {
    claims := &Claims{
        UserID: user.ID,
        Email:  user.Email,
        Name:   user.Name,
        Roles:  user.GetRoles(),
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(TokenExpiration).Unix(),
            IssuedAt:  time.Now().Unix(),
            NotBefore: time.Now().Unix(),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
```

#### Token验证

```go
// 中间件验证
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 从header获取token
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            response.Error(c, "缺少认证信息")
            c.Abort()
            return
        }
        
        // 解析Bearer token
        token := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := ValidateToken(token)
        if err != nil {
            response.Error(c, "Token无效或已过期")
            c.Abort()
            return
        }
        
        // 保存用户信息
        c.Set("user_id", claims.UserID)
        c.Set("user_roles", claims.Roles)
        c.Next()
    }
}

// 验证Token
func ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }
    
    return claims, nil
}
```

### RBAC权限系统

#### 权限模型

```go
// 三层权限结构
type Role struct {
    ID          int64
    Name        string              // admin, user, moderator
    Permissions []*Permission       // 关联权限
}

type Permission struct {
    ID          int64
    Name        string              // topics.create, comments.delete
    Description string
}

type UserRole struct {
    ID     int64
    UserID int64
    RoleID int64
}
```

#### 权限检查

```go
// 方式1: 注解式权限检查
func (uc *TopicController) Create(c *gin.Context) {
    // 检查权限
    if !uc.hasPermission(c, "topics.create") {
        response.Error(c, "没有创建话题的权限")
        return
    }
    
    // 创建逻辑
}

// 方式2: 策略类检查
func (p *TopicPolicy) Create(c *gin.Context, topic *Topic) bool {
    userID := c.GetInt64("user_id")
    
    // 管理员可以创建任何话题
    if p.isAdmin(c) {
        return true
    }
    
    // 普通用户只能创建自己的话题
    return topic.UserID == userID
}

// 方式3: 权限中间件
func CanPerform(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        roles := c.GetStringSlice("user_roles")
        
        // 检查角色是否有权限
        if !hasPermission(roles, permission) {
            response.Error(c, "权限不足")
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

#### 常用权限列表

```
用户权限:
├── users.create          - 注册账户
├── users.read            - 查看个人信息
├── users.update          - 修改个人信息
└── users.delete          - 删除账户

话题权限:
├── topics.create         - 创建话题
├── topics.read           - 查看话题
├── topics.update         - 编辑话题
└── topics.delete         - 删除话题

评论权限:
├── comments.create       - 创建评论
├── comments.read         - 查看评论
├── comments.update       - 编辑评论
└── comments.delete       - 删除评论

管理权限:
├── admin.users           - 管理用户
├── admin.topics          - 管理话题
├── admin.comments        - 管理评论
├── admin.roles           - 管理角色
└── admin.permissions     - 管理权限
```

---

## 🛡️ 数据保护

### 密码安全

```go
// bcrypt加密密码
import "golang.org/x/crypto/bcrypt"

// 用户注册时
hashedPassword, err := bcrypt.GenerateFromPassword(
    []byte(password), 
    bcrypt.DefaultCost,  // 成本因子
)
user.Password = string(hashedPassword)

// 登录验证
err := bcrypt.CompareHashAndPassword(
    []byte(user.Password),
    []byte(loginPassword),
)
if err != nil {
    // 密码不匹配
}

// 密码策略
const (
    MinPasswordLength = 8
    MaxPasswordLength = 128
)

// 验证密码强度
func ValidatePasswordStrength(password string) error {
    if len(password) < MinPasswordLength {
        return errors.New("密码长度不能少于8位")
    }
    
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasDigit := regexp.MustCompile(`\d`).MatchString(password)
    
    if !hasUpper || !hasLower || !hasDigit {
        return errors.New("密码必须包含大小写字母和数字")
    }
    
    return nil
}
```

### 敏感数据加密

```go
import "crypto/aes"

// 配置加密密钥
const (
    EncryptionKey = "your-32-byte-key" // 32字节(256位)
)

// 加密敏感字段
func EncryptSensitive(data string) (string, error) {
    cipher, err := aes.NewCipher([]byte(EncryptionKey))
    if err != nil {
        return "", err
    }
    
    // 使用GCM模式
    gcm, err := cipher.NewGCM()
    if err != nil {
        return "", err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }
    
    ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// 解密敏感字段
func DecryptSensitive(encrypted string) (string, error) {
    ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
    if err != nil {
        return "", err
    }
    
    cipher, err := aes.NewCipher([]byte(EncryptionKey))
    if err != nil {
        return "", err
    }
    
    gcm, err := cipher.NewGCM()
    if err != nil {
        return "", err
    }
    
    nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    
    return string(plaintext), err
}
```

### 数据库安全

```go
// ✅ 使用参数化查询防止SQL注入
user, err := repo.GetByEmail(ctx, email)  // ✅ 正确
// 内部实现:
db.Where("email = ?", email).First(&user)

// ❌ 避免字符串拼接
// db.Where(fmt.Sprintf("email = '%s'", email)).First(&user)  // 危险！

// 限制查询权限
// 数据库用户应只有必要的权限
// - API用户: SELECT, INSERT, UPDATE(自己数据)
// - 备份用户: SELECT only
// - 迁移用户: DDL权限
```

---

## 🚫 内容安全

### XSS防护

```go
import (
    "github.com/microcosm-cc/bluemonday"
    "html"
)

// 配置XSS清理策略
var strictPolicy = bluemonday.StrictPolicy()
var uGCPolicy = bluemonday.UGCPolicy()

// 方式1: 严格清理（仅允许文本）
func SanitizeStrict(input string) string {
    return strictPolicy.Sanitize(input)
}

// 方式2: UGC清理（允许安全的HTML标签）
func SanitizeUGC(input string) string {
    return uGCPolicy.Sanitize(input)
}

// 方式3: HTML转义（显示用户输入的HTML）
func EscapeHTML(input string) string {
    return html.EscapeString(input)
}

// 在创建话题/评论时
func (s *TopicService) Create(ctx context.Context, req *CreateTopicRequest) error {
    // 清理用户输入
    req.Title = SanitizeStrict(req.Title)
    req.Body = SanitizeUGC(req.Body)
    
    // 继续创建逻辑...
}
```

### CSRF防护

```go
import "github.com/ugorji/go/codec"

// 使用CSRF Token
func CSRFMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // GET请求生成Token
        if c.Request.Method == "GET" {
            token := generateCSRFToken()
            c.SetCookie("csrf_token", token, 3600, "/", "", false, true)
            c.Set("csrf_token", token)
            c.Next()
            return
        }
        
        // POST/PUT/DELETE请求验证Token
        if c.Request.Method != "GET" {
            token := c.PostForm("csrf_token")
            if token == "" {
                token = c.GetHeader("X-CSRF-Token")
            }
            
            cookieToken, _ := c.Cookie("csrf_token")
            if token != cookieToken {
                response.Error(c, "CSRF验证失败")
                c.Abort()
                return
            }
        }
        
        c.Next()
    }
}
```

### 敏感词过滤

```go
import "github.com/importantimport/senswords"

// 初始化敏感词库
var sensitiveWords = senswords.New()

// 从文件加载敏感词库
func LoadSensitiveWords(filePath string) error {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return err
    }
    
    words := strings.Split(string(data), "\n")
    sensitiveWords.Add(words...)
    return nil
}

// 检查和过滤敏感词
func CheckSensitiveWords(text string) (bool, error) {
    // 检查是否包含敏感词
    if sensitiveWords.Contains(text) {
        return true, nil
    }
    return false, nil
}

func FilterSensitiveWords(text string) string {
    // 替换敏感词为***
    return sensitiveWords.Replace(text, '*')
}

// 在创建话题/评论时
func (s *TopicService) Create(ctx context.Context, req *CreateTopicRequest) error {
    // 检查敏感词
    if hasSensitiveWords, _ := CheckSensitiveWords(req.Body); hasSensitiveWords {
        return errors.New("内容包含违禁词，请修改")
    }
    
    // 创建逻辑...
}
```

---

## 🔒 API安全

### 速率限制

```go
import "github.com/gin-contrib/ratelimit"

// 配置限流
func RateLimitMiddleware() gin.HandlerFunc {
    return ratelimit.NewRateLimiter(
        ratelimit.FixedWindowLimiter(
            100,               // 最多100个请求
            time.Minute,       // 在1分钟内
        ),
    )
}

// 特定端点限流
func LoginRateLimit() gin.HandlerFunc {
    return ratelimit.NewRateLimiter(
        ratelimit.FixedWindowLimiter(
            5,                 // 最多5次尝试
            time.Minute,       // 在1分钟内
        ),
    )
}

// 路由应用
func SetupRoutes(r *gin.Engine) {
    // 全局限流
    r.Use(RateLimitMiddleware())
    
    // 登录端点特殊限流
    r.POST("/api/auth/login", LoginRateLimit(), handler.Login)
}
```

### 跨域安全

```go
import "github.com/gin-contrib/cors"

// CORS配置
config := cors.Config{
    AllowOrigins:     []string{"https://example.com", "https://app.example.com"},  // 明确的源
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Content-Type", "Authorization"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}

r.Use(cors.New(config))

// ❌ 避免
cors.Config{
    AllowAllOrigins: true,  // 危险！允许所有源
}
```

### API密钥安全

```go
// API密钥验证
func APIKeyAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        apiKey := c.GetHeader("X-API-Key")
        if apiKey == "" {
            response.Error(c, "缺少API密钥")
            c.Abort()
            return
        }
        
        // 验证API密钥
        app, err := appService.ValidateAPIKey(c, apiKey)
        if err != nil {
            response.Error(c, "无效的API密钥")
            c.Abort()
            return
        }
        
        c.Set("app_id", app.ID)
        c.Next()
    }
}

// 路由应用
api := r.Group("/api/v1")
api.Use(APIKeyAuth())
api.GET("/data", handler.GetData)
```

---

## 🏗️ 基础设施安全

### HTTPS/TLS

```go
// 强制HTTPS
func HTTPSRedirect() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Request.Header.Get("X-Forwarded-Proto") != "https" {
            c.Redirect(http.StatusMovedPermanently, 
                fmt.Sprintf("https://%s%s", c.Request.Host, c.Request.RequestURI))
            return
        }
        c.Next()
    }
}

// 安全的Header设置
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 防止点击劫持
        c.Header("X-Frame-Options", "DENY")
        
        // 防止MIME嗅探
        c.Header("X-Content-Type-Options", "nosniff")
        
        // 启用XSS防护
        c.Header("X-XSS-Protection", "1; mode=block")
        
        // 内容安全策略
        c.Header("Content-Security-Policy", "default-src 'self'")
        
        // 强制HSTS
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        
        c.Next()
    }
}
```

### 依赖安全

```bash
# 检查依赖漏洞
go list -json -m all | nancy sleuth

# 更新依赖
go get -u ./...

# 审计依赖
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

---

## 📊 审计与监控

### 操作日志

```go
// 记录所有关键操作
type AuditLog struct {
    ID        int64
    UserID    int64
    Action    string           // create, update, delete
    Resource  string           // topic, comment
    ResourceID int64
    Changes   map[string]interface{}  // 变更内容
    IPAddress string
    UserAgent string
    Status    string           // success, failed
    CreatedAt time.Time
}

// 记录操作
func (s *AuditService) Log(ctx context.Context, log *AuditLog) error {
    log.UserID = ctx.Value("user_id").(int64)
    log.IPAddress = c.ClientIP()
    log.UserAgent = c.Request.UserAgent()
    log.CreatedAt = time.Now()
    
    return s.repo.Create(ctx, log)
}

// 中间件自动记录
func AuditMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        startTime := time.Now()
        
        c.Next()
        
        // 只记录修改操作
        if c.Request.Method != "GET" {
            auditService.Log(c, &AuditLog{
                Action:    c.Request.Method,
                Resource:  c.Request.URL.Path,
                Status:    strconv.Itoa(c.Writer.Status()),
            })
        }
    }
}
```

### 安全监控

```go
// 监控异常行为
type SecurityAlert struct {
    Type      string    // login_failure, permission_denied, rate_limit
    UserID    int64
    Message   string
    Severity  string    // info, warning, critical
    CreatedAt time.Time
}

// 监控登录失败
func (s *AuthService) Login(email, password string) error {
    attempts, _ := cache.Get(fmt.Sprintf("login_attempt:%s", email))
    
    if attempts > 5 {
        // 触发告警
        alertService.Send(&SecurityAlert{
            Type:     "login_failure",
            Message:  fmt.Sprintf("用户 %s 多次登录失败", email),
            Severity: "warning",
        })
        return errors.New("账户已锁定，请稍后重试")
    }
    
    // 验证登录...
}

// 监控权限拒绝
func (s *AuthService) CheckPermission(userID int64, permission string) bool {
    if !hasPermission(userID, permission) {
        alertService.Send(&SecurityAlert{
            Type:     "permission_denied",
            UserID:   userID,
            Message:  fmt.Sprintf("用户尝试访问无权限的资源: %s", permission),
            Severity: "warning",
        })
        return false
    }
    return true
}
```

---

## 🚨 应急响应

### 安全事件处理流程

```
检测事件
  ↓
评估严重程度
  ↓
  ├─ 低: 记录日志
  ├─ 中: 通知团队
  └─ 高: 立即升级
  ↓
隔离影响
  ├─ 禁用账户
  ├─ 撤销Token
  └─ 备份数据
  ↓
修复问题
  ├─ 代码修复
  ├─ 依赖更新
  └─ 配置调整
  ↓
恢复服务
  ├─ 验证修复
  ├─ 逐步恢复
  └─ 监控状态
  ↓
事后分析
  ├─ 原因分析
  ├─ 流程改进
  └─ 知识共享
```

### 常见安全事件

| 事件类型 | 症状 | 响应 |
|---------|------|------|
| **账户被黑** | 异常登录、未授权操作 | 重置密码、撤销Token、审计日志 |
| **数据泄露** | 数据库被访问、敏感信息公开 | 隔离系统、通知用户、法律咨询 |
| **DDoS攻击** | 大量请求、服务不可用 | 启用CDN、限流、地理限制 |
| **代码漏洞** | 崩溃、不正常行为 | 紧急修复、灰度发布、监控 |
| **依赖漏洞** | 已知漏洞CVE | 紧急更新、安全补丁、扫描 |

---

## ✅ 安全检查清单

部署前完成以下检查：

- [ ] 修改所有默认密码和密钥
- [ ] 启用HTTPS/TLS
- [ ] 配置防火墙规则
- [ ] 设置CORS允许列表
- [ ] 启用日志审计
- [ ] 配置备份和恢复
- [ ] 运行安全扫描工具
- [ ] 进行渗透测试
- [ ] 建立应急响应计划
- [ ] 培训团队安全意识

---

## 📚 相关文档

- [权限系统详解](03_RBAC.md) - RBAC设计和实现
- [生产部署](09_PRODUCTION.md) - 安全的部署配置

---

**安全级别**: 🔒 高  
**最后更新**: 2026年1月1日  
*由GoHub Security Team维护* ✨
