package limitx

import (
	"context"
	"time"
)

// Limiter 限流器接口，支持不同的底层实现（go-redis、go-zero redis 等）。
type Limiter interface {
	Take(ctx context.Context, key string) (Result, error)
}

// Result 是一次限流判定的结果。
//
// 除了"允不允许"，它还要给出余量与窗口重置时长：调用方据此在响应上回传余量，
// 使客户端能提前降速，而不是撞到 429 才知道。
type Result struct {
	State     int           // Allowed / HitQuota / OverQuota
	Limit     int           // 窗口内的配额上限
	Remaining int           // 判定后本窗口剩余可申请次数
	Reset     time.Duration // 最早一条记录滑出窗口、配额开始回落的剩余时长
}

const (
	Allowed   = 1
	HitQuota  = 2
	OverQuota = 0
)
