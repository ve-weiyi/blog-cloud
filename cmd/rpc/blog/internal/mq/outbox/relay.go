package outbox

import (
	"context"
	"errors"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/vkit/adapter/mqx"
)

const (
	defaultInterval     = time.Second
	defaultBatch        = 50
	defaultMaxRetries   = 10
	defaultBaseDelay    = time.Second
	defaultMaxDelay     = 5 * time.Minute
	defaultStuckAfter   = 5 * time.Minute
	defaultRetention    = 7 * 24 * time.Hour
	defaultStatsEvery   = time.Minute
	defaultCleanupEvery = time.Hour
	defaultResetEvery   = time.Minute
)

// Metrics 记录 outbox 积压状态，由调用方接入 Prometheus 等。
type Metrics interface {
	OutboxStats(pending, dead int64, oldestAge time.Duration)
}

// RelayOption 配置 relay 的可选依赖。
type RelayOption func(*Relay)

// WithMetrics 注入积压指标。
func WithMetrics(metrics Metrics) RelayOption {
	return func(r *Relay) {
		if metrics != nil {
			r.metrics = metrics
		}
	}
}

// Relay 把 outbox 中的待发布消息投递到 MQ。
//
// 投递语义为至少一次：发布成功但标记前崩溃、或标记失败，都会在
// ResetStuck 之后重投；消费方按 event_id 幂等。
type Relay struct {
	store   Store
	pub     mqx.Publisher
	metrics Metrics

	interval     time.Duration
	batch        int
	maxRetries   int
	baseDelay    time.Duration
	maxDelay     time.Duration
	stuckAfter   time.Duration
	retention    time.Duration
	statsEvery   time.Duration
	cleanupEvery time.Duration

	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}

	startOnce sync.Once
	stopOnce  sync.Once

	lastStats   time.Time
	lastCleanup time.Time
	lastReset   time.Time
}

// NewRelay 创建 relay。
func NewRelay(store Store, pub mqx.Publisher, opts ...RelayOption) *Relay {
	relay := &Relay{
		store:        store,
		pub:          pub,
		interval:     defaultInterval,
		batch:        defaultBatch,
		maxRetries:   defaultMaxRetries,
		baseDelay:    defaultBaseDelay,
		maxDelay:     defaultMaxDelay,
		stuckAfter:   defaultStuckAfter,
		retention:    defaultRetention,
		statsEvery:   defaultStatsEvery,
		cleanupEvery: defaultCleanupEvery,
	}
	for _, opt := range opts {
		opt(relay)
	}
	return relay
}

// Start 启动后台投递循环。store 或 publisher 缺失时不启动，只记录错误。
func (r *Relay) Start() {
	if r == nil || r.store == nil {
		logx.Error("[outbox] relay not started: store is nil")
		return
	}
	if r.pub == nil {
		logx.Error("[outbox] relay not started: publisher is nil")
		return
	}
	r.startOnce.Do(func() {
		r.ctx, r.cancel = context.WithCancel(context.Background())
		r.done = make(chan struct{})
		go r.loop()
	})
}

// Stop 停止投递循环，并等待当前批次退出。
func (r *Relay) Stop(ctx context.Context) {
	if r == nil {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}
	r.stopOnce.Do(func() {
		if r.cancel != nil {
			r.cancel()
		}
		if r.done == nil {
			return
		}
		select {
		case <-r.done:
		case <-ctx.Done():
		}
	})
}

func (r *Relay) loop() {
	defer close(r.done)

	r.tick()
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-r.ctx.Done():
			return
		case <-ticker.C:
			r.tick()
		}
	}
}

func (r *Relay) tick() {
	if err := r.RunOnce(r.ctx); err != nil && r.ctx.Err() == nil {
		logx.Errorf("[outbox] relay: %v", err)
	}
	r.cleanup(r.ctx)
	r.reportStats(r.ctx)
}

// RunOnce 认领并发布一批消息，返回本轮遇到的首个基础设施错误。
func (r *Relay) RunOnce(ctx context.Context) error {
	if r.pub == nil {
		return mqx.ErrUnavailable
	}
	now := time.Now()
	if now.Sub(r.lastReset) >= defaultResetEvery {
		if _, err := r.store.ResetStuck(ctx, now.Add(-r.stuckAfter)); err != nil {
			return err
		}
		r.lastReset = now
	}

	rows, err := r.store.Claim(ctx, r.batch)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}

	publishedByToken := make(map[string][]int64, 1)
	for _, row := range rows {
		msg, err := row.Message()
		if err != nil {
			r.fail(ctx, row, mqx.Permanent(err))
			continue
		}

		pubCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		_, err = r.pub.Publish(pubCtx, &msg)
		cancel()
		if err != nil {
			r.fail(ctx, row, err)
			continue
		}
		publishedByToken[row.LockToken] = append(publishedByToken[row.LockToken], row.ID)
	}

	for token, ids := range publishedByToken {
		if err := r.store.MarkPublished(ctx, ids, token); err != nil {
			return err
		}
	}
	return nil
}

func (r *Relay) fail(ctx context.Context, row Record, cause error) {
	// 坏消息直接死亡；broker 不可用则无限退避重试，不把整批事件判死。
	if mqx.IsPermanent(cause) || errors.Is(cause, mqx.ErrBadMessage) {
		r.markDead(ctx, row, cause)
		return
	}

	retries := row.RetryCount + 1
	if !errors.Is(cause, mqx.ErrUnavailable) && retries >= r.maxRetries {
		r.markDead(ctx, row, cause)
		return
	}

	delay := r.baseDelay << (retries - 1)
	if delay <= 0 || delay > r.maxDelay {
		delay = r.maxDelay
	}
	delay = time.Duration(float64(delay) * (0.8 + 0.4*rand.Float64()))

	next := time.Now().Add(delay)
	if err := r.store.MarkFailed(ctx, row.ID, row.LockToken, retries, next, cause); err != nil {
		logx.Errorf("[outbox] mark failed: id=%d error=%v", row.ID, err)
	}
}

func (r *Relay) markDead(ctx context.Context, row Record, cause error) {
	if err := r.store.MarkDead(ctx, row.ID, row.LockToken, cause); err != nil {
		logx.Errorf("[outbox] mark dead failed: id=%d error=%v", row.ID, err)
	}
	logx.Errorf("[outbox] event dead: id=%d event_id=%s retries=%d error=%v",
		row.ID, row.EventID, row.RetryCount+1, cause)
}

func (r *Relay) cleanup(ctx context.Context) {
	now := time.Now()
	if now.Sub(r.lastCleanup) < r.cleanupEvery {
		return
	}
	r.lastCleanup = now

	deleted, err := r.store.DeletePublished(ctx, now.Add(-r.retention))
	if err != nil {
		logx.Errorf("[outbox] cleanup: %v", err)
		return
	}
	if deleted > 0 {
		logx.Infof("[outbox] cleaned %d published events", deleted)
	}
}

func (r *Relay) reportStats(ctx context.Context) {
	now := time.Now()
	if now.Sub(r.lastStats) < r.statsEvery {
		return
	}
	r.lastStats = now

	stats, err := r.store.Stats(ctx)
	if err != nil {
		logx.Errorf("[outbox] stats: %v", err)
		return
	}
	var oldestAge time.Duration
	if stats.Oldest != nil {
		oldestAge = now.Sub(*stats.Oldest)
	}
	if r.metrics != nil {
		r.metrics.OutboxStats(stats.Pending, stats.Dead, oldestAge)
	}
	if stats.Dead > 0 {
		logx.Errorf("[outbox] dead events: %d", stats.Dead)
	}
	if stats.Pending > 100 {
		logx.Errorf("[outbox] pending backlog: %d", stats.Pending)
	}
	if stats.Oldest != nil && now.Sub(*stats.Oldest) > 5*time.Minute {
		logx.Errorf("[outbox] oldest pending: %s", now.Sub(*stats.Oldest))
	}
}
