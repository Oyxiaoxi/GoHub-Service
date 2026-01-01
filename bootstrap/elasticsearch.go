package bootstrap

import (
	"context"
	"time"

	"go.uber.org/zap"

	"GoHub-Service/pkg/config"
	"GoHub-Service/pkg/elasticsearch"
	"GoHub-Service/pkg/logger"
)

// LoadElasticsearch 加载Elasticsearch客户端
func LoadElasticsearch() {
	esConfig := config.GetBool("elasticsearch.enabled", false)
	if !esConfig {
		logger.Info("Elasticsearch is disabled")
		return
	}

	// 使用默认地址
	addresses := []string{"http://localhost:9200"}
	
	timeout := config.GetInt("elasticsearch.timeout", 30000)

	client, err := elasticsearch.NewClient(addresses)
	if err != nil {
		logger.Warn("Failed to create Elasticsearch client", zap.Error(err))
		return
	}

	// 检查健康状态
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()

	if ok, err := client.Health(ctx); !ok {
		logger.Warn("Elasticsearch health check failed", zap.Error(err))
		return
	}

	// 创建索引
	indexManager := elasticsearch.NewIndexManager(client)

	if err := indexManager.CreateTopicIndex(ctx); err != nil {
		logger.Warn("Failed to create Elasticsearch index", zap.Error(err))
		return
	}

	logger.Info("Elasticsearch initialized successfully")
}
