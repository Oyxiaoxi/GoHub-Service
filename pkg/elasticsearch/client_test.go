package elasticsearch

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClient 集成测试用例
type TestClient struct {
	*Client
	indexName string
	t         *testing.T
}

// setupTest 初始化测试环境
func setupTest(t *testing.T) *TestClient {
	// 连接到本地Elasticsearch (确保Docker已启动)
	client, err := NewClient([]string{"http://localhost:9200"})
	require.NoError(t, err, "failed to create elasticsearch client")

	indexName := fmt.Sprintf("test-topics-%d", time.Now().UnixNano())

	return &TestClient{
		Client:    client,
		indexName: indexName,
		t:         t,
	}
}

// cleanup 清理测试环境
func (tc *TestClient) cleanup(ctx context.Context) {
	// 删除测试索引
	tc.DeleteIndex(ctx, tc.indexName)
}

// TestHealth 测试ES健康检查
func TestHealth(t *testing.T) {
	tc := setupTest(t)
	defer tc.cleanup(context.Background())

	ctx := context.Background()
	healthy, err := tc.Health(ctx)

	assert.NoError(t, err, "健康检查应该成功")
	assert.True(t, healthy, "Elasticsearch应该是健康的")
}

// TestCreateIndex 测试索引创建
func TestCreateIndex(t *testing.T) {
	tc := setupTest(t)
	defer tc.cleanup(context.Background())

	ctx := context.Background()

	// 创建索引
	im := NewIndexManager(tc.Client)
	err := im.CreateIndex(ctx, tc.indexName)
	assert.NoError(t, err, "创建索引应该成功")

	// 验证索引存在
	exists, err := im.IndexExists(ctx, tc.indexName)
	assert.NoError(t, err, "检查索引存在应该成功")
	assert.True(t, exists, "索引应该存在")
}

// TestDeleteIndex 测试索引删除
func TestDeleteIndex(t *testing.T) {
	tc := setupTest(t)

	ctx := context.Background()
	im := NewIndexManager(tc.Client)

	// 创建索引
	err := im.CreateIndex(ctx, tc.indexName)
	require.NoError(t, err)

	// 删除索引
	err = tc.DeleteIndex(ctx, tc.indexName)
	assert.NoError(t, err, "删除索引应该成功")

	// 验证索引不存在
	exists, err := im.IndexExists(ctx, tc.indexName)
	assert.NoError(t, err)
	assert.False(t, exists, "索引应该不存在")
}

// TestIndexTopic 测试话题索引
func TestIndexTopic(t *testing.T) {
	tc := setupTest(t)
	defer tc.cleanup(context.Background())

	ctx := context.Background()

	// 创建索引
	im := NewIndexManager(tc.Client)
	err := im.CreateIndex(ctx, tc.indexName)
	require.NoError(t, err)

	// 索引话题
	topicData := map[string]interface{}{
		"id":             int64(1),
		"title":          "如何学习Go语言",
		"body":           "Go语言是一门现代编程语言...",
		"user_id":        int64(100),
		"category_id":    int64(1),
		"view_count":     100,
		"like_count":     20,
		"comment_count":  5,
		"created_at":     time.Now(),
		"status":         "active",
	}

	err = tc.IndexDocument(ctx, tc.indexName, "1", topicData)
	assert.NoError(t, err, "索引话题应该成功")

	// 等待索引刷新
	time.Sleep(1 * time.Second)

	// 验证文档存在
	doc, err := tc.GetDocument(ctx, tc.indexName, "1")
	assert.NoError(t, err, "获取文档应该成功")
	assert.NotNil(t, doc)
}

// TestBulkIndex 测试批量索引
func TestBulkIndex(t *testing.T) {
	tc := setupTest(t)
	defer tc.cleanup(context.Background())

	ctx := context.Background()

	// 创建索引
	im := NewIndexManager(tc.Client)
	err := im.CreateIndex(ctx, tc.indexName)
	require.NoError(t, err)

	// 准备批量数据
	docs := []map[string]interface{}{
		{
			"id":          int64(1),
			"title":       "Go语言教程",
			"body":        "第一章 入门...",
			"user_id":     int64(100),
			"category_id": int64(1),
			"view_count":  100,
			"created_at":  time.Now(),
			"status":      "active",
		},
		{
			"id":          int64(2),
			"title":       "Rust入门指南",
			"body":        "Rust是一门系统编程语言...",
			"user_id":     int64(101),
			"category_id": int64(2),
			"view_count":  50,
			"created_at":  time.Now(),
			"status":      "active",
		},
		{
			"id":          int64(3),
			"title":       "Python数据分析",
			"body":        "使用pandas进行数据分析...",
			"user_id":     int64(102),
			"category_id": int64(3),
			"view_count":  200,
			"created_at":  time.Now(),
			"status":      "active",
		},
	}

	// 批量索引
	err = tc.BulkIndex(ctx, tc.indexName, docs)
	assert.NoError(t, err, "批量索引应该成功")

	// 等待索引刷新
	time.Sleep(1 * time.Second)

	// 验证文档数量
	count, err := tc.CountDocuments(ctx, tc.indexName)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count, "应该有3个文档")
}

// TestSearch 测试搜索功能
func TestSearch(t *testing.T) {
	tc := setupTest(t)
	defer tc.cleanup(context.Background())

	ctx := context.Background()

	// 创建索引并索引数据
	im := NewIndexManager(tc.Client)
	err := im.CreateIndex(ctx, tc.indexName)
	require.NoError(t, err)

	docs := []map[string]interface{}{
		{
			"id":    int64(1),
			"title": "Go语言高性能编程",
			"body":  "并发编程、内存管理、性能优化...",
		},
		{
			"id":    int64(2),
			"title": "Go语言入门教程",
			"body":  "基础语法、goroutine、channel...",
		},
		{
			"id":    int64(3),
			"title": "Python数据科学",
			"body":  "numpy、pandas、scikit-learn...",
		},
	}

	err = tc.BulkIndex(ctx, tc.indexName, docs)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	// 搜索Go语言相关内容
	ss := NewSearchService(tc.Client)
	results, err := ss.Search(ctx, tc.indexName, map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "Go语言",
			},
		},
	})

	assert.NoError(t, err, "搜索应该成功")
	assert.Greater(t, len(results), 0, "应该找到搜索结果")
	assert.Equal(t, 2, len(results), "应该找到2个Go语言相关的文档")
}

// TestSearchWithFilter 测试带过滤的搜索
func TestSearchWithFilter(t *testing.T) {
	tc := setupTest(t)
	defer tc.cleanup(context.Background())

	ctx := context.Background()

	// 设置索引
	im := NewIndexManager(tc.Client)
	err := im.CreateIndex(ctx, tc.indexName)
	require.NoError(t, err)

	docs := []map[string]interface{}{
		{
			"id":          int64(1),
			"title":       "高性能Go",
			"category_id": int64(1),
			"view_count":  1000,
		},
		{
			"id":          int64(2),
			"title":       "入门Go",
			"category_id": int64(2),
			"view_count":  100,
		},
		{
			"id":          int64(3),
			"title":       "高性能Python",
			"category_id": int64(1),
			"view_count":  500,
		},
	}

	err = tc.BulkIndex(ctx, tc.indexName, docs)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	// 搜索Go且view_count>500
	ss := NewSearchService(tc.Client)
	results, err := ss.Search(ctx, tc.indexName, map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"title": "Go",
						},
					},
				},
				"filter": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							"view_count": map[string]interface{}{
								"gte": 500,
							},
						},
					},
				},
			},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, len(results), "应该找到1个符合条件的文档")
}

// TestAggregation 测试聚合功能
func TestAggregation(t *testing.T) {
	tc := setupTest(t)
	defer tc.cleanup(context.Background())

	ctx := context.Background()

	// 设置索引
	im := NewIndexManager(tc.Client)
	err := im.CreateIndex(ctx, tc.indexName)
	require.NoError(t, err)

	docs := []map[string]interface{}{
		{
			"id":          int64(1),
			"title":       "文章1",
			"category_id": int64(1),
			"view_count":  100,
		},
		{
			"id":          int64(2),
			"title":       "文章2",
			"category_id": int64(1),
			"view_count":  200,
		},
		{
			"id":          int64(3),
			"title":       "文章3",
			"category_id": int64(2),
			"view_count":  300,
		},
	}

	err = tc.BulkIndex(ctx, tc.indexName, docs)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	// 按category_id聚合
	ss := NewSearchService(tc.Client)
	results, err := ss.Aggregate(ctx, tc.indexName, map[string]interface{}{
		"aggs": map[string]interface{}{
			"categories": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "category_id",
				},
			},
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, results)
}

// TestUpdateDocument 测试文档更新
func TestUpdateDocument(t *testing.T) {
	tc := setupTest(t)
	defer tc.cleanup(context.Background())

	ctx := context.Background()

	// 创建索引
	im := NewIndexManager(tc.Client)
	err := im.CreateIndex(ctx, tc.indexName)
	require.NoError(t, err)

	// 索引初始文档
	initialDoc := map[string]interface{}{
		"id":         int64(1),
		"title":      "原始标题",
		"view_count": 10,
	}

	err = tc.IndexDocument(ctx, tc.indexName, "1", initialDoc)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	// 更新文档
	updateData := map[string]interface{}{
		"view_count": 20,
		"title":      "更新后的标题",
	}

	err = tc.UpdateDocument(ctx, tc.indexName, "1", updateData)
	assert.NoError(t, err, "更新文档应该成功")

	time.Sleep(1 * time.Second)

	// 验证更新
	doc, err := tc.GetDocument(ctx, tc.indexName, "1")
	assert.NoError(t, err)
	assert.NotNil(t, doc)
}

// TestDeleteDocument 测试文档删除
func TestDeleteDocument(t *testing.T) {
	tc := setupTest(t)
	defer tc.cleanup(context.Background())

	ctx := context.Background()

	// 创建索引并索引文档
	im := NewIndexManager(tc.Client)
	err := im.CreateIndex(ctx, tc.indexName)
	require.NoError(t, err)

	doc := map[string]interface{}{
		"id":    int64(1),
		"title": "待删除的文档",
	}

	err = tc.IndexDocument(ctx, tc.indexName, "1", doc)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	// 删除文档
	err = tc.DeleteDocument(ctx, tc.indexName, "1")
	assert.NoError(t, err, "删除文档应该成功")

	time.Sleep(1 * time.Second)

	// 验证文档已删除
	resultDoc, err := tc.GetDocument(ctx, tc.indexName, "1")
	assert.NoError(t, err)
	assert.Nil(t, resultDoc, "文档应该被删除")
}

// TestSearchSuggestions 测试搜索建议
func TestSearchSuggestions(t *testing.T) {
	tc := setupTest(t)
	defer tc.cleanup(context.Background())

	ctx := context.Background()

	// 创建索引
	im := NewIndexManager(tc.Client)
	err := im.CreateIndex(ctx, tc.indexName)
	require.NoError(t, err)

	// 索引示例数据
	docs := []map[string]interface{}{
		{"id": int64(1), "title": "Golang教程"},
		{"id": int64(2), "title": "Go并发编程"},
		{"id": int64(3), "title": "Go语言最佳实践"},
	}

	err = tc.BulkIndex(ctx, tc.indexName, docs)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	// 搜索建议
	ss := NewSearchService(tc.Client)
	suggestions, err := ss.Suggest(ctx, tc.indexName, "go")

	assert.NoError(t, err, "获取建议应该成功")
	assert.Greater(t, len(suggestions), 0, "应该返回建议")
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	tc := setupTest(t)

	ctx := context.Background()

	// 尝试访问不存在的索引
	_, err := tc.GetDocument(ctx, "non-existent-index", "1")
	assert.Error(t, err, "访问不存在的索引应该返回错误")

	// 尝试删除不存在的索引
	err = tc.DeleteIndex(ctx, "non-existent-index")
	assert.NoError(t, err, "删除不存在的索引不应该报错")
}

// BenchmarkSearch 搜索性能基准
func BenchmarkSearch(b *testing.B) {
	tc := setupTest(&testing.T{})
	defer tc.cleanup(context.Background())

	ctx := context.Background()

	// 设置索引和数据
	im := NewIndexManager(tc.Client)
	im.CreateIndex(ctx, tc.indexName)

	// 准备1000个文档
	docs := make([]map[string]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		docs[i] = map[string]interface{}{
			"id":    int64(i),
			"title": fmt.Sprintf("文档%d Go语言内容", i),
			"body":  "Go语言是一门现代编程语言，具有高性能和并发能力...",
		}
	}

	tc.BulkIndex(ctx, tc.indexName, docs)
	time.Sleep(1 * time.Second)

	ss := NewSearchService(tc.Client)

	b.ResetTimer()

	// 基准测试搜索性能
	for i := 0; i < b.N; i++ {
		ss.Search(ctx, tc.indexName, map[string]interface{}{
			"query": map[string]interface{}{
				"match": map[string]interface{}{
					"title": "Go语言",
				},
			},
		})
	}
}
