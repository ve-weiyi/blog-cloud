package outbox

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ve-weiyi/vkit/adapter/mqx"
)

// TableName 是事务性 outbox 表名。
const TableName = "t_outbox_event"

const (
	StatusPending    = "pending"
	StatusPublishing = "publishing"
	StatusPublished  = "published"
	StatusDead       = "dead"
)

// Record 是 outbox 表的一行。
//
// 索引标签与建表语句一一对应：AutoMigrate 据此创建的表必须能直接承载 relay。
type Record struct {
	ID            int64      `gorm:"column:id;type:bigint;primaryKey;autoIncrement" json:"id"`
	EventID       string     `gorm:"column:event_id;type:varchar(64);not null;default:'';uniqueIndex:uk_event_id" json:"event_id"`
	Destination   string     `gorm:"column:destination;type:varchar(128);not null;default:''" json:"destination"`
	RoutingKey    string     `gorm:"column:routing_key;type:varchar(128);not null;default:''" json:"routing_key"`
	MessageKey    string     `gorm:"column:message_key;type:varchar(255);not null;default:''" json:"message_key"`
	Headers       string     `gorm:"column:headers;type:text" json:"headers"`
	Payload       []byte     `gorm:"column:payload;type:longblob;not null" json:"payload"`
	Status        string     `gorm:"column:status;type:varchar(16);not null;default:'pending';index:idx_status_next,priority:1;index:idx_published,priority:1" json:"status"`
	RetryCount    int        `gorm:"column:retry_count;type:int;not null;default:0" json:"retry_count"`
	NextAttemptAt time.Time  `gorm:"column:next_attempt_at;type:datetime(3);not null;index:idx_status_next,priority:2" json:"next_attempt_at"`
	LockedAt      *time.Time `gorm:"column:locked_at;type:datetime(3)" json:"locked_at"`
	LockToken     string     `gorm:"column:lock_token;type:varchar(64);not null;default:''" json:"lock_token"`
	LastError     string     `gorm:"column:last_error;type:text" json:"last_error"`
	OccurredAt    time.Time  `gorm:"column:occurred_at;type:datetime(3);not null" json:"occurred_at"`
	CreatedAt     time.Time  `gorm:"column:created_at;type:datetime(3);not null" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;type:datetime(3);not null" json:"updated_at"`
	PublishedAt   *time.Time `gorm:"column:published_at;type:datetime(3);index:idx_published,priority:2" json:"published_at"`
}

// TableName 实现 gorm.Tabler。
func (Record) TableName() string { return TableName }

// Message 把表行还原为待发布消息。
func (r Record) Message() (mqx.Message, error) {
	if r.Destination == "" {
		return mqx.Message{}, fmt.Errorf("outbox: destination is required")
	}

	headers := make(map[string]string, 8)
	if r.Headers != "" {
		if err := json.Unmarshal([]byte(r.Headers), &headers); err != nil {
			return mqx.Message{}, fmt.Errorf("outbox: decode headers: %w", err)
		}
	}

	timestamp := r.OccurredAt
	if timestamp.IsZero() {
		timestamp = r.CreatedAt
	}
	return mqx.Message{
		ID:          r.EventID,
		Destination: r.Destination,
		RoutingKey:  r.RoutingKey,
		Key:         r.MessageKey,
		Headers:     headers,
		Body:        r.Payload,
		Timestamp:   timestamp,
	}, nil
}
