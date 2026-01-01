package elasticsearch

import (
	"context"
	"fmt"
	"log"
	"time"

	"GoHub-Service/app/models/topic"
	"GoHub-Service/pkg/database"
)

// SyncService 数据同步服务 - MySQL → Elasticsearch
type SyncService struct {
	client        *Client
	indexManager  *IndexManager
	searchService *SearchService
}

// NewSyncService 创建数据同步服务
func NewSyncService(client *Client) *SyncService {
	return &SyncService{
		client:        client,
		indexManager:  NewIndexManager(client),
		searchService: NewSearchService(client),
	}
}

// SyncAllTopics 同步所有话题到Elasticsearch
func (s *SyncService) SyncAllTopics(ctx context.Context, batchSize int) error {
	log.Println("Starting full topic sync to Elasticsearch...")

	// 第1步: 创建新索引
	if err := s.indexManager.CreateTopicIndex(ctx); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	// 第2步: 分批查询并同步
	var offset int
	totalProcessed := 0
	totalErrors := 0

	for {
		var topics []topic.Topic
		result := database.DB.WithContext(ctx).
			Offset(offset).
			Limit(batchSize).
			Find(&topics)

		if result.Error != nil {
			return fmt.Errorf("failed to fetch topics: %w", result.Error)
		}

		if len(topics) == 0 {
			break
		}

		// 批量索引
		topicMaps := make([]map[string]interface{}, 0, len(topics))
		for _, t := range topics {
			topicMap := s.topicToMap(t)
			topicMaps = append(topicMaps, topicMap)
		}

		if err := s.client.BulkIndex(ctx, topicMaps); err != nil {
			log.Printf("Batch error at offset %d: %v", offset, err)
			totalErrors += len(topicMaps)
		} else {
			totalProcessed += len(topicMaps)
		}

		offset += batchSize
	}

	log.Printf("✓ Topic sync completed: %d processed, %d errors", totalProcessed, totalErrors)
	return nil
}

// SyncTopicIncremental 增量同步 - 同步最近修改的话题
func (s *SyncService) SyncTopicIncremental(ctx context.Context, sinceMinutes int) error {
	log.Printf("Starting incremental sync (last %d minutes)...", sinceMinutes)

	since := time.Now().Add(-time.Duration(sinceMinutes) * time.Minute)

	var topics []topic.Topic
	result := database.DB.WithContext(ctx).
		Where("updated_at >= ?", since).
		Find(&topics)

	if result.Error != nil {
		return fmt.Errorf("failed to fetch recent topics: %w", result.Error)
	}

	if len(topics) == 0 {
		log.Println("No topics to sync")
		return nil
	}

	// 批量索引
	topicMaps := make([]map[string]interface{}, 0, len(topics))
	for _, t := range topics {
		topicMap := s.topicToMap(t)
		topicMaps = append(topicMaps, topicMap)
	}

	if err := s.client.BulkIndex(ctx, topicMaps); err != nil {
		return fmt.Errorf("failed to sync topics: %w", err)
	}

	log.Printf("✓ Incremental sync completed: %d topics", len(topicMaps))
	return nil
}

// IndexSingleTopic 索引单个话题
func (s *SyncService) IndexSingleTopic(ctx context.Context, topicID uint64) error {
	var t topic.Topic
	result := database.DB.WithContext(ctx).First(&t, topicID)

	if result.Error != nil {
		return fmt.Errorf("failed to fetch topic: %w", result.Error)
	}

	topicMap := s.topicToMap(t)
	return s.client.IndexTopic(ctx, topicMap)
}

// RemoveSingleTopic 删除单个话题索引
func (s *SyncService) RemoveSingleTopic(ctx context.Context, topicID uint64) error {
	return s.client.DeleteTopic(ctx, fmt.Sprintf("%d", topicID))
}

// topicToMap 将话题模型转为搜索文档
func (s *SyncService) topicToMap(t topic.Topic) map[string]interface{} {
	return map[string]interface{}{
		"id":             t.ID,
		"title":          t.Title,
		"content":        t.Body,
		"description":    t.Description,
		"category_id":    t.CategoryID,
		"user_id":        t.UserID,
		"status":         "published", // 假设已发布
		"created_at":     t.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updated_at":     t.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		"likes_count":    t.LikesCount,
		"views_count":    t.ViewsCount,
		"comments_count": t.CommentsCount,
	}
}

// ReindexTopics 重建索引 - 删除旧索引并创建新索引
func (s *SyncService) ReindexTopics(ctx context.Context, batchSize int) error {
	log.Println("Starting full reindex...")

	// 删除旧索引
	if err := s.indexManager.DeleteTopicIndex(ctx); err != nil {
		log.Printf("Warning: failed to delete old index: %v", err)
	}

	// 完整同步
	return s.SyncAllTopics(ctx, batchSize)
}

// GetSyncStatus 获取同步状态
func (s *SyncService) GetSyncStatus(ctx context.Context) (map[string]interface{}, error) {
	// 获取MySQL中话题总数
	var mysqlCount int64
	database.DB.WithContext(ctx).Model(&topic.Topic{}).Count(&mysqlCount)

	// 获取ES中索引的话题数
	esCount, err := s.searchService.CountTopics(ctx, SearchRequest{
		Query: "*",
	})
	if err != nil {
		esCount = -1
	}

	return map[string]interface{}{
		"mysql_total": mysqlCount,
		"es_indexed":  esCount,
		"synced":      mysqlCount == esCount,
		"sync_status": func() string {
			if mysqlCount == esCount {
				return "in-sync"
			}
			return "out-of-sync"
		}(),
	}, nil
}
