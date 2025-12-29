// Package cache Topic缓存层
package cache

import (
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
func (tc *TopicCache) GetByID(id string) (*topic.Topic, error) {
	key := tc.cacheKeyPrefix + id
	
	// 从缓存获取
	data := cache.Get(key)
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
func (tc *TopicCache) Set(topicModel *topic.Topic) error {
	key := tc.cacheKeyPrefix + fmt.Sprintf("%d", topicModel.ID)
	
	data, err := json.Marshal(topicModel)
	if err != nil {
		return err
	}
	
	cache.Set(key, string(data), tc.cacheTime)
	return nil
}

// Delete 删除话题缓存
func (tc *TopicCache) Delete(id string) error {
	key := tc.cacheKeyPrefix + id
	cache.Forget(key)
	return nil
}

// GetListCacheKey 获取列表缓存key
func (tc *TopicCache) GetListCacheKey(c *gin.Context) string {
	page := c.DefaultQuery("page", "1")
	perPage := c.DefaultQuery("per_page", "10")
	return tc.cacheKeyPrefix + "list:page:" + page + ":per_page:" + perPage
}

// GetList 从缓存获取话题列表
func (tc *TopicCache) GetList(c *gin.Context) ([]topic.Topic, bool) {
	key := tc.GetListCacheKey(c)
	
	data := cache.Get(key)
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
func (tc *TopicCache) SetList(c *gin.Context, topics []topic.Topic) error {
	key := tc.GetListCacheKey(c)
	
	data, err := json.Marshal(topics)
	if err != nil {
		return err
	}
	
	cache.Set(key, string(data), 10*time.Minute)
	return nil
}

// ClearList 清除话题列表缓存
func (tc *TopicCache) ClearList() error {
	// 简单处理，清空所有缓存
	cache.Flush()
	return nil
}
