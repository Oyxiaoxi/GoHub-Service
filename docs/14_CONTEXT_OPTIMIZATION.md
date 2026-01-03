# 🔄 Context 传递优化指南

**最后更新**: 2026年1月3日 | **版本**: v2.1

> 📚 **文档导航**: [返回文档中心](00_INDEX.md) | [性能优化](07_PERFORMANCE.md) | [开发规范](05_DEVELOPMENT.md)

---

## 优化概述

本次优化主要解决项目中 Context 传递不规范的问题，实现：

- ✅ **请求级别超时控制** - 每个请求都可以设置独立的超时时间
- ✅ **链路追踪支持** - 通过 Context 传递 TraceID、RequestID 等信息
- ✅ **优雅的取消机制** - 支持请求取消和资源释放
- ✅ **标准化的 Context 使用** - 统一的 Context 创建和传递方式

---

## 问题分析

### 1. 原有问题

#### ❌ Redis 使用固定 Context
```go
// 问题代码
type RedisClient struct {
    Client  *redis.Client
    Context context.Context  // 固定的 context.Background()
}

func (rds RedisClient) Set(key string, value interface{}, expiration time.Duration) bool {
    if err := rds.Client.Set(rds.Context, key, value, expiration).Err(); err != nil {
        return false
    }
    return true
}
```

**问题**：
- 无法传递请求级别的超时
- 无法实现请求取消
- 无法进行链路追踪

#### ❌ Repository 层缺少 Context 参数
```go
// 问题代码
func (r *commentRepository) GetByID(id string) (*comment.Comment, error) {
    // 无法传递 context，无法控制超时
    return database.DB.First(&comment, id).Error
}
```

**问题**：
- 数据库查询无法设置超时
- 无法取消长时间运行的查询
- 缺少请求追踪信息

---

## 优化方案

### 1. Context 助手工具包

创建 `pkg/ctx/context.go` 提供统一的 Context 管理：

```go
package ctx

import (
    "context"
    "time"
    "github.com/gin-gonic/gin"
)

// 从 Gin Context 获取请求 Context
func FromGinContext(c *gin.Context) context.Context {
    return c.Request.Context()
}

// 创建带超时的 Context
func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
    if timeout == 0 {
        timeout = 30 * time.Second // 默认30秒
    }
    return context.WithTimeout(parent, timeout)
}

// 添加请求ID
func WithRequestID(ctx context.Context, requestID string) context.Context {
    return context.WithValue(ctx, RequestIDKey, requestID)
}

// 添加用户ID
func WithUserID(ctx context.Context, userID string) context.Context {
    return context.WithValue(ctx, UserIDKey, userID)
}

// 添加链路追踪ID
func WithTraceID(ctx context.Context, traceID string) context.Context {
    return context.WithValue(ctx, TraceIDKey, traceID)
}
```

### 2. Redis 客户端优化

#### ✅ 移除固定 Context，方法级传递

```go
// 优化后
type RedisClient struct {
    Client *redis.Client  // 移除 Context 字段
}

// 所有方法都接收 context 参数
func (rds RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) bool {
    if err := rds.Client.Set(ctx, key, value, expiration).Err(); err != nil {
        logger.ErrorString("Redis", "Set", err.Error())
        return false
    }
    return true
}

func (rds RedisClient) Get(ctx context.Context, key string) string {
    result, err := rds.Client.Get(ctx, key).Result()
    if err != nil {
        if err != redis.Nil {
            logger.ErrorString("Redis", "Get", err.Error())
        }
        return ""
    }
    return result
}
```

**优势**：
- ✅ 每个 Redis 操作都可以独立设置超时
- ✅ 支持请求级别的取消
- ✅ 可以传递链路追踪信息

### 3. Repository 层优化

#### ✅ 所有方法添加 Context 参数

```go
// 优化后的接口定义
type CommentRepository interface {
    GetByID(ctx context.Context, id string) (*comment.Comment, error)
    List(ctx context.Context, c *gin.Context, perPage int) ([]comment.Comment, *paginator.Paging, error)
    Create(ctx context.Context, comment *comment.Comment) error
    Update(ctx context.Context, comment *comment.Comment) error
    Delete(ctx context.Context, id string) error
    // ... 其他方法
}

// 实现示例
func (r *commentRepository) GetByID(ctx context.Context, id string) (*comment.Comment, error) {
    var commentModel comment.Comment
    
    // 使用带超时的 context 进行查询
    if err := database.DB.WithContext(ctx).
        Select("id", "topic_id", "user_id", "content", "like_count", "created_at").
        First(&commentModel, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, err
    }
    return &commentModel, nil
}
```

### 4. Service 层优化

#### ✅ 从 Gin Context 创建请求 Context

```go
// 优化后
func (s *CommentService) GetByID(c *gin.Context, id string) (*CommentResponseDTO, *apperrors.AppError) {
    // 从 Gin Context 创建请求 Context
    ctx := ctx.FromGinContext(c)
    
    // 添加请求信息
    ctx = ctx.WithRequestID(ctx, c.GetString("request_id"))
    ctx = ctx.WithUserID(ctx, auth.CurrentUID(c))
    
    // 创建带超时的 context
    ctx, cancel := ctx.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    // 传递 context 到 Repository 层
    commentModel, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, apperrors.DatabaseError("获取评论", err)
    }
    
    return s.toResponseDTO(commentModel), nil
}
```

---

## 使用示例

### 示例 1: Controller 中使用 Context

```go
func (ctrl *CommentsController) Show(c *gin.Context) {
    // 从 Gin Context 获取请求 Context
    ctx := ctx.FromGinContext(c)
    
    // 添加请求追踪信息
    ctx = ctx.WithRequestID(ctx, c.GetString("request_id"))
    ctx = ctx.WithUserID(ctx, auth.CurrentUID(c))
    
    // 调用 Service 层
    comment, err := ctrl.commentService.GetByID(ctx, c.Param("id"))
    if err != nil {
        response.ApiError(c, 500, err.Code, err.Message)
        return
    }
    
    response.Data(c, comment)
}
```

### 示例 2: Repository 中使用 Context

```go
func (r *commentRepository) GetByID(ctx context.Context, id string) (*comment.Comment, error) {
    var commentModel comment.Comment
    
    // GORM 使用 WithContext 传递 context
    if err := database.DB.WithContext(ctx).
        Select("id", "content", "user_id").
        First(&commentModel, id).Error; err != nil {
        return nil, err
    }
    
    return &commentModel, nil
}
```

### 示例 3: Redis 中使用 Context

```go
func (c *CommentCache) GetByID(ctx context.Context, id string) (*comment.Comment, error) {
    key := c.cacheKeyPrefix + id
    
    // Redis 操作传递 context
    data := redis.Redis.Get(ctx, key)
    if data == "" {
        return nil, nil
    }
    
    var commentModel comment.Comment
    err := json.Unmarshal([]byte(data), &commentModel)
    if err != nil {
        return nil, err
    }
    
    return &commentModel, nil
}

func (c *CommentCache) Set(ctx context.Context, commentModel *comment.Comment) error {
    key := c.cacheKeyPrefix + fmt.Sprintf("%d", commentModel.ID)
    
    data, err := json.Marshal(commentModel)
    if err != nil {
        return err
    }
    
    // Redis Set 传递 context
    redis.Redis.Set(ctx, key, string(data), c.cacheTime)
    return nil
}
```

### 示例 4: 带超时控制的查询

```go
func (s *CommentService) ListByTopicID(c *gin.Context, topicID string, perPage int) (*CommentListResponseDTO, *apperrors.AppError) {
    ctx := ctx.FromGinContext(c)
    
    // 设置 3 秒超时
    ctx, cancel := ctx.WithTimeout(ctx, 3*time.Second)
    defer cancel()
    
    // 查询会在 3 秒后自动取消
    comments, paging, err := s.repo.ListByTopicID(ctx, c, topicID, perPage)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            return nil, apperrors.TimeoutError("查询超时")
        }
        return nil, apperrors.DatabaseError("获取评论列表", err)
    }
    
    return &CommentListResponseDTO{
        Comments: s.toResponseDTOList(comments),
        Paging:   paging,
    }, nil
}
```

---

## 迁移指南

### 步骤 1: 更新 Redis 调用

```go
// 旧代码
redis.Redis.Set(key, value, expiration)
data := redis.Redis.Get(key)

// 新代码 - 传递 context
ctx := ctx.FromGinContext(c)
redis.Redis.Set(ctx, key, value, expiration)
data := redis.Redis.Get(ctx, key)
```

### 步骤 2: 更新 Repository 调用

```go
// 旧代码
comment, err := repo.GetByID(id)

// 新代码 - 传递 context
ctx := ctx.FromGinContext(c)
comment, err := repo.GetByID(ctx, id)
```

### 步骤 3: 更新 GORM 查询

```go
// 旧代码
database.DB.First(&model, id)

// 新代码 - 使用 WithContext
database.DB.WithContext(ctx).First(&model, id)
```

---

## 最佳实践

### 1. Context 创建

✅ **推荐**：从 Gin Context 获取
```go
ctx := ctx.FromGinContext(c)
```

❌ **避免**：在业务代码中使用 context.Background()
```go
ctx := context.Background() // 仅用于初始化和测试
```

### 2. Context 传递

✅ **推荐**：作为第一个参数传递
```go
func GetByID(ctx context.Context, id string) (*Model, error)
```

❌ **避免**：将 Context 存储在结构体中
```go
type Service struct {
    ctx context.Context // 不要这样做
}
```

### 3. 超时设置

✅ **推荐**：根据操作类型设置合理的超时
```go
// 查询操作：3-5秒
ctx, cancel := ctx.WithTimeout(ctx, 3*time.Second)
defer cancel()

// 写入操作：5-10秒
ctx, cancel := ctx.WithTimeout(ctx, 10*time.Second)
defer cancel()
```

### 4. 错误处理

✅ **推荐**：检查超时错误
```go
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        return apperrors.TimeoutError("操作超时")
    }
    if errors.Is(err, context.Canceled) {
        return apperrors.CanceledError("操作已取消")
    }
    return apperrors.DatabaseError("数据库错误", err)
}
```

### 5. Cancel 函数

✅ **推荐**：总是调用 cancel
```go
ctx, cancel := ctx.WithTimeout(parent, timeout)
defer cancel() // 确保资源释放
```

---

## 性能影响

### 优化前后对比

| 场景 | 优化前 | 优化后 | 改进 |
|------|--------|--------|------|
| **请求超时控制** | 无 | 支持 | ✅ 新增 |
| **请求取消** | 无 | 支持 | ✅ 新增 |
| **链路追踪** | 无 | 支持 | ✅ 新增 |
| **性能开销** | 0 | <1ms | ✅ 可忽略 |

### Context 传递的额外开销

- **内存开销**: 每个 Context ~48 字节
- **CPU开销**: 创建和传递 <1 微秒
- **结论**: 性能影响可忽略不计，收益远大于开销

---

## 注意事项

### 1. Context 使用规范

- ✅ Context 应该作为函数的第一个参数
- ✅ Context 只能传递，不能存储在结构体中
- ✅ Context 的值应该是请求范围的，不要用于传递可选参数
- ✅ 使用 context.TODO() 标记需要添加 Context 的地方

### 2. 超时时间设置

- 数据库查询: 3-5 秒
- Redis 操作: 1-3 秒
- HTTP 请求: 10-30 秒
- 批量操作: 根据数据量调整

### 3. 避免的反模式

❌ **不要在结构体中存储 Context**
```go
type Service struct {
    ctx context.Context // 错误！
}
```

❌ **不要传递 nil Context**
```go
DoSomething(nil) // 错误！应该传递 context.Background() 或 context.TODO()
```

❌ **不要忽略 cancel 函数**
```go
ctx, _ := context.WithTimeout(parent, timeout) // 错误！会导致资源泄漏
```

---

## 测试建议

### 1. 超时测试

```go
func TestGetByID_Timeout(t *testing.T) {
    repo := NewCommentRepository()
    
    // 创建一个立即超时的 context
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
    defer cancel()
    
    time.Sleep(10 * time.Millisecond) // 确保超时
    
    _, err := repo.GetByID(ctx, "1")
    assert.Error(t, err)
    assert.True(t, errors.Is(err, context.DeadlineExceeded))
}
```

### 2. 取消测试

```go
func TestGetByID_Cancel(t *testing.T) {
    repo := NewCommentRepository()
    
    ctx, cancel := context.WithCancel(context.Background())
    
    // 在查询前取消
    cancel()
    
    _, err := repo.GetByID(ctx, "1")
    assert.Error(t, err)
    assert.True(t, errors.Is(err, context.Canceled))
}
```

---

## 相关文档

- [性能优化指南](07_PERFORMANCE.md) - 整体性能优化策略
- [开发规范](05_DEVELOPMENT.md) - 编码最佳实践
- [数据库优化](13_DATABASE_OPTIMIZATION.md) - 数据库查询优化
- [监控告警](11_MONITORING.md) - 性能监控和告警

---

**[⬆️ 返回顶部](#-context-传递优化指南)** | **[📚 返回文档中心](00_INDEX.md)**
