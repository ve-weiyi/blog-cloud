package rabbitmq

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

type ForeverBackoff struct {
	bkf    backoff.BackOff
	ticker *backoff.Ticker
}

// NewForeverBackoff 创建一个永不停止的 backoff
func NewForeverBackoff() *ForeverBackoff {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 0 // 一直重试
	b.MaxInterval = 10 * time.Second
	b.InitialInterval = 100 * time.Millisecond

	return &ForeverBackoff{
		bkf:    b,
		ticker: backoff.NewTicker(b),
	}
}

func (b *ForeverBackoff) Reset() {
	b.bkf.Reset()
}

func (b *ForeverBackoff) Wait() {
	<-b.ticker.C
}
