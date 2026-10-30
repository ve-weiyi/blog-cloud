package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TNotificationModel = (*defaultTNotificationModel)(nil)

type (
	// 接口定义
	TNotificationModel interface {
		Table[TNotification]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TNotificationModel)
		// add extra method in here
		FindOneByRecipientIdDedupeKey(ctx context.Context, recipientId string, dedupeKey string) (out *TNotification, err error)
	}

	// 表字段定义
	TNotification struct {
		Id          int64      `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:ID" json:"id"`                                                                                                                                               // ID
		EventId     int64      `gorm:"column:event_id;type:bigint;not null;index:idx_event,priority:1;default:0;comment:通知事件ID" json:"event_id"`                                                                                                               // 通知事件ID
		RecipientId string     `gorm:"column:recipient_id;type:varchar(64);not null;uniqueIndex:uk_recipient_dedupe,priority:1;index:idx_recipient_scene_group,priority:1;index:idx_recipient_read,priority:1;default:'';comment:收件人用户ID" json:"recipient_id"` // 收件人用户ID
		Scene       string     `gorm:"column:scene;type:varchar(48);not null;index:idx_recipient_scene_group,priority:2;index:idx_recipient_read,priority:3;index:idx_scene_read_created,priority:1;default:'';comment:场景标识" json:"scene"`                     // 场景标识
		ObjectType  string     `gorm:"column:object_type;type:varchar(32);not null;index:idx_object,priority:1;default:'';comment:对象类型：article | talk | comment | user | announcement" json:"object_type"`                                                     // 对象类型：article | talk | comment | user | announcement
		ObjectId    int64      `gorm:"column:object_id;type:bigint;not null;index:idx_object,priority:2;default:0;comment:对象ID" json:"object_id"`                                                                                                              // 对象ID
		GroupKey    string     `gorm:"column:group_key;type:varchar(191);not null;index:idx_recipient_scene_group,priority:3;index:idx_recipient_read,priority:4;default:'';comment:展示分组键；写入时物化，不是计数" json:"group_key"`                                        // 展示分组键；写入时物化，不是计数
		DedupeKey   string     `gorm:"column:dedupe_key;type:varchar(191);not null;uniqueIndex:uk_recipient_dedupe,priority:2;default:'';comment:去重键：可分组场景含触发者，不可分组场景含事件" json:"dedupe_key"`                                                                   // 去重键：可分组场景含触发者，不可分组场景含事件
		ReadAt      *time.Time `gorm:"column:read_at;type:datetime;index:idx_recipient_read,priority:2;index:idx_scene_read_created,priority:2;comment:已读时间" json:"read_at"`                                                                                   // 已读时间
		CreatedAt   time.Time  `gorm:"column:created_at;type:datetime;not null;index:idx_recipient_scene_group,priority:4;index:idx_scene_read_created,priority:3;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                                   // 创建时间
	}

	// 接口实现
	defaultTNotificationModel struct {
		baseModel[TNotification]
	}
)

func NewTNotificationModel(db *gorm.DB) TNotificationModel {
	return &defaultTNotificationModel{newBaseModel[TNotification](db, "`t_notification`")}
}

// 在事务中操作
func (m *defaultTNotificationModel) WithTx(tx *gorm.DB) (out TNotificationModel) {
	return &defaultTNotificationModel{newBaseModel[TNotification](tx, m.table)}
}

// add extra method in here
func (m *defaultTNotificationModel) FindOneByRecipientIdDedupeKey(ctx context.Context, recipientId string, dedupeKey string) (out *TNotification, err error) {
	db := m.scope(ctx)

	err = db.Where("`recipient_id` = ? and `dedupe_key` = ?", recipientId, dedupeKey).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
