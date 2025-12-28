# 日志管理文档

> 创建时间：2025年12月29日  
> 最后更新：2025年12月29日 v1.0  
> 状态：核心功能已完成

---

## 📋 目录

- [概述](#概述)
- [日志架构](#日志架构)
- [已实现功能](#已实现功能)
- [日志配置](#日志配置)
- [使用指南](#使用指南)
- [日志分级策略](#日志分级策略)
- [日志轮转和归档](#日志轮转和归档)
- [上下文日志追踪](#上下文日志追踪)
- [最佳实践](#最佳实践)
- [生产环境配置](#生产环境配置)
- [日志集中收集](#日志集中收集)

---

## 概述

GoHub-Service 使用 Uber Zap 作为高性能日志库，结合 Lumberjack 实现日志轮转。支持结构化日志、分级输出、自动轮转、上下文追踪等企业级特性。

### 核心特性

- ✅ **结构化日志**: JSON格式（生产）/ Console格式（开发）
- ✅ **分级输出**: Debug、Info、Warn、Error四个级别
- ✅ **自动轮转**: 按大小、按时间、按日期分割日志文件
- ✅ **上下文追踪**: RequestID、ErrorType、StackTrace自动记录
- ✅ **性能优化**: 零内存分配、异步写入
- ✅ **灵活配置**: 环境变量配置，支持热更新

---

## 日志架构

```
应用层
    ↓
Logger包 (pkg/logger/)
    ├── logger.go        - 全局Logger初始化
    ├── context.go       - 上下文日志追踪
    └── gorm_logger.go   - GORM日志适配器
    ↓
Zap核心
    ├── Encoder (JSON/Console)
    ├── WriteSyncer (Stdout/File)
    └── LogLevel (Debug/Info/Warn/Error)
    ↓
Lumberjack (日志轮转)
    ├── 按大小分割
    ├── 按时间归档
    └── 自动清理
    ↓
文件系统 (storage/logs/)
```

---

## 已实现功能

### ✅ 基础功能

1. **结构化日志输出**
   - 文件: `pkg/logger/logger.go`
   - 基于Zap的高性能结构化日志
   - JSON格式（生产环境）
   - Console格式（本地开发，带颜色高亮）

2. **日志分级输出**
   - 配置文件: `config/log.go`
   - 支持4个级别: debug, info, warn, error
   - 环境变量: `LOG_LEVEL`
   - 从低到高级别过滤

3. **日志轮转和归档**
   - 使用lumberjack.v2实现
   - 按大小轮转 (默认64MB)
   - 按时间归档 (默认保留30天)
   - 按日期分文件 (single/daily模式)
   - 自动清理过期日志 (默认保留5个文件)
   - 可选压缩功能

4. **上下文日志追踪**
   - 文件: `pkg/logger/context.go`
   - `LogErrorWithContext()`: 错误日志 + 上下文
   - `LogWithRequestID()`: 通用日志 + RequestID
   - 自动包含: RequestID, ErrorType, StackTrace

### ⏳ 待实现功能

- [ ] 日志集中收集 (ELK/Loki)
- [ ] 日志性能监控
- [ ] 慢查询日志独立输出
- [ ] 日志采样（高流量场景）

---

## 日志配置

### 配置文件

**位置**: `config/log.go`

```go
config.Add("log", func() map[string]interface{} {
    return map[string]interface{}{
        // 日志级别: debug, info, warn, error
        "level": config.Env("LOG_LEVEL", "debug"),
        
        // 日志类型: single (独立文件), daily (按日期)
        "type": config.Env("LOG_TYPE", "single"),
        
        // 日志文件路径
        "filename": config.Env("LOG_NAME", "storage/logs/logs.log"),
        
        // 每个日志文件最大尺寸 (MB)
        "max_size": config.Env("LOG_MAX_SIZE", 64),
        
        // 最多保存日志文件数 (0为不限)
        "max_backup": config.Env("LOG_MAX_BACKUP", 5),
        
        // 最多保存天数 (0为不删除)
        "max_age": config.Env("LOG_MAX_AGE", 30),
        
        // 是否压缩归档日志
        "compress": config.Env("LOG_COMPRESS", false),
    }
})
```

### 环境变量

**.env 配置示例**:

```bash
# 开发环境
LOG_LEVEL=debug
LOG_TYPE=single
LOG_NAME=storage/logs/logs.log
LOG_MAX_SIZE=64
LOG_MAX_BACKUP=5
LOG_MAX_AGE=30
LOG_COMPRESS=false

# 生产环境
LOG_LEVEL=error
LOG_TYPE=daily
LOG_NAME=storage/logs/app.log
LOG_MAX_SIZE=100
LOG_MAX_BACKUP=10
LOG_MAX_AGE=90
LOG_COMPRESS=true
```

---

## 使用指南

### 1. 基础日志

```go
import "GoHub-Service/pkg/logger"

// Debug 级别
logger.Debug("调试信息", 
    zap.String("user_id", "123"),
    zap.Int("count", 10),
)

// Info 级别
logger.Info("用户登录成功", 
    zap.String("username", "admin"),
    zap.String("ip", "192.168.1.1"),
)

// Warn 级别
logger.Warn("配置项缺失，使用默认值",
    zap.String("config_key", "max_connections"),
)

// Error 级别
logger.Error("数据库连接失败",
    zap.Error(err),
    zap.String("database", "mysql"),
)
```

### 2. 上下文日志（带RequestID）

```go
import (
    "GoHub-Service/pkg/logger"
    "github.com/gin-gonic/gin"
)

func SomeHandler(c *gin.Context) {
    // 通用上下文日志
    logger.LogWithRequestID(c, "info", "开始处理请求",
        zap.String("path", c.Request.URL.Path),
        zap.String("method", c.Request.Method),
    )
    
    // 错误上下文日志
    if err := doSomething(); err != nil {
        logger.LogErrorWithContext(c, err, "操作失败",
            zap.String("operation", "create_topic"),
        )
    }
}
```

### 3. Service层日志

```go
func (s *TopicService) CreateTopic(dto TopicCreateDTO) (*Topic, error) {
    logger.Info("创建话题",
        zap.String("title", dto.Title),
        zap.Uint64("user_id", dto.UserID),
    )
    
    topic, err := s.repository.Create(&models.Topic{
        Title: dto.Title,
        Body:  dto.Body,
    })
    
    if err != nil {
        logger.Error("话题创建失败",
            zap.Error(err),
            zap.String("title", dto.Title),
        )
        return nil, err
    }
    
    logger.Info("话题创建成功",
        zap.Uint64("topic_id", topic.ID),
    )
    return topic, nil
}
```

### 4. Controller层日志

```go
func (ctrl *TopicsController) Store(c *gin.Context) {
    var dto services.TopicCreateDTO
    if err := c.ShouldBindJSON(&dto); err != nil {
        logger.LogErrorWithContext(c, err, "请求参数验证失败")
        response.ValidationError(c, err)
        return
    }
    
    topic, err := ctrl.service.CreateTopic(dto)
    if err != nil {
        logger.LogErrorWithContext(c, err, "创建话题失败",
            zap.String("title", dto.Title),
        )
        response.ApiError(c, err)
        return
    }
    
    logger.LogWithRequestID(c, "info", "话题创建成功",
        zap.Uint64("topic_id", topic.ID),
    )
    response.Created(c, topic)
}
```

---

## 日志分级策略

### 级别说明

| 级别 | 用途 | 示例 | 生产环境 |
|-----|------|------|---------|
| **Debug** | 详细调试信息 | HTTP请求参数、SQL查询、函数调用 | ❌ 不推荐 |
| **Info** | 业务运行日志 | 用户登录、订单创建、任务完成 | ✅ 可选 |
| **Warn** | 需要关注的信息 | 配置项缺失、使用降级方案、性能警告 | ✅ 推荐 |
| **Error** | 错误信息 | 数据库错误、API调用失败、Panic | ✅ 必须 |

### 环境配置建议

#### 开发环境

```bash
LOG_LEVEL=debug
```

记录所有级别日志，便于调试。

#### 测试环境

```bash
LOG_LEVEL=info
```

记录业务日志和错误，用于验证功能。

#### 生产环境

```bash
LOG_LEVEL=error  # 或 warn
```

仅记录错误和警告，减少日志量，提升性能。

---

## 日志轮转和归档

### 按大小轮转

当单个日志文件达到 `max_size` (MB) 时自动切割：

```bash
LOG_MAX_SIZE=64  # 64MB切割一次
```

**示例文件**:
```
storage/logs/logs.log           # 当前文件
storage/logs/logs-2024-12-29.log   # 归档文件1
storage/logs/logs-2024-12-28.log   # 归档文件2
```

### 按时间归档

保留最近 `max_age` 天的日志：

```bash
LOG_MAX_AGE=30  # 保留30天
```

超过30天的日志文件自动删除。

### 按日期分文件

设置 `LOG_TYPE=daily` 每天生成一个日志文件：

```bash
LOG_TYPE=daily
LOG_NAME=storage/logs/app.log
```

**生成文件**:
```
storage/logs/app-2024-12-29.log
storage/logs/app-2024-12-28.log
storage/logs/app-2024-12-27.log
```

### 文件数量限制

保留最近 `max_backup` 个备份文件：

```bash
LOG_MAX_BACKUP=5  # 保留5个备份
```

配合 `max_age` 使用，哪个先达到就先删除。

### 日志压缩

归档日志自动压缩（.gz格式）：

```bash
LOG_COMPRESS=true
```

**压缩后文件**:
```
storage/logs/logs-2024-12-29.log.gz
```

节省磁盘空间，但查看需要解压。

---

## 上下文日志追踪

### LogErrorWithContext

**用途**: 记录错误日志，自动包含上下文信息

**函数签名**:
```go
func LogErrorWithContext(c *gin.Context, err error, message string, fields ...zap.Field)
```

**自动包含**:
- RequestID (X-Request-ID)
- ErrorType (AppError类型)
- ErrorCode (错误码)
- ErrorDetails (错误详情)
- StackTrace (堆栈追踪)

**使用示例**:
```go
if err := someOperation(); err != nil {
    logger.LogErrorWithContext(c, err, "操作失败",
        zap.String("operation", "update_user"),
        zap.Uint64("user_id", userID),
    )
}
```

**日志输出**:
```json
{
  "level": "ERROR",
  "time": "2024-12-29 15:30:45",
  "caller": "controllers/users_controller.go:42",
  "message": "操作失败",
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "error_type": "Database",
  "error_code": 5001,
  "error_details": {"query": "UPDATE users SET..."},
  "stack_trace": "goroutine 1 [running]:\n...",
  "operation": "update_user",
  "user_id": 123
}
```

### LogWithRequestID

**用途**: 记录通用日志，自动包含RequestID

**函数签名**:
```go
func LogWithRequestID(c *gin.Context, level string, message string, fields ...zap.Field)
```

**级别**: debug, info, warn, error

**使用示例**:
```go
logger.LogWithRequestID(c, "info", "用户操作记录",
    zap.String("action", "view_topic"),
    zap.Uint64("topic_id", topicID),
)
```

**日志输出**:
```json
{
  "level": "INFO",
  "time": "2024-12-29 15:30:45",
  "message": "用户操作记录",
  "request_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "action": "view_topic",
  "topic_id": 456
}
```

---

## 最佳实践

### 1. 日志级别选择

```go
// ✅ 正确：Debug用于详细调试
logger.Debug("SQL查询",
    zap.String("query", "SELECT * FROM users WHERE id = ?"),
    zap.Any("params", []interface{}{123}),
)

// ✅ 正确：Info用于业务流程
logger.Info("用户登录",
    zap.String("username", username),
    zap.String("ip", clientIP),
)

// ✅ 正确：Warn用于需要关注的情况
logger.Warn("缓存未命中，使用数据库查询",
    zap.String("cache_key", key),
)

// ✅ 正确：Error用于错误
logger.Error("数据库连接失败",
    zap.Error(err),
    zap.Int("retry_count", retries),
)
```

### 2. 结构化字段

```go
// ✅ 推荐：使用结构化字段
logger.Info("创建订单",
    zap.Uint64("order_id", orderID),
    zap.Float64("amount", 99.99),
    zap.String("status", "pending"),
)

// ❌ 不推荐：使用字符串拼接
logger.Info(fmt.Sprintf("创建订单 ID=%d amount=%.2f status=%s", 
    orderID, 99.99, "pending"))
```

### 3. 敏感信息脱敏

```go
// ✅ 正确：脱敏密码
logger.Info("用户注册",
    zap.String("email", email),
    zap.String("password", "******"),  // 脱敏
)

// ❌ 错误：记录明文密码
logger.Info("用户注册",
    zap.String("email", email),
    zap.String("password", password),  // 安全风险！
)
```

### 4. 错误上下文

```go
// ✅ 推荐：使用LogErrorWithContext
logger.LogErrorWithContext(c, err, "创建用户失败",
    zap.String("email", email),
)

// 可选：普通Error日志
logger.Error("创建用户失败",
    zap.Error(err),
    zap.String("email", email),
)
```

### 5. 避免过度日志

```go
// ❌ 不推荐：循环内打印日志
for _, item := range items {
    logger.Debug("处理项目", zap.Any("item", item))  // 可能产生大量日志
}

// ✅ 推荐：批量记录
logger.Debug("处理批量项目", 
    zap.Int("count", len(items)),
    zap.Any("first_item", items[0]),
)
```

---

## 生产环境配置

### 推荐配置

```bash
# 日志级别：仅记录错误
LOG_LEVEL=error

# 按日期分文件
LOG_TYPE=daily

# 日志路径
LOG_NAME=/var/log/gohub/app.log

# 每个文件最大100MB
LOG_MAX_SIZE=100

# 保留最近10个备份
LOG_MAX_BACKUP=10

# 保留90天
LOG_MAX_AGE=90

# 启用压缩
LOG_COMPRESS=true
```

### 文件权限

```bash
# 创建日志目录
mkdir -p /var/log/gohub
chown app:app /var/log/gohub
chmod 755 /var/log/gohub

# 设置日志文件权限
chmod 644 /var/log/gohub/*.log
```

### 日志监控

使用 `logrotate` 作为备用方案：

```bash
# /etc/logrotate.d/gohub
/var/log/gohub/*.log {
    daily
    rotate 30
    compress
    delaycompress
    notifempty
    create 0644 app app
    sharedscripts
    postrotate
        systemctl reload gohub || true
    endscript
}
```

---

## 日志集中收集

### ELK Stack 方案

#### 1. 安装 Filebeat

```bash
# 下载安装
curl -L -O https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-8.x.x-linux-x86_64.tar.gz
tar xzvf filebeat-8.x.x-linux-x86_64.tar.gz
```

#### 2. 配置 Filebeat

**filebeat.yml**:
```yaml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/gohub/*.log
  json.keys_under_root: true
  json.add_error_key: true
  fields:
    app: gohub
    env: production

output.elasticsearch:
  hosts: ["localhost:9200"]
  index: "gohub-logs-%{+yyyy.MM.dd}"

setup.template.name: "gohub-logs"
setup.template.pattern: "gohub-logs-*"
```

#### 3. 启动 Filebeat

```bash
./filebeat -e -c filebeat.yml
```

#### 4. Kibana 可视化

访问 `http://localhost:5601` 创建索引模式 `gohub-logs-*`

### Grafana Loki 方案

#### 1. 安装 Promtail

```bash
curl -LO https://github.com/grafana/loki/releases/download/v2.x.x/promtail-linux-amd64.zip
unzip promtail-linux-amd64.zip
```

#### 2. 配置 Promtail

**promtail-config.yml**:
```yaml
server:
  http_listen_port: 9080

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://localhost:3100/loki/api/v1/push

scrape_configs:
  - job_name: gohub
    static_configs:
      - targets:
          - localhost
        labels:
          app: gohub
          env: production
          __path__: /var/log/gohub/*.log
    pipeline_stages:
      - json:
          expressions:
            level: level
            time: time
            message: message
            request_id: request_id
```

#### 3. 启动 Promtail

```bash
./promtail -config.file=promtail-config.yml
```

#### 4. Grafana 查询

Loki 数据源查询示例：
```
{app="gohub", level="error"} |= "database"
```

---

## 故障排查

### 日志不输出

1. 检查日志级别配置
2. 检查日志文件权限
3. 检查磁盘空间

```bash
# 检查配置
echo $LOG_LEVEL

# 检查权限
ls -la storage/logs/

# 检查磁盘
df -h
```

### 日志文件过大

1. 降低日志级别（debug → info → error）
2. 减小 `max_size`
3. 减少 `max_backup` 数量
4. 启用压缩 `compress=true`

### 日志丢失

1. 检查 `max_age` 和 `max_backup` 配置
2. 检查磁盘空间是否充足
3. 检查是否有外部logrotate清理

---

## 相关文档

- [OPTIMIZATION_PLAN.md](../OPTIMIZATION_PLAN.md) - 优化计划
- [API_SECURITY.md](API_SECURITY.md) - API安全文档
- [PERFORMANCE_OPTIMIZATION.md](PERFORMANCE_OPTIMIZATION.md) - 性能优化

---

## 版本历史

| 版本 | 日期 | 变更内容 |
|-----|------|---------|
| v1.0 | 2025-12-29 | 初始版本，文档化现有日志系统 |

---

**维护者**: GoHub-Service Team  
**最后审核**: 2025-12-29
