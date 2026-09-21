package limitx

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// slidingWindowLua 是滑动窗口限流脚本：每次申请作为一条带时间戳的记录放进 ZSET，
// 先清掉窗口外的记录，再按窗口内的记录数判定。
//
// 用记录数而非固定窗口计数器，是为了消除窗口交界处的突刺——固定窗口在交界处
// 可放行接近两倍配额。配额是每分钟个位数，按记录存储的代价可忽略。
//
// 时钟取自调用方（与仓库内既有的时效判定一致），多副本部署下依赖机器时钟同步。
//
// KEYS[1] 限流键
// ARGV[1] 当前时间（毫秒）
// ARGV[2] 窗口长度（毫秒）
// ARGV[3] 窗口内配额
// ARGV[4] 本次申请的记录标识（须唯一，重复标识会被 ZADD 覆盖而少计一次）
// 返回 {状态, 剩余次数, 重置毫秒数}
const slidingWindowLua = `
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local quota = tonumber(ARGV[3])

redis.call("ZREMRANGEBYSCORE", KEYS[1], 0, now - window)
local used = redis.call("ZCARD", KEYS[1])

local reset = window
local oldest = redis.call("ZRANGE", KEYS[1], 0, 0, "WITHSCORES")
if oldest[2] then
    reset = window - (now - tonumber(oldest[2]))
end

if used >= quota then
    return {0, 0, reset}
end

redis.call("ZADD", KEYS[1], now, ARGV[4])
redis.call("PEXPIRE", KEYS[1], window)

local state = 1
if used + 1 == quota then
    state = 2
end

return {state, quota - used - 1, reset}
`

// evalClient 抽象出限流所需的最小 Redis 能力，便于替换实现与测试。
type evalClient interface {
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
}

// PeriodLimit 基于 Redis 的滑动窗口限流器，支持任意时间窗口。
type PeriodLimit struct {
	period    int
	quota     int
	keyPrefix string
	client    evalClient
}

func NewPeriodLimit(period, quota int, client *redis.Client, keyPrefix string) *PeriodLimit {
	return &PeriodLimit{
		period:    period,
		quota:     quota,
		keyPrefix: keyPrefix,
		client:    client,
	}
}

// Take 申请一次配额，返回 Allowed / HitQuota / OverQuota，并给出余量与窗口重置时长。
func (l *PeriodLimit) Take(ctx context.Context, key string) (Result, error) {
	now := time.Now()
	window := time.Duration(l.period) * time.Second

	reply, err := l.client.Eval(ctx, slidingWindowLua,
		[]string{l.keyPrefix + key},
		now.UnixMilli(),
		window.Milliseconds(),
		l.quota,
		strconv.FormatInt(now.UnixNano(), 10),
	).Int64Slice()
	if err != nil {
		return Result{}, err
	}
	if len(reply) != 3 {
		return Result{}, fmt.Errorf("unexpected eval response length: %d", len(reply))
	}

	return Result{
		State:     int(reply[0]),
		Limit:     l.quota,
		Remaining: int(reply[1]),
		Reset:     time.Duration(reply[2]) * time.Millisecond,
	}, nil
}
