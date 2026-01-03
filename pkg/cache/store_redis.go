package cache

import (
	"context"
	"GoHub-Service/pkg/config"
	"GoHub-Service/pkg/redis"
	"time"
)

// RedisStore 实现 cache.Store interface
type RedisStore struct {
	RedisClient *redis.RedisClient
	KeyPrefix   string
}

func NewRedisStore(address string, username string, password string, db int) *RedisStore {
	rs := &RedisStore{}
	rs.RedisClient = redis.NewClient(address, username, password, db)
	rs.KeyPrefix = config.GetString("app.name") + ":cache:"
	return rs
}

func (s *RedisStore) Set(ctx context.Context, key string, value string, expireTime time.Duration) {
	s.RedisClient.Set(ctx, s.KeyPrefix+key, value, expireTime)
}

func (s *RedisStore) Get(ctx context.Context, key string) string {
	return s.RedisClient.Get(ctx, s.KeyPrefix+key)
}

func (s *RedisStore) Has(ctx context.Context, key string) bool {
	return s.RedisClient.Has(ctx, s.KeyPrefix+key)
}

func (s *RedisStore) Forget(ctx context.Context, key string) {
	s.RedisClient.Del(ctx, s.KeyPrefix+key)
}

func (s *RedisStore) Forever(ctx context.Context, key string, value string) {
	s.RedisClient.Set(ctx, s.KeyPrefix+key, value, 0)
}

func (s *RedisStore) Flush(ctx context.Context) {
	s.RedisClient.FlushDB(ctx)
}

func (s *RedisStore) Increment(ctx context.Context, parameters ...interface{}) {
	s.RedisClient.Increment(ctx, parameters...)
}

func (s *RedisStore) Decrement(ctx context.Context, parameters ...interface{}) {
	s.RedisClient.Decrement(ctx, parameters...)
}

func (s *RedisStore) IsAlive(ctx context.Context) error {
	return s.RedisClient.Ping(ctx)
}
