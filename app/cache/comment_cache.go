// Package cache Comment缓存层
package cache

import (
	"GoHub-Service/app/models/comment"
	"GoHub-Service/pkg/cache"
	"encoding/json"
	"fmt"
	"time"
)

// CommentCache 评论缓存
type CommentCache struct {
	cacheKeyPrefix      string
	topicCacheKeyPrefix string
	cacheTime           time.Duration
}

// NewCommentCache 创建评论缓存实例，使用分级 TTL
func NewCommentCache() *CommentCache {
	tier := GetEntityTier("comment")
	return &CommentCache{
		cacheKeyPrefix:      "comment:",
		topicCacheKeyPrefix: "comment:topic:",
		cacheTime:           tier.TTL,
	}
}

// GetByID 从缓存获取评论
func (cc *CommentCache) GetByID(id string) (*comment.Comment, error) {
	key := cc.cacheKeyPrefix + id

	// 从缓存获取
	data := cache.Get(key)
	dataStr, ok := data.(string)
	if !ok || dataStr == "" {
		return nil, nil
	}

	var commentModel comment.Comment
	err := json.Unmarshal([]byte(dataStr), &commentModel)
	if err != nil {
		return nil, nil
	}

	return &commentModel, nil
}

// Set 设置评论缓存
func (cc *CommentCache) Set(commentModel *comment.Comment) error {
	key := cc.cacheKeyPrefix + fmt.Sprintf("%d", commentModel.ID)

	data, err := json.Marshal(commentModel)
	if err != nil {
		return err
	}

	cache.Set(key, string(data), cc.cacheTime)
	return nil
}

// Invalidate 删除评论缓存
func (cc *CommentCache) Invalidate(id string) error {
	key := cc.cacheKeyPrefix + id
	cache.Forget(key)
	return nil
}

// GetByTopicID 从缓存获取话题的评论列表
func (cc *CommentCache) GetByTopicID(topicID string) ([]comment.Comment, error) {
	key := cc.topicCacheKeyPrefix + topicID

	data := cache.Get(key)
	dataStr, ok := data.(string)
	if !ok || dataStr == "" {
		return nil, nil
	}

	var comments []comment.Comment
	err := json.Unmarshal([]byte(dataStr), &comments)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// SetByTopicID 设置话题的评论列表缓存
func (cc *CommentCache) SetByTopicID(topicID string, comments []comment.Comment) error {
	key := cc.topicCacheKeyPrefix + topicID

	data, err := json.Marshal(comments)
	if err != nil {
		return err
	}

	cache.Set(key, string(data), cc.cacheTime)
	return nil
}

// InvalidateByTopicID 删除话题的评论列表缓存
func (cc *CommentCache) InvalidateByTopicID(topicID string) error {
	key := cc.topicCacheKeyPrefix + topicID
	cache.Forget(key)
	return nil
}

// InvalidateAll 清除所有评论相关缓存
func (cc *CommentCache) InvalidateAll() error {
	// 这里可以实现更复杂的缓存清除逻辑
	// 比如使用 Redis SCAN 命令查找所有相关 key
	return nil
}
