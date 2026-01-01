package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"GoHub-Service/pkg/elasticsearch"
)

// 这是一个测试脚本，用于初始化和测试Elasticsearch集成

func main() {
	// 1. 创建ES客户端
	addresses := []string{"http://localhost:9200"}
	client, err := elasticsearch.NewClient(addresses)
	if err != nil {
		log.Fatalf("Failed to create ES client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 2. 检查ES健康状态
	health, err := client.Health(ctx)
	if err != nil {
		log.Fatalf("Failed to check ES health: %v", err)
	}
	fmt.Printf("ES Health: %v\n", health)

	// 3. 创建索引管理器
	indexManager := elasticsearch.NewIndexManager(client)

	// 4. 创建话题索引
	if err := indexManager.CreateTopicIndex(ctx); err != nil {
		log.Fatalf("Failed to create topic index: %v", err)
	}
	fmt.Println("✓ Topic index created successfully")

	// 5. 索引测试数据
	testTopics := []map[string]interface{}{
		{
			"id":             1,
			"title":          "Golang 最佳实践",
			"content":        "本文介绍Golang的最佳实践...",
			"description":    "Golang编程最佳实践指南",
			"category_id":    1,
			"user_id":        1,
			"created_at":     "2026-01-01T10:00:00Z",
			"likes_count":    100,
			"views_count":    5000,
			"comments_count": 20,
			"status":         "published",
		},
		{
			"id":             2,
			"title":          "Elasticsearch 搜索引擎",
			"content":        "Elasticsearch是一个强大的搜索引擎...",
			"description":    "Elasticsearch入门指南",
			"category_id":    2,
			"user_id":        2,
			"created_at":     "2026-01-02T10:00:00Z",
			"likes_count":    80,
			"views_count":    4000,
			"comments_count": 15,
			"status":         "published",
		},
	}

	if err := client.BulkIndex(ctx, testTopics); err != nil {
		log.Fatalf("Failed to bulk index topics: %v", err)
	}
	fmt.Println("✓ Test data indexed successfully")

	// 6. 刷新索引
	if err := indexManager.RefreshIndex(ctx, "gohub-topics"); err != nil {
		log.Fatalf("Failed to refresh index: %v", err)
	}
	fmt.Println("✓ Index refreshed")

	// 7. 搜索测试
	searchService := elasticsearch.NewSearchService(client)
	results, total, err := searchService.SearchTopics(ctx, elasticsearch.SearchRequest{
		Query:    "Golang",
		Page:     1,
		PageSize: 10,
		SortBy:   "relevance",
	})
	if err != nil {
		log.Fatalf("Failed to search: %v", err)
	}

	fmt.Printf("\n✓ Search results (total: %d):\n", total)
	for _, result := range results {
		fmt.Printf("  - [%d] %s (score: %.2f)\n", result.ID, result.Title, result.Score)
	}

	// 8. 获取索引统计
	stats, err := indexManager.GetIndexStats(ctx, "gohub-topics")
	if err != nil {
		log.Fatalf("Failed to get stats: %v", err)
	}
	fmt.Printf("\n✓ Index Stats:\n%+v\n", stats)

	fmt.Println("\n✓ All tests passed!")
}
