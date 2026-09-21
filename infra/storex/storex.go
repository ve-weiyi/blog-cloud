package storex

import (
	"context"
	"time"
)

// ---- KVStore ------------------------------------------------------------------

// KVStore 纯键值存储，单 key 单 value。
type KVStore interface {
	// Set 写入 key。ttl <= 0 表示不过期。
	Set(ctx context.Context, key, value string, ttl time.Duration) error

	// Get 读取 key。第二个返回值区分「没写过」与「写了个空串」——
	// 两者都表现为空字符串，光看 value 分不出来。
	Get(ctx context.Context, key string) (string, bool, error)

	// Delete 删除一个或多个 key。首参必填，保证不会一个 key 都不传。
	Delete(ctx context.Context, key string, keys ...string) error

	// Keys 返回匹配 pattern 的全部 key。
	//
	// 代价随整个 keyspace 规模线性增长（底层是 SCAN 游标遍历，不是按前缀索引），
	// 因此只应出现在管理类路径上，不要放进请求热路径 —— 需要按维度取一组 key 时，
	// 优先用 MemberStore 直接记下那份清单。
	Keys(ctx context.Context, pattern string) ([]string, error)
}

// ---- MemberStore --------------------------------------------------------------

// MemberStore 单 key 多 member 集合存储，member 之间独立 TTL。
type MemberStore interface {
	// Add 加入一个 member。ttl <= 0 表示不过期，与 KVStore 一致。
	Add(ctx context.Context, key, member string, ttl time.Duration) error

	// Members 返回全部未过期的 member。
	Members(ctx context.Context, key string) ([]string, error)

	// Has 判断 member 是否存在且未过期。
	Has(ctx context.Context, key, member string) (bool, error)

	// Remove 移除指定 member；不指定 member 时删除整个 key。
	Remove(ctx context.Context, key string, members ...string) error
}
