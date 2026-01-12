package tokenx

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// RedisStore Redis 存储实现
type RedisStore struct {
	client *redis.Redis
	prefix string
}

// NewRedisStore 创建 Redis 存储实例
func NewRedisStore(client *redis.Redis, prefix string) *RedisStore {
	return &RedisStore{
		client: client,
		prefix: prefix,
	}
}

// key 添加前缀
func (s *RedisStore) key(k string) string {
	if s.prefix == "" {
		return k
	}
	return s.prefix + k
}

func (s *RedisStore) Set(key string, value string, expireSeconds int) error {
	return s.client.Setex(s.key(key), value, expireSeconds)
}

func (s *RedisStore) Get(key string) (string, error) {
	val, err := s.client.Get(s.key(key))
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (s *RedisStore) Delete(key string) error {
	_, err := s.client.Del(s.key(key))
	return err
}

func (s *RedisStore) Exists(key string) (bool, error) {
	exists, err := s.client.Exists(s.key(key))
	return exists, err
}

func (s *RedisStore) SetExpire(key string, expireSeconds int) error {
	return s.client.Expire(s.key(key), expireSeconds)
}
