package storex

import (
	"context"
	"time"
)

// ---- SetOption ----------------------------------------------------------------

// SetOption 写操作配置。
type SetOption func(*SetOptions)

// SetOptions 写操作选项。
type SetOptions struct {
	TTL    time.Duration
	Prefix string
}

// WithTTL 设置过期时间。不支持 TTL 的后端可忽略。
func WithTTL(d time.Duration) SetOption {
	return func(o *SetOptions) {
		o.TTL = d
	}
}

// WithPrefix 统一 key 前缀。
func WithPrefix(prefix string) SetOption {
	return func(o *SetOptions) {
		o.Prefix = prefix
	}
}

// ---- KVStore ------------------------------------------------------------------

// KVStore 纯键值存储，单 key 单 value。
type KVStore interface {
	Set(ctx context.Context, key, value string, opts ...SetOption) error
	Get(ctx context.Context, key string) (string, bool, error)
	Delete(ctx context.Context, key string) error
}

// ---- MemberStore --------------------------------------------------------------

// MemberStore 单 key 多 member 集合存储，member 之间独立 TTL。
type MemberStore interface {
	Add(ctx context.Context, key, member string, opts ...SetOption) error
	Members(ctx context.Context, key string) ([]string, error)
	Has(ctx context.Context, key, member string) (bool, error)
	Remove(ctx context.Context, key string, members ...string) error
}
