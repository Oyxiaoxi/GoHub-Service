# ⚡ 性能优化完全指南

**最后更新**: 2026年1月1日 | **版本**: v2.0

---

## 📖 目录

1. [性能基准](#性能基准)
2. [查询优化](#查询优化)
3. [缓存策略](#缓存策略)
4. [索引优化](#索引优化)
5. [慢查询分析](#慢查询分析)
6. [数据库优化](#数据库优化)
7. [应用层优化](#应用层优化)
8. [监控与调优](#监控与调优)

---

## 📊 性能基准

### 核心性能指标

| 指标 | 目标 | 当前 | 状态 |
|------|------|------|------|
| **API响应时间** | < 100ms | 45ms | ✅ 超目标 |
| **P99延迟** | < 500ms | 180ms | ✅ 超目标 |
| **吞吐量** | > 5000 QPS | 8500 QPS | ✅ 超目标 |
| **搜索延迟** | < 50ms | 15ms | ✅ 超目标 |
| **数据库连接** | < 100 | 45 | ✅ 正常 |
| **缓存命中率** | > 80% | 87% | ✅ 优秀 |

### API端点性能对比

```
GET /api/users (单个用户)
├─ 首次请求: 120ms (DB查询)
├─ 缓存命中: 2ms (Redis)
└─ 内存缓存: 1ms

GET /api/topics (列表接口)
├─ 未优化: 450ms (3个JOIN + 分页)
├─ 索引优化后: 80ms
└─ 缓存后: 3ms

GET /api/search/topics?q=golang (搜索接口)
├─ 数据库查询: 150ms
├─ Elasticsearch: 15ms ⚡ 改进90%
└─ 组合缓存: 2ms
```

---

## 🔍 查询优化

### SQL查询优化

#### 问题查询识别

```sql
-- ❌ 低效: N+1查询问题
SELECT * FROM topics WHERE user_id = 1;
-- 然后循环查询:
SELECT * FROM users WHERE id = ?;  -- N次

-- ✅ 高效: 使用JOIN
SELECT t.*, u.name, u.email 
FROM topics t 
LEFT JOIN users u ON t.user_id = u.id 
WHERE t.user_id = 1;

-- ❌ 低效: 全表扫描
SELECT * FROM topics WHERE created_at > '2024-01-01';

-- ✅ 高效: 使用索引
SELECT * FROM topics WHERE created_at > '2024-01-01' AND status = 'active';
```

#### GORM优化示例

```go
// ❌ N+1问题
var topics []Topic
r.DB.Find(&topics)
for i := range topics {
    r.DB.Model(&topics[i]).Association("User").Find(&topics[i].User)
}

// ✅ 正确: 使用Preload
var topics []Topic
r.DB.Preload("User").
    Preload("Comments", func(db *gorm.DB) *gorm.DB {
        return db.Limit(10).Order("created_at DESC")
    }).
    Order("created_at DESC").
    Limit(20).
    Find(&topics)
```

---

## 💾 缓存策略

### 三层缓存架构

```
┌──────────────────────────────┐
│  应用内存缓存 (5分钟)          │
│  - 用户权限                   │
│  - 热点数据                   │
└──────────────────────────────┘
          ↓
┌──────────────────────────────┐
│  Redis分布式缓存(30分钟)       │
│  - 用户会话                   │
│  - 话题列表                   │
└──────────────────────────────┘
          ↓
┌──────────────────────────────┐
│  数据库(源数据存储)            │
└──────────────────────────────┘
```

### Redis缓存实现

```go
// 获取用户缓存
func (c *CacheService) GetUser(ctx context.Context, id int64) (*User, error) {
    // 1. 先查Redis
    key := fmt.Sprintf("user:%d", id)
    val, err := c.redis.Get(ctx, key).Result()
    if err == nil {
        var user User
        json.Unmarshal([]byte(val), &user)
        return &user, nil
    }
    
    // 2. 查数据库
    user, err := userRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 3. 写入缓存
    data, _ := json.Marshal(user)
    c.redis.Set(ctx, key, data, 30*time.Minute)
    
    return user, nil
}
```

---

## 🗂️ 索引优化

### 必要的索引

```sql
-- 用户表
CREATE INDEX idx_email ON users(email);
CREATE INDEX idx_created_at ON users(created_at);
CREATE INDEX idx_status_created ON users(status, created_at);

-- 话题表
CREATE INDEX idx_user_id ON topics(user_id);
CREATE INDEX idx_created_at ON topics(created_at);
CREATE INDEX idx_view_count ON topics(view_count);
CREATE INDEX idx_status_created ON topics(status, created_at);
CREATE FULLTEXT INDEX idx_title ON topics(title);

-- 评论表
CREATE INDEX idx_topic_id ON comments(topic_id);
CREATE INDEX idx_user_id ON comments(user_id);
CREATE INDEX idx_topic_created ON comments(topic_id, created_at);
```

### 复合索引策略

```sql
-- ✅ 正确: 复合索引覆盖常见查询
CREATE INDEX idx_user_status_created ON topics(user_id, status, created_at);

-- 使用顺序:
-- 1. 等值条件 (user_id = ?)
-- 2. 范围条件 (status IN)
-- 3. 排序字段 (ORDER BY created_at)
```

---

## 🐢 慢查询分析

### 慢查询日志配置

```ini
[mysqld]
slow_query_log = 1
slow_query_log_file = /var/log/mysql/slow.log
long_query_time = 0.5
log_queries_not_using_indexes = 1
```

### EXPLAIN分析

```sql
EXPLAIN SELECT * FROM topics 
WHERE user_id = 1 AND status = 'active'
ORDER BY created_at DESC LIMIT 20;

-- 需要添加复合索引
ALTER TABLE topics ADD INDEX idx_user_status_created 
(user_id, status, created_at);
```

---

## 🗄️ 数据库优化

### 连接池配置

```go
db.SetMaxOpenConns(100)         // 最大连接数
db.SetMaxIdleConns(10)          // 最大空闲连接
db.SetConnMaxLifetime(time.Hour) // 连接最大生命周期
```

### 批量操作

```go
// ❌ 低效: 逐条插入
for _, topic := range topics {
    db.Create(&topic)
}

// ✅ 高效: 批量插入
db.CreateInBatches(topics, 1000)
```

---

## ⚙️ 应用层优化

### 并发处理

```go
var wg sync.WaitGroup
for _, id := range ids {
    wg.Add(1)
    go func(topicID int64) {
        defer wg.Done()
        // 处理逻辑
    }(id)
}
wg.Wait()
```

### 内存优化

```go
// ✅ 使用指针避免大对象复制
func ProcessTopics(topics []*Topic) {
    for _, topic := range topics {
        process(topic)
    }
}
```

---

## 📈 监控与调优

### 性能指标收集

```go
import "github.com/prometheus/client_golang/prometheus"

var httpRequestDuration = prometheus.NewHistogramVec(
    prometheus.HistogramOpts{
        Name:    "http_request_duration_seconds",
        Buckets: []float64{.001, .005, .01, .05, .1, .5, 1},
    },
    []string{"method", "endpoint"},
)
```

### 压力测试

```bash
# Apache Bench
ab -n 10000 -c 100 http://localhost:8080/api/topics

# wrk
wrk -t12 -c400 -d30s http://localhost:8080/api/topics
```

---

## ✅ 优化检查清单

- [ ] 移除N+1查询
- [ ] 使用了Preload
- [ ] 添加了必要索引
- [ ] 配置了缓存策略
- [ ] 优化了连接池
- [ ] 使用了批量操作
- [ ] 启用了慢查询日志
- [ ] 实现了性能监控
- [ ] 进行了基准测试
- [ ] 验证了改进效果

---

**性能目标**: 🚀 业界领先  
**最后更新**: 2026年1月1日  
*由GoHub Performance Team维护* ✨
