// Package cache 缓存Key集中管理
package cache

import (
	"fmt"
	"time"
)

// KeyManager 缓存Key管理器
type KeyManager struct {
	prefix string
	ttl    time.Duration
}

// 缓存Key常量定义
const (
	// 实体缓存前缀
	PrefixComment  = "comment:"
	PrefixTopic    = "topic:"
	PrefixUser     = "user:"
	PrefixCategory = "category:"
	PrefixLink     = "link:"
	PrefixRole     = "role:"
	PrefixPermission = "permission:"
	
	// 列表缓存前缀
	PrefixCommentList  = "comment:list:"
	PrefixTopicList    = "topic:list:"
	PrefixUserList     = "user:list:"
	PrefixCategoryList = "category:list:"
	
	// 关联缓存前缀
	PrefixCommentByTopic = "comment:topic:"
	PrefixTopicByUser    = "topic:user:"
	PrefixTopicByCategory = "topic:category:"
	
	// 统计缓存前缀
	PrefixCommentCount = "count:comment:"
	PrefixTopicCount   = "count:topic:"
	PrefixViewCount    = "count:view:"
	PrefixLikeCount    = "count:like:"
	
	// 空值缓存标记（防穿透）
	EmptyValue = "EMPTY_VALUE"
	EmptyTTL   = 5 * time.Minute
	
	// 锁前缀
	PrefixLock = "lock:"
)

// 默认TTL配置
const (
	TTLShort  = 5 * time.Minute   // 短期缓存：5分钟
	TTLMedium = 30 * time.Minute  // 中期缓存：30分钟
	TTLLong   = 2 * time.Hour     // 长期缓存：2小时
	TTLDay    = 24 * time.Hour    // 日缓存：24小时
)

// NewKeyManager 创建Key管理器
func NewKeyManager(prefix string, ttl time.Duration) *KeyManager {
	return &KeyManager{
		prefix: prefix,
		ttl:    ttl,
	}
}

// CommentKeys 评论Key管理
type CommentKeys struct{}

func (k *CommentKeys) ByID(id string) string {
	return PrefixComment + id
}

func (k *CommentKeys) ListByTopic(topicID string) string {
	return PrefixCommentByTopic + topicID
}

func (k *CommentKeys) ListByPage(page, perPage int) string {
	return fmt.Sprintf("%spage:%d:per_page:%d", PrefixCommentList, page, perPage)
}

func (k *CommentKeys) CountByTopic(topicID string) string {
	return PrefixCommentCount + "topic:" + topicID
}

func (k *CommentKeys) Lock(id string) string {
	return PrefixLock + PrefixComment + id
}

// TopicKeys 话题Key管理
type TopicKeys struct{}

func (k *TopicKeys) ByID(id string) string {
	return PrefixTopic + id
}

func (k *TopicKeys) ListByCategory(categoryID string, page, perPage int) string {
	return fmt.Sprintf("%scategory:%s:page:%d:per_page:%d", 
		PrefixTopicByCategory, categoryID, page, perPage)
}

func (k *TopicKeys) ListByUser(userID string, page, perPage int) string {
	return fmt.Sprintf("%suser:%s:page:%d:per_page:%d", 
		PrefixTopicByUser, userID, page, perPage)
}

func (k *TopicKeys) ListByPage(page, perPage int) string {
	return fmt.Sprintf("%spage:%d:per_page:%d", PrefixTopicList, page, perPage)
}

func (k *TopicKeys) ViewCount(id string) string {
	return PrefixViewCount + "topic:" + id
}

func (k *TopicKeys) LikeCount(id string) string {
	return PrefixLikeCount + "topic:" + id
}

func (k *TopicKeys) Lock(id string) string {
	return PrefixLock + PrefixTopic + id
}

// UserKeys 用户Key管理
type UserKeys struct{}

func (k *UserKeys) ByID(id string) string {
	return PrefixUser + id
}

func (k *UserKeys) ByEmail(email string) string {
	return PrefixUser + "email:" + email
}

func (k *UserKeys) ByPhone(phone string) string {
	return PrefixUser + "phone:" + phone
}

func (k *UserKeys) Lock(id string) string {
	return PrefixLock + PrefixUser + id
}

// CategoryKeys 分类Key管理
type CategoryKeys struct{}

func (k *CategoryKeys) ByID(id string) string {
	return PrefixCategory + id
}

func (k *CategoryKeys) List() string {
	return PrefixCategoryList + "all"
}

func (k *CategoryKeys) Lock(id string) string {
	return PrefixLock + PrefixCategory + id
}

// LinkKeys 链接Key管理
type LinkKeys struct{}

func (k *LinkKeys) ByID(id string) string {
	return PrefixLink + id
}

func (k *LinkKeys) List() string {
	return "link:list:all"
}

// 全局Key管理器实例
var (
	Comment  = &CommentKeys{}
	Topic    = &TopicKeys{}
	User     = &UserKeys{}
	Category = &CategoryKeys{}
	Link     = &LinkKeys{}
)

// BuildKey 构建完整的缓存Key（带前缀和命名空间）
func BuildKey(namespace, key string) string {
	return namespace + ":" + key
}

// BuildListKey 构建列表缓存Key
func BuildListKey(namespace string, page, perPage int) string {
	return fmt.Sprintf("%s:list:page:%d:per_page:%d", namespace, page, perPage)
}

// BuildCountKey 构建计数缓存Key
func BuildCountKey(resource, id string) string {
	return fmt.Sprintf("count:%s:%s", resource, id)
}

// IsEmptyValue 检查是否为空值标记
func IsEmptyValue(value interface{}) bool {
	if str, ok := value.(string); ok {
		return str == EmptyValue
	}
	return false
}
