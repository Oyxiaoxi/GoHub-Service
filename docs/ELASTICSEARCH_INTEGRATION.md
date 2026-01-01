# GoHub-Service Elasticsearch 集成方案

## 1. 概述

### 1.1 现状分析
- **当前搜索**: 基于MySQL全文搜索，响应时间 **150ms+**
- **性能瓶颈**: 大数据量下扫描行数多(5000+)，无相关性排序
- **目标**: 集成Elasticsearch实现 **15ms+** 高性能搜索

### 1.2 目标与价值
- ✅ 搜索响应时间降低 **90%** (150ms → 15ms)
- ✅ 支持相关性排序、模糊搜索、聚合统计
- ✅ 支持百万级文档，无性能衰减
- ✅ 提升用户搜索体验

---

## 2. 架构设计

### 2.1 系统架构

```
┌──────────────────────────────────────────────────────┐
│                   API 层                              │
│         (搜索请求处理)                               │
└────────────┬─────────────────────────────────────────┘
             │
    ┌────────┴────────┐
    ▼                 ▼
┌────────────┐   ┌──────────────┐
│  MySQL DB  │   │ Elasticsearch│
│  (主数据)  │   │  (搜索索引)  │
└────────────┘   └──────────────┘
    △                 △
    │                 │
    └──────┬──────────┘
           │
    ┌──────▼───────┐
    │ 消息队列/同步 │
    │  (Kafka/定时) │
    └───────────────┘
```

### 2.2 数据流

```
1. 创建/更新话题
   └─> MySQL 保存
   └─> 发送事件到Kafka/消息队列
   └─> Elasticsearch索引器消费事件
   └─> 更新ES索引

2. 搜索请求
   └─> API 接收搜索请求
   └─> Elasticsearch 查询
   └─> 返回结果(ID列表)
   └─> 从MySQL补充详细数据(可选)
```

---

## 3. 技术方案

### 3.1 依赖包

```bash
# Elasticsearch Go客户端
go get github.com/elastic/go-elasticsearch/v8

# JSON处理(已有)
go get encoding/json

# 日志(已有)
```

### 3.2 核心组件设计

#### 3.2.1 Elasticsearch客户端封装

```go
// pkg/elasticsearch/client.go
package elasticsearch

import (
    "context"
    "github.com/elastic/go-elasticsearch/v8"
)

type ESClient struct {
    client *elasticsearch.Client
}

// 初始化
func NewESClient(addresses []string) (*ESClient, error) {
    cfg := elasticsearch.Config{
        Addresses: addresses, // ["http://localhost:9200"]
    }
    client, err := elasticsearch.NewClient(cfg)
    if err != nil {
        return nil, err
    }
    
    return &ESClient{client: client}, nil
}

// 健康检查
func (ec *ESClient) Health(ctx context.Context) (bool, error) {
    res, err := ec.client.Info()
    if err != nil {
        return false, err
    }
    defer res.Body.Close()
    return res.StatusCode == 200, nil
}
```

#### 3.2.2 索引管理

```go
// pkg/elasticsearch/index.go
package elasticsearch

type IndexManager struct {
    client *ESClient
}

// 创建索引
func (im *IndexManager) CreateTopicIndex(ctx context.Context) error {
    mapping := map[string]interface{}{
        "settings": map[string]interface{}{
            "number_of_shards": 3,
            "number_of_replicas": 1,
            "analysis": map[string]interface{}{
                "analyzer": map[string]interface{}{
                    "ik_analyzer": map[string]interface{}{
                        "type": "custom",
                        "tokenizer": "ik_max_word",
                        "filter": []string{"lowercase"},
                    },
                },
            },
        },
        "mappings": map[string]interface{}{
            "properties": map[string]interface{}{
                "id": map[string]interface{}{
                    "type": "keyword",
                },
                "title": map[string]interface{}{
                    "type": "text",
                    "analyzer": "ik_analyzer",
                },
                "content": map[string]interface{}{
                    "type": "text",
                    "analyzer": "ik_analyzer",
                },
                "category_id": map[string]interface{}{
                    "type": "keyword",
                },
                "user_id": map[string]interface{}{
                    "type": "keyword",
                },
                "created_at": map[string]interface{}{
                    "type": "date",
                },
                "likes_count": map[string]interface{}{
                    "type": "integer",
                },
                "views_count": map[string]interface{}{
                    "type": "integer",
                },
                "status": map[string]interface{}{
                    "type": "keyword",
                },
            },
        },
    }
    
    // 创建索引逻辑...
    return nil
}
```

#### 3.2.3 搜索服务

```go
// app/services/search_service.go
package services

type SearchService struct {
    esClient *elasticsearch.ESClient
    repo     repository.TopicRepository
}

type SearchRequest struct {
    Query       string
    CategoryID  string
    Page        int
    PageSize    int
    SortBy      string // "relevance", "latest", "popular"
}

type SearchResult struct {
    ID        int64
    Title     string
    Content   string
    Score     float64
    CreatedAt time.Time
}

// 搜索话题
func (ss *SearchService) SearchTopics(ctx context.Context, req SearchRequest) ([]SearchResult, int64, error) {
    // 构建ES查询
    query := map[string]interface{}{
        "bool": map[string]interface{}{
            "must": []map[string]interface{}{
                {
                    "multi_match": map[string]interface{}{
                        "query": req.Query,
                        "fields": []string{"title^2", "content"}, // title权重2倍
                    },
                },
            },
            "filter": []map[string]interface{}{
                {
                    "term": map[string]interface{}{
                        "status": "published",
                    },
                },
            },
        },
    }
    
    // 添加分类过滤
    if req.CategoryID != "" {
        filter := query["bool"].(map[string]interface{})["filter"].([]map[string]interface{})
        filter = append(filter, map[string]interface{}{
            "term": map[string]interface{}{
                "category_id": req.CategoryID,
            },
        })
    }
    
    // 排序
    sort := []map[string]interface{}{
        {
            "_score": map[string]string{"order": "desc"},
        },
        {
            "created_at": map[string]string{"order": "desc"},
        },
    }
    
    // 执行搜索...
    results := []SearchResult{}
    return results, 0, nil
}
```

#### 3.2.4 索引同步

```go
// app/jobs/elasticsearch_sync.go
package jobs

type ESIndexSyncJob struct {
    esService *services.SearchService
    topicRepo repository.TopicRepository
}

// 定时全量同步(凌晨2点)
func (job *ESIndexSyncJob) FullSync(ctx context.Context) error {
    // 1. 获取所有发布的话题
    topics, err := job.topicRepo.GetPublishedTopics(ctx)
    if err != nil {
        return err
    }
    
    // 2. 批量索引到ES
    for _, topic := range topics {
        err := job.esService.IndexTopic(ctx, topic)
        if err != nil {
            // 记录错误，继续处理
            log.Errorf("Failed to index topic %d: %v", topic.ID, err)
        }
    }
    
    return nil
}

// 增量同步(消息队列)
func (job *ESIndexSyncJob) IncrementalSync(ctx context.Context, event TopicEvent) error {
    switch event.Type {
    case "created", "updated":
        return job.esService.IndexTopic(ctx, event.Topic)
    case "deleted":
        return job.esService.DeleteTopic(ctx, event.TopicID)
    }
    return nil
}
```

---

## 4. 实施步骤

### 4.1 第1周：环境准备与集成

#### Day 1-2：安装与配置
- [ ] Docker启动Elasticsearch 8.x
- [ ] 配置ES集群(推荐3节点)
- [ ] 配置IK中文分词器
- [ ] 验证ES健康状态

```bash
# Docker启动ES
docker run -d \
  -p 9200:9200 \
  -e "discovery.type=single-node" \
  -e "xpack.security.enabled=false" \
  docker.elastic.co/elasticsearch/elasticsearch:8.0.0
```

#### Day 3-4：Go客户端集成
- [ ] 添加elasticsearch依赖
- [ ] 实现ESClient封装
- [ ] 实现索引管理接口
- [ ] 单元测试

#### Day 5：搜索服务开发
- [ ] 实现SearchService
- [ ] 实现各类型搜索(标题、内容、组合)
- [ ] 集成排序与分页
- [ ] 性能测试

### 4.2 第2周：数据同步与优化

#### Day 1-2：索引同步机制
- [ ] 实现全量同步Job
- [ ] 实现增量同步(Kafka/消息队列)
- [ ] 处理同步失败场景
- [ ] 监控同步状态

#### Day 3-4：API集成
- [ ] 修改搜索API端点
- [ ] 实现请求参数映射
- [ ] 返回结果格式化
- [ ] 向前兼容性处理

#### Day 5：灰度发布
- [ ] 配置灰度规则(10% → 50% → 100%)
- [ ] 监控搜索质量指标
- [ ] 收集用户反馈
- [ ] 性能对标

### 4.3 第3周：性能优化与监控

#### Day 1-2：性能优化
- [ ] 调整shard/replica配置
- [ ] 优化查询模板
- [ ] 实现查询缓存
- [ ] 批量索引优化

#### Day 3-4：监控告警
- [ ] Prometheus集成
- [ ] 关键指标监控(搜索延迟、索引大小)
- [ ] 告警规则配置
- [ ] 日志采集

#### Day 5：文档与培训
- [ ] 技术文档编写
- [ ] 运维文档编写
- [ ] 团队培训
- [ ] 故障处理手册

---

## 5. 配置示例

### 5.1 环境配置

```yaml
# config/elasticsearch.go
type ElasticsearchConfig struct {
    Enabled   bool
    Addresses []string      // ["http://localhost:9200"]
    Username  string        // 可选
    Password  string        // 可选
    Timeout   time.Duration // 30s
    
    // 索引配置
    IndexName        string // "gohub-topics"
    Shards           int    // 3
    Replicas         int    // 1
    RefreshInterval  string // "1s"
    
    // 性能配置
    MaxRetries       int           // 3
    BulkSize         int           // 1000
    FlushInterval    time.Duration // 5s
    ConcurrentBulks  int           // 5
}
```

### 5.2 应用配置

```go
// bootstrap/elasticsearch.go
func init() {
    container.Singleton("elasticsearch", func(config *config.Config) (*elasticsearch.ESClient, error) {
        esConfig := config.GetElasticsearchConfig()
        
        client, err := elasticsearch.NewESClient(esConfig.Addresses)
        if err != nil {
            return nil, err
        }
        
        // 创建索引
        indexManager := elasticsearch.NewIndexManager(client)
        indexManager.CreateTopicIndex(context.Background())
        
        return client, nil
    })
}
```

---

## 6. API 集成示例

### 6.1 搜索端点修改

```go
// routes/search.go
func SearchTopics(c *gin.Context) {
    query := c.DefaultQuery("q", "")
    categoryID := c.DefaultQuery("category_id", "")
    page := c.DefaultQuery("page", "1")
    pageSize := c.DefaultQuery("page_size", "10")
    sortBy := c.DefaultQuery("sort", "relevance") // relevance, latest, popular
    
    req := services.SearchRequest{
        Query:      query,
        CategoryID: categoryID,
        Page:       page,
        PageSize:   pageSize,
        SortBy:     sortBy,
    }
    
    results, total, err := searchService.SearchTopics(c, req)
    if err != nil {
        response.Fail(c, "Search failed", err)
        return
    }
    
    response.Success(c, map[string]interface{}{
        "topics": results,
        "total":  total,
    })
}
```

### 6.2 请求/响应示例

```bash
# 搜索请求
GET /api/v1/search/topics?q=golang&category_id=1&sort=relevance&page=1&page_size=10

# 响应
{
  "code": 200,
  "data": {
    "topics": [
      {
        "id": 123,
        "title": "Golang 最佳实践",
        "content": "...",
        "score": 9.5,
        "created_at": "2026-01-01T10:00:00Z"
      }
    ],
    "total": 156
  }
}
```

---

## 7. 性能对标

### 7.1 性能指标对比

| 指标 | MySQL | Elasticsearch | 提升 |
|-----|--------|---------------|------|
| **搜索延迟** | 150ms | 15ms | **90%** ↓ |
| **吞吐量(QPS)** | 100 | 1000+ | **10倍** ↑ |
| **支持文档数** | 百万 | 数十亿 | **显著** ↑ |
| **相关性排序** | ❌ | ✅ | - |
| **模糊搜索** | ❌ | ✅ | - |
| **聚合统计** | ❌ | ✅ | - |

### 7.2 基准测试结果(预期)

```
搜索"golang"(100万文档):
  - MySQL: 150ms (全表扫描)
  - ES: 12-15ms (倒排索引)
  
搜索"高并发"(同步1000请求):
  - MySQL: 严重堵塞
  - ES: 平均响应 20ms
  
中文分词精度:
  - MySQL: 无法精准分词
  - ES + IK: 95%+ 准确率
```

---

## 8. 监控与告警

### 8.1 关键指标

| 指标 | 告警阈值 | 备注 |
|-----|---------|------|
| 搜索P99延迟 | > 50ms | 性能异常 |
| 索引失败率 | > 1% | 同步问题 |
| 磁盘使用率 | > 80% | 容量预警 |
| ES集群健康 | RED | 立即告警 |

### 8.2 Prometheus指标

```go
// 搜索延迟直方图
searchLatencyHistogram := prometheus.NewHistogramVec(
    prometheus.HistogramOpts{
        Name: "search_latency_seconds",
        Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1},
    },
    []string{"query_type"},
)

// 索引操作计数
indexOperationCounter := prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Name: "es_index_operations_total",
    },
    []string{"operation", "status"},
)
```

---

## 9. 故障处理

### 9.1 常见问题

| 问题 | 原因 | 解决方案 |
|-----|------|---------|
| 索引失败 | 网络问题/ES宕机 | 重试机制、死信队列 |
| 搜索结果不准 | 索引不同步 | 全量重新索引 |
| 性能下降 | 分片不均 | rebalance分片 |
| 磁盘爆满 | 索引过多 | 删除过期索引 |

### 9.2 应急预案

```go
// 快速回滚到MySQL搜索
type SearchService struct {
    useES bool // 配置开关
}

func (ss *SearchService) SearchTopics(ctx context.Context, req SearchRequest) ([]SearchResult, error) {
    if !ss.useES || ss.esClient == nil {
        return ss.searchFromMySQL(ctx, req) // 降级
    }
    
    results, err := ss.searchFromES(ctx, req)
    if err != nil {
        log.Warnf("ES search failed, fallback to MySQL: %v", err)
        return ss.searchFromMySQL(ctx, req) // 自动降级
    }
    
    return results, nil
}
```

---

## 10. 工作量评估

| 任务 | 工作量 | 人力 |
|-----|--------|------|
| 环境搭建 | 5h | 1人 |
| Go客户端集成 | 8h | 1人 |
| 搜索服务开发 | 16h | 1人 |
| 数据同步机制 | 12h | 1人 |
| API集成 | 6h | 1人 |
| 性能测试 | 8h | 1人 |
| 监控与告警 | 4h | 1人 |
| 文档与培训 | 5h | 1人 |
| **总计** | **64h** | **1-2人** |

---

## 11. 成功标准

✅ **功能目标**:
- 搜索响应时间 < 20ms
- 支持中文分词、模糊搜索、排序
- 支持百万级文档

✅ **质量目标**:
- 测试覆盖率 > 80%
- 性能基准测试通过
- 故障恢复时间 < 5分钟

✅ **运维目标**:
- 自动化索引管理
- 完整的监控告警
- 详细的运维文档

---

## 12. 参考资源

- [Elasticsearch官方文档](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html)
- [Go Elasticsearch客户端](https://github.com/elastic/go-elasticsearch)
- [IK分词器](https://github.com/medcl/elasticsearch-analysis-ik)
- [Kibana可视化](https://www.elastic.co/guide/en/kibana/current/index.html)

---

**方案版本**: v1.0  
**更新日期**: 2026年1月1日  
**负责人**: GoHub开发团队
