package storex

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// RedisStore 基于 Redis 的 KVStore 实现。
type RedisStore struct {
	client *redis.Client
}

// NewRedisStore 创建 RedisStore。
func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

func (s *RedisStore) Set(ctx context.Context, key, value string, opts ...SetOption) error {
	o := resolveOpts(opts)
	return s.client.Set(ctx, o.Prefix+key, value, o.TTL).Err()
}

func (s *RedisStore) Get(ctx context.Context, key string) (string, bool, error) {
	val, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return val, true, nil
}

func (s *RedisStore) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}
