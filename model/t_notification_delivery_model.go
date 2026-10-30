package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TNotificationDeliveryModel = (*defaultTNotificationDeliveryModel)(nil)

type (
	// 接口定义
	TNotificationDeliveryModel interface {
		Table[TNotificationDelivery]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TNotificationDeliveryModel)
		// add extra method in here
		FindOneByNotificationIdChannel(ctx context.Context, notificationId int64, channel string) (out *TNotificationDelivery, err error)
	}

	// 表字段定义
	TNotificationDelivery struct {
		Id             int64      `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:ID" json:"id"`                                                                                                                   // ID
		NotificationId int64      `gorm:"column:notification_id;type:bigint;not null;uniqueIndex:uk_notification_channel,priority:1;default:0;comment:收件箱条目ID" json:"notification_id"`                                                // 收件箱条目ID
		Channel        string     `gorm:"column:channel;type:varchar(16);not null;uniqueIndex:uk_notification_channel,priority:2;default:'';comment:渠道：email | sms | push" json:"channel"`                                            // 渠道：email | sms | push
		Recipient      string     `gorm:"column:recipient;type:varchar(255);not null;default:'';comment:收件地址：邮箱 / 手机号 / 设备端点" json:"recipient"`                                                                                       // 收件地址：邮箱 / 手机号 / 设备端点
		Status         string     `gorm:"column:status;type:varchar(16);not null;index:idx_status_retry,priority:1;index:idx_status_updated,priority:1;default:pending;comment:状态：pending | sent | failed | cancelled" json:"status"` // 状态：pending | sent | failed | cancelled
		TemplateCode   string     `gorm:"column:template_code;type:varchar(64);not null;default:'';comment:模板编码" json:"template_code"`                                                                                                // 模板编码
		Content        *string    `gorm:"column:content;type:text;comment:实际发送内容" json:"content"`                                                                                                                                     // 实际发送内容
		ProviderMsgId  string     `gorm:"column:provider_msg_id;type:varchar(128);not null;default:'';comment:渠道侧消息标识" json:"provider_msg_id"`                                                                                        // 渠道侧消息标识
		ErrorMsg       string     `gorm:"column:error_msg;type:varchar(512);not null;default:'';comment:失败原因" json:"error_msg"`                                                                                                       // 失败原因
		Attempts       int64      `gorm:"column:attempts;type:int;not null;default:0;comment:尝试次数" json:"attempts"`                                                                                                                   // 尝试次数
		NextRetryAt    *time.Time `gorm:"column:next_retry_at;type:datetime;index:idx_status_retry,priority:2;comment:下次重试时间" json:"next_retry_at"`                                                                                   // 下次重试时间
		SentAt         *time.Time `gorm:"column:sent_at;type:datetime;comment:发送时间" json:"sent_at"`                                                                                                                                   // 发送时间
		CreatedAt      time.Time  `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                                                                                          // 创建时间
		UpdatedAt      time.Time  `gorm:"column:updated_at;type:datetime;not null;index:idx_status_updated,priority:2;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                                                      // 更新时间
	}

	// 接口实现
	defaultTNotificationDeliveryModel struct {
		baseModel[TNotificationDelivery]
	}
)

func NewTNotificationDeliveryModel(db *gorm.DB) TNotificationDeliveryModel {
	return &defaultTNotificationDeliveryModel{newBaseModel[TNotificationDelivery](db, "`t_notification_delivery`")}
}

// 在事务中操作
func (m *defaultTNotificationDeliveryModel) WithTx(tx *gorm.DB) (out TNotificationDeliveryModel) {
	return &defaultTNotificationDeliveryModel{newBaseModel[TNotificationDelivery](tx, m.table)}
}

// add extra method in here
func (m *defaultTNotificationDeliveryModel) FindOneByNotificationIdChannel(ctx context.Context, notificationId int64, channel string) (out *TNotificationDelivery, err error) {
	db := m.scope(ctx)

	err = db.Where("`notification_id` = ? and `channel` = ?", notificationId, channel).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
