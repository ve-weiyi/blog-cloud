package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// redis缓存
type RedisCache struct {
	rdsDb *redis.Client
}

func NewRedisCache(rdb *redis.Client) *RedisCache {
	return &RedisCache{
		rdsDb: rdb,
	}
}

func (s *RedisCache) Get(key string) (string, error) {
	str, err := s.rdsDb.Get(context.Background(), key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}
	return str, nil
}
func (s *RedisCache) Set(key string, val string, exp int64) error {
	res := s.rdsDb.Set(context.Background(), key, val, time.Second*time.Duration(exp))
	return res.Err()
}

func (s *RedisCache) Inc(key string) error {
	st := `local current = redis.call('incr',KEYS[1]);
local t = redis.call('ttl',KEYS[1]); 
if t == -1 then  
	redis.call('expire',KEYS[1],ARGV[1]) 
end;
return current`
	err := s.rdsDb.Eval(context.Background(), st, []string{key}, 10*60).Err()
	return err
}

func (s *RedisCache) Del(key string) error {
	res := s.rdsDb.Del(context.Background(), key)
	return res.Err()
}

func (s *RedisCache) HsetNx(key, filed string, val interface{}) (bool, error) {
	res := s.rdsDb.HSetNX(context.Background(), key, filed, val)
	return res.Result()
}

func (s *RedisCache) HGetAll(key string) (map[string]string, error) {
	res := s.rdsDb.HGetAll(context.Background(), key)
	return res.Result()
}

func (s *RedisCache) HDel(key string, fields ...string) error {
	res := s.rdsDb.HDel(context.Background(), key, fields...)
	return res.Err()
}

func (s *RedisCache) HGet(key string, field string) (string, error) {
	res := s.rdsDb.HGet(context.Background(), key, field)
	return res.Result()
}
func (s *RedisCache) HVals(key string) ([]string, error) {
	strs, err := s.rdsDb.HVals(context.Background(), key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return []string{}, err
	}
	return strs, nil
}
