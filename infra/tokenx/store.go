package tokenx

import (
	"context"
	"time"
)

// Store 纯 KV 存储接口，与业务解耦。
type Store interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error

	// Get key 不存在或已过期时返回空字符串，不返回 error
	Get(ctx context.Context, key string) (string, error)

	Delete(ctx context.Context, key ...string) error
	Keys(ctx context.Context, pattern string) ([]string, error)
}
