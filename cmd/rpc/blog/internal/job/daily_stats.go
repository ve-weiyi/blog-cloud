package job

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/model"
)

const (
	// dailyStatsCatchUp 启动后的补跑延迟，避开启动流量。
	dailyStatsCatchUp = 30 * time.Second
	// dailyStatsBackfillDays 启动时回填的天数，覆盖短期停机。
	dailyStatsBackfillDays = 7
	// dailyStatsHour / dailyStatsMinute 每日执行时刻（本地时区）。
	dailyStatsHour   = 0
	dailyStatsMinute = 5
)

// DailyStats 补齐 t_daily_stats 的用户指标。
//
// UV/PV 由访问日志实时累加；new_users / total_users / active_users 需要按天
// 聚合。任务幂等：同一天重复执行只覆盖这三列，不碰 UV/PV。
type DailyStats struct {
	db *gorm.DB

	cancel context.CancelFunc
	done   chan struct{}
}

// NewDailyStats 创建每日统计任务。
func NewDailyStats(db *gorm.DB) *DailyStats {
	return &DailyStats{db: db}
}

// Start 启动后台循环；重复调用无副作用。
func (d *DailyStats) Start() {
	if d.cancel != nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	d.cancel = cancel
	d.done = make(chan struct{})

	go func() {
		defer close(d.done)

		// 启动补跑最近几天的数据；停机跨过零点时仍能补齐
		timer := time.NewTimer(dailyStatsCatchUp)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
		}
		Guard("daily-stats", func() {
			if err := d.runRecent(ctx, dailyStatsBackfillDays); err != nil {
				logx.WithContext(ctx).Errorf("每日统计补跑失败: %v", err)
			}
		})

		for {
			timer.Reset(time.Until(nextDaily(time.Now(), dailyStatsHour, dailyStatsMinute)))
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
			}
			Guard("daily-stats", func() {
				if err := d.RunOnce(ctx, time.Now().AddDate(0, 0, -1)); err != nil {
					logx.WithContext(ctx).Errorf("每日统计失败: %v", err)
				}
			})
		}
	}()
}

// Stop 停止循环并等待在途批次结束；nil 安全。
func (d *DailyStats) Stop(ctx context.Context) {
	if d == nil || d.cancel == nil {
		return
	}
	d.cancel()
	select {
	case <-d.done:
	case <-ctx.Done():
	}
}

// RunOnce 聚合指定自然日的用户指标并写入 t_daily_stats。
func (d *DailyStats) RunOnce(ctx context.Context, day time.Time) error {
	day = time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := day.AddDate(0, 0, 1)

	release, ok, err := AcquireLock(ctx, d.db, "blog:daily-stats")
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	defer release()

	userTable := model.NewTUserModel(d.db).TableName()
	visitTable := model.NewTVisitLogModel(d.db).TableName()
	statsTable := model.NewTDailyStatsModel(d.db).TableName()

	newUsers, err := d.count(ctx, userTable, "created_at >= ? AND created_at < ?", day, end)
	if err != nil {
		return fmt.Errorf("job: count new users: %w", err)
	}
	totalUsers, err := d.count(ctx, userTable, "created_at < ?", end)
	if err != nil {
		return fmt.Errorf("job: count total users: %w", err)
	}
	// 空 user_id 是游客，不算活跃用户
	activeUsers, err := d.countDistinct(ctx, visitTable, "user_id", "created_at >= ? AND created_at < ? AND user_id <> ''", day, end)
	if err != nil {
		return fmt.Errorf("job: count active users: %w", err)
	}

	// 行不存在时插入；UV/PV 为 0，累计值沿用前一天的快照。
	var prev struct {
		TotalUV int64 `gorm:"column:total_uv_count"`
		TotalPV int64 `gorm:"column:total_pv_count"`
	}
	if err := d.db.WithContext(ctx).Table(statsTable).
		Select("total_uv_count, total_pv_count").
		Where("date < ?", day).
		Order("date DESC").
		Limit(1).
		Scan(&prev).Error; err != nil {
		return fmt.Errorf("job: read previous totals: %w", err)
	}

	sql := fmt.Sprintf(`INSERT INTO %s
		(date, new_users, total_users, active_users, uv_count, pv_count, total_uv_count, total_pv_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, 0, 0, ?, ?, NOW(3), NOW(3))
		ON DUPLICATE KEY UPDATE
			new_users = VALUES(new_users),
			total_users = VALUES(total_users),
			active_users = VALUES(active_users),
			updated_at = NOW(3)`, statsTable)
	if result := d.db.WithContext(ctx).Exec(sql, day, newUsers, totalUsers, activeUsers, prev.TotalUV, prev.TotalPV); result.Error != nil {
		return fmt.Errorf("job: upsert daily stats: %w", result.Error)
	}

	logx.WithContext(ctx).Infof(
		"每日统计：%s 新增用户 %d、累计用户 %d、活跃用户 %d",
		day.Format("2006-01-02"), newUsers, totalUsers, activeUsers)
	return nil
}

// runRecent 回填最近 days 天（不含今天）。
func (d *DailyStats) runRecent(ctx context.Context, days int) error {
	now := time.Now()
	for i := days; i >= 1; i-- {
		if err := d.RunOnce(ctx, now.AddDate(0, 0, -i)); err != nil {
			return err
		}
	}
	return nil
}

func (d *DailyStats) count(ctx context.Context, table, conditions string, args ...interface{}) (int64, error) {
	var n int64
	if err := d.db.WithContext(ctx).Table(table).Where(conditions, args...).Count(&n).Error; err != nil {
		return 0, err
	}
	return n, nil
}

func (d *DailyStats) countDistinct(ctx context.Context, table, column, conditions string, args ...interface{}) (int64, error) {
	var n int64
	if err := d.db.WithContext(ctx).Table(table).Where(conditions, args...).Distinct(column).Count(&n).Error; err != nil {
		return 0, err
	}
	return n, nil
}

// nextDaily 返回 now 之后最近的一个 hour:minute。
func nextDaily(now time.Time, hour, minute int) time.Time {
	next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
	if !next.After(now) {
		next = next.AddDate(0, 0, 1)
	}
	return next
}
