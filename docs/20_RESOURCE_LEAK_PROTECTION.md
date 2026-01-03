# 20. 资源泄漏防护 (Resource Leak Protection)

## 版本

- **版本**: v2.9
- **更新日期**: 2026-01-03
- **作者**: GoHub-Service Team

## 概述

本文档介绍 GoHub-Service 的资源泄漏防护方案，涵盖：

1. **HTTP Response Body 管理**：确保 HTTP 连接正确关闭
2. **数据库事务处理**：防止事务泄漏和连接占用
3. **Goroutine 管理**：防止 goroutine 无限制增长
4. **Context 取消**：防止 context 泄漏导致的 goroutine 泄漏
5. **资源追踪工具**：检测和报告潜在的资源泄漏

## 资源泄漏风险

### 常见泄漏场景

| 资源类型 | 泄漏场景 | 影响 | 检测难度 |
|---------|---------|------|---------|
| HTTP Body | 忘记调用 `Close()` | 连接池耗尽 | ⭐⭐⭐ |
| 数据库事务 | 未 Commit/Rollback | 连接占用 | ⭐⭐⭐⭐ |
| Goroutine | 无限制创建 | 内存耗尽 | ⭐⭐⭐⭐⭐ |
| Context | 未调用 `cancel()` | Goroutine 泄漏 | ⭐⭐⭐⭐ |
| 文件句柄 | 未关闭文件 | 句柄耗尽 | ⭐⭐ |

### 项目现状

✅ **已正确处理**：
- HTTP Response Body 已使用 `defer res.Body.Close()`（pkg/elasticsearch/client.go）
- 数据库事务使用 GORM 的 `Transaction()` 方法（自动处理 Rollback）
- 文件操作使用 `defer file.Close()`（pkg/file/file.go）

⚠️ **需要加强**：
- 缺少 goroutine 池管理机制
- 缺少 context 取消的统一防护
- 缺少资源泄漏检测工具
- 缺少自定义资源的 Rollback 处理示例

## 解决方案

### 1. 资源管理工具包 (pkg/resource)

提供统一的资源管理工具：

#### SafeClose - 安全关闭资源

```go
import "GoHub-Service/pkg/resource"

// 使用示例
resp, err := http.Get("https://api.example.com/data")
if err != nil {
    return err
}
defer resource.SafeClose(resp.Body, logger.Logger)

// SafeClose 会：
// 1. 检查 closer 是否为 nil
// 2. 捕获 Close() 可能的 panic
// 3. 记录错误日志
```

**优势**：
- 防止 panic 传播
- 统一错误处理
- 自动日志记录

#### Tracker - 资源泄漏追踪

```go
tracker := resource.NewTracker(logger.Logger)

// 追踪资源
resourceID := fmt.Sprintf("conn-%p", conn)
tracker.Track(resourceID, "database.Connection")

// 使用资源...

// 释放资源
defer tracker.Untrack(resourceID)

// 定期检查泄漏（超过 5 分钟未释放）
tracker.Report(5 * time.Minute)
```

**功能**：
- 追踪资源创建时间
- 记录资源类型和堆栈
- 定期检测超时未释放的资源
- 自动生成泄漏报告

#### GoRoutinePool - Goroutine 池

```go
pool := resource.NewGoRoutinePool(20, logger.Logger) // 最多 20 个并发
defer pool.Shutdown(30 * time.Second)

// 提交任务
for i := 0; i < 10000; i++ {
    taskID := i
    if err := pool.Submit(func() {
        // 执行任务
        processTask(taskID)
    }); err != nil {
        return err
    }
}
```

**特性**：
- 限制并发 goroutine 数量
- 自动捕获任务 panic
- 优雅关闭机制
- 防止 goroutine 泄漏

#### TransactionGuard - 事务守卫

```go
tx, err := db.Begin()
if err != nil {
    return err
}

guard := resource.NewTransactionGuard(tx, logger.Logger)
defer guard.Release() // 自动回滚未提交的事务

// 执行操作
_, err = tx.Exec("INSERT INTO users (name) VALUES (?)", "alice")
if err != nil {
    return err
}

// 提交事务
if err := tx.Commit(); err != nil {
    return err
}
guard.Commit() // 标记已提交
```

**保护机制**：
- 自动检测未提交/未回滚的事务
- `Release()` 时自动回滚
- 记录泄漏堆栈
- 防止重复回滚

#### ContextGuard - Context 守卫

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
guard := resource.NewContextGuard(ctx, cancel, logger.Logger)
defer guard.Release() // 自动调用 cancel

// 使用 context...
data, err := fetchData(ctx)
if err != nil {
    return err
}

guard.Cancel() // 手动取消（可选）
```

**作用**：
- 确保 `cancel()` 被调用
- 防止 context 泄漏
- 记录未取消的 context
- 调试 goroutine 泄漏

## 最佳实践

### 1. HTTP 请求处理

#### ❌ 错误示例

```go
func fetchData() ([]byte, error) {
    resp, err := http.Get("https://api.example.com/data")
    if err != nil {
        return nil, err
    }
    // 忘记关闭 Body
    return io.ReadAll(resp.Body)
}
```

#### ✅ 正确示例

```go
func fetchData() ([]byte, error) {
    resp, err := http.Get("https://api.example.com/data")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close() // 确保关闭
    
    return io.ReadAll(resp.Body)
}
```

#### ⭐ 推荐示例

```go
func fetchData() ([]byte, error) {
    resp, err := http.Get("https://api.example.com/data")
    if err != nil {
        return nil, err
    }
    defer resource.SafeClose(resp.Body, logger.Logger) // 安全关闭
    
    return io.ReadAll(resp.Body)
}
```

### 2. 数据库事务处理

#### ❌ 错误示例（原始 sql.DB）

```go
func createUser(db *sql.DB, name string) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    
    _, err = tx.Exec("INSERT INTO users (name) VALUES (?)", name)
    if err != nil {
        return err // 泄漏：未回滚
    }
    
    return tx.Commit()
}
```

#### ✅ 推荐示例（GORM）

```go
func createUser(db *gorm.DB, user *User) error {
    // GORM Transaction 自动处理 Commit/Rollback
    return db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(user).Error; err != nil {
            return err // 自动回滚
        }
        
        settings := &UserSettings{UserID: user.ID}
        if err := tx.Create(settings).Error; err != nil {
            return err // 自动回滚
        }
        
        return nil // 自动提交
    })
}
```

#### ⭐ 复杂场景示例（sql.DB + TransactionGuard）

```go
func complexTransaction(db *sql.DB) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    
    guard := resource.NewTransactionGuard(tx, logger.Logger)
    defer guard.Release() // 确保回滚
    
    // 步骤 1
    if _, err := tx.Exec("INSERT INTO users (name) VALUES (?)", "alice"); err != nil {
        return err
    }
    
    // 步骤 2
    if _, err := tx.Exec("INSERT INTO logs (action) VALUES (?)", "user_created"); err != nil {
        return err
    }
    
    // 提交
    if err := tx.Commit(); err != nil {
        return err
    }
    
    guard.Commit()
    return nil
}
```

### 3. Goroutine 管理

#### ❌ 错误示例

```go
func processUsers(users []User) {
    for _, user := range users {
        go func(u User) {
            // 处理用户
            time.Sleep(1 * time.Second)
            fmt.Printf("Processed user %d\n", u.ID)
        }(user)
        // 可能创建数千个 goroutine
    }
}
```

#### ✅ 正确示例

```go
func processUsers(users []User) error {
    pool := resource.NewGoRoutinePool(10, logger.Logger) // 限制 10 个并发
    defer pool.Shutdown(30 * time.Second)
    
    for _, user := range users {
        u := user
        if err := pool.Submit(func() {
            // 处理用户
            time.Sleep(1 * time.Second)
            fmt.Printf("Processed user %d\n", u.ID)
        }); err != nil {
            return err
        }
    }
    
    return nil
}
```

### 4. Context 超时控制

#### ❌ 错误示例

```go
func queryData() (*Data, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    // 忘记调用 cancel
    
    return db.QueryContext(ctx, "SELECT * FROM data")
}
```

#### ✅ 正确示例

```go
func queryData() (*Data, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel() // 确保调用
    
    return db.QueryContext(ctx, "SELECT * FROM data")
}
```

#### ⭐ 推荐示例

```go
func queryData() (*Data, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    guard := resource.NewContextGuard(ctx, cancel, logger.Logger)
    defer guard.Release()
    
    data, err := db.QueryContext(ctx, "SELECT * FROM data")
    if err != nil {
        return nil, err
    }
    
    guard.Cancel() // 手动取消（可选）
    return data, nil
}
```

### 5. 资源泄漏检测

#### 生产环境监控

```go
func main() {
    tracker := resource.NewTracker(logger.Logger)
    
    // 启动定期检查（每分钟）
    go func() {
        ticker := time.NewTicker(1 * time.Minute)
        defer ticker.Stop()
        
        for range ticker.C {
            // 报告超过 5 分钟未释放的资源
            tracker.Report(5 * time.Minute)
            
            // 记录资源统计
            logger.Info("资源统计",
                zap.Int("active_resources", tracker.Count()),
            )
        }
    }()
    
    // 应用逻辑...
}
```

#### 开发环境检测

```go
func TestResourceLeak(t *testing.T) {
    tracker := resource.NewTracker(logger.Logger)
    
    // 模拟资源使用
    tracker.Track("res-1", "http.Response")
    time.Sleep(100 * time.Millisecond)
    
    // 检查泄漏
    leaked := tracker.Check(50 * time.Millisecond)
    assert.Empty(t, leaked, "发现资源泄漏")
    
    // 清理
    tracker.Untrack("res-1")
}
```

## 性能影响

### 开销对比

| 工具 | 内存开销 | CPU 开销 | 场景适用性 |
|-----|---------|---------|----------|
| SafeClose | ~48 字节 | <1μs | 所有场景 ✅ |
| Tracker | ~200 字节/资源 | <5μs/操作 | 开发/测试 ⭐⭐⭐ |
| GoRoutinePool | ~1KB + 8KB/goroutine | <10μs/提交 | 高并发场景 ✅ |
| TransactionGuard | ~64 字节 | <1μs | sql.DB 事务 ⭐⭐ |
| ContextGuard | ~48 字节 | <1μs | 所有 context ✅ |

### 性能测试结果

```bash
# Benchmark 结果
BenchmarkSafeClose-8          10000000    120 ns/op    48 B/op    1 allocs/op
BenchmarkTracker-8             1000000   1200 ns/op   240 B/op    4 allocs/op
BenchmarkGoRoutinePool-8        500000   2500 ns/op   512 B/op    2 allocs/op
```

**结论**：
- SafeClose/ContextGuard：零性能影响，推荐全局使用
- GoRoutinePool：轻微开销，高并发场景必用
- Tracker：适中开销，建议开发/测试环境启用，生产环境可选

## 迁移指南

### 第一阶段：关键路径（必须）

1. **HTTP 请求**：将 `defer resp.Body.Close()` 替换为 `defer resource.SafeClose(resp.Body, logger.Logger)`
2. **高并发任务**：使用 `GoRoutinePool` 替代直接创建 goroutine
3. **Context 超时**：使用 `ContextGuard` 确保 cancel 被调用

### 第二阶段：全面覆盖（推荐）

1. **所有 HTTP 请求**：统一使用 SafeClose
2. **所有 goroutine**：通过 GoRoutinePool 管理
3. **所有 context**：使用 ContextGuard 防护

### 第三阶段：监控检测（进阶）

1. **生产环境**：启用 Tracker 定期检查（阈值 5 分钟）
2. **测试环境**：强制检测资源泄漏（阈值 10 秒）
3. **CI/CD**：集成泄漏检测测试

### 向后兼容性

✅ **完全兼容**：
- 新增工具包，不修改现有代码
- 可逐步迁移，无破坏性变更
- GORM Transaction 继续使用原有方式
- 现有 `defer Close()` 代码可保留

## 故障排查

### Q1: 如何诊断 goroutine 泄漏？

**A**: 使用 pprof + Tracker 定位：

```bash
# 1. 获取 goroutine profile
curl http://localhost:6060/debug/pprof/goroutine > goroutine.txt

# 2. 分析 goroutine 数量
go tool pprof -top goroutine.txt

# 3. 查看 Tracker 报告
# 日志中搜索 "检测到可能的资源泄漏"
```

### Q2: GoRoutinePool 队列满了怎么办？

**A**: 调整池大小或增加超时：

```go
// 方案 1：增大池容量
pool := resource.NewGoRoutinePool(50, logger.Logger) // 20 → 50

// 方案 2：增加提交超时（修改源码）
// 当前超时：5 秒
```

### Q3: Tracker 误报资源泄漏？

**A**: 调整检测阈值或排除特定资源：

```go
// 方案 1：增大阈值
tracker.Report(10 * time.Minute) // 5 分钟 → 10 分钟

// 方案 2：排除长期资源（修改源码）
// 在 Check() 中添加资源类型过滤
```

### Q4: TransactionGuard 与 GORM 冲突？

**A**: GORM 推荐使用 `Transaction()` 方法，无需 TransactionGuard：

```go
// ✅ GORM 推荐
db.Transaction(func(tx *gorm.DB) error {
    // 自动处理 Commit/Rollback
})

// ⚠️ 仅在使用 sql.DB 时使用 TransactionGuard
tx, _ := sqlDB.Begin()
guard := resource.NewTransactionGuard(tx, logger.Logger)
defer guard.Release()
```

## 相关文档

- [05. 开发规范](./05_DEVELOPMENT.md) - 编码规范
- [07. 性能优化](./07_PERFORMANCE.md) - 性能监控
- [14. Context 优化](./14_CONTEXT_OPTIMIZATION.md) - Context 传递
- [18. 连接池优化](./18_DATABASE_POOL_OPTIMIZATION.md) - 数据库连接管理

## 更新历史

- **v2.9 (2026-01-03)**: 初始版本
  - 实现 SafeClose 安全关闭工具
  - 实现 Tracker 资源泄漏检测
  - 实现 GoRoutinePool goroutine 池
  - 实现 TransactionGuard 事务守卫
  - 实现 ContextGuard context 守卫
  - 提供完整使用示例和迁移指南
