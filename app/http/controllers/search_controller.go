package controllers

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"GoHub-Service/pkg/elasticsearch"
	"GoHub-Service/pkg/response"
)

// SearchController 搜索控制器
type SearchController struct {
	searchService *elasticsearch.SearchService
}

// NewSearchController 创建搜索控制器
func NewSearchController(searchService *elasticsearch.SearchService) *SearchController {
	return &SearchController{
		searchService: searchService,
	}
}

// SearchTopics 搜索话题
// @Summary 搜索话题
// @Description 全文搜索话题，支持分类筛选、排序等
// @Tags Search
// @Accept json
// @Produce json
// @Param query query string false "搜索关键词"
// @Param category_id query int false "分类ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param sort_by query string false "排序方式: relevance|latest|popular" default(relevance)
// @Success 200 {object} map[string]interface{}
// @Router /api/search/topics [get]
func (sc *SearchController) SearchTopics(c *gin.Context) {
	// 获取请求参数
	query := c.DefaultQuery("query", "")
	categoryID := c.DefaultQuery("category_id", "")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	sortBy := c.DefaultQuery("sort_by", "relevance")

	// 参数转换
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	if pageInt < 1 {
		pageInt = 1
	}
	if pageSizeInt < 1 || pageSizeInt > 100 {
		pageSizeInt = 20
	}

	// 构建搜索请求
	req := elasticsearch.SearchRequest{
		Query:      query,
		CategoryID: categoryID,
		Page:       pageInt,
		PageSize:   pageSizeInt,
		SortBy:     sortBy,
	}

	// 执行搜索
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	results, total, err := sc.searchService.SearchTopics(ctx, req)
	if err != nil {
		response.ApiError(c, 500, response.CodeServerError, "搜索失败")
		return
	}

	// 返回结果
	response.ApiSuccess(c, gin.H{
		"items":     results,
		"total":     total,
		"page":      pageInt,
		"page_size": pageSizeInt,
	})
}

// SearchSuggestions 搜索建议/自动完成
// @Summary 搜索建议
// @Description 获取搜索建议用于自动完成功能
// @Tags Search
// @Accept json
// @Produce json
// @Param query query string true "搜索前缀"
// @Param limit query int false "返回数量" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /api/search/suggestions [get]
func (sc *SearchController) SearchSuggestions(c *gin.Context) {
	query := c.Query("query")
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	if limit < 1 || limit > 50 {
		limit = 10
	}

	if query == "" {
		response.ApiError(c, 400, response.CodeInvalidParams, "query parameter is required")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	suggestions, err := sc.searchService.SuggestTopics(ctx, query, limit)
	if err != nil {
		response.ApiError(c, 500, response.CodeServerError, "获取建议失败")
		return
	}

	response.ApiSuccess(c, gin.H{
		"suggestions": suggestions,
	})
}

// GetHotTopics 获取热门话题
// @Summary 热门话题
// @Description 获取最近7天的热门话题
// @Tags Search
// @Accept json
// @Produce json
// @Param limit query int false "返回数量" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /api/search/hot-topics [get]
func (sc *SearchController) GetHotTopics(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	if limit < 1 || limit > 50 {
		limit = 10
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	topics, err := sc.searchService.GetHotTopics(ctx, limit)
	if err != nil {
		response.ApiError(c, 500, response.CodeServerError, "获取热门话题失败")
		return
	}

	response.ApiSuccess(c, gin.H{
		"items": topics,
	})
}

// IndexTopicForSearch 索引话题用于搜索 (内部接口)
// 当话题创建或更新时调用此接口
func (sc *SearchController) IndexTopicForSearch(c *gin.Context, topic map[string]interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sc.searchService.IndexTopic(ctx, topic); err != nil {
		// 日志记录但不中断主流程
		c.Error(err)
		return
	}
}

// RemoveTopicFromSearch 从搜索索引中移除话题 (内部接口)
// 当话题删除时调用此接口
func (sc *SearchController) RemoveTopicFromSearch(c *gin.Context, topicID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sc.searchService.RemoveTopic(ctx, topicID); err != nil {
		// 日志记录但不中断主流程
		c.Error(err)
		return
	}
}
