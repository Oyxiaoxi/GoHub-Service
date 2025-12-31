# MySQL 慢查询日志分析指南

## 概述

本文档介绍如何使用 GoHub-Service 内置的慢查询日志分析工具，帮助识别和优化数据库性能瓶颈。

## 目录

- [快速开始](#快速开始)
- [功能特性](#功能特性)
- [使用方法](#使用方法)
- [分析报告](#分析报告)
- [优化建议](#优化建议)
- [最佳实践](#最佳实践)

---

## 快速开始

### 1. 启用MySQL慢查询日志

在MySQL配置文件 `/etc/mysql/my.cnf` 中添加：

```ini
[mysqld]
# 启用慢查询日志
slow_query_log = 1

# 慢查询日志文件路径
slow_query_log_file = /var/log/mysql/slow-query.log

# 慢查询阈值（秒）
long_query_time = 0.5

# 记录未使用索引的查询
log_queries_not_using_indexes = 1

# 限制每分钟记录的未使用索引的查询数量
log_throttle_queries_not_using_indexes = 10
```

重启MySQL使配置生效：
```bash
sudo systemctl restart mysql
```

### 2. 运行分析工具

```bash
# 基础分析
./main slowlog -f /var/log/mysql/slow-query.log

# 显示Top 20慢查询
./main slowlog -f /var/log/mysql/slow-query.log -t 20

# 只显示超过1秒的查询
./main slowlog -f /var/log/mysql/slow-query.log -s 1.0

# 显示详细信息
./main slowlog -f /var/log/mysql/slow-query.log -d

# 导出为CSV
./main slowlog -f /var/log/mysql/slow-query.log -e csv

# 导出为JSON
./main slowlog -f /var/log/mysql/slow-query.log -e json
```

---

## 功能特性

### 1. 自动化分析

- ✅ 自动解析MySQL慢查询日志文件
- ✅ 提取查询语句、执行时间、扫描行数等关键指标
- ✅ 对查询语句进行标准化（移除具体参数）
- ✅ 按总耗时自动排序

### 2. 统计信息

- **总体统计**: 总查询数、总耗时、平均耗时、扫描行数
- **单条查询**: 执行次数、总耗时、平均/最大/最小耗时、扫描/返回行数
- **百分比分析**: 计算每条查询占总耗时的百分比

### 3. 智能优化建议

系统会根据以下指标自动生成优化建议：

- **扫描行数**: 
  - >10K 行: 高优先级，建议添加索引
  - >1K 行: 中优先级，检查索引使用情况

- **选择性比例**:
  - 扫描与返回行数比例 >100:1: 建议优化索引选择性

- **执行时间**:
  - >1秒: 紧急优化
  - >500ms: 需要优化

- **查询频率**:
  - >1000次: 建议添加缓存

- **SQL语法问题**:
  - 使用 `SELECT *`
  - 使用前导通配符 `LIKE '%...`
  - 使用多个 `OR` 而非 `IN`
  - 缺少 `LIMIT` 限制

### 4. 导出功能

支持将分析结果导出为：
- **CSV格式**: 方便在Excel中进一步分析
- **JSON格式**: 便于程序化处理和集成

---

## 使用方法

### 命令行参数

| 参数 | 短选项 | 说明 | 默认值 | 示例 |
|-----|-------|------|-------|------|
| `--file` | `-f` | 慢查询日志文件路径 | 必填 | `-f /var/log/mysql/slow.log` |
| `--top` | `-t` | 显示Top N慢查询 | 10 | `-t 20` |
| `--threshold` | `-s` | 过滤阈值（秒） | 0 | `-s 1.0` |
| `--export` | `-e` | 导出格式（csv/json） | - | `-e csv` |
| `--detail` | `-d` | 显示完整SQL | false | `-d` |

### 使用场景

#### 场景1: 日常性能监控

每天定时分析慢查询日志：

```bash
#!/bin/bash
# /usr/local/bin/daily-slowlog-check.sh

LOG_FILE="/var/log/mysql/slow-query.log"
EXPORT_DIR="/var/log/gohub/slowlog-reports"
DATE=$(date +%Y%m%d)

mkdir -p $EXPORT_DIR

# 分析并导出报告
/path/to/main slowlog -f $LOG_FILE -t 20 -e csv > "$EXPORT_DIR/report_$DATE.txt"

# 如果有超过2秒的查询，发送告警
CRITICAL=$(grep "Very slow query (>1s)" "$EXPORT_DIR/report_$DATE.txt" | wc -l)
if [ $CRITICAL -gt 0 ]; then
    echo "发现 $CRITICAL 个严重慢查询" | mail -s "慢查询告警" admin@example.com
fi
```

设置定时任务：
```bash
# 每天凌晨1点执行
0 1 * * * /usr/local/bin/daily-slowlog-check.sh
```

#### 场景2: 问题排查

当系统性能下降时：

```bash
# 1. 先看最慢的查询
./main slowlog -f /var/log/mysql/slow-query.log -t 5 -d

# 2. 查看所有超过1秒的查询
./main slowlog -f /var/log/mysql/slow-query.log -s 1.0 -t 50

# 3. 导出详细报告给DBA
./main slowlog -f /var/log/mysql/slow-query.log -e json -t 100
```

#### 场景3: 上线前性能测试

```bash
# 1. 清空慢查询日志
mysql -e "SET GLOBAL slow_query_log = 'OFF'; TRUNCATE TABLE mysql.slow_log; SET GLOBAL slow_query_log = 'ON';"

# 2. 执行压力测试
ab -n 10000 -c 100 http://localhost:3000/api/v1/topics

# 3. 分析慢查询
./main slowlog -f /var/log/mysql/slow-query.log -d -e csv
```

---

## 分析报告

### 报告示例

```
Analyzing slow query log: /var/log/mysql/slow-query.log
================================================================================

📊 Summary Statistics
--------------------------------------------------------------------------------
Total slow queries: 523
Threshold: >= 0.50 seconds
Total execution time: 342.56 seconds
Average execution time: 0.6550 seconds
Total rows examined: 2.5M
Total rows sent: 15.2K
Average rows examined per query: 4781

🔍 Top 10 Slow Queries (by total time)
================================================================================

#1 Query Pattern:
--------------------------------------------------------------------------------
Count: 89 times
Total Time: 127.4523 seconds (37.2% of total)
Avg Time: 1.4321 seconds
Max Time: 3.5612 seconds
Min Time: 0.5234 seconds
Avg Rows Examined: 35421
Avg Rows Sent: 20

Query:
SELECT topics.*, users.name, users.avatar, categories.name AS category_name
FROM topics
LEFT JOIN users ON users.id = topics.user_id
LEFT JOIN categories ON categories.id = topics.category_id
WHERE topics.category_id = ?
ORDER BY topics.created_at DESC
LIMIT ?

💡 Optimization Suggestions:
  • High rows examined (>10K). Consider adding indexes.
  • Poor selectivity (1771:1 ratio). Add more selective indexes.
  • Very slow query (>1s). Urgent optimization needed.
  • High frequency query (89 times). Caching recommended.

#2 Query Pattern:
--------------------------------------------------------------------------------
Count: 156 times
Total Time: 89.2341 seconds (26.0% of total)
Avg Time: 0.5720 seconds
Max Time: 1.2345 seconds
Min Time: 0.3120 seconds
Avg Rows Examined: 8945
Avg Rows Sent: 15

Query:
SELECT comments.*, users.name, users.avatar
FROM comments
LEFT JOIN users ON users.id = comments.user_id
WHERE comments.topic_id = ?
ORDER BY comments.created_at DESC
LIMIT ?

💡 Optimization Suggestions:
  • Moderate rows examined (>1K). Check if indexes are being used.
  • High frequency query (156 times). Caching recommended.
  • Slow query (>500ms). Consider optimization.

...

Analysis completed!
Results exported to: slowlog_analysis_20251231_143025.csv
```

### CSV导出格式

```csv
Rank,Query,Count,Total Time,Avg Time,Max Time,Min Time,Avg Rows Examined,Avg Rows Sent
1,"SELECT topics.* FROM topics WHERE category_id = ? ORDER BY created_at DESC",89,127.4523,1.4321,3.5612,0.5234,35421,20
2,"SELECT comments.* FROM comments WHERE topic_id = ? ORDER BY created_at DESC",156,89.2341,0.5720,1.2345,0.3120,8945,15
```

### JSON导出格式

```json
[
  {
    "rank": 1,
    "query": "SELECT topics.* FROM topics WHERE category_id = ? ORDER BY created_at DESC",
    "count": 89,
    "total_time": 127.4523,
    "avg_time": 1.4321,
    "max_time": 3.5612,
    "min_time": 0.5234,
    "avg_rows_examined": 35421,
    "avg_rows_sent": 20
  }
]
```

---

## 优化建议

### 1. 添加索引

当看到 "High rows examined" 建议时：

```sql
-- 查看当前索引
SHOW INDEX FROM topics;

-- 添加复合索引
CREATE INDEX idx_topics_category_created ON topics(category_id, created_at DESC);

-- 验证索引效果
EXPLAIN SELECT * FROM topics 
WHERE category_id = 1 
ORDER BY created_at DESC 
LIMIT 20;
```

### 2. 优化查询选择性

当扫描/返回比例过高时：

```sql
-- 不好: 扫描10000行返回10行
SELECT * FROM topics WHERE created_at > '2025-01-01';

-- 更好: 添加更多过滤条件
SELECT * FROM topics 
WHERE created_at > '2025-01-01' 
  AND category_id = 1 
  AND status = 'published';

-- 最好: 使用覆盖索引
CREATE INDEX idx_topics_filter 
ON topics(category_id, status, created_at) 
INCLUDE (id, title, user_id);
```

### 3. 避免全表扫描

```sql
-- 不好: 前导通配符导致索引失效
SELECT * FROM topics WHERE title LIKE '%keyword%';

-- 更好: 使用全文搜索
ALTER TABLE topics ADD FULLTEXT INDEX ft_title (title);
SELECT * FROM topics WHERE MATCH(title) AGAINST('keyword');

-- 最好: 使用Elasticsearch
```

### 4. 使用缓存

当查询频率很高时：

```go
// 在应用层添加缓存
func (s *TopicService) GetTopicsByCategory(categoryID string) ([]Topic, error) {
    // 尝试从缓存获取
    cacheKey := fmt.Sprintf("topics:category:%s", categoryID)
    if cached, err := cache.Get(cacheKey); err == nil {
        return cached, nil
    }

    // 从数据库查询
    topics, err := s.repo.GetByCategory(categoryID)
    if err != nil {
        return nil, err
    }

    // 写入缓存（5分钟过期）
    cache.Set(cacheKey, topics, 5*time.Minute)

    return topics, nil
}
```

### 5. 优化JOIN查询

```sql
-- 不好: 多次JOIN导致笛卡尔积
SELECT * FROM topics t
LEFT JOIN users u1 ON t.user_id = u1.id
LEFT JOIN comments c ON t.id = c.topic_id
LEFT JOIN users u2 ON c.user_id = u2.id;

-- 更好: 分开查询 + 应用层聚合
SELECT * FROM topics WHERE id IN (1,2,3);
SELECT * FROM comments WHERE topic_id IN (1,2,3);
```

---

## 最佳实践

### 1. 慢查询阈值设置

不同场景的建议阈值：

| 场景 | 阈值 | 说明 |
|-----|------|------|
| 开发环境 | 0.1秒 | 尽早发现问题 |
| 测试环境 | 0.2秒 | 接近生产环境 |
| 生产环境 | 0.5秒 | 关注严重问题 |
| 高性能API | 0.1秒 | 严格要求 |

### 2. 日志轮转

慢查询日志会快速增长，需要配置日志轮转：

```bash
# /etc/logrotate.d/mysql-slow
/var/log/mysql/slow-query.log {
    daily
    rotate 30
    missingok
    compress
    delaycompress
    notifempty
    create 640 mysql mysql
    postrotate
        test -x /usr/bin/mysqladmin && \
        /usr/bin/mysqladmin flush-logs
    endscript
}
```

### 3. 定期清理

避免慢查询日志占用过多磁盘空间：

```bash
#!/bin/bash
# 每周清理一次30天前的日志
find /var/log/mysql/ -name "slow-query.log.*" -mtime +30 -delete
```

### 4. 监控告警

结合监控系统：

```yaml
# Prometheus监控规则
groups:
  - name: mysql_slow_queries
    rules:
      - alert: HighSlowQueryRate
        expr: rate(mysql_global_status_slow_queries[5m]) > 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "MySQL慢查询率过高"
          description: "最近5分钟慢查询率超过10次/秒"
```

### 5. 持续优化流程

1. **每日分析**: 定时生成慢查询报告
2. **周报总结**: 汇总本周Top慢查询
3. **优化验证**: 添加索引后对比效果
4. **文档记录**: 记录优化措施和效果

---

## 常见问题

### Q1: 日志文件太大，分析很慢怎么办？

A: 可以先截取最近的部分日志：

```bash
# 只分析最后10000行
tail -n 10000 /var/log/mysql/slow-query.log > /tmp/recent-slow.log
./main slowlog -f /tmp/recent-slow.log
```

### Q2: 如何只分析特定时间段的日志？

A: 使用时间范围截取：

```bash
# 分析今天的慢查询
sed -n '/^# Time: 2025-12-31/,/^# Time: 2026-01-01/p' \
  /var/log/mysql/slow-query.log > /tmp/today-slow.log
./main slowlog -f /tmp/today-slow.log
```

### Q3: 分析工具不识别日志格式？

A: 确保MySQL慢查询日志格式正确：

```sql
-- 检查日志格式
SHOW VARIABLES LIKE 'log_output';

-- 应该是 'FILE' 而不是 'TABLE'
SET GLOBAL log_output = 'FILE';
```

### Q4: 如何集成到CI/CD流程？

A: 在部署前运行性能测试和慢查询分析：

```yaml
# .github/workflows/performance-test.yml
- name: Performance Test
  run: |
    # 启动测试环境
    docker-compose up -d
    
    # 运行压力测试
    ab -n 1000 -c 10 http://localhost:3000/api/v1/topics
    
    # 分析慢查询
    docker exec mysql cat /var/log/mysql/slow-query.log > slow.log
    ./main slowlog -f slow.log -s 0.5
    
    # 如果有超过1秒的查询，测试失败
    if ./main slowlog -f slow.log -s 1.0 | grep -q "Total slow queries: [1-9]"; then
      echo "发现严重慢查询！"
      exit 1
    fi
```

---

## 相关工具

### MySQL官方工具

- **mysqldumpslow**: MySQL内置的慢查询分析工具
  ```bash
  mysqldumpslow -s t -t 10 /var/log/mysql/slow-query.log
  ```

- **pt-query-digest**: Percona Toolkit的慢查询分析工具
  ```bash
  pt-query-digest /var/log/mysql/slow-query.log
  ```

### 对比优势

| 特性 | GoHub slowlog | mysqldumpslow | pt-query-digest |
|-----|--------------|---------------|-----------------|
| **安装** | 内置，无需安装 | 内置 | 需要安装 |
| **性能** | 快速 | 较慢 | 中等 |
| **优化建议** | ✅ 智能建议 | ❌ 无 | ⚠️ 基础建议 |
| **导出格式** | CSV, JSON | 文本 | HTML, JSON |
| **易用性** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| **可定制** | ✅ 源码可改 | ❌ 固定 | ⚠️ 插件 |

---

## 总结

慢查询日志分析是数据库性能优化的重要手段。通过：

1. ✅ 定期分析慢查询日志
2. ✅ 根据建议添加索引
3. ✅ 优化高频慢查询
4. ✅ 使用缓存减少数据库压力
5. ✅ 持续监控和改进

可以显著提升系统性能，改善用户体验。

---

**版本**: 1.0.0  
**更新时间**: 2025-12-31  
**维护者**: GoHub Development Team
