package model

import (
	"time"

	"gorm.io/gorm"
)

var _ TLoginLogModel = (*defaultTLoginLogModel)(nil)

type (
	// 接口定义
	TLoginLogModel interface {
		Table[TLoginLog]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TLoginLogModel)
	}

	// 表字段定义
	TLoginLog struct {
		Id         int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:日志ID" json:"id"`                                                        // 日志ID
		UserId     string    `gorm:"column:user_id;type:varchar(64);not null;index:idx_user_id,priority:1;default:'';comment:用户id" json:"user_id"`                      // 用户id
		DeviceId   string    `gorm:"column:device_id;type:varchar(64);not null;default:'';comment:设备id" json:"device_id"`                                               // 设备id
		EventType  int64     `gorm:"column:event_type;type:tinyint;not null;default:1;comment:事件类型 1登录 2登出" json:"event_type"`                                          // 事件类型 1登录 2登出
		LoginType  string    `gorm:"column:login_type;type:varchar(20);not null;default:'';comment:登录类型：password-密码登录 sms-短信登录（登出事件为空）" json:"login_type"`              // 登录类型：password-密码登录 sms-短信登录（登出事件为空）
		LogoutType string    `gorm:"column:logout_type;type:varchar(32);not null;default:'';comment:登出类型：手动登出/超时/强制下线（登录事件为空）" json:"logout_type"`                      // 登出类型：手动登出/超时/强制下线（登录事件为空）
		Status     int64     `gorm:"column:status;type:tinyint;not null;default:1;comment:登录状态：0-失败 1-成功（登出事件恒为 1）" json:"status"`                                      // 登录状态：0-失败 1-成功（登出事件恒为 1）
		FailReason string    `gorm:"column:fail_reason;type:varchar(500);not null;default:'';comment:失败原因" json:"fail_reason"`                                          // 失败原因
		CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null;index:idx_created_at,priority:1;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
		UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                                 // 更新时间
	}

	// 接口实现
	defaultTLoginLogModel struct {
		baseModel[TLoginLog]
	}
)

func NewTLoginLogModel(db *gorm.DB) TLoginLogModel {
	return &defaultTLoginLogModel{newBaseModel[TLoginLog](db, "`t_login_log`")}
}

// 在事务中操作
func (m *defaultTLoginLogModel) WithTx(tx *gorm.DB) (out TLoginLogModel) {
	return &defaultTLoginLogModel{newBaseModel[TLoginLog](tx, m.table)}
}
