// Package metrics 把运行时事件接入 Prometheus。
package metrics

import (
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/metric"

	"github.com/ve-weiyi/vkit/adapter/mqx"
)

const (
	metricNamespace = "blog"
	metricSubsystem = "mqx"
)

// durationBuckets 覆盖由毫秒级内存投递到分钟级外呼的耗时区间。
var durationBuckets = []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5, 10, 30}

// MQX 实现 mqx.Metrics。
//
// 指标注册到默认 registry，只在 go-zero prometheus agent 启动后记录；
// agent 未启用时是空操作，不影响调用路径。
type MQX struct {
	publishDuration metric.HistogramVec
	publishTotal    metric.CounterVec
	consumeDuration metric.HistogramVec
	consumeTotal    metric.CounterVec
	stoppedTotal    metric.CounterVec
}

var (
	_ mqx.Metrics = (*MQX)(nil)

	mqxOnce    sync.Once
	mqxMetrics *MQX
)

// NewMQX 返回进程级指标集。Prometheus 按指标名注册，重复构造会 panic，
// 因此同一个进程内只创建一次。
func NewMQX() *MQX {
	mqxOnce.Do(func() {
		mqxMetrics = &MQX{
			publishDuration: metric.NewHistogramVec(&metric.HistogramVecOpts{
				Namespace: metricNamespace,
				Subsystem: metricSubsystem,
				Name:      "publish_duration_seconds",
				Help:      "消息发布耗时",
				Labels:    []string{"destination", "routing_key", "result"},
				Buckets:   durationBuckets,
			}),
			publishTotal: metric.NewCounterVec(&metric.CounterVecOpts{
				Namespace: metricNamespace,
				Subsystem: metricSubsystem,
				Name:      "publish_total",
				Help:      "消息发布次数",
				Labels:    []string{"destination", "routing_key", "result"},
			}),
			consumeDuration: metric.NewHistogramVec(&metric.HistogramVecOpts{
				Namespace: metricNamespace,
				Subsystem: metricSubsystem,
				Name:      "consume_duration_seconds",
				Help:      "单次 handler 执行耗时（含进程内重试的每次尝试）",
				Labels:    []string{"destination", "routing_key", "group", "result"},
				Buckets:   durationBuckets,
			}),
			consumeTotal: metric.NewCounterVec(&metric.CounterVecOpts{
				Namespace: metricNamespace,
				Subsystem: metricSubsystem,
				Name:      "consume_total",
				Help:      "handler 执行次数",
				Labels:    []string{"destination", "routing_key", "group", "result"},
			}),
			stoppedTotal: metric.NewCounterVec(&metric.CounterVecOpts{
				Namespace: metricNamespace,
				Subsystem: metricSubsystem,
				Name:      "subscription_stopped_total",
				Help:      "订阅终止次数",
				Labels:    []string{"destination", "group", "result"},
			}),
		}
	})
	return mqxMetrics
}

// PublishCompleted 记录一次发布。
func (m *MQX) PublishCompleted(destination, routingKey string, dur time.Duration, err error) {
	result := resultOf(err)
	m.publishDuration.ObserveFloat(dur.Seconds(), destination, routingKey, result)
	m.publishTotal.Inc(destination, routingKey, result)
}

// ConsumeCompleted 记录一次 handler 执行。
func (m *MQX) ConsumeCompleted(destination, routingKey, group string, _ int, dur time.Duration, err error) {
	result := resultOf(err)
	m.consumeDuration.ObserveFloat(dur.Seconds(), destination, routingKey, group, result)
	m.consumeTotal.Inc(destination, routingKey, group, result)
}

// SubscriptionStopped 记录一次订阅终止。
func (m *MQX) SubscriptionStopped(destination, group string, err error) {
	m.stoppedTotal.Inc(destination, group, resultOf(err))
}

func resultOf(err error) string {
	if err != nil {
		return "error"
	}
	return "ok"
}
