package model

import (
	"time"

	"gorm.io/gorm"
)

var _ TAiSessionModel = (*defaultTAiSessionModel)(nil)

type (
	// 接口定义
	TAiSessionModel interface {
		Table[TAiSession]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TAiSessionModel)
	}

	// 表字段定义
	TAiSession struct {
		Id            int64      `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键" json:"id"`                                  // 主键
		UserId        string     `gorm:"column:user_id;type:varchar(64);not null;index:idx_user,priority:1;default:'';comment:所属用户" json:"user_id"` // 所属用户
		Title         string     `gorm:"column:title;type:varchar(128);not null;default:'';comment:标题" json:"title"`                                // 标题
		Model         string     `gorm:"column:model;type:varchar(64);not null;default:'';comment:模型" json:"model"`                                 // 模型
		MessageCount  int64      `gorm:"column:message_count;type:int;not null;default:0;comment:消息条数" json:"message_count"`                        // 消息条数
		LastMessageAt *time.Time `gorm:"column:last_message_at;type:datetime;index:idx_user,priority:2;comment:最近消息时间" json:"last_message_at"`      // 最近消息时间
		CreatedAt     time.Time  `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`         // 创建时间
		UpdatedAt     time.Time  `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`         // 更新时间
	}

	// 接口实现
	defaultTAiSessionModel struct {
		baseModel[TAiSession]
	}
)

func NewTAiSessionModel(db *gorm.DB) TAiSessionModel {
	return &defaultTAiSessionModel{newBaseModel[TAiSession](db, "`t_ai_session`")}
}

// 在事务中操作
func (m *defaultTAiSessionModel) WithTx(tx *gorm.DB) (out TAiSessionModel) {
	return &defaultTAiSessionModel{newBaseModel[TAiSession](tx, m.table)}
}
