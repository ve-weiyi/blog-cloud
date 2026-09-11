package limitx

import "context"

// Limiter 限流器接口，支持不同的底层实现（go-redis、go-zero redis 等）。
type Limiter interface {
	Take(ctx context.Context, key string) (int, error)
}

const (
	Allowed   = 1
	HitQuota  = 2
	OverQuota = 0
)
