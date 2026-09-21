package storex

import (
	"context"
	"time"

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

func (s *RedisStore) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return s.client.Set(ctx, key, value, ttl).Err()
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

func (s *RedisStore) Delete(ctx context.Context, key string, keys ...string) error {
	return s.client.Del(ctx, append([]string{key}, keys...)...).Err()
}

func (s *RedisStore) Keys(ctx context.Context, pattern string) ([]string, error) {
	var keys []string

	iter := s.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}
