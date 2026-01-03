// Package cache 缓存防护策略
package cache

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	// ErrCacheMiss 缓存未命中
	ErrCacheMiss = errors.New("cache miss")
	
	// ErrEmptyValue 空值（防穿透标记）
	ErrEmptyValue = errors.New("empty value cached")
)

// GuardConfig 防护配置
type GuardConfig struct {
	// 是否启用缓存穿透防护
	EnablePenetrationGuard bool
	
	// 是否启用缓存雪崩防护
	EnableAvalancheGuard bool
	
	// 是否启用缓存击穿防护（使用singleflight）
	EnableBreakdownGuard bool
	
	// 空值缓存时间（防穿透）
	EmptyValueTTL time.Duration
	
	// 过期时间随机范围（防雪崩，秒）
	RandomExpireRange int
	
	// 互斥锁超时时间（防击穿）
	MutexTimeout time.Duration
}

// DefaultGuardConfig 默认防护配置
func DefaultGuardConfig() *GuardConfig {
	return &GuardConfig{
		EnablePenetrationGuard: true,
		EnableAvalancheGuard:   true,
		EnableBreakdownGuard:   true,
		EmptyValueTTL:          EmptyTTL,
		RandomExpireRange:      300, // 5分钟随机范围
		MutexTimeout:           5 * time.Second,
	}
}

// CacheGuard 缓存防护器
type CacheGuard struct {
	config *GuardConfig
	locks  sync.Map // 用于防击穿的互斥锁
}

// NewCacheGuard 创建缓存防护器
func NewCacheGuard(config *GuardConfig) *CacheGuard {
	if config == nil {
		config = DefaultGuardConfig()
	}
	return &CacheGuard{
		config: config,
	}
}

// GetWithProtection 带防护的缓存获取
// 集成：防穿透 + 防雪崩 + 防击穿
func (g *CacheGuard) GetWithProtection(
	ctx context.Context,
	key string,
	ttl time.Duration,
	fetchFunc func() (interface{}, error),
) (interface{}, error) {
	// 1. 尝试从缓存获取
	value := Get(ctx, key)
	
	// 1.1 防穿透：检查是否为空值标记
	if g.config.EnablePenetrationGuard && IsEmptyValue(value) {
		return nil, ErrEmptyValue
	}
	
	// 1.2 缓存命中
	if value != nil {
		return value, nil
	}
	
	// 2. 缓存未命中，防击穿：使用互斥锁
	if g.config.EnableBreakdownGuard {
		return g.getWithMutex(ctx, key, ttl, fetchFunc)
	}
	
	// 3. 不使用互斥锁，直接查询
	return g.fetchAndCache(ctx, key, ttl, fetchFunc)
}

// getWithMutex 使用互斥锁防止缓存击穿
func (g *CacheGuard) getWithMutex(
	ctx context.Context,
	key string,
	ttl time.Duration,
	fetchFunc func() (interface{}, error),
) (interface{}, error) {
	// 获取或创建锁
	lockKey := PrefixLock + key
	lockInterface, _ := g.locks.LoadOrStore(lockKey, &sync.Mutex{})
	lock := lockInterface.(*sync.Mutex)
	
	// 加锁
	lock.Lock()
	defer lock.Unlock()
	
	// 双重检查：再次尝试从缓存获取
	value := Get(ctx, key)
	if value != nil {
		return value, nil
	}
	
	// 查询并缓存
	return g.fetchAndCache(ctx, key, ttl, fetchFunc)
}

// fetchAndCache 查询数据并缓存
func (g *CacheGuard) fetchAndCache(
	ctx context.Context,
	key string,
	ttl time.Duration,
	fetchFunc func() (interface{}, error),
) (interface{}, error) {
	// 执行查询函数
	value, err := fetchFunc()
	if err != nil {
		return nil, err
	}
	
	// 防穿透：如果查询结果为空，缓存空值标记
	if value == nil && g.config.EnablePenetrationGuard {
		Set(ctx, key, EmptyValue, g.config.EmptyValueTTL)
		return nil, ErrEmptyValue
	}
	
	// 防雪崩：添加随机过期时间
	finalTTL := ttl
	if g.config.EnableAvalancheGuard {
		finalTTL = g.addRandomExpire(ttl)
	}
	
	// 缓存数据
	Set(ctx, key, value, finalTTL)
	return value, nil
}

// addRandomExpire 添加随机过期时间（防雪崩）
func (g *CacheGuard) addRandomExpire(ttl time.Duration) time.Duration {
	if g.config.RandomExpireRange <= 0 {
		return ttl
	}
	
	// 添加随机秒数
	randomSeconds := rand.Intn(g.config.RandomExpireRange)
	return ttl + time.Duration(randomSeconds)*time.Second
}

// SetWithExpireJitter 设置缓存并添加过期时间抖动（防雪崩）
func (g *CacheGuard) SetWithExpireJitter(
	ctx context.Context,
	key string,
	value interface{},
	ttl time.Duration,
) {
	finalTTL := ttl
	if g.config.EnableAvalancheGuard {
		finalTTL = g.addRandomExpire(ttl)
	}
	Set(ctx, key, value, finalTTL)
}

// CacheNullValue 缓存空值（防穿透）
func (g *CacheGuard) CacheNullValue(ctx context.Context, key string) {
	if !g.config.EnablePenetrationGuard {
		return
	}
	Set(ctx, key, EmptyValue, g.config.EmptyValueTTL)
}

// DeleteNullValue 删除空值标记
func (g *CacheGuard) DeleteNullValue(ctx context.Context, key string) {
	Forget(ctx, key)
}

// 全局防护器实例
var (
	globalGuard     *CacheGuard
	globalGuardOnce sync.Once
)

// GetGuard 获取全局防护器
func GetGuard() *CacheGuard {
	globalGuardOnce.Do(func() {
		globalGuard = NewCacheGuard(DefaultGuardConfig())
	})
	return globalGuard
}

// SafeGet 安全获取缓存（带防护）
func SafeGet(
	ctx context.Context,
	key string,
	ttl time.Duration,
	fetchFunc func() (interface{}, error),
) (interface{}, error) {
	return GetGuard().GetWithProtection(ctx, key, ttl, fetchFunc)
}

// SafeSet 安全设置缓存（带过期时间抖动）
func SafeSet(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	GetGuard().SetWithExpireJitter(ctx, key, value, ttl)
}

// SetEmptyCache 设置空缓存（防穿透）
func SetEmptyCache(ctx context.Context, key string) {
	GetGuard().CacheNullValue(ctx, key)
}

// WarmupCache 缓存预热
type WarmupCache struct {
	warmupFuncs map[string]func(context.Context) error
	mu          sync.RWMutex
}

var (
	warmupInstance *WarmupCache
	warmupOnce     sync.Once
)

// GetWarmup 获取预热实例
func GetWarmup() *WarmupCache {
	warmupOnce.Do(func() {
		warmupInstance = &WarmupCache{
			warmupFuncs: make(map[string]func(context.Context) error),
		}
	})
	return warmupInstance
}

// Register 注册预热函数
func (w *WarmupCache) Register(name string, fn func(context.Context) error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.warmupFuncs[name] = fn
}

// Warmup 执行所有预热函数
func (w *WarmupCache) Warmup(ctx context.Context) error {
	w.mu.RLock()
	defer w.mu.RUnlock()
	
	for name, fn := range w.warmupFuncs {
		if err := fn(ctx); err != nil {
			return fmt.Errorf("warmup %s failed: %w", name, err)
		}
	}
	return nil
}

// WarmupOne 执行指定的预热函数
func (w *WarmupCache) WarmupOne(ctx context.Context, name string) error {
	w.mu.RLock()
	fn, exists := w.warmupFuncs[name]
	w.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("warmup function %s not found", name)
	}
	
	return fn(ctx)
}

// List 列出所有已注册的预热函数
func (w *WarmupCache) List() []string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	
	names := make([]string, 0, len(w.warmupFuncs))
	for name := range w.warmupFuncs {
		names = append(names, name)
	}
	return names
}

// RegisterWarmup 注册预热函数（便捷方法）
func RegisterWarmup(name string, fn func(context.Context) error) {
	GetWarmup().Register(name, fn)
}

// WarmupAll 执行所有预热（便捷方法）
func WarmupAll(ctx context.Context) error {
	return GetWarmup().Warmup(ctx)
}
