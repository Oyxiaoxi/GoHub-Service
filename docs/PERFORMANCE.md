# ⚡ 性能优化

缓存策略、数据库优化和性能最佳实践。

## 1. 缓存策略

### 多层缓存架构

```
┌─────────────────────────────┐
│   缓存层 1: 内存缓存          │ (毫秒级)
│   - 最热数据                  │
│   - TTL: 5-10 分钟            │
└────────────┬──────────────────┘
             │ 未命中
┌────────────▼──────────────────┐
│   缓存层 2: Redis             │ (10毫秒级)
│   - 热数据                     │
│   - TTL: 30-60 分钟            │
└────────────┬──────────────────┘
             │ 未命中
┌────────────▼──────────────────┐
│   缓存层 3: 数据库            │ (100毫秒级)
│   - 真实数据源                 │
└─────────────────────────────┘
```

### Redis 缓存使用

文件位置: `app/cache/`

```go
package cache

import (
    "GoHub-Service/pkg/redis"
    "time"
)

type TopicCache struct {
    redis *redis.Client
}

func NewTopicCache(redis *redis.Client) *TopicCache {
    return &TopicCache{redis: redis}
}

// 缓存键命名规范
// 格式: <module>:<action>:<param1>:<param2>
// 示例: topics:popular:page:1, topics:list:category:1:page:2

const (
    TopicListKey    = "topics:list"
    TopicDetailKey  = "topics:detail"
    PopularTopics   = "topics:popular"
)

// 获取缓存
func (c *TopicCache) GetList(page, pageSize int) interface{} {
    key := fmt.Sprintf("%s:page:%d:size:%d", TopicListKey, page, pageSize)
    val, _ := c.redis.Get(key)
    return val
}

// 设置缓存
func (c *TopicCache) SetList(page, pageSize int, data interface{}) {
    key := fmt.Sprintf("%s:page:%d:size:%d", TopicListKey, page, pageSize)
    c.redis.Set(key, data, 30*time.Minute)
}

// 删除缓存
func (c *TopicCache) Clear() {
    c.redis.Delete(TopicListKey + ":*")
}

// 缓存失效
func (c *TopicCache) Invalidate(id uint) {
    key := fmt.Sprintf("%s:%d", TopicDetailKey, id)
    c.redis.Delete(key)
}
```

### 缓存失效策略

```go
// Strategy 1: 主动失效 (推荐)
func (s *TopicService) Update(topic *models.Topic) error {
    // 1. 更新数据库
    if err := s.repo.Update(topic); err != nil {
        return err
    }
    
    // 2. 清除相关缓存
    s.cache.InvalidateTopic(topic.ID)
    s.cache.InvalidateCategory(topic.CategoryID)
    
    return nil
}

// Strategy 2: TTL 过期
func (c *TopicCache) SetList(data interface{}) {
    // 30 分钟后自动过期
    c.redis.Set("topics:list", data, 30*time.Minute)
}

// Strategy 3: 分布式失效
func (s *TopicService) Delete(id uint) error {
    // 发送缓存失效消息到消息队列
    s.queue.Publish("cache:invalidate", map[string]interface{}{
        "type": "topic",
        "id":   id,
    })
    
    return s.repo.Delete(id)
}
```

## 2. 数据库优化

### 索引优化

✅ **应该创建索引的字段**:
- 外键字段 (FK)
- 频繁用于 WHERE 的字段
- 频繁用于 ORDER BY 的字段
- 频繁用于 JOIN 的字段

```sql
-- 创建索引
CREATE INDEX idx_user_id ON topics(user_id);
CREATE INDEX idx_category_id ON topics(category_id);
CREATE INDEX idx_created_at ON topics(created_at);

-- 复合索引（字段顺序很重要）
CREATE INDEX idx_user_category ON topics(user_id, category_id);

-- 删除不必要的索引
DROP INDEX idx_unused ON topics;
```

❌ **不应该创建索引的字段**:
- 值分布少的字段 (如 status: 1/0)
- 更新频繁的字段
- 低基数字段 (如 gender: M/F)

### 查询优化

```go
// ❌ N+1 查询问题
topics, _ := repo.FindAll()
for _, topic := range topics {
    user, _ := repo.GetUser(topic.UserID)  // 执行 N 次查询！
}

// ✅ 使用预加载
topics, _ := repo.FindAllWithUser()  // 单次查询

// 实现方式
func (r *TopicRepository) FindAllWithUser() ([]models.Topic, error) {
    var topics []models.Topic
    // 使用 Preload 预加载关联数据
    err := r.db.Preload("User").Find(&topics).Error
    return topics, err
}
```

### 分页优化

```go
// ❌ 不好的分页（计算 offset）
SELECT * FROM topics OFFSET 10000 LIMIT 10;  // 跳过 10000 行

// ✅ 游标分页（更高效）
SELECT * FROM topics WHERE id > 999 LIMIT 10;
```

### 批量操作

```go
// ❌ 逐条插入
for _, topic := range topics {
    db.Create(&topic)  // 执行 N 次查询
}

// ✅ 批量插入
db.CreateInBatches(topics, 1000)  // 执行 1 次查询
```

## 3. 查询优化建议

### SELECT 优化

```go
// ❌ 查询所有字段
var topics []models.Topic
db.Find(&topics)

// ✅ 只查询需要的字段
var topics []struct {
    ID        uint
    Title     string
    CreatedAt time.Time
}
db.Model(&models.Topic{}).
    Select("id", "title", "created_at").
    Find(&topics)
```

### WHERE 优化

```go
// ✅ 使用索引字段
db.Where("user_id = ?", userID).Find(&topics)

// ❌ 函数查询（无法使用索引）
db.Where("YEAR(created_at) = ?", 2024).Find(&topics)

// ✅ 改用日期范围
start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
end := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
db.Where("created_at >= ? AND created_at < ?", start, end).Find(&topics)
```

### JOIN 优化

```go
// ✅ 主动 JOIN 而不是 N+1
var results []struct {
    TopicID  uint
    Title    string
    Username string
}
db.Table("topics").
    Select("topics.id as topic_id, topics.title, users.username").
    Joins("JOIN users ON users.id = topics.user_id").
    Find(&results)
```

## 4. 连接池优化

### 数据库连接池配置

文件位置: `config/database.go`

```go
type DatabaseConfig struct {
    // 连接池大小
    MaxOpenConns    int  // 最多开放连接数（推荐: CPU数 * 2 + 额外连接）
    MaxIdleConns    int  // 最多空闲连接数（推荐: MaxOpenConns 的 50%）
    ConnMaxLifetime time.Duration  // 连接最大存活时间
    ConnMaxIdleTime time.Duration  // 连接最大空闲时间
}

// 使用
sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(25)      // MySQL: 25-100
sqlDB.SetMaxIdleConns(5)       // 5-10
sqlDB.SetConnMaxLifetime(time.Hour)
sqlDB.SetConnMaxIdleTime(10 * time.Minute)
```

### Redis 连接池配置

文件位置: `config/redis.go`

```go
type RedisConfig struct {
    // 连接池大小
    PoolSize    int  // 连接池大小（推荐: 10-50）
    MinIdleConn int  // 最少空闲连接数
}
```

## 5. 慢查询分析

### 启用慢查询日志

```sql
-- MySQL
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 0.1;  -- 100ms
SET GLOBAL log_queries_not_using_indexes = 'ON';
```

### 查看慢查询

```sql
-- 查看是否有索引未使用
SELECT * FROM slow_log;

-- 使用 EXPLAIN 分析查询
EXPLAIN SELECT * FROM topics WHERE user_id = 1 AND category_id = 2;
```

### EXPLAIN 结果分析

| 字段 | 含义 | 优化目标 |
|------|------|---------|
| type | 访问方式 | ALL < index < range < ref < eq_ref < const |
| key | 使用的索引 | 应该不为 NULL |
| rows | 扫描的行数 | 越少越好 |
| Extra | 额外信息 | 避免 Using filesort、Using temporary |

## 6. 监控指标

### 关键性能指标 (KPI)

```go
// 响应时间
- API 平均响应时间: < 100ms
- 数据库查询时间: < 50ms
- Redis 查询时间: < 10ms

// 错误率
- API 5xx 错误率: < 0.1%
- 数据库连接错误: < 0.01%

// 资源使用
- 内存使用: < 80%
- CPU 使用: < 70%
- 磁盘 I/O: < 60%

// 缓存指标
- 缓存命中率: > 80%
- Redis 连接数: < 最大数的 50%
```

### 性能测试

```bash
# 压力测试
ab -n 10000 -c 100 http://localhost:3000/api/v1/topics

# 使用 Apache Bench
ab -c 100 -n 10000 -H "Authorization: Bearer <token>" http://localhost:3000/api/v1/topics

# 使用 wrk
wrk -t12 -c400 -d30s http://localhost:3000/api/v1/topics
```

## 7. 优化检查清单

- [ ] 添加必要的数据库索引
- [ ] 避免 N+1 查询（使用 Preload）
- [ ] 使用分页而不是加载所有数据
- [ ] 配置数据库连接池
- [ ] 实施多层缓存策略
- [ ] 使用 Redis 缓存热数据
- [ ] 批量操作而不是逐个操作
- [ ] 只查询需要的字段
- [ ] 使用查询结果集代替关联加载
- [ ] 定期分析慢查询日志
- [ ] 启用缓存预热
- [ ] 监控缓存命中率
- [ ] 优化索引字段选择
- [ ] 避免在索引字段上使用函数
- [ ] 使用异步处理耗时操作

## 常见性能问题

### 问题 1: 缓存击穿
```
症状: 缓存过期，大量请求直接打到数据库
解决: 使用缓存预热或布隆过滤器

func (s *Service) GetWithLock(id uint) {
    key := fmt.Sprintf("lock:get:%d", id)
    if !lock.Acquire(key, 3*time.Second) {
        return  // 只有一个请求查询数据库
    }
    defer lock.Release(key)
    
    // 查询数据库并缓存
}
```

### 问题 2: 缓存雪崩
```
症状: 大量缓存同时过期，请求集中打到数据库
解决: 使用随机 TTL

ttl := time.Duration(30+rand.Intn(10)) * time.Minute
cache.Set(key, value, ttl)
```

### 问题 3: 缓存一致性
```
症状: 缓存和数据库数据不一致
解决: 主动更新缓存

func (s *Service) Update(entity) {
    db.Save(entity)
    cache.Invalidate(entity.ID)  // 清除缓存
}
```

---

更多信息请查看 [ARCHITECTURE.md](./ARCHITECTURE.md)
