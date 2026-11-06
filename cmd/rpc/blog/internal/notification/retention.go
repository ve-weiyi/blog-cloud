package notification

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/job"
	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
)

const (
	// retentionInterval 清理周期。
	retentionInterval = 24 * time.Hour
	// retentionFirstDelay 启动后的首轮延迟：既不与启动流量抢资源，也保证
	// 频繁重启的进程仍能执行清理。
	retentionFirstDelay = 5 * time.Minute
	// retentionBatch 每批删除的条目数。条目与其投递在同一事务内删除。
	retentionBatch = 500
	// retentionMaxRounds 单场景单轮最多处理的批数，避免首次执行长时间占用数据库。
	retentionMaxRounds = 200
	// retentionEventAge 终态事件在无引用条目后再保留一段时间，供排障与对账。
	retentionEventAge = 7 * 24 * time.Hour
)

// RetentionWorker 按场景保留期清理收件箱条目与终态事件。
//
// 保留期由场景目录声明：已读先于未读清理，未读的保留期更长。清理必须走
// (scene, read_at, created_at) 索引并分批执行，否则表越大锁越久。
type RetentionWorker struct {
	svc  *Service
	loop *job.Periodic
}

// NewRetentionWorker 创建保留期清理任务。
func NewRetentionWorker(svc *Service, interval time.Duration) *RetentionWorker {
	if interval <= 0 {
		interval = retentionInterval
	}
	w := &RetentionWorker{svc: svc}
	w.loop = job.NewPeriodic("notify-retention", interval, retentionFirstDelay, func(ctx context.Context) {
		// 清理必须单活：多副本并发删同一批只会互相等锁
		release, ok, err := job.AcquireLock(ctx, w.svc.db, "blog:notify:retention")
		if err != nil {
			logx.WithContext(ctx).Errorf("通知保留清理：获取任务锁失败: %v", err)
			return
		}
		if !ok {
			return
		}
		defer release()

		entries, events, err := w.svc.RunRetention(ctx, time.Now(), retentionBatch)
		if err != nil {
			logx.WithContext(ctx).Errorf("通知保留清理失败: %v", err)
			return
		}
		if entries > 0 || events > 0 {
			logx.WithContext(ctx).Infof("通知保留清理：条目 %d 条、事件 %d 条", entries, events)
		}
	})
	return w
}

// Start 启动后台循环；重复调用无副作用。首轮在启动 5 分钟后执行，
// 之后按 interval 重复；清理按批推进且幂等，重启触发的额外一轮不会放大副作用。
func (w *RetentionWorker) Start() {
	w.loop.Start()
}

// Stop 停止循环并等待在途批次结束；nil 安全。
func (w *RetentionWorker) Stop(ctx context.Context) {
	if w == nil || w.loop == nil {
		return
	}
	w.loop.Stop(ctx)
}

// RunRetention 执行一轮清理，返回删除的条目数与事件数。
func (s *Service) RunRetention(ctx context.Context, now time.Time, batch int) (int64, int64, error) {
	return s.runRetention(ctx, now, batch, notify.All())
}

// runRetention 按给定场景执行清理；场景列表由调用方提供，便于用固定目录测试。
func (s *Service) runRetention(ctx context.Context, now time.Time, batch int, scenes []notify.Scene) (int64, int64, error) {
	if batch <= 0 {
		batch = retentionBatch
	}

	// 排序只影响日志与测试的稳定性，复制一份避免改写调用方切片
	scenes = append([]notify.Scene(nil), scenes...)
	sort.Slice(scenes, func(i, j int) bool { return scenes[i].Key < scenes[j].Key })

	var entries int64
	for _, scene := range scenes {
		if scene.ReadTTL > 0 {
			n, err := s.drainScene(ctx, scene.Key, true, now.Add(-scene.ReadTTL), batch)
			if err != nil {
				return entries, 0, fmt.Errorf("notification: retention scene %s read: %w", scene.Key, err)
			}
			entries += n
		}
		if scene.UnreadTTL > 0 {
			n, err := s.drainScene(ctx, scene.Key, false, now.Add(-scene.UnreadTTL), batch)
			if err != nil {
				return entries, 0, fmt.Errorf("notification: retention scene %s unread: %w", scene.Key, err)
			}
			entries += n
		}
	}

	var events int64
	cutoff := now.Add(-retentionEventAge)
	for round := 0; round < retentionMaxRounds; round++ {
		n, err := s.pruneEvents(ctx, cutoff, batch)
		if err != nil {
			return entries, events, fmt.Errorf("notification: prune events: %w", err)
		}
		events += n
		if n < int64(batch) {
			break
		}
	}
	return entries, events, nil
}

// drainScene 分批删除某场景下已过期的条目，直到不足一批或达到批数上限。
func (s *Service) drainScene(ctx context.Context, scene string, read bool, cutoff time.Time, batch int) (int64, error) {
	var deleted int64
	for round := 0; round < retentionMaxRounds; round++ {
		n, err := s.deleteExpiredBatch(ctx, scene, read, cutoff, batch)
		if err != nil {
			return deleted, err
		}
		deleted += n
		if n < int64(batch) {
			return deleted, nil
		}
	}
	logx.WithContext(ctx).Errorf(
		"通知保留清理：场景 %s（read=%t）达到单轮上限 %d 批，剩余留待下一轮",
		scene, read, retentionMaxRounds)
	return deleted, nil
}

// deleteExpiredBatch 删除一批过期条目及其投递记录。
func (s *Service) deleteExpiredBatch(ctx context.Context, scene string, read bool, cutoff time.Time, limit int) (int64, error) {
	ids := make([]int64, 0, limit)
	q := s.db.WithContext(ctx).Table(s.notifs.TableName()).
		Select("id").
		Where("scene = ?", scene).
		Limit(limit)
	// 与索引 (scene, read_at, created_at) 同序，避免为取最早一批而对整段过期区间排序
	if read {
		q = q.Where("read_at IS NOT NULL AND read_at < ?", cutoff).Order("read_at").Order("id")
	} else {
		q = q.Where("read_at IS NULL AND created_at < ?", cutoff).Order("created_at").Order("id")
	}
	if err := q.Pluck("id", &ids).Error; err != nil {
		return 0, fmt.Errorf("notification: query expired entries: %w", err)
	}
	if len(ids) == 0 {
		return 0, nil
	}

	var deleted int64
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := s.deliveries.WithTx(tx).DeleteBatch(ctx, "notification_id IN ?", ids); err != nil {
			return fmt.Errorf("notification: delete deliveries of expired entries: %w", err)
		}
		n, err := s.notifs.WithTx(tx).DeleteBatch(ctx, "id IN ?", ids)
		if err != nil {
			return fmt.Errorf("notification: delete expired entries: %w", err)
		}
		deleted = n
		return nil
	})
	return deleted, err
}

// pruneEvents 回收已终止且已无引用条目的旧事件。
func (s *Service) pruneEvents(ctx context.Context, cutoff time.Time, limit int) (int64, error) {
	ids := make([]int64, 0, limit)
	notExists := fmt.Sprintf(
		"NOT EXISTS (SELECT 1 FROM %s WHERE event_id = %s.id)",
		s.notifs.TableName(), s.events.TableName())
	err := s.db.WithContext(ctx).Table(s.events.TableName()).
		Select("id").
		Where("status IN ?", []string{StatusDone, StatusFailed, StatusTerminated}).
		Where("updated_at < ?", cutoff).
		Where(notExists).
		Order("id").
		Limit(limit).
		Pluck("id", &ids).Error
	if err != nil {
		return 0, fmt.Errorf("notification: query terminal events: %w", err)
	}
	if len(ids) == 0 {
		return 0, nil
	}
	return s.events.DeleteBatch(ctx, "id IN ?", ids)
}
