package storex

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisMemberStore 基于 Redis ZSET 的 MemberStore 实现。
type RedisMemberStore struct {
	client *redis.Client
}

// NewRedisMemberStore 创建 RedisMemberStore。
func NewRedisMemberStore(client *redis.Client) *RedisMemberStore {
	return &RedisMemberStore{client: client}
}

func (s *RedisMemberStore) Add(ctx context.Context, key, member string, opts ...SetOption) error {
	o := resolveOpts(opts)
	k := o.Prefix + key
	score := float64(time.Now().Unix()) + o.TTL.Seconds()

	err := s.client.ZAdd(ctx, k, redis.Z{Score: score, Member: member}).Err()
	if err != nil {
		return err
	}
	if o.TTL > 0 {
		return s.client.Expire(ctx, k, o.TTL).Err()
	}
	return nil
}

func (s *RedisMemberStore) Members(ctx context.Context, key string) ([]string, error) {
	now := float64(time.Now().Unix())

	// 清理已过期 member
	if err := s.client.ZRemRangeByScore(ctx, key, "-inf", formatFloat(now)).Err(); err != nil {
		return nil, err
	}

	vals, err := s.client.ZRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return vals, nil
}

func (s *RedisMemberStore) Has(ctx context.Context, key, member string) (bool, error) {
	score, err := s.client.ZScore(ctx, key, member).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	now := float64(time.Now().Unix())
	if score <= now {
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

func resolveOpts(opts []SetOption) SetOptions {
	o := SetOptions{}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
