// Package cache 缓存工具类，可以缓存各种类型包括 struct 对象
package cache

import (
	"context"
	"encoding/json"
	"GoHub-Service/pkg/logger"
	"sync"
	"time"

	"github.com/spf13/cast"
)

type CacheService struct {
	Store Store
}

var once sync.Once
var Cache *CacheService

func InitWithCacheStore(store Store) {
	once.Do(func() {
		Cache = &CacheService{
			Store: store,
		}
	})
}

func Set(ctx context.Context, key string, obj interface{}, expireTime time.Duration) {
	b, err := json.Marshal(&obj)
	logger.LogIf(err)
	Cache.Store.Set(ctx, key, string(b), expireTime)
}

func Get(ctx context.Context, key string) interface{} {
	stringValue := Cache.Store.Get(ctx, key)
	var wanted interface{}
	err := json.Unmarshal([]byte(stringValue), &wanted)
	logger.LogIf(err)
	return wanted
}

func Has(ctx context.Context, key string) bool {
	return Cache.Store.Has(ctx, key)
}

// GetObject 应该传地址，用法如下:
//     model := user.User{}
//     cache.GetObject("key", &model)
func GetObject(ctx context.Context, key string, wanted interface{}) {
	val := Cache.Store.Get(ctx, key)
	if len(val) > 0 {
		err := json.Unmarshal([]byte(val), &wanted)
		logger.LogIf(err)
	}
}

func GetString(ctx context.Context, key string) string {
	return cast.ToString(Get(ctx, key))
}

func GetBool(ctx context.Context, key string) bool {
	return cast.ToBool(Get(ctx, key))
}

func GetInt(ctx context.Context, key string) int {
	return cast.ToInt(Get(ctx, key))
}

func GetInt32(ctx context.Context, key string) int32 {
	return cast.ToInt32(Get(ctx, key))
}

func GetInt64(ctx context.Context, key string) int64 {
	return cast.ToInt64(Get(ctx, key))
}

func GetUint(ctx context.Context, key string) uint {
	return cast.ToUint(Get(ctx, key))
}

func GetUint32(ctx context.Context, key string) uint32 {
	return cast.ToUint32(Get(ctx, key))
}

func GetUint64(ctx context.Context, key string) uint64 {
	return cast.ToUint64(Get(ctx, key))
}

func GetFloat64(ctx context.Context, key string) float64 {
	return cast.ToFloat64(Get(ctx, key))
}

func GetTime(ctx context.Context, key string) time.Time {
	return cast.ToTime(Get(ctx, key))
}

func GetDuration(ctx context.Context, key string) time.Duration {
	return cast.ToDuration(Get(ctx, key))
}

func GetIntSlice(ctx context.Context, key string) []int {
	return cast.ToIntSlice(Get(ctx, key))
}

func GetStringSlice(ctx context.Context, key string) []string {
	return cast.ToStringSlice(Get(ctx, key))
}

func GetStringMap(ctx context.Context, key string) map[string]interface{} {
	return cast.ToStringMap(Get(ctx, key))
}

func GetStringMapString(ctx context.Context, key string) map[string]string {
	return cast.ToStringMapString(Get(ctx, key))
}

func GetStringMapStringSlice(ctx context.Context, key string) map[string][]string {
	return cast.ToStringMapStringSlice(Get(ctx, key))
}

func Forget(ctx context.Context, key string) {
	Cache.Store.Forget(ctx, key)
}

func Forever(ctx context.Context, key string, value string) {
	Cache.Store.Set(ctx, key, value, 0)
}

func Flush(ctx context.Context) {
	Cache.Store.Flush(ctx)
}

func Increment(ctx context.Context, parameters ...interface{}) {
	Cache.Store.Increment(ctx, parameters...)
}

func Decrement(ctx context.Context, parameters ...interface{}) {
	Cache.Store.Decrement(ctx, parameters...)
}

func IsAlive(ctx context.Context) error {
	return Cache.Store.IsAlive(ctx)
}
