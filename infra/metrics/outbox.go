package metrics

import (
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/metric"
)

// Outbox 记录 outbox 积压状态。
type Outbox struct {
	pending metric.GaugeVec
	dead    metric.GaugeVec
	oldest  metric.GaugeVec
}

var (
	outboxOnce    sync.Once
	outboxMetrics *Outbox
)

// NewOutbox 返回进程级 outbox 指标集。Prometheus 按指标名注册，重复构造会 panic，
// 因此同一个进程内只创建一次。
func NewOutbox() *Outbox {
	outboxOnce.Do(func() {
		outboxMetrics = &Outbox{
			pending: metric.NewGaugeVec(&metric.GaugeVecOpts{
				Namespace: metricNamespace,
				Subsystem: "outbox",
				Name:      "pending",
				Help:      "待发布的 outbox 事件数",
			}),
			dead: metric.NewGaugeVec(&metric.GaugeVecOpts{
				Namespace: metricNamespace,
				Subsystem: "outbox",
				Name:      "dead",
				Help:      "标记为 dead、等待人工处理的 outbox 事件数",
			}),
			oldest: metric.NewGaugeVec(&metric.GaugeVecOpts{
				Namespace: metricNamespace,
				Subsystem: "outbox",
				Name:      "oldest_pending_seconds",
				Help:      "最老的待发布事件距今的秒数",
			}),
		}
	})
	return outboxMetrics
}

// OutboxStats 记录一次积压快照。
func (o *Outbox) OutboxStats(pending, dead int64, oldestAge time.Duration) {
	o.pending.Set(float64(pending))
	o.dead.Set(float64(dead))
	o.oldest.Set(oldestAge.Seconds())
}
