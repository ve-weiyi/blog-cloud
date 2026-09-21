package storex

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisMemberStore 基于 Redis ZSET 的 MemberStore 实现。
// member 的过期时间记在各自的 score 上，因此互不影响。
type RedisMemberStore struct {
	client *redis.Client
}

// NewRedisMemberStore 创建 RedisMemberStore。
func NewRedisMemberStore(client *redis.Client) *RedisMemberStore {
	return &RedisMemberStore{client: client}
}

func (s *RedisMemberStore) Add(ctx context.Context, key, member string, ttl time.Duration) error {
	if err := s.client.ZAdd(ctx, key, redis.Z{Score: expiryScore(ttl), Member: member}).Err(); err != nil {
		return err
	}
	if ttl > 0 {
		return s.client.Expire(ctx, key, ttl).Err()
	}

	return nil
}

func (s *RedisMemberStore) Members(ctx context.Context, key string) ([]string, error) {
	// 清理已过期 member
	if err := s.client.ZRemRangeByScore(ctx, key, "-inf", formatFloat(float64(time.Now().Unix()))).Err(); err != nil {
		return nil, err
	}

	return s.client.ZRange(ctx, key, 0, -1).Result()
}

func (s *RedisMemberStore) Has(ctx context.Context, key, member string) (bool, error) {
	score, err := s.client.ZScore(ctx, key, member).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if score <= float64(time.Now().Unix()) {
		// 惰性清理
		_ = s.client.ZRem(ctx, key, member).Err()
		return false, nil
	}

	return true, nil
}

func (s *RedisMemberStore) Remove(ctx context.Context, key string, members ...string) error {
	if len(members) == 0 {
		return s.client.Del(ctx, key).Err()
	}

	return s.client.ZRem(ctx, key, members).Err()
}

// ---- helpers ----------------------------------------------------------------

// expiryScore 把 TTL 换算成 ZSET 的过期时间戳。
// ttl <= 0 表示不过期，用 +Inf 表示 —— Redis 的 score 允许取无穷大。
func expiryScore(ttl time.Duration) float64 {
	if ttl <= 0 {
		return math.Inf(1)
	}

	return float64(time.Now().Unix()) + ttl.Seconds()
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
