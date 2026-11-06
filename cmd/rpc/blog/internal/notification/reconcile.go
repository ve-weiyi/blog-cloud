package notification

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/job"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/model"
)

const (
	// reconcileStaleAfter 超过该时长仍停在待扇出/扇出中的事件视为停滞。
	// 必须大于一批扇出的正常耗时，避免与正常唤醒抢同一批。
	reconcileStaleAfter = 5 * time.Minute
	// reconcileBatch 单轮最多补发的唤醒消息数。
	reconcileBatch = 100
	// reconcileInterval 对账周期。
	reconcileInterval = time.Minute
)

// ReconcileWorker 周期性补发丢失的扇出唤醒消息。
//
// 事件表兼作工作单元：正常路径靠 outbox 里的唤醒消息推进。消息丢失或进程在
// 批次之间退出时，事件会停在 pending/fanning，没有别的任务会推进它。
type ReconcileWorker struct {
	svc  *Service
	loop *job.Periodic
}

// NewReconcileWorker 创建扇出对账任务。
func NewReconcileWorker(svc *Service, interval time.Duration) *ReconcileWorker {
	if interval <= 0 {
		interval = reconcileInterval
	}
	w := &ReconcileWorker{svc: svc}
	w.loop = job.NewPeriodic("notify-reconcile", interval, 0, func(ctx context.Context) {
		// 多副本时只让一个实例补发；锁拿不到说明别的实例正在做
		release, ok, err := job.AcquireLock(ctx, w.svc.db, "blog:notify:reconcile")
		if err != nil {
			logx.WithContext(ctx).Errorf("通知扇出对账：获取任务锁失败: %v", err)
			return
		}
		if !ok {
			return
		}
		defer release()

		n, err := w.svc.ReconcileStalled(ctx, time.Now(), reconcileStaleAfter, reconcileBatch)
		if err != nil {
			logx.WithContext(ctx).Errorf("通知扇出对账失败: %v", err)
			return
		}
		if n > 0 {
			logx.WithContext(ctx).Infof("通知扇出对账：补发 %d 条唤醒", n)
		}
	})
	return w
}

// Start 启动后台循环；重复调用无副作用。
func (w *ReconcileWorker) Start() {
	w.loop.Start()
}

// Stop 停止循环并等待在途批次结束；nil 安全。
func (w *ReconcileWorker) Stop(ctx context.Context) {
	if w == nil || w.loop == nil {
		return
	}
	w.loop.Stop(ctx)
}

// ReconcileStalled 为停滞事件补发扇出唤醒，返回补发条数。
//
// 只补发唤醒、不复位游标：Fanout 每批独立事务并持久化进度，继续处理即可。
// 重放会从头扫描受众，让扇出期间新注册的用户收到旧通知。
func (s *Service) ReconcileStalled(ctx context.Context, now time.Time, staleAfter time.Duration, limit int) (int, error) {
	if s.outbox == nil {
		return 0, errors.New("notification: outbox is not configured")
	}
	if limit <= 0 {
		limit = reconcileBatch
	}

	rows := make([]*model.TNotificationEvent, 0, limit)
	if err := s.db.WithContext(ctx).Table(s.events.TableName()).
		Where("status IN ? AND updated_at < ?", []string{StatusPending, StatusFanning}, now.Add(-staleAfter)).
		Order("id ASC").
		Limit(limit).
		Find(&rows).Error; err != nil {
		return 0, fmt.Errorf("notification: query stalled events: %w", err)
	}

	nudged := 0
	for _, ev := range rows {
		msg, err := mq.FanoutMessage.Message(&mq.NotificationFanoutEvent{EventID: ev.Id})
		if err != nil {
			logx.WithContext(ctx).Errorf("通知扇出对账：构造事件 %d 的唤醒失败: %v", ev.Id, err)
			continue
		}
		// 对账没有业务写入，直接以 db 作为执行句柄；outbox 的单条 INSERT 本身是原子的
		if err := s.outbox.Enqueue(ctx, s.db, msg); err != nil {
			logx.WithContext(ctx).Errorf("通知扇出对账：事件 %d 入队失败: %v", ev.Id, err)
			continue
		}

		// 刷新 updated_at：同一事件在下一轮之前不再重复补发，避免 MQ 长期不可用时
		// outbox 堆积重复唤醒。扇出进度仍由 Fanout 自己维护。
		if _, err := s.events.UpdateFields(ctx,
			map[string]interface{}{"updated_at": now}, "id = ?", ev.Id); err != nil {
			logx.WithContext(ctx).Errorf("通知扇出对账：刷新事件 %d 时间失败: %v", ev.Id, err)
		}
		nudged++
	}

	// 全部入队失败说明 outbox 或序列化出了系统性问题；回传错误，
	// 否则调用方会把「一条都没补上」当成「没有停滞事件」。
	if len(rows) > 0 && nudged == 0 {
		return 0, fmt.Errorf("notification: none of %d stalled events could be enqueued", len(rows))
	}
	return nudged, nil
}
