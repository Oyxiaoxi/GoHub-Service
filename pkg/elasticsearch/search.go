package elasticsearch

import (
	"context"
	"fmt"
)

// SearchService 搜索服务
type SearchService struct {
	client *Client
}

// NewSearchService 创建搜索服务
func NewSearchService(client *Client) *SearchService {
	return &SearchService{client: client}
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Query      string
	CategoryID string
	Page       int
	PageSize   int
	SortBy     string // "relevance", "latest", "popular"
}

// SearchResult 搜索结果
type SearchResult struct {
	ID            int64   `json:"id"`
	Title         string  `json:"title"`
	Content       string  `json:"content"`
	Description   string  `json:"description"`
	CategoryID    int64   `json:"category_id"`
	UserID        int64   `json:"user_id"`
	CreatedAt     string  `json:"created_at"`
	LikesCount    int     `json:"likes_count"`
	ViewsCount    int     `json:"views_count"`
	CommentsCount int     `json:"comments_count"`
	Score         float64 `json:"score"`
}

// SearchTopics 搜索话题
func (ss *SearchService) SearchTopics(ctx context.Context, req SearchRequest) ([]SearchResult, int64, error) {
	// 构建查询
	query := ss.buildQuery(req)

	// 执行搜索
	results, err := ss.client.Search(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// 解析结果
	topics, total := ss.parseResults(results)

	return topics, total, nil
}

// buildQuery 构建ES查询
func (ss *SearchService) buildQuery(req SearchRequest) map[string]interface{} {
	query := map[string]interface{}{
		"size": req.PageSize,
		"from": (req.Page - 1) * req.PageSize,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"multi_match": map[string]interface{}{
							"query":  req.Query,
							"fields": []string{"title^2", "content", "description"},
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
		},
	}

	// 添加分类过滤
	if req.CategoryID != "" {
		filters := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]map[string]interface{})
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"category_id": req.CategoryID,
			},
		})
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filters
	}

	// 添加排序
	sort := ss.buildSort(req.SortBy)
	if len(sort) > 0 {
		query["sort"] = sort
	}

	return query
}

// buildSort 构建排序参数
func (ss *SearchService) buildSort(sortBy string) []interface{} {
	switch sortBy {
	case "latest":
		return []interface{}{
			map[string]interface{}{
				"created_at": map[string]string{"order": "desc"},
			},
		}
	case "popular":
		return []interface{}{
			map[string]interface{}{
				"likes_count": map[string]string{"order": "desc"},
			},
			map[string]interface{}{
				"views_count": map[string]string{"order": "desc"},
			},
		}
	default: // relevance
		return []interface{}{
			map[string]interface{}{
				"_score": map[string]string{"order": "desc"},
			},
			map[string]interface{}{
				"created_at": map[string]string{"order": "desc"},
			},
		}
	}
}

// parseResults 解析ES搜索结果
func (ss *SearchService) parseResults(results map[string]interface{}) ([]SearchResult, int64) {
	hits := results["hits"].(map[string]interface{})
	total := int64(0)

	if totalInfo, ok := hits["total"].(map[string]interface{}); ok {
		total = int64(totalInfo["value"].(float64))
	} else if totalInfo, ok := hits["total"].(float64); ok {
		total = int64(totalInfo)
	}

	topics := []SearchResult{}
	hitsList := hits["hits"].([]interface{})

	for _, hit := range hitsList {
		hitData := hit.(map[string]interface{})
		source := hitData["_source"].(map[string]interface{})
		score := hitData["_score"]

		topic := SearchResult{
			ID:          int64(source["id"].(float64)),
			Title:       source["title"].(string),
			CategoryID:  int64(source["category_id"].(float64)),
			UserID:      int64(source["user_id"].(float64)),
			CreatedAt:   source["created_at"].(string),
			LikesCount:  int(source["likes_count"].(float64)),
			ViewsCount:  int(source["views_count"].(float64)),
		}

		if content, ok := source["content"].(string); ok {
			topic.Content = content
		}
		if description, ok := source["description"].(string); ok {
			topic.Description = description
		}
		if commentsCount, ok := source["comments_count"].(float64); ok {
			topic.CommentsCount = int(commentsCount)
		}
		if scoreVal, ok := score.(float64); ok {
			topic.Score = scoreVal
		}

		topics = append(topics, topic)
	}

	return topics, total
}

// SuggestTopics 搜索建议/自动完成
func (ss *SearchService) SuggestTopics(ctx context.Context, prefix string, limit int) ([]string, error) {
	query := map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			"titles": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "title.keyword",
					"size":  limit,
					"order": map[string]string{
						"_count": "desc",
					},
				},
			},
		},
		"query": map[string]interface{}{
			"prefix": map[string]interface{}{
				"title": prefix,
			},
		},
	}

	results, err := ss.client.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	suggests := []string{}
	aggs := results["aggregations"].(map[string]interface{})
	titles := aggs["titles"].(map[string]interface{})
	buckets := titles["buckets"].([]interface{})

	for _, bucket := range buckets {
		bucketData := bucket.(map[string]interface{})
		suggests = append(suggests, bucketData["key"].(string))
	}

	return suggests, nil
}

// GetHotTopics 获取热门话题
func (ss *SearchService) GetHotTopics(ctx context.Context, limit int) ([]SearchResult, error) {
	query := map[string]interface{}{
		"size": limit,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"status": "published",
						},
					},
				},
			},
		},
		"sort": []interface{}{
			map[string]interface{}{
				"likes_count": map[string]string{"order": "desc"},
			},
			map[string]interface{}{
				"views_count": map[string]string{"order": "desc"},
			},
		},
	}

	results, err := ss.client.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	topics, _ := ss.parseResults(results)
	return topics, nil
}

// CountTopics 统计搜索结果数量
func (ss *SearchService) CountTopics(ctx context.Context, req SearchRequest) (int64, error) {
	// 移除分页参数，只查询数量
	query := ss.buildQuery(req)
	delete(query, "size")
	delete(query, "from")

	results, err := ss.client.Search(ctx, query)
	if err != nil {
		return 0, err
	}

	hits := results["hits"].(map[string]interface{})
	if totalInfo, ok := hits["total"].(map[string]interface{}); ok {
		return int64(totalInfo["value"].(float64)), nil
	} else if totalInfo, ok := hits["total"].(float64); ok {
		return int64(totalInfo), nil
	}

	return 0, fmt.Errorf("failed to parse total count")
}

// IndexTopic 索引单个话题文档 (由控制器调用)
func (ss *SearchService) IndexTopic(ctx context.Context, topic map[string]interface{}) error {
	return ss.client.IndexTopic(ctx, topic)
}

// RemoveTopic 从搜索索引删除话题 (由控制器调用)
func (ss *SearchService) RemoveTopic(ctx context.Context, topicID string) error {
	return ss.client.DeleteTopic(ctx, topicID)
}
