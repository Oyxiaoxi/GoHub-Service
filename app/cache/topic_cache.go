// Package cache Topic缓存层
package cache

import (
	"context"
	"GoHub-Service/app/models/topic"
	"GoHub-Service/pkg/cache"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// TopicCache 话题缓存
type TopicCache struct {
	cacheKeyPrefix string
	cacheTime      time.Duration
}

// NewTopicCache 创建话题缓存实例，使用分级 TTL
func NewTopicCache() *TopicCache {
	tier := GetEntityTier("topic")
	return &TopicCache{
		cacheKeyPrefix: "topic:",
		cacheTime:      tier.TTL,
	}
}

// GetByID 从缓存获取话题
func (tc *TopicCache) GetByID(ctx context.Context, id string) (*topic.Topic, error) {
	key := tc.cacheKeyPrefix + id
	
	// 从缓存获取
	data := cache.Get(ctx, key)
	dataStr, ok := data.(string)
	if !ok || dataStr == "" {
		return nil, nil
	}
	
	var topicModel topic.Topic
	err := json.Unmarshal([]byte(dataStr), &topicModel)
	if err != nil {
		return nil, nil
	}
	
	return &topicModel, nil
}

// Set 设置话题缓存
func (tc *TopicCache) Set(ctx context.Context, topicModel *topic.Topic) error {
	key := tc.cacheKeyPrefix + fmt.Sprintf("%d", topicModel.ID)
	
	data, err := json.Marshal(topicModel)
	if err != nil {
		return err
	}
	
	cache.Set(ctx, key, string(data), tc.cacheTime)
	return nil
}

// Delete 删除话题缓存
func (tc *TopicCache) Delete(ctx context.Context, id string) error {
	key := tc.cacheKeyPrefix + id
	cache.Forget(ctx, key)
	return nil
}

// GetListCacheKey 获取列表缓存key
func (tc *TopicCache) GetListCacheKey(c *gin.Context) string {
	page := c.DefaultQuery("page", "1")
	perPage := c.DefaultQuery("per_page", "10")
	return tc.cacheKeyPrefix + "list:page:" + page + ":per_page:" + perPage
}

// GetList 从缓存获取话题列表
func (tc *TopicCache) GetList(ctx context.Context, c *gin.Context) ([]topic.Topic, bool) {
	key := tc.GetListCacheKey(c)
	
	data := cache.Get(ctx, key)
	dataStr, ok := data.(string)
	if !ok || dataStr == "" {
		return nil, false
	}
	
	var topics []topic.Topic
	err := json.Unmarshal([]byte(dataStr), &topics)
	if err != nil {
		return nil, false
	}
	
	return topics, true
}

// SetList 设置话题列表缓存
func (tc *TopicCache) SetList(ctx context.Context, c *gin.Context, topics []topic.Topic) error {
	key := tc.GetListCacheKey(c)
	
	data, err := json.Marshal(topics)
	if err != nil {
		return err
	}
	
	cache.Set(ctx, key, string(data), 10*time.Minute)
	return nil
}

// ClearList 清除话题列表缓存
func (tc *TopicCache) ClearList(ctx context.Context) error {
	// 简单处理，清空所有缓存
	cache.Flush(ctx)
	return nil
}
