# GoHub-Service 数据库连接池配置与监控指南

## 目录

- [概述](#概述)
- [连接池配置](#连接池配置)
- [监控功能](#监控功能)
- [性能优化](#性能优化)
- [故障排查](#故障排查)
- [最佳实践](#最佳实践)

---

## 概述

### 优化目标

本文档介绍数据库连接池的合理配置和实时监控方案，帮助优化数据库连接管理，提升系统性能和稳定性。

### 关键改进

- ✅ **优化连接池配置**：max_open_connections 从 25 提升到 100
- ✅ **添加实时监控**：连接池统计、健康检查、指标采集
- ✅ **智能建议系统**：根据运行状态自动推荐配置调整
- ✅ **Prometheus 集成**：支持标准指标格式导出

---

## 连接池配置

### 1. 配置参数说明

#### config/database.go

```go
"mysql": map[string]interface{}{
    // 最大打开连接数
    // - 当前配置：100
    // - 推荐范围：100-200（根据业务负载调整）
    // - 计算公式：(2 × CPU核心数) + 磁盘数量
    // - 注意：设置过高会消耗过多资源，过低会导致等待
    "max_open_connections": config.Env("DB_MAX_OPEN_CONNECTIONS", 100),
    
    // 最大空闲连接数
    // - 当前配置：25
    // - 推荐范围：max_open_connections 的 1/4 到 1/2
    // - 作用：保持一定数量的空闲连接，减少连接建立开销
    // - 注意：设置过高会浪费资源，过低会频繁创建连接
    "max_idle_connections": config.Env("DB_MAX_IDLE_CONNECTIONS", 25),
    
    // 连接最大生命周期（秒）
    // - 当前配置：600 秒（10 分钟）
    // - 推荐范围：5-30 分钟
    // - 作用：防止连接长时间占用导致的问题
    // - 注意：MySQL wait_timeout 默认 8 小时，应设置小于此值
    "max_life_seconds": config.Env("DB_MAX_LIFE_SECONDS", 10*60),
},
```

### 2. 环境变量配置

#### .env 文件

```bash
# 数据库连接池配置
DB_MAX_OPEN_CONNECTIONS=100      # 最大打开连接数（默认：100）
DB_MAX_IDLE_CONNECTIONS=25       # 最大空闲连接数（默认：25）
DB_MAX_LIFE_SECONDS=600          # 连接最大生命周期秒数（默认：600）
```

### 3. 不同场景的推荐配置

#### 场景 1：低并发应用（< 100 QPS）

```bash
DB_MAX_OPEN_CONNECTIONS=50
DB_MAX_IDLE_CONNECTIONS=10
DB_MAX_LIFE_SECONDS=600
```

#### 场景 2：中等并发应用（100-1000 QPS）

```bash
DB_MAX_OPEN_CONNECTIONS=100
DB_MAX_IDLE_CONNECTIONS=25
DB_MAX_LIFE_SECONDS=600
```

#### 场景 3：高并发应用（> 1000 QPS）

```bash
DB_MAX_OPEN_CONNECTIONS=200
DB_MAX_IDLE_CONNECTIONS=50
DB_MAX_LIFE_SECONDS=1200
```

#### 场景 4：突发流量应用

```bash
DB_MAX_OPEN_CONNECTIONS=150
DB_MAX_IDLE_CONNECTIONS=40
DB_MAX_LIFE_SECONDS=900
```

---

## 监控功能

### 1. 统计信息

#### API 端点：`GET /api/v1/monitor/database/stats`

**返回数据：**

```json
{
  "max_open_connections": 100,
  "open_connections": 45,
  "in_use": 38,
  "idle": 7,
  "wait_count": 1250,
  "wait_duration": "2.5s",
  "max_idle_closed": 120,
  "max_lifetime_closed": 580,
  "max_idle_time_closed": 45,
  "utilization_rate": 38.0,
  "idle_rate": 15.56,
  "avg_wait_duration": "2ms"
}
```

**字段说明：**

| 字段 | 说明 | 正常范围 |
|------|------|---------|
| max_open_connections | 最大连接数 | 配置值 |
| open_connections | 当前打开连接 | < max_open_connections |
| in_use | 使用中连接 | 根据负载 |
| idle | 空闲连接 | 10-30% 的 open_connections |
| utilization_rate | 使用率（%） | < 80% |
| idle_rate | 空闲率（%） | 20-40% |
| wait_count | 等待次数 | 越少越好 |
| avg_wait_duration | 平均等待时间 | < 100ms |

### 2. 健康检查

#### API 端点：`GET /api/v1/monitor/database/health`

**返回数据：**

```json
{
  "healthy": true,
  "warnings": [
    "连接池运行正常"
  ]
}
```

**健康检查规则：**

1. ❌ 连接使用率 > 80% → 需要增加最大连接数
2. ❌ 等待次数 > 1000 → 连接不足
3. ❌ 平均等待时间 > 100ms → 性能瓶颈
4. ⚠️ 空闲率 > 50% → 可考虑降低空闲连接数
5. ⚠️ 生命周期关闭 > 10000 → 可延长 max_life_seconds

### 3. 指标导出

#### API 端点：`GET /api/v1/monitor/database/metrics`

**Prometheus 格式指标：**

```json
{
  "db_max_open_connections": 100,
  "db_open_connections": 45,
  "db_in_use_connections": 38,
  "db_idle_connections": 7,
  "db_wait_count_total": 1250,
  "db_wait_duration_seconds": 2.5,
  "db_max_idle_closed_total": 120,
  "db_max_lifetime_closed_total": 580,
  "db_max_idle_time_closed_total": 45,
  "db_utilization_rate": 38.0,
  "db_idle_rate": 15.56,
  "db_avg_wait_duration_seconds": 0.002
}
```

### 4. 配置建议

#### API 端点：`GET /api/v1/monitor/database/recommend`

**智能推荐示例：**

```json
{
  "status": "连接池配置合理，无需调整"
}
```

**或需要优化时：**

```json
{
  "max_open_connections": "当前: 100, 建议: 200 (使用率 85.00% 过高)",
  "performance": "平均等待 150ms，建议增加连接数或优化查询"
}
```

---

## 性能优化

### 1. 连接数计算公式

#### 基础公式

```
最大连接数 = (2 × CPU核心数) + 磁盘数量
```

**示例：**
- 4 核 CPU + 2 块磁盘 = 10 个连接（最小配置）
- 8 核 CPU + 4 块磁盘 = 20 个连接（推荐）

#### 实际调整

1. **监控使用率**：如果经常 > 80%，增加连接数
2. **等待时间**：如果 > 100ms，增加连接数
3. **业务特点**：长查询多 → 增加连接数；短查询多 → 适中即可

### 2. 优化检查清单

#### ✅ 配置优化

- [ ] 最大连接数是否合理（根据 CPU 核心数计算）
- [ ] 空闲连接数是否为最大连接数的 1/4 - 1/2
- [ ] 连接生命周期是否 < MySQL wait_timeout
- [ ] 是否配置了连接超时（ConnMaxIdleTime）

#### ✅ 监控告警

- [ ] 是否监控连接使用率
- [ ] 是否监控等待次数和时间
- [ ] 是否配置使用率 > 80% 的告警
- [ ] 是否配置等待时间 > 100ms 的告警

#### ✅ 应用层优化

- [ ] 是否使用了连接池（不手动管理连接）
- [ ] 查询是否添加了超时控制
- [ ] 事务是否及时提交或回滚
- [ ] 是否避免了长时间占用连接

### 3. 常见性能问题

#### 问题 1：连接频繁等待

**症状：**
- `wait_count` 持续增长
- `avg_wait_duration` > 100ms
- API 响应时间变慢

**原因：**
- 最大连接数设置过低
- 慢查询占用连接时间过长
- 事务未及时提交

**解决方案：**
```bash
# 1. 增加最大连接数
DB_MAX_OPEN_CONNECTIONS=200

# 2. 优化慢查询
# 检查慢查询日志，添加索引

# 3. 设置查询超时
db.SetMaxIdleTime(30 * time.Second)
```

#### 问题 2：空闲连接过多

**症状：**
- `idle_rate` > 50%
- `open_connections` 远大于 `in_use`
- 资源浪费

**原因：**
- 最大空闲连接数设置过高
- 流量波动大

**解决方案：**
```bash
# 减少最大空闲连接数
DB_MAX_IDLE_CONNECTIONS=15
```

#### 问题 3：连接频繁创建销毁

**症状：**
- `max_idle_closed` 持续增长
- CPU 使用率高
- 网络开销大

**原因：**
- 最大空闲连接数设置过低
- 连接生命周期过短

**解决方案：**
```bash
# 增加空闲连接数
DB_MAX_IDLE_CONNECTIONS=40

# 延长生命周期
DB_MAX_LIFE_SECONDS=1200
```

---

## 故障排查

### 1. 诊断流程

```
1. 访问健康检查接口
   GET /api/v1/monitor/database/health

2. 查看统计信息
   GET /api/v1/monitor/database/stats

3. 获取配置建议
   GET /api/v1/monitor/database/recommend

4. 根据建议调整配置

5. 重启应用并持续监控
```

### 2. 常见错误

#### 错误 1：too many connections

**错误信息：**
```
Error 1040: Too many connections
```

**原因：**
- 应用连接数超过 MySQL max_connections
- 连接未正常释放

**解决方案：**
```sql
-- 1. 临时增加 MySQL 最大连接数
SET GLOBAL max_connections = 500;

-- 2. 永久修改 my.cnf
[mysqld]
max_connections = 500

-- 3. 检查应用配置
# 确保 max_open_connections < MySQL max_connections
DB_MAX_OPEN_CONNECTIONS=400
```

#### 错误 2：connection refused

**错误信息：**
```
dial tcp 127.0.0.1:3306: connect: connection refused
```

**原因：**
- MySQL 服务未启动
- 网络不通
- 防火墙阻止

**排查步骤：**
```bash
# 1. 检查 MySQL 是否运行
systemctl status mysql

# 2. 检查端口是否监听
netstat -tlnp | grep 3306

# 3. 检查防火墙
iptables -L -n | grep 3306
```

#### 错误 3：database is locked

**错误信息：**
```
database is locked
```

**原因：**
- SQLite 并发写入冲突
- 事务未提交

**解决方案：**
```go
// 使用 WAL 模式
db.Exec("PRAGMA journal_mode=WAL")

// 设置超时
db.Exec("PRAGMA busy_timeout=5000")
```

---

## 最佳实践

### 1. 配置管理

#### ✅ 推荐做法

```bash
# 1. 使用环境变量管理配置
DB_MAX_OPEN_CONNECTIONS=100

# 2. 不同环境使用不同配置
# .env.development
DB_MAX_OPEN_CONNECTIONS=20

# .env.production
DB_MAX_OPEN_CONNECTIONS=200

# 3. 配置注释清晰
# 最大连接数：建议为 CPU 核心数的 2-3 倍
DB_MAX_OPEN_CONNECTIONS=100
```

#### ❌ 避免做法

```go
// ❌ 硬编码连接数
db.SetMaxOpenConns(100)

// ❌ 所有环境使用相同配置

// ❌ 没有注释说明配置原因
```

### 2. 监控告警

#### Prometheus 告警规则示例

```yaml
groups:
  - name: database_pool
    rules:
      # 连接使用率过高
      - alert: DatabasePoolHighUtilization
        expr: db_utilization_rate > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "数据库连接池使用率过高"
          description: "当前使用率 {{ $value }}%"
      
      # 等待时间过长
      - alert: DatabasePoolHighWaitTime
        expr: db_avg_wait_duration_seconds > 0.1
        for: 3m
        labels:
          severity: critical
        annotations:
          summary: "数据库连接等待时间过长"
          description: "平均等待时间 {{ $value }}s"
      
      # 连接不足
      - alert: DatabasePoolHighWaitCount
        expr: rate(db_wait_count_total[5m]) > 100
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "数据库连接等待频繁"
          description: "等待频率 {{ $value }}/s"
```

### 3. 代码规范

#### ✅ 正确的连接使用

```go
// 1. 使用 context 控制超时
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err := db.WithContext(ctx).Where("id = ?", id).First(&user).Error

// 2. 事务及时提交
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()  // 发生 panic 时回滚
    }
}()

if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}

tx.Commit()  // 成功时提交

// 3. 使用连接池，不手动管理连接
// ✅ 通过 GORM 自动管理
db.Find(&users)

// ❌ 不要手动获取连接
sqlDB, _ := db.DB()
conn, _ := sqlDB.Conn(ctx)  // 避免
```

#### ❌ 错误的连接使用

```go
// ❌ 忘记提交/回滚事务
tx := db.Begin()
tx.Create(&user)
// 忘记 Commit 或 Rollback，连接一直占用

// ❌ 长时间占用连接
db.Where("id = ?", id).First(&user)
time.Sleep(10 * time.Minute)  // 连接被长时间占用

// ❌ 不使用 context 超时控制
db.Raw("SELECT SLEEP(60)").Scan(&result)  // 可能长时间阻塞
```

### 4. 性能优化建议

#### 数据库层面

```sql
-- 1. 优化 MySQL 配置
[mysqld]
max_connections = 500
wait_timeout = 28800
interactive_timeout = 28800
max_connect_errors = 100

-- 2. 监控慢查询
slow_query_log = 1
slow_query_log_file = /var/log/mysql/slow.log
long_query_time = 1

-- 3. 优化索引
SHOW INDEX FROM users;
EXPLAIN SELECT * FROM users WHERE email = 'test@example.com';
```

#### 应用层面

```go
// 1. 使用连接池统计优化
stats := database.GetStats()
if stats.UtilizationRate > 80 {
    // 记录日志，考虑增加连接数
    logger.Warn("Database pool utilization high", "rate", stats.UtilizationRate)
}

// 2. 定期检查健康状态
go func() {
    ticker := time.NewTicker(5 * time.Minute)
    for range ticker.C {
        healthy, warnings := database.CheckHealth()
        if !healthy {
            for _, warning := range warnings {
                logger.Warn("Database pool warning", "msg", warning)
            }
        }
    }
}()

// 3. 批量操作使用事务
tx := db.Begin()
for _, user := range users {
    tx.Create(&user)
}
tx.Commit()
```

---

## 监控仪表板

### Grafana 配置示例

```json
{
  "dashboard": {
    "title": "GoHub-Service 数据库监控",
    "panels": [
      {
        "title": "连接池使用率",
        "targets": [
          {
            "expr": "db_utilization_rate",
            "legendFormat": "使用率 %"
          }
        ]
      },
      {
        "title": "连接数统计",
        "targets": [
          {
            "expr": "db_max_open_connections",
            "legendFormat": "最大连接数"
          },
          {
            "expr": "db_open_connections",
            "legendFormat": "打开连接数"
          },
          {
            "expr": "db_in_use_connections",
            "legendFormat": "使用中"
          },
          {
            "expr": "db_idle_connections",
            "legendFormat": "空闲"
          }
        ]
      },
      {
        "title": "等待统计",
        "targets": [
          {
            "expr": "rate(db_wait_count_total[5m])",
            "legendFormat": "等待频率/s"
          },
          {
            "expr": "db_avg_wait_duration_seconds",
            "legendFormat": "平均等待时间s"
          }
        ]
      }
    ]
  }
}
```

---

## 总结

### 优化成果

| 指标 | 优化前 | 优化后 | 改善 |
|------|-------|-------|------|
| 最大连接数 | 25 | 100 | +300% |
| 最大空闲数 | 100 | 25 | -75% |
| 连接监控 | ❌ | ✅ | 新增 |
| 健康检查 | ❌ | ✅ | 新增 |
| 配置建议 | ❌ | ✅ | 新增 |

### 核心价值

- ✅ **性能提升**：连接数从 25 提升到 100，满足高并发需求
- ✅ **资源优化**：空闲连接从 100 降到 25，节省资源
- ✅ **实时监控**：4 个监控接口，全方位掌握连接池状态
- ✅ **智能运维**：自动分析并推荐配置调整
- ✅ **Prometheus 集成**：标准指标格式，便于集成监控系统

---

## 参考资源

- [MySQL Connection Pool Best Practices](https://dev.mysql.com/doc/refman/8.0/en/connection-pooling.html)
- [Go database/sql Package](https://pkg.go.dev/database/sql)
- [GORM Performance](https://gorm.io/docs/performance.html)
- [Prometheus Monitoring](https://prometheus.io/docs/introduction/overview/)

---

**文档版本：** v1.0  
**最后更新：** 2026-01-03  
**维护者：** GoHub-Service Team
