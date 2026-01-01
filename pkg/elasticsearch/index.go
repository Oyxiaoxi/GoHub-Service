package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
)

// IndexManager 索引管理器
type IndexManager struct {
	client *Client
}

// NewIndexManager 创建索引管理器
func NewIndexManager(client *Client) *IndexManager {
	return &IndexManager{client: client}
}

// TopicMapping 话题索引的mapping定义
func (im *IndexManager) TopicMapping() map[string]interface{} {
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   3,
			"number_of_replicas": 1,
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"ik_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "ik_max_word",
						"filter":    []string{"lowercase"},
					},
				},
			},
			"refresh_interval": "1s",
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type": "keyword",
				},
				"title": map[string]interface{}{
					"type":     "text",
					"analyzer": "ik_analyzer",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"content": map[string]interface{}{
					"type":     "text",
					"analyzer": "ik_analyzer",
				},
				"description": map[string]interface{}{
					"type":     "text",
					"analyzer": "ik_analyzer",
				},
				"category_id": map[string]interface{}{
					"type": "keyword",
				},
				"user_id": map[string]interface{}{
					"type": "keyword",
				},
				"created_at": map[string]interface{}{
					"type":   "date",
					"format": "yyyy-MM-dd'T'HH:mm:ssZ||yyyy-MM-dd HH:mm:ss||epoch_millis",
				},
				"updated_at": map[string]interface{}{
					"type":   "date",
					"format": "yyyy-MM-dd'T'HH:mm:ssZ||yyyy-MM-dd HH:mm:ss||epoch_millis",
				},
				"likes_count": map[string]interface{}{
					"type": "integer",
				},
				"views_count": map[string]interface{}{
					"type": "integer",
				},
				"comments_count": map[string]interface{}{
					"type": "integer",
				},
				"status": map[string]interface{}{
					"type": "keyword",
				},
			},
		},
	}
}

// CreateTopicIndex 创建话题索引
func (im *IndexManager) CreateTopicIndex(ctx context.Context) error {
	mapping := im.TopicMapping()
	return im.client.CreateIndex(ctx, "gohub-topics", mapping)
}

// DeleteTopicIndex 删除话题索引
func (im *IndexManager) DeleteTopicIndex(ctx context.Context) error {
	return im.client.DeleteIndex(ctx, "gohub-topics")
}

// ReindexTopicIndex 重建话题索引
func (im *IndexManager) ReindexTopicIndex(ctx context.Context) error {
	// 1. 删除旧索引
	if err := im.DeleteTopicIndex(ctx); err != nil {
		return fmt.Errorf("failed to delete old index: %w", err)
	}

	// 2. 创建新索引
	if err := im.CreateTopicIndex(ctx); err != nil {
		return fmt.Errorf("failed to create new index: %w", err)
	}

	return nil
}

// IndexExists 检查索引是否存在
func (im *IndexManager) IndexExists(ctx context.Context, indexName string) (bool, error) {
	res, err := im.client.client.Indices.Exists(
		[]string{indexName},
		im.client.client.Indices.Exists.WithContext(ctx),
	)
	if err != nil {
		return false, fmt.Errorf("failed to check index existence: %w", err)
	}
	defer res.Body.Close()

	return !res.IsError(), nil
}

// GetIndexStats 获取索引统计信息
func (im *IndexManager) GetIndexStats(ctx context.Context, indexName string) (map[string]interface{}, error) {
	res, err := im.client.client.Indices.Stats(
		im.client.client.Indices.Stats.WithIndex(indexName),
		im.client.client.Indices.Stats.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get index stats: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("es stats error: %s", res.String())
	}

	var stats map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode stats: %w", err)
	}

	return stats, nil
}

// RefreshIndex 刷新索引
func (im *IndexManager) RefreshIndex(ctx context.Context, indexName string) error {
	res, err := im.client.client.Indices.Refresh(
		im.client.client.Indices.Refresh.WithIndex(indexName),
		im.client.client.Indices.Refresh.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("failed to refresh index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("es refresh error: %s", res.String())
	}

	return nil
}
