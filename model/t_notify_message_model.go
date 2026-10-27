package model

import (
	"time"

	"gorm.io/gorm"
)

var _ TNotifyMessageModel = (*defaultTNotifyMessageModel)(nil)

type (
	// 接口定义
	TNotifyMessageModel interface {
		Table[TNotifyMessage]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TNotifyMessageModel)
	}

	// 表字段定义
	TNotifyMessage struct {
		Id          int64      `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:ID" json:"id"`                                                                                 // ID
		Title       string     `gorm:"column:title;type:varchar(256);not null;default:'';comment:消息标题" json:"title"`                                                                             // 消息标题
		Content     *string    `gorm:"column:content;type:text;comment:消息内容" json:"content"`                                                                                                     // 消息内容
		Category    string     `gorm:"column:category;type:varchar(50);not null;index:idx_category,priority:1;default:system;comment:分类：system | maintenance | update | remind" json:"category"` // 分类：system | maintenance | update | remind
		Level       string     `gorm:"column:level;type:varchar(20);not null;default:info;comment:等级：info | warning | error" json:"level"`                                                       // 等级：info | warning | error
		TargetType  string     `gorm:"column:target_type;type:varchar(20);not null;default:all;comment:目标类型：all | user_ids" json:"target_type"`                                                  // 目标类型：all | user_ids
		TargetIds   *string    `gorm:"column:target_ids;type:json;comment:目标用户ID列表" json:"target_ids"`                                                                                           // 目标用户ID列表
		Status      string     `gorm:"column:status;type:varchar(20);not null;index:idx_status,priority:1;default:draft;comment:状态：draft | published | revoked" json:"status"`                   // 状态：draft | published | revoked
		PublishedAt *time.Time `gorm:"column:published_at;type:datetime;index:idx_published_at,priority:1;comment:发布时间" json:"published_at"`                                                     // 发布时间
		PublishedBy string     `gorm:"column:published_by;type:varchar(64);not null;default:'';comment:发布人" json:"published_by"`                                                                 // 发布人
		CreatedAt   time.Time  `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                                                        // 创建时间
		UpdatedAt   time.Time  `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                                                        // 更新时间
	}

	// 接口实现
	defaultTNotifyMessageModel struct {
		baseModel[TNotifyMessage]
	}
)

func NewTNotifyMessageModel(db *gorm.DB) TNotifyMessageModel {
	return &defaultTNotifyMessageModel{newBaseModel[TNotifyMessage](db, "`t_notify_message`")}
}

// 在事务中操作
func (m *defaultTNotifyMessageModel) WithTx(tx *gorm.DB) (out TNotifyMessageModel) {
	return &defaultTNotifyMessageModel{newBaseModel[TNotifyMessage](tx, m.table)}
}
