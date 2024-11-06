package captcha

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// 验证码存储
type RedisStore struct {
	Redis      *redis.Client // 缓存 15分钟
	Expiration time.Duration // 过期时间
	PrefixKey  string        // 前缀
	Context    context.Context
}

func NewRedisStore(rd *redis.Client) *RedisStore {
	return &RedisStore{
		Expiration: 15 * 60 * time.Second,
		PrefixKey:  "captcha:",
		Redis:      rd,
		Context:    context.Background(),
	}
}

func (rs *RedisStore) Set(key string, value string) error {
	cacheKey := rs.PrefixKey + key
	err := rs.Redis.Set(rs.Context, cacheKey, value, rs.Expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rs *RedisStore) Get(key string, clear bool) string {
	cacheKey := rs.PrefixKey + key
	val, err := rs.Redis.Get(rs.Context, cacheKey).Result()
	if err != nil {
		return ""
	}
	if clear {
		err := rs.Redis.Del(rs.Context, cacheKey).Err()
		if err != nil {
			return ""
		}
	}
	return val
}

func (rs *RedisStore) Verify(key, answer string, clear bool) bool {
	v := rs.Get(key, clear)
	if v == "" {
		return false
	}
	log.Printf("RedisStore Verify. key:%v,answer:%v,v:%v,clear:%v", key, answer, v, clear)
	return v == answer
}
