package model

import (
	"time"

	"gorm.io/gorm"
)

var _ TNotificationEventModel = (*defaultTNotificationEventModel)(nil)

type (
	// 接口定义
	TNotificationEventModel interface {
		Table[TNotificationEvent]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TNotificationEventModel)
	}

	// 表字段定义
	TNotificationEvent struct {
		Id            int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:ID" json:"id"`                                                                                            // ID
		Scene         string    `gorm:"column:scene;type:varchar(48);not null;default:'';comment:场景标识" json:"scene"`                                                                                         // 场景标识
		ActorId       string    `gorm:"column:actor_id;type:varchar(64);not null;default:'';comment:触发者用户ID；系统事件为空" json:"actor_id"`                                                                         // 触发者用户ID；系统事件为空
		ObjectType    string    `gorm:"column:object_type;type:varchar(32);not null;index:idx_object,priority:1;default:'';comment:对象类型" json:"object_type"`                                                 // 对象类型
		ObjectId      int64     `gorm:"column:object_id;type:bigint;not null;index:idx_object,priority:2;default:0;comment:对象ID" json:"object_id"`                                                           // 对象ID
		AudienceType  string    `gorm:"column:audience_type;type:varchar(16);not null;default:'';comment:受众类型：all | user_ids；空值非法，扇出时直接失败" json:"audience_type"`                                             // 受众类型：all | user_ids；空值非法，扇出时直接失败
		AudienceIds   *string   `gorm:"column:audience_ids;type:json;comment:指定受众的用户ID列表" json:"audience_ids"`                                                                                               // 指定受众的用户ID列表
		Payload       *string   `gorm:"column:payload;type:json;comment:渲染载荷；公告事件只存内容引用" json:"payload"`                                                                                                     // 渲染载荷；公告事件只存内容引用
		Status        string    `gorm:"column:status;type:varchar(16);not null;index:idx_status_updated,priority:1;default:pending;comment:状态：pending | fanning | done | failed | terminated" json:"status"` // 状态：pending | fanning | done | failed | terminated
		AudienceTotal int64     `gorm:"column:audience_total;type:bigint;not null;default:0;comment:受众总数，扇出开始时确定，之后不再变化" json:"audience_total"`                                                              // 受众总数，扇出开始时确定，之后不再变化
		Processed     int64     `gorm:"column:processed;type:bigint;not null;default:0;comment:已处理收件人数，与受众总数构成进度" json:"processed"`                                                                          // 已处理收件人数，与受众总数构成进度
		Cursor        int64     `gorm:"column:cursor;type:bigint;not null;default:0;comment:扇出游标" json:"cursor"`                                                                                             // 扇出游标
		CreatedAt     time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                                                                   // 创建时间
		UpdatedAt     time.Time `gorm:"column:updated_at;type:datetime;not null;index:idx_status_updated,priority:2;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                               // 更新时间
	}

	// 接口实现
	defaultTNotificationEventModel struct {
		baseModel[TNotificationEvent]
	}
)

func NewTNotificationEventModel(db *gorm.DB) TNotificationEventModel {
	return &defaultTNotificationEventModel{newBaseModel[TNotificationEvent](db, "`t_notification_event`")}
}

// 在事务中操作
func (m *defaultTNotificationEventModel) WithTx(tx *gorm.DB) (out TNotificationEventModel) {
	return &defaultTNotificationEventModel{newBaseModel[TNotificationEvent](tx, m.table)}
}
