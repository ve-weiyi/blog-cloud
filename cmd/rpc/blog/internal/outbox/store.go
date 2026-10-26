package outbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ve-weiyi/vkit/adapter/mqx"
)

// Store 是 outbox 的存储契约，便于 relay 用内存实现测试。
type Store interface {
	// Enqueue 在调用方事务内写入一条待发布消息。
	Enqueue(ctx context.Context, tx *gorm.DB, msg mqx.Message) error
	// Claim 认领一批到期消息，并标记为 publishing。
	Claim(ctx context.Context, limit int) ([]Record, error)
	MarkPublished(ctx context.Context, ids []int64, token string) error
	MarkFailed(ctx context.Context, id int64, token string, retries int, nextAttempt time.Time, cause error) error
	MarkDead(ctx context.Context, id int64, token string, cause error) error
	ResetStuck(ctx context.Context, olderThan time.Time) (int64, error)
	DeletePublished(ctx context.Context, olderThan time.Time) (int64, error)
	Stats(ctx context.Context) (Stats, error)
}

// Stats 是 relay 用于日志告警的概览。
type Stats struct {
	Pending int64
	Oldest  *time.Time
	Dead    int64
}

// GORMStore 是 Store 的 MySQL 实现。
type GORMStore struct {
	db *gorm.DB
}

var _ Store = (*GORMStore)(nil)

// NewGORMStore 创建 outbox 仓储。
func NewGORMStore(db *gorm.DB) *GORMStore {
	return &GORMStore{db: db}
}

// Ping 校验 outbox 表可用；失败通常意味着迁移脚本尚未执行。
func (s *GORMStore) Ping(ctx context.Context) error {
	if _, err := s.Stats(ctx); err != nil {
		return fmt.Errorf("outbox unavailable, run protocol/sql/2026-10-26-t_outbox_event.sql: %w", err)
	}
	return nil
}

// Enqueue 在调用方事务内写入 outbox。
func (s *GORMStore) Enqueue(ctx context.Context, tx *gorm.DB, msg mqx.Message) error {
	if tx == nil {
		return fmt.Errorf("outbox: transaction is required")
	}
	if msg.ID == "" {
		return fmt.Errorf("outbox: event id is required")
	}
	if msg.Destination == "" {
		return fmt.Errorf("outbox: destination is required")
	}

	headers, err := json.Marshal(msg.Headers)
	if err != nil {
		return fmt.Errorf("outbox: encode headers: %w", err)
	}
	now := time.Now()
	occurredAt := msg.Timestamp
	if occurredAt.IsZero() {
		occurredAt = now
	}

	row := &Record{
		EventID:       msg.ID,
		Destination:   msg.Destination,
		RoutingKey:    msg.RoutingKey,
		MessageKey:    msg.Key,
		Headers:       string(headers),
		Payload:       msg.Body,
		Status:        StatusPending,
		NextAttemptAt: now,
		OccurredAt:    occurredAt,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := tx.WithContext(ctx).Create(row).Error; err != nil {
		return fmt.Errorf("outbox: enqueue: %w", err)
	}
	return nil
}

// Claim 认领到期消息。MySQL 8 的 SKIP LOCKED 支持多副本 relay 并行认领。
func (s *GORMStore) Claim(ctx context.Context, limit int) ([]Record, error) {
	if limit <= 0 {
		limit = 1
	}

	var rows []Record
	now := time.Now()
	token := uuid.NewString()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("status = ? AND next_attempt_at <= ?", StatusPending, now).
			Order("id").
			Limit(limit).
			Find(&rows).Error; err != nil {
			return err
		}
		if len(rows) == 0 {
			return nil
		}

		ids := make([]int64, len(rows))
		for i := range rows {
			ids[i] = rows[i].ID
		}
		return tx.Model(&Record{}).Where("id IN ?", ids).Updates(map[string]any{
			"status":     StatusPublishing,
			"locked_at":  now,
			"lock_token": token,
			"updated_at": now,
		}).Error
	})
	if err != nil {
		return nil, fmt.Errorf("outbox: claim: %w", err)
	}
	for i := range rows {
		rows[i].LockToken = token
	}
	return rows, nil
}

// MarkPublished 标记一批消息已发布。
func (s *GORMStore) MarkPublished(ctx context.Context, ids []int64, token string) error {
	if len(ids) == 0 {
		return nil
	}
	if token == "" {
		return fmt.Errorf("outbox: lock token is required")
	}
	now := time.Now()
	if err := s.db.WithContext(ctx).Model(&Record{}).
		Where("id IN ? AND lock_token = ?", ids, token).
		Updates(map[string]any{
			"status":       StatusPublished,
			"published_at": now,
			"locked_at":    nil,
			"lock_token":   "",
			"last_error":   "",
			"updated_at":   now,
		}).Error; err != nil {
		return fmt.Errorf("outbox: mark published: %w", err)
	}
	return nil
}

// MarkFailed 记录一次可重试失败并安排下次尝试。
func (s *GORMStore) MarkFailed(ctx context.Context, id int64, token string, retries int, nextAttempt time.Time, cause error) error {
	if token == "" {
		return fmt.Errorf("outbox: lock token is required")
	}
	errMsg := ""
	if cause != nil {
		errMsg = cause.Error()
	}
	if err := s.db.WithContext(ctx).Model(&Record{}).
		Where("id = ? AND lock_token = ?", id, token).
		Updates(map[string]any{
			"status":          StatusPending,
			"retry_count":     retries,
			"next_attempt_at": nextAttempt,
			"locked_at":       nil,
			"lock_token":      "",
			"last_error":      errMsg,
			"updated_at":      time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("outbox: mark failed: %w", err)
	}
	return nil
}

// MarkDead 标记为不可重试，等待人工处理。
func (s *GORMStore) MarkDead(ctx context.Context, id int64, token string, cause error) error {
	if token == "" {
		return fmt.Errorf("outbox: lock token is required")
	}
	errMsg := ""
	if cause != nil {
		errMsg = cause.Error()
	}
	if err := s.db.WithContext(ctx).Model(&Record{}).
		Where("id = ? AND lock_token = ?", id, token).
		Updates(map[string]any{
			"status":     StatusDead,
			"locked_at":  nil,
			"lock_token": "",
			"last_error": errMsg,
			"updated_at": time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("outbox: mark dead: %w", err)
	}
	return nil
}

// ResetStuck 把认领后超时未完成的消息放回待发布队列。
func (s *GORMStore) ResetStuck(ctx context.Context, olderThan time.Time) (int64, error) {
	res := s.db.WithContext(ctx).Model(&Record{}).
		Where("status = ? AND locked_at < ?", StatusPublishing, olderThan).
		Updates(map[string]any{
			"status":     StatusPending,
			"locked_at":  nil,
			"lock_token": "",
			"updated_at": time.Now(),
		})
	if res.Error != nil {
		return 0, fmt.Errorf("outbox: reset stuck: %w", res.Error)
	}
	return res.RowsAffected, nil
}

// DeletePublished 清理保留期之外的已发布消息。
func (s *GORMStore) DeletePublished(ctx context.Context, olderThan time.Time) (int64, error) {
	res := s.db.WithContext(ctx).
		Where("status = ? AND published_at < ?", StatusPublished, olderThan).
		Delete(&Record{})
	if res.Error != nil {
		return 0, fmt.Errorf("outbox: delete published: %w", res.Error)
	}
	return res.RowsAffected, nil
}

// Stats 返回积压概览；oldest 为 nil 表示没有待发布消息。
func (s *GORMStore) Stats(ctx context.Context) (Stats, error) {
	var stats Stats
	if err := s.db.WithContext(ctx).Model(&Record{}).
		Where("status = ?", StatusPending).Count(&stats.Pending).Error; err != nil {
		return Stats{}, fmt.Errorf("outbox: count pending: %w", err)
	}

	var oldest Record
	err := s.db.WithContext(ctx).Select("next_attempt_at").
		Where("status = ?", StatusPending).
		Order("next_attempt_at").
		First(&oldest).Error
	if err == nil {
		t := oldest.NextAttemptAt
		stats.Oldest = &t
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return Stats{}, fmt.Errorf("outbox: oldest pending: %w", err)
	}

	if err := s.db.WithContext(ctx).Model(&Record{}).
		Where("status = ?", StatusDead).Count(&stats.Dead).Error; err != nil {
		return Stats{}, fmt.Errorf("outbox: count dead: %w", err)
	}
	return stats, nil
}
