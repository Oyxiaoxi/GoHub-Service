// Package cache 缓存配置与分级策略
package cache

import "time"

// CacheTier 缓存层级配置
type CacheTier struct {
	Name string
	TTL  time.Duration
}

// 预定义的缓存层级
var (
	// HotDataTier 热点数据：高频访问，长 TTL
	HotDataTier = CacheTier{
		Name: "hot",
		TTL:  60 * time.Minute,
	}

	// WarmDataTier 温数据：中等访问频率
	WarmDataTier = CacheTier{
		Name: "warm",
		TTL:  30 * time.Minute,
	}

	// ColdDataTier 冷数据：低频访问，短 TTL
	ColdDataTier = CacheTier{
		Name: "cold",
		TTL:  10 * time.Minute,
	}

	// ListDataTier 列表数据：快速变化，短 TTL
	ListDataTier = CacheTier{
		Name: "list",
		TTL:  10 * time.Minute,
	}
)

// GetTierByAccessFrequency 根据访问频率获取缓存层级
func GetTierByAccessFrequency(accessCount int) CacheTier {
	if accessCount > 1000 {
		return HotDataTier
	} else if accessCount > 100 {
		return WarmDataTier
	}
	return ColdDataTier
}

// EntityCacheTiers 实体类型的默认缓存层级
var EntityCacheTiers = map[string]CacheTier{
	"topic":    WarmDataTier, // 话题：中等热度
	"user":     HotDataTier,  // 用户：高频访问
	"category": HotDataTier,  // 分类：高频访问
	"link":     ColdDataTier, // 链接：低频访问
}

// GetEntityTier 获取实体的缓存层级
func GetEntityTier(entityType string) CacheTier {
	if tier, ok := EntityCacheTiers[entityType]; ok {
		return tier
	}
	return WarmDataTier // 默认中等层级
}
