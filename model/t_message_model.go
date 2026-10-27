package model

import (
	"time"

	"gorm.io/gorm"
)

var _ TMessageModel = (*defaultTMessageModel)(nil)

type (
	// 接口定义
	TMessageModel interface {
		Table[TMessage]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TMessageModel)
	}

	// 表字段定义
	TMessage struct {
		Id             int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键id" json:"id"`                        // 主键id
		UserId         string    `gorm:"column:user_id;type:varchar(64);not null;default:'';comment:用户id" json:"user_id"`                   // 用户id
		DeviceId       string    `gorm:"column:device_id;type:varchar(64);not null;default:'';comment:终端id" json:"device_id"`               // 终端id
		MessageContent string    `gorm:"column:message_content;type:varchar(255);not null;default:'';comment:留言内容" json:"message_content"`  // 留言内容
		Status         int64     `gorm:"column:status;type:tinyint;not null;default:0;comment:状态:0正常 1编辑 2撤回 3删除" json:"status"`            // 状态:0正常 1编辑 2撤回 3删除
		CreatedAt      time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:发布时间" json:"created_at"` // 发布时间
		UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
	}

	// 接口实现
	defaultTMessageModel struct {
		baseModel[TMessage]
	}
)

func NewTMessageModel(db *gorm.DB) TMessageModel {
	return &defaultTMessageModel{newBaseModel[TMessage](db, "`t_message`")}
}

// 在事务中操作
func (m *defaultTMessageModel) WithTx(tx *gorm.DB) (out TMessageModel) {
	return &defaultTMessageModel{newBaseModel[TMessage](tx, m.table)}
}
