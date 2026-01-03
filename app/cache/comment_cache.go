// Package cache Comment缓存层
package cache

import (
	"context"
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
func (cc *CommentCache) GetByID(ctx context.Context, id string) (*comment.Comment, error) {
	key := cc.cacheKeyPrefix + id

	// 从缓存获取
	data := cache.Get(ctx, key)
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
func (cc *CommentCache) Set(ctx context.Context, commentModel *comment.Comment) error {
	key := cc.cacheKeyPrefix + fmt.Sprintf("%d", commentModel.ID)

	data, err := json.Marshal(commentModel)
	if err != nil {
		return err
	}

	cache.Set(ctx, key, string(data), cc.cacheTime)
	return nil
}

// Invalidate 删除评论缓存
func (cc *CommentCache) Invalidate(ctx context.Context, id string) error {
	key := cc.cacheKeyPrefix + id
	cache.Forget(ctx, key)
	return nil
}

// GetByTopicID 从缓存获取话题的评论列表
func (cc *CommentCache) GetByTopicID(ctx context.Context, topicID string) ([]comment.Comment, error) {
	key := cc.topicCacheKeyPrefix + topicID

	data := cache.Get(ctx, key)
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
func (cc *CommentCache) SetByTopicID(ctx context.Context, topicID string, comments []comment.Comment) error {
	key := cc.topicCacheKeyPrefix + topicID

	data, err := json.Marshal(comments)
	if err != nil {
		return err
	}

	cache.Set(ctx, key, string(data), cc.cacheTime)
	return nil
}

// InvalidateByTopicID 删除话题的评论列表缓存
func (cc *CommentCache) InvalidateByTopicID(ctx context.Context, topicID string) error {
	key := cc.topicCacheKeyPrefix + topicID
	cache.Forget(ctx, key)
	return nil
}

// InvalidateAll 清除所有评论相关缓存
func (cc *CommentCache) InvalidateAll(ctx context.Context) error {
	// 这里可以实现更复杂的缓存清除逻辑
	// 比如使用 Redis SCAN 命令查找所有相关 key
	return nil
}
