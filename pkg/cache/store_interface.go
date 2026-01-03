package cache

import (
	"context"
	"time"
)

type Store interface {
	Set(ctx context.Context, key string, value string, expireTime time.Duration)
	Get(ctx context.Context, key string) string
	Has(ctx context.Context, key string) bool
	Forget(ctx context.Context, key string)
	Forever(ctx context.Context, key string, value string)
	Flush(ctx context.Context)

	IsAlive(ctx context.Context) error

	// Increment 当参数只有 1 个时，为 key，增加 1。
	// 当参数有 2 个时，第一个参数为 key ，第二个参数为要增加的值 int64 类型。
	Increment(ctx context.Context, parameters ...interface{})

	// Decrement 当参数只有 1 个时，为 key，减去 1。
	// 当参数有 2 个时，第一个参数为 key ，第二个参数为要减去的值 int64 类型。
	Decrement(ctx context.Context, parameters ...interface{})
}
