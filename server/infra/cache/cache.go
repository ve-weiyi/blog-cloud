package cache

// 提取RedisCache所有方法到接口
type Cache interface {
	Get(key string) (string, error)
	Set(key string, val string, exp int64) error
	Inc(key string) error
	Del(key string) error
	HsetNx(key, filed string, val interface{}) (bool, error)
	HGetAll(key string) (map[string]string, error)
	HDel(key string, fields ...string) error
	HGet(key string, field string) (string, error)
	HVals(key string) ([]string, error)
}
