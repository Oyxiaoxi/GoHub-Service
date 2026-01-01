package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// Client 封装的ES客户端
type Client struct {
	client *elasticsearch.Client
}

// NewClient 创建新的ES客户端
func NewClient(addresses []string) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	return &Client{client: client}, nil
}

// Health 检查ES集群健康状态
func (c *Client) Health(ctx context.Context) (bool, error) {
	req := esapi.InfoRequest{}
	res, err := req.Do(ctx, c.client)
	if err != nil {
		return false, fmt.Errorf("failed to check es health: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return false, fmt.Errorf("es health check failed: %s", res.String())
	}

	return true, nil
}

// Close 关闭客户端
func (c *Client) Close() error {
	return nil // ES Go客户端无需显式关闭
}

// IndexTopic 索引话题文档
func (c *Client) IndexTopic(ctx context.Context, topic map[string]interface{}) error {
	topicID := fmt.Sprintf("%v", topic["id"])

	body, err := json.Marshal(topic)
	if err != nil {
		return fmt.Errorf("failed to marshal topic: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      "gohub-topics",
		DocumentID: topicID,
		Body:       bytes.NewReader(body),
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to index topic: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("es index error: %s", res.String())
	}

	return nil
}

// DeleteTopic 删除话题索引
func (c *Client) DeleteTopic(ctx context.Context, topicID string) error {
	req := esapi.DeleteRequest{
		Index:      "gohub-topics",
		DocumentID: topicID,
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to delete topic: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("es delete error: %s", res.String())
	}

	return nil
}

// BulkIndex 批量索引文档
func (c *Client) BulkIndex(ctx context.Context, topics []map[string]interface{}) error {
	if len(topics) == 0 {
		return nil
	}

	for _, topic := range topics {
		body, err := json.Marshal(topic)
		if err != nil {
			log.Printf("Failed to marshal topic: %v", err)
			continue
		}

		req := esapi.IndexRequest{
			Index:      "gohub-topics",
			DocumentID: fmt.Sprintf("%v", topic["id"]),
			Body:       bytes.NewReader(body),
		}

		res, err := req.Do(ctx, c.client)
		if err != nil {
			log.Printf("Failed to index topic: %v", err)
		}
		if res != nil {
			res.Body.Close()
		}
	}

	return nil
}

// Search 搜索话题
func (c *Client) Search(ctx context.Context, query map[string]interface{}) (map[string]interface{}, error) {
	body, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	req := esapi.SearchRequest{
		Index: []string{"gohub-topics"},
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("es search error: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// CreateIndex 创建索引
func (c *Client) CreateIndex(ctx context.Context, indexName string, mapping map[string]interface{}) error {
	body, err := json.Marshal(mapping)
	if err != nil {
		return fmt.Errorf("failed to marshal mapping: %w", err)
	}

	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 400 { // 400 = 索引已存在
		return fmt.Errorf("es create index error: %s", res.String())
	}

	return nil
}

// DeleteIndex 删除索引
func (c *Client) DeleteIndex(ctx context.Context, indexName string) error {
	req := esapi.IndicesDeleteRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to delete index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("es delete index error: %s", res.String())
	}

	return nil
}
