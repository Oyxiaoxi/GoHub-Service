# Elasticsearch 集成快速开始指南

## 1. 启动Elasticsearch集群

```bash
# 使用Docker Compose启动3节点ES集群 + Kibana
docker-compose -f docker-compose.elasticsearch.yml up -d

# 验证集群状态
curl http://localhost:9200/_cluster/health?pretty

# 查看Kibana
open http://localhost:5601
```

## 2. 验证安装

```bash
# 检查集群节点
curl http://localhost:9200/_nodes?pretty

# 检查索引
curl http://localhost:9200/_cat/indices?v
```

## 3. Go代码集成

### 3.1 添加依赖

```bash
go get github.com/elastic/go-elasticsearch/v8
```

### 3.2 配置Elasticsearch

编辑 `config/app.yaml`：

```yaml
elasticsearch:
  enabled: true
  addresses:
    - "http://localhost:9200"
  timeout: 30s
  index_name: "gohub-topics"
  shards: 3
  replicas: 1
```

### 3.3 初始化客户端

在 `bootstrap/elasticsearch.go` 中：

```go
package bootstrap

import (
    "context"
    "gohub/pkg/elasticsearch"
    "gohub/pkg/logger"
)

func init() {
    container.Singleton("elasticsearch", func(config *config.Config) (*elasticsearch.Client, error) {
        esConfig := config.GetElasticsearchConfig()
        
        if !esConfig.Enabled {
            logger.Warn("Elasticsearch is disabled")
            return nil, nil
        }
        
        client, err := elasticsearch.NewClient(esConfig.Addresses)
        if err != nil {
            return nil, err
        }
        
        // 创建索引
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        
        indexManager := elasticsearch.NewIndexManager(client)
        if err := indexManager.CreateTopicIndex(ctx); err != nil {
            logger.Warn("Failed to create topic index: %v", err)
        }
        
        logger.Info("Elasticsearch client initialized")
        return client, nil
    })
}
```

## 4. 测试集成

```bash
# 运行测试脚本
go run scripts/test_elasticsearch.go

# 预期输出
# ✓ ES Health: true
# ✓ Topic index created successfully
# ✓ Test data indexed successfully
# ✓ Index refreshed
# ✓ Search results (total: 1):
#   - [1] Golang 最佳实践 (score: 5.23)
# ✓ All tests passed!
```

## 5. 在API中集成搜索

### 5.1 创建搜索路由

```go
// routes/search.go
func SearchTopics(c *gin.Context) {
    query := c.DefaultQuery("q", "")
    categoryID := c.DefaultQuery("category_id", "")
    page := c.GetInt("page", 1)
    pageSize := c.GetInt("page_size", 10)
    sortBy := c.DefaultQuery("sort", "relevance")
    
    if query == "" {
        c.JSON(400, gin.H{"error": "Query is required"})
        return
    }
    
    req := elasticsearch.SearchRequest{
        Query:      query,
        CategoryID: categoryID,
        Page:       page,
        PageSize:   pageSize,
        SortBy:     sortBy,
    }
    
    searchService := container.Make("elasticsearch.search").(*elasticsearch.SearchService)
    results, total, err := searchService.SearchTopics(c, req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{
        "topics": results,
        "total":  total,
    })
}
```

### 5.2 测试搜索API

```bash
# 搜索话题
curl "http://localhost:8000/api/v1/search/topics?q=golang&page=1&page_size=10"

# 带分类过滤
curl "http://localhost:8000/api/v1/search/topics?q=golang&category_id=1&sort=latest"

# 预期响应
{
  "topics": [
    {
      "id": 1,
      "title": "Golang 最佳实践",
      "content": "...",
      "score": 9.5,
      "created_at": "2026-01-01T10:00:00Z"
    }
  ],
  "total": 156
}
```

## 6. 数据同步

### 6.1 话题创建时同步

在 `services/topic_service.go` 中：

```go
func (s *TopicService) CreateTopic(ctx context.Context, req CreateTopicRequest) (*Topic, error) {
    // 1. 保存到MySQL
    topic, err := s.repo.Create(ctx, req)
    if err != nil {
        return nil, err
    }
    
    // 2. 异步索引到ES
    go func() {
        esClient := container.Make("elasticsearch").(*elasticsearch.Client)
        esClient.IndexTopic(context.Background(), map[string]interface{}{
            "id":             topic.ID,
            "title":          topic.Title,
            "content":        topic.Content,
            "category_id":    topic.CategoryID,
            "user_id":        topic.UserID,
            "created_at":     topic.CreatedAt.Format(time.RFC3339),
            "likes_count":    0,
            "views_count":    0,
            "comments_count": 0,
            "status":         "published",
        })
    }()
    
    return topic, nil
}
```

### 6.2 话题更新时同步

```go
func (s *TopicService) UpdateTopic(ctx context.Context, id int64, req UpdateTopicRequest) error {
    // 1. 更新MySQL
    if err := s.repo.Update(ctx, id, req); err != nil {
        return err
    }
    
    // 2. 异步更新ES
    go func() {
        topic, _ := s.repo.GetByID(ctx, id)
        esClient := container.Make("elasticsearch").(*elasticsearch.Client)
        esClient.IndexTopic(context.Background(), topicToESDoc(topic))
    }()
    
    return nil
}
```

## 7. 监控与维护

### 7.1 监控集群健康

```bash
# 实时监控
watch -n 1 'curl -s http://localhost:9200/_cluster/health?pretty'

# 获取节点信息
curl http://localhost:9200/_nodes/stats?pretty
```

### 7.2 索引管理

```bash
# 获取索引大小
curl http://localhost:9200/_cat/indices/gohub-topics?v

# 强制刷新索引
curl -X POST http://localhost:9200/gohub-topics/_refresh

# 获取索引映射
curl http://localhost:9200/gohub-topics/_mapping?pretty

# 删除索引
curl -X DELETE http://localhost:9200/gohub-topics
```

## 8. 常见问题

### Q: Elasticsearch启动失败？
**A**: 检查内存设置，ES需要最少512MB。修改 `docker-compose.elasticsearch.yml` 中的 `ES_JAVA_OPTS`。

### Q: 搜索结果为空？
**A**: 确保：
1. 索引已创建：`curl http://localhost:9200/_cat/indices?v`
2. 数据已索引：`curl http://localhost:9200/gohub-topics/_count`
3. 刷新索引：`curl -X POST http://localhost:9200/gohub-topics/_refresh`

### Q: 性能很慢？
**A**: 检查：
1. ES集群健康状态：`curl http://localhost:9200/_cluster/health?pretty`
2. 垃圾回收情况：`curl http://localhost:9200/_nodes/stats/jvm?pretty`
3. 增加分片数或replicas

## 9. 下一步

- [ ] 配置IK中文分词器
- [ ] 实现自动同步机制（Kafka/消息队列）
- [ ] 集成Prometheus监控
- [ ] 性能压力测试
- [ ] 灰度发布计划

## 10. 相关文档

- [Elasticsearch集成方案](./ELASTICSEARCH_INTEGRATION.md)
- [性能优化指南](./PERFORMANCE.md)
- [故障排查手册](./TROUBLESHOOTING.md)
