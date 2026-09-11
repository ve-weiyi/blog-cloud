package limitx

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const periodLua = `
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local current = redis.call("INCRBY", KEYS[1], 1)
if current == 1 then
    redis.call("expire", KEYS[1], window)
end
if current < limit then
    return 1
elseif current == limit then
    return 2
else
    return 0
end
`

// PeriodLimit 基于 Redis 的滑动窗口限流器，支持任意时间窗口。
type PeriodLimit struct {
	period    int
	quota     int
	keyPrefix string
	client    *redis.Client
}

func NewPeriodLimit(period, quota int, client *redis.Client, keyPrefix string) *PeriodLimit {
	return &PeriodLimit{
		period:    period,
		quota:     quota,
		keyPrefix: keyPrefix,
		client:    client,
	}
}

// Take 申请一次配额，返回 Allowed / HitQuota / OverQuota。
func (l *PeriodLimit) Take(ctx context.Context, key string) (int, error) {
	resp, err := l.client.Eval(ctx, periodLua, []string{l.keyPrefix + key}, l.quota, l.period).Result()
	if err != nil {
		return 0, err
	}
	code, ok := resp.(int64)
	if !ok {
		return 0, fmt.Errorf("unexpected eval response type: %T", resp)
	}
	return int(code), nil
}
