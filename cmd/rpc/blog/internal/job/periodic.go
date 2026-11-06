package job

import (
	"context"
	"time"
)

// Periodic 是一个受 Guard 保护的周期任务。
//
// Start 与 Stop 由 Runner 串行调用，不对并发启停做额外保护。
type Periodic struct {
	name       string
	interval   time.Duration
	firstDelay time.Duration
	run        func(context.Context)

	cancel context.CancelFunc
	done   chan struct{}
}

// NewPeriodic 创建周期任务。firstDelay 为 0 时首轮在 interval 之后执行；
// run 内的 panic 由 Guard 捕获，循环继续。
func NewPeriodic(name string, interval, firstDelay time.Duration, run func(context.Context)) *Periodic {
	if interval <= 0 {
		interval = time.Minute
	}
	if firstDelay <= 0 {
		firstDelay = interval
	}
	return &Periodic{
		name:       name,
		interval:   interval,
		firstDelay: firstDelay,
		run:        run,
	}
}

// Start 启动循环；重复调用无副作用。
func (p *Periodic) Start() {
	if p.cancel != nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	p.done = make(chan struct{})

	go func() {
		defer close(p.done)
		timer := time.NewTimer(p.firstDelay)
		defer timer.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				Guard(p.name, func() { p.run(ctx) })
				timer.Reset(p.interval)
			}
		}
	}()
}

// Stop 停止循环并等待在途任务结束；nil 安全，重复调用无副作用。
func (p *Periodic) Stop(ctx context.Context) {
	if p == nil || p.cancel == nil {
		return
	}
	p.cancel()
	select {
	case <-p.done:
	case <-ctx.Done():
	}
}
