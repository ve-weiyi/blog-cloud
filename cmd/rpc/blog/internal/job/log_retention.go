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
	// logRetentionInterval 清理周期。
	logRetentionInterval = 24 * time.Hour
	// logRetentionFirstDelay 启动后的首轮延迟，避开启动流量。
	logRetentionFirstDelay = 10 * time.Minute
	// logRetentionBatch 每批删除行数。
	logRetentionBatch = 2000
	// logRetentionMaxRounds 单表单轮最多删除的批数，避免首次执行长时间占用数据库。
	logRetentionMaxRounds = 100
)

const day = 24 * time.Hour

// logPolicy 一张日志表的保留策略。
//
// 保留窗口是产品决定，不是技术参数：操作日志用于追责，保留最久；
// 访问日志量最大、排障价值随时间快速下降，保留最短。
type logPolicy struct {
	table string
	keep  time.Duration
}

// LogRetention 按表清理超过保留期的日志。
type LogRetention struct {
	db   *gorm.DB
	loop *Periodic
}

// NewLogRetention 创建日志保留清理任务。
func NewLogRetention(db *gorm.DB, interval time.Duration) *LogRetention {
	if interval <= 0 {
		interval = logRetentionInterval
	}
	r := &LogRetention{db: db}
	r.loop = NewPeriodic("log-retention", interval, logRetentionFirstDelay, func(ctx context.Context) {
		// 多副本时只让一个实例删，避免同一批行互相等锁
		release, ok, err := AcquireLock(ctx, db, "blog:log-retention")
		if err != nil {
			logx.WithContext(ctx).Errorf("日志保留清理：获取任务锁失败: %v", err)
			return
		}
		if !ok {
			return
		}
		defer release()

		deleted, err := r.RunOnce(ctx, time.Now(), logPolicies(db), logRetentionBatch)
		if err != nil {
			logx.WithContext(ctx).Errorf("日志保留清理失败: %v", err)
			return
		}
		if deleted > 0 {
			logx.WithContext(ctx).Infof("日志保留清理：删除 %d 行", deleted)
		}
	})
	return r
}

// Start 启动后台循环；重复调用无副作用。
func (r *LogRetention) Start() {
	r.loop.Start()
}

// Stop 停止循环并等待在途批次结束；nil 安全。
func (r *LogRetention) Stop(ctx context.Context) {
	if r == nil || r.loop == nil {
		return
	}
	r.loop.Stop(ctx)
}

// RunOnce 按策略清理一轮，返回删除的总行数。
func (r *LogRetention) RunOnce(ctx context.Context, now time.Time, policies []logPolicy, batch int) (int64, error) {
	if batch <= 0 {
		batch = logRetentionBatch
	}

	var total int64
	for _, policy := range policies {
		deleted, err := r.drainTable(ctx, policy.table, now.Add(-policy.keep), batch)
		if err != nil {
			return total, err
		}
		total += deleted
		if deleted > 0 {
			logx.WithContext(ctx).Infof("日志保留清理：%s 删除 %d 行", policy.table, deleted)
		}
	}
	return total, nil
}

// drainTable 分批删除一张表里超过 cutoff 的行。
func (r *LogRetention) drainTable(ctx context.Context, table string, cutoff time.Time, batch int) (int64, error) {
	var deleted int64
	for round := 0; round < logRetentionMaxRounds; round++ {
		// 表名来自模型常量，不是用户输入；时间与批大小走参数绑定。
		result := r.db.WithContext(ctx).Exec(
			fmt.Sprintf("DELETE FROM %s WHERE created_at < ? ORDER BY created_at LIMIT ?", table),
			cutoff, batch)
		if result.Error != nil {
			return deleted, fmt.Errorf("job: delete %s: %w", table, result.Error)
		}
		deleted += result.RowsAffected
		if result.RowsAffected < int64(batch) {
			return deleted, nil
		}
	}
	logx.WithContext(ctx).Errorf(
		"日志保留清理：%s 达到单轮上限 %d 批，剩余留待下一轮", table, logRetentionMaxRounds)
	return deleted, nil
}

// logPolicies 返回各日志表的保留策略。
func logPolicies(db *gorm.DB) []logPolicy {
	return []logPolicy{
		{table: model.NewTLoginLogModel(db).TableName(), keep: 180 * day},
		{table: model.NewTOperationLogModel(db).TableName(), keep: 365 * day},
		{table: model.NewTVisitLogModel(db).TableName(), keep: 90 * day},
		{table: model.NewTUploadLogModel(db).TableName(), keep: 180 * day},
	}
}
