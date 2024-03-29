package cache

import "github.com/orca-zhang/ecache"

// 内存缓存
type MemoryCache struct {
	defaultInstance ecache.Cache
	instances       map[int64]ecache.Cache
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{}
}

func (s *MemoryCache) Get(key string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *MemoryCache) Set(key string, val string, exp int64) error {
	//TODO implement me
	panic("implement me")
}

func (s *MemoryCache) Inc(key string) error {
	//TODO implement me
	panic("implement me")
}

func (s *MemoryCache) Del(key string) error {
	//TODO implement me
	panic("implement me")
}

func (s *MemoryCache) HsetNx(key, filed string, val interface{}) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s *MemoryCache) HGetAll(key string) (map[string]string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *MemoryCache) HDel(key string, fields ...string) error {
	//TODO implement me
	panic("implement me")
}

func (s *MemoryCache) HGet(key string, field string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *MemoryCache) HVals(key string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}
