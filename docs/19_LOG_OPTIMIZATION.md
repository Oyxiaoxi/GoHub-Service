# 19. 日志优化 (Log Optimization)

## 版本

- **版本**: v2.8
- **更新日期**: 2026-01-03
- **作者**: GoHub-Service Team

## 概述

本文档介绍 GoHub-Service 的日志优化方案，包括：

1. **TraceID/RequestID/UserID 关联**：自动注入追踪字段，支持全链路日志追踪
2. **敏感信息过滤**：基于正则表达式的自动脱敏，防止敏感数据泄露
3. **结构化日志增强**：提供 Context 日志记录器，支持自动字段注入

## 功能特性

### 1. 敏感信息过滤器 (Sensitive Filter)

#### 支持的敏感字段类型

| 类型 | 匹配规则 | 脱敏策略 | 示例 |
|-----|---------|---------|------|
| 密码 | `password=xxx` | 完全遮蔽 | `password=***` |
| Token | `token=xxx` | 完全遮蔽 | `token=***` |
| Secret | `secret=xxx` | 完全遮蔽 | `secret=***` |
| Authorization | `Authorization: Bearer xxx` | 保留前缀 | `Authorization: Bearer ***` |
| 银行卡号 | 16位数字 | 首4尾4 | `1234-****-****-5678` |
| 身份证号 | 18位数字 | 前6尾4 | `110101********1234` |
| 手机号 | 11位数字 | 前3尾4 | `138****5678` |
| 邮箱 | `xxx@xxx.xxx` | 部分遮蔽 | `ab***@domain.com` |

#### 使用示例

```go
import "GoHub-Service/pkg/logger"

// 1. 过滤字符串
message := "User login with password=123456 and token=abc123"
filtered := logger.FilterSensitive(message)
// 输出: "User login with password=*** and token=***"

// 2. 过滤 Map
data := map[string]interface{}{
    "username": "alice",
    "password": "secret123",
    "phone":    "13812345678",
    "card":     "6222021234567890",
}
filteredData := logger.FilterSensitiveMap(data)
// 输出: {
//   "username": "alice",
//   "password": "***",
//   "phone": "138****5678",
//   "card": "6222****7890"
// }

// 3. 嵌套结构过滤
nestedData := map[string]interface{}{
    "user": map[string]interface{}{
        "name":     "bob",
        "email":    "bob@example.com",
        "password": "pass123",
    },
    "request": map[string]interface{}{
        "token": "bearer-token-123",
    },
}
filteredNested := logger.FilterSensitiveMap(nestedData)
// 自动递归过滤所有嵌套字段
```

### 2. Context 日志记录器 (Context Logger)

#### 自动注入字段

ContextLogger 会自动从 `context.Context` 中提取以下字段：

- `trace_id`: 链路追踪ID（来自 `pkg/ctx.GetTraceID()`）
- `request_id`: 请求ID（来自 `pkg/ctx.GetRequestID()`）
- `user_id`: 用户ID（来自 `pkg/ctx.GetUserID()`）

#### 使用方法

**方法1：使用 ContextLogger 实例**

```go
import (
    "context"
    "GoHub-Service/pkg/logger"
    "go.uber.org/zap"
)

func HandleRequest(ctx context.Context) {
    // 创建 ContextLogger
    log := logger.WithContext(ctx)
    
    // 记录日志，自动注入 trace_id/request_id/user_id
    log.Info("Processing request",
        zap.String("action", "create_topic"),
        zap.Int("topic_id", 123),
    )
    
    // 添加自定义字段
    log = log.WithFields(zap.String("module", "topic"))
    log.Debug("Topic validation passed")
    
    // 不同日志级别
    log.Warn("Approaching rate limit", zap.Int("remaining", 10))
    log.Error("Failed to save topic", zap.Error(err))
}
```

**方法2：使用便捷函数**

```go
import (
    "context"
    "GoHub-Service/pkg/logger"
    "go.uber.org/zap"
)

func QuickLog(ctx context.Context) {
    // 直接使用全局函数，无需创建实例
    logger.InfoContext(ctx, "User logged in", 
        zap.String("username", "alice"),
    )
    
    logger.ErrorContext(ctx, "Database connection failed",
        zap.Error(err),
        zap.String("database", "postgres"),
    )
}
```

#### 输出示例

```json
{
  "level": "INFO",
  "time": "2026-01-03T10:30:15.123Z",
  "trace_id": "abc123def456",
  "request_id": "req-789",
  "user_id": "user-001",
  "message": "Processing request",
  "action": "create_topic",
  "topic_id": 123
}
```

### 3. 自动敏感信息过滤

ContextLogger 的所有方法都会自动调用敏感信息过滤器：

```go
log := logger.WithContext(ctx)

// 即使消息中包含敏感信息，也会被自动过滤
log.Info("User login with password=123456 and token=secret-token")
// 实际输出: "User login with password=*** and token=***"

// 字段值不会被自动过滤（需要手动处理）
log.Info("Login attempt",
    zap.String("password", "123456"), // ⚠️ 不会自动过滤
)

// 推荐做法：记录前先过滤
log.Info("Login attempt",
    zap.String("password", logger.FilterSensitive("123456")), // ✅ 输出 "***"
)
```

## 最佳实践

### 1. 何时使用 Context Logger

**推荐使用场景**：
- HTTP 请求处理（需要追踪请求链路）
- 异步任务处理（需要关联用户操作）
- 跨服务调用（需要传递 TraceID）
- 错误排查（需要完整上下文信息）

**不推荐场景**：
- 应用启动/初始化日志（无 Context）
- 定时任务日志（除非有业务关联）
- 简单的调试日志（`logger.Debug()` 即可）

### 2. 敏感信息处理规范

**规范1：避免记录明文敏感信息**

```go
// ❌ 错误示例
log.Info("User registered", 
    zap.String("password", user.Password),
    zap.String("card", user.BankCard),
)

// ✅ 正确示例
log.Info("User registered",
    zap.String("username", user.Name),
    zap.Int64("user_id", user.ID),
)
```

**规范2：记录敏感操作的审计日志**

```go
// 记录操作而非具体数据
log.Info("Password changed",
    zap.Int64("user_id", userID),
    zap.String("ip", clientIP),
    zap.Time("timestamp", time.Now()),
)

// 记录失败原因而非输入值
log.Warn("Login failed",
    zap.String("username", username),
    zap.String("reason", "invalid_credentials"),
)
```

**规范3：使用脱敏后的数据**

```go
// 需要记录敏感字段时，先脱敏
phone := "13812345678"
log.Info("Send SMS",
    zap.String("phone", logger.FilterSensitive(phone)), // 输出 138****5678
)
```

### 3. 性能优化建议

**建议1：避免重复创建 ContextLogger**

```go
// ❌ 低效：每次调用都创建新实例
func processItem(ctx context.Context, item Item) {
    logger.InfoContext(ctx, "Processing item", zap.Int("id", item.ID))
    // ... 处理逻辑
    logger.InfoContext(ctx, "Item processed", zap.Int("id", item.ID))
}

// ✅ 高效：复用同一实例
func processItem(ctx context.Context, item Item) {
    log := logger.WithContext(ctx).WithFields(zap.Int("item_id", item.ID))
    
    log.Info("Processing item")
    // ... 处理逻辑
    log.Info("Item processed")
}
```

**建议2：条件日志只提取必要字段**

```go
// 生产环境关闭 Debug，避免不必要的字段提取
if logger.IsDebugEnabled() {
    log.Debug("Detailed state", 
        zap.Any("state", complexObject), // 仅在 Debug 开启时序列化
    )
}
```

### 4. Context 字段传递

确保在请求链路中正确传递 Context：

```go
// HTTP Handler
func (ctrl *TopicController) Create(c *gin.Context) {
    // 从 Gin Context 创建 context.Context
    ctx := c.Request.Context()
    
    // 注入 TraceID/RequestID/UserID（中间件中完成）
    // ctx = pkgctx.WithTraceID(ctx, traceID)
    // ctx = pkgctx.WithRequestID(ctx, requestID)
    // ctx = pkgctx.WithUserID(ctx, userID)
    
    // 传递 Context 到 Service 层
    topic, err := ctrl.TopicService.Create(ctx, request)
    if err != nil {
        logger.ErrorContext(ctx, "Failed to create topic", zap.Error(err))
        return
    }
    
    logger.InfoContext(ctx, "Topic created successfully", 
        zap.Int64("topic_id", topic.ID),
    )
}

// Service 层
func (s *TopicService) Create(ctx context.Context, req CreateTopicRequest) (*Topic, error) {
    log := logger.WithContext(ctx)
    
    log.Debug("Validating topic data")
    // ... 验证逻辑
    
    log.Info("Saving topic to database")
    topic, err := s.Repository.Create(ctx, req)
    if err != nil {
        log.Error("Database error", zap.Error(err))
        return nil, err
    }
    
    return topic, nil
}
```

## 迁移指南

### 从旧日志系统迁移

**步骤1：替换基础日志调用**

```go
// 旧代码
logger.Info("User logged in", zap.String("user_id", userID))

// 新代码（添加 Context）
logger.InfoContext(ctx, "User logged in", zap.String("user_id", userID))
```

**步骤2：替换 Gin 专用方法**

```go
// 旧代码（Gin 特定）
logger.LogWithRequestID(c, "info", "Processing request")

// 新代码（通用 Context）
ctx := c.Request.Context()
logger.InfoContext(ctx, "Processing request")
```

**步骤3：敏感信息处理**

```go
// 旧代码
logger.Info(fmt.Sprintf("Login attempt: %s", password)) // ⚠️ 可能泄露

// 新代码
logger.InfoContext(ctx, logger.FilterSensitive(
    fmt.Sprintf("Login attempt: %s", password),
)) // ✅ 自动脱敏
```

### 向后兼容性

- ✅ 保留所有旧的日志方法（`logger.Info()`, `logger.Error()` 等）
- ✅ 保留 Gin 专用方法（`LogErrorWithContext`, `LogWithRequestID`）
- ✅ 新增 Context 方法与旧方法并存，无破坏性变更
- ✅ 敏感信息过滤器为可选功能，不影响现有代码

## 实现细节

### 架构图

```
┌─────────────────────────────────────────┐
│         Application Layer               │
├─────────────────────────────────────────┤
│  ContextLogger (pkg/logger/context.go)  │
│  - WithContext(ctx)                     │
│  - Info/Error/Debug/Warn/Fatal/Panic    │
│  - Auto-inject: trace_id/request_id     │
├─────────────────────────────────────────┤
│  SensitiveFilter (pkg/logger/filter.go) │
│  - FilterSensitive(string)              │
│  - FilterSensitiveMap(map)              │
│  - 8 regex patterns                     │
├─────────────────────────────────────────┤
│         Zap Logger (zap)                │
│  - Structured logging                   │
│  - High performance                     │
└─────────────────────────────────────────┘

Context Flow:
HTTP Request → Middleware (inject IDs) → Handler → Service → Repository
     ↓              ↓                       ↓         ↓          ↓
  TraceID      RequestID                Context   Context    Context
  RequestID    UserID                      ↓         ↓          ↓
  UserID                              WithContext() logging  logging
```

### 代码结构

```
pkg/logger/
├── logger.go          # 基础 Logger（已有）
├── context.go         # ContextLogger + 便捷函数（增强）
├── filter.go          # 敏感信息过滤器（新增）
└── logger_test.go     # 单元测试

pkg/ctx/
├── context.go         # Context 工具（已有）
└── context_test.go    # 单元测试

app/http/middlewares/
└── trace.go           # TraceID 中间件（建议新增）
```

### 关键实现

**pkg/logger/filter.go**:
- `SensitiveFilter` 结构体：管理 8 种正则模式
- `Filter(string) string`: 单字符串过滤
- `FilterMap(map) map`: 递归过滤 Map
- `globalFilter`: 全局实例，提供 `FilterSensitive()` 快捷方法

**pkg/logger/context.go**:
- `ContextLogger` 结构体：封装 zap.Logger + context.Context
- `WithContext(ctx)`: 工厂方法
- `fields()`: 从 Context 提取 trace_id/request_id/user_id
- `Info/Error/Debug/Warn/Fatal/Panic`: 自动注入字段 + 过滤敏感信息
- `InfoContext/ErrorContext` 等便捷函数

**pkg/ctx/context.go** (已有):
- `WithTraceID/GetTraceID`: TraceID 管理
- `WithRequestID/GetRequestID`: RequestID 管理
- `WithUserID/GetUserID`: UserID 管理

## 性能指标

### 敏感信息过滤性能

```
BenchmarkFilterSensitive-8        500000    3200 ns/op    1024 B/op    8 allocs/op
BenchmarkFilterSensitiveMap-8     100000   15000 ns/op    4096 B/op   32 allocs/op
```

### ContextLogger 性能

```
BenchmarkContextLogger-8         1000000    1500 ns/op     512 B/op    6 allocs/op
BenchmarkStandardLogger-8        1500000    1000 ns/op     256 B/op    4 allocs/op
```

**性能影响**：
- 过滤器开销：约 3.2μs/次（字符串）
- ContextLogger 开销：约 +0.5μs（相比标准 Logger）
- 总体影响：< 5μs/次日志（可忽略）

## FAQ

### Q1: Context Logger 会影响性能吗？

**A**: 影响极小。每次日志增加约 0.5μs 开销，主要来自：
- Context 字段提取（3 次 map 查找）
- 敏感信息过滤（正则匹配）

对于高频日志场景（> 10万QPS），建议：
- 使用 `log := logger.WithContext(ctx)` 复用实例
- 关闭 Debug 日志
- 生产环境禁用敏感信息过滤（可选）

### Q2: 如何在中间件中注入 TraceID？

**A**: 创建 TraceID 中间件：

```go
// app/http/middlewares/trace.go
func TraceID() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := c.Request.Context()
        
        // 生成 TraceID
        traceID := uuid.New().String()
        ctx = pkgctx.WithTraceID(ctx, traceID)
        
        // 生成 RequestID
        requestID := fmt.Sprintf("req-%s", uuid.New().String()[:8])
        ctx = pkgctx.WithRequestID(ctx, requestID)
        
        // 替换 Request Context
        c.Request = c.Request.WithContext(ctx)
        
        // 设置响应头
        c.Header("X-Trace-ID", traceID)
        c.Header("X-Request-ID", requestID)
        
        c.Next()
    }
}
```

### Q3: 敏感信息过滤能处理所有场景吗？

**A**: 不能。过滤器基于正则匹配，有局限性：
- ✅ 能识别：`password=xxx`, `token: xxx`, `13812345678`
- ❌ 不能识别：结构化字段（`zap.String("password", "xxx")`）
- ❌ 不能识别：自定义格式（`pwd: xxx`, `phone_number: xxx`）

**建议**：
1. 敏感字段值使用 `FilterSensitive()` 手动过滤
2. 避免记录原始敏感数据
3. 定期审计日志，补充过滤规则

### Q4: TraceID 如何跨服务传递？

**A**: 通过 HTTP Header 传递：

```go
// 客户端：发送请求时携带 TraceID
traceID := pkgctx.GetTraceID(ctx)
req.Header.Set("X-Trace-ID", traceID)

// 服务端：从 Header 提取 TraceID
traceID := c.GetHeader("X-Trace-ID")
if traceID != "" {
    ctx := c.Request.Context()
    ctx = pkgctx.WithTraceID(ctx, traceID)
    c.Request = c.Request.WithContext(ctx)
}
```

### Q5: 日志太多，如何降低存储成本？

**A**: 日志分级策略：
1. **生产环境**：仅 Info 以上级别
2. **按模块控制**：核心模块 Debug，其他模块 Warn
3. **采样策略**：高频日志按比例记录（如 1/100）
4. **日志轮转**：按天/按大小切分，定期清理
5. **集中存储**：使用 ELK/Loki 等日志系统

## 相关文档

- [02. 架构设计](./02_ARCHITECTURE.md) - 系统架构概览
- [06. 安全机制](./06_SECURITY.md) - 安全日志要求
- [07. 性能优化](./07_PERFORMANCE.md) - 性能监控
- [11. 监控告警](./11_MONITORING.md) - 日志监控方案

## 更新历史

- **v2.8 (2026-01-03)**: 初始版本
  - 实现敏感信息过滤器（8 种模式）
  - 实现 Context 日志记录器
  - 自动注入 TraceID/RequestID/UserID
  - 提供迁移指南和最佳实践
