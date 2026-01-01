package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

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

// IndexDocument 索引单个文档到指定索引
func (c *Client) IndexDocument(ctx context.Context, indexName, docID string, doc map[string]interface{}) error {
	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: docID,
		Body:       bytes.NewReader(body),
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("es index error: %s", res.String())
	}

	return nil
}

// GetDocument 获取单个文档
func (c *Client) GetDocument(ctx context.Context, indexName, docID string) (map[string]interface{}, error) {
	req := esapi.GetRequest{
		Index:      indexName,
		DocumentID: docID,
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 404 {
			return nil, nil // 文档不存在
		}
		return nil, fmt.Errorf("es get error: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// 提取source部分
	if source, ok := result["_source"]; ok {
		if sourceMap, ok := source.(map[string]interface{}); ok {
			return sourceMap, nil
		}
	}

	return nil, nil
}

// UpdateDocument 更新文档
func (c *Client) UpdateDocument(ctx context.Context, indexName, docID string, doc map[string]interface{}) error {
	updateBody := map[string]interface{}{
		"doc": doc,
	}

	body, err := json.Marshal(updateBody)
	if err != nil {
		return fmt.Errorf("failed to marshal update: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      indexName,
		DocumentID: docID,
		Body:       bytes.NewReader(body),
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("es update error: %s", res.String())
	}

	return nil
}

// DeleteDocument 删除单个文档
func (c *Client) DeleteDocument(ctx context.Context, indexName, docID string) error {
	req := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: docID,
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("es delete error: %s", res.String())
	}

	return nil
}

// CountDocuments 统计索引中的文档数
func (c *Client) CountDocuments(ctx context.Context, indexName string) (int64, error) {
	req := esapi.CountRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, fmt.Errorf("es count error: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if count, ok := result["count"].(float64); ok {
		return int64(count), nil
	}

	return 0, nil
}

// BulkIndex 批量索引文档（支持[]map或map[string]map）
func (c *Client) BulkIndex(ctx context.Context, indexName string, docs interface{}) error {
	var buf bytes.Buffer

	switch docsTyped := docs.(type) {
	case map[string]map[string]interface{}:
		// 处理map[string]map[string]interface{}
		for docID, doc := range docsTyped {
			meta := []byte(fmt.Sprintf(`{"index":{"_index":"%s","_id":"%s"}}`, indexName, docID))
			buf.Write(meta)
			buf.WriteString("\n")

			docBytes, err := json.Marshal(doc)
			if err != nil {
				return fmt.Errorf("failed to marshal document: %w", err)
			}
			buf.Write(docBytes)
			buf.WriteString("\n")
		}

	case []map[string]interface{}:
		// 处理[]map[string]interface{}
		for i, doc := range docsTyped {
			docID := fmt.Sprintf("%d", i)
			meta := []byte(fmt.Sprintf(`{"index":{"_index":"%s","_id":"%s"}}`, indexName, docID))
			buf.Write(meta)
			buf.WriteString("\n")

			docBytes, err := json.Marshal(doc)
			if err != nil {
				return fmt.Errorf("failed to marshal document: %w", err)
			}
			buf.Write(docBytes)
			buf.WriteString("\n")
		}

	default:
		return fmt.Errorf("unsupported docs type: must be map[string]map[string]interface{} or []map[string]interface{}")
	}

	req := esapi.BulkRequest{
		Body: &buf,
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to bulk index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("es bulk index error: %s", res.String())
	}

	// 检查响应中是否有错误
	var bulkResp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&bulkResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if hasErrors, ok := bulkResp["errors"].(bool); ok && hasErrors {
		return fmt.Errorf("some documents failed to index")
	}

	return nil
}

// Search 执行搜索查询
func (c *Client) Search(ctx context.Context, query map[string]interface{}) (map[string]interface{}, error) {
	body, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	req := esapi.SearchRequest{
		Body: bytes.NewReader(body),
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

// Aggregate 执行聚合查询
func (c *Client) Aggregate(ctx context.Context, indexName string, aggs map[string]interface{}) (map[string]interface{}, error) {
	query := map[string]interface{}{
		"size": 0,
		"aggs": aggs,
	}

	body, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	req := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("es aggregate error: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// Suggest 搜索建议查询
func (c *Client) Suggest(ctx context.Context, indexName string, text string) (map[string]interface{}, error) {
	query := map[string]interface{}{
		"suggest": map[string]interface{}{
			"text": text,
			"title-suggest": map[string]interface{}{
				"completion": map[string]interface{}{
					"field": "title.completion",
				},
			},
		},
	}

	body, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	req := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return nil, fmt.Errorf("failed to suggest: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("es suggest error: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// IndexExists 检查索引是否存在
func (c *Client) IndexExists(ctx context.Context, indexName string) (bool, error) {
	res, err := c.client.Indices.Exists(
		[]string{indexName},
		c.client.Indices.Exists.WithContext(ctx),
	)
	if err != nil {
		return false, fmt.Errorf("failed to check index existence: %w", err)
	}
	defer res.Body.Close()

	return !res.IsError(), nil
}
