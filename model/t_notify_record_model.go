package model

import (
	"time"

	"gorm.io/gorm"
)

var _ TNotifyRecordModel = (*defaultTNotifyRecordModel)(nil)

type (
	// 接口定义
	TNotifyRecordModel interface {
		Table[TNotifyRecord]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TNotifyRecordModel)
	}

	// 表字段定义
	TNotifyRecord struct {
		Id           int64      `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:ID" json:"id"`                                                                                        // ID
		MessageId    int64      `gorm:"column:message_id;type:bigint;not null;index:idx_message_id,priority:1;default:0;comment:关联消息ID" json:"message_id"`                                               // 关联消息ID
		Channel      string     `gorm:"column:channel;type:varchar(16);not null;index:idx_channel,priority:1;default:inbox;comment:渠道：inbox | sms | email" json:"channel"`                               // 渠道：inbox | sms | email
		Recipient    string     `gorm:"column:recipient;type:varchar(128);not null;index:idx_recipient,priority:1;default:'';comment:接收者（user_id | mobile | email）" json:"recipient"`                    // 接收者（user_id | mobile | email）
		TemplateCode string     `gorm:"column:template_code;type:varchar(64);not null;default:'';comment:模板编码" json:"template_code"`                                                                     // 模板编码
		Content      *string    `gorm:"column:content;type:text;comment:实际发送内容" json:"content"`                                                                                                          // 实际发送内容
		Status       string     `gorm:"column:status;type:varchar(16);not null;index:idx_status,priority:1;default:unread;comment:状态：inbox: unread|read / sms/email: pending|sent|failed" json:"status"` // 状态：inbox: unread|read / sms/email: pending|sent|failed
		BizId        string     `gorm:"column:biz_id;type:varchar(64);not null;index:idx_biz_id,priority:1;default:'';comment:业务幂等键（email/sms 为事件ID，inbox 为消息ID:接收者）" json:"biz_id"`                     // 业务幂等键（email/sms 为事件ID，inbox 为消息ID:接收者）
		ErrorMsg     string     `gorm:"column:error_msg;type:varchar(512);not null;default:'';comment:失败原因" json:"error_msg"`                                                                            // 失败原因
		ReadAt       *time.Time `gorm:"column:read_at;type:datetime;comment:阅读时间（inbox）" json:"read_at"`                                                                                                 // 阅读时间（inbox）
		SentAt       *time.Time `gorm:"column:sent_at;type:datetime;comment:发送时间（sms/email）" json:"sent_at"`                                                                                             // 发送时间（sms/email）
		CreatedAt    time.Time  `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                                                               // 创建时间
	}

	// 接口实现
	defaultTNotifyRecordModel struct {
		baseModel[TNotifyRecord]
	}
)

func NewTNotifyRecordModel(db *gorm.DB) TNotifyRecordModel {
	return &defaultTNotifyRecordModel{newBaseModel[TNotifyRecord](db, "`t_notify_record`")}
}

// 在事务中操作
func (m *defaultTNotifyRecordModel) WithTx(tx *gorm.DB) (out TNotifyRecordModel) {
	return &defaultTNotifyRecordModel{newBaseModel[TNotifyRecord](tx, m.table)}
}
