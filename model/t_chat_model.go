package model

import (
	"time"

	"gorm.io/gorm"
)

var _ TChatModel = (*defaultTChatModel)(nil)

type (
	// 接口定义
	TChatModel interface {
		Table[TChat]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TChatModel)
	}

	// 表字段定义
	TChat struct {
		Id        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键" json:"id"`                                 // 主键
		UserId    string    `gorm:"column:user_id;type:varchar(64);not null;index:idx_uid,priority:1;default:'';comment:用户id" json:"user_id"` // 用户id
		DeviceId  string    `gorm:"column:device_id;type:varchar(64);not null;default:'';comment:设备id" json:"device_id"`                      // 设备id
		Nickname  string    `gorm:"column:nickname;type:varchar(64);not null;default:'';comment:昵称" json:"nickname"`                          // 昵称
		Avatar    string    `gorm:"column:avatar;type:varchar(255);not null;default:'';comment:头像" json:"avatar"`                             // 头像
		IpAddress string    `gorm:"column:ip_address;type:varchar(64);not null;default:'';comment:用户ip 127.0.0.1" json:"ip_address"`          // 用户ip 127.0.0.1
		IpSource  string    `gorm:"column:ip_source;type:varchar(128);not null;default:'';comment:用户地址 广东省深圳市" json:"ip_source"`              // 用户地址 广东省深圳市
		Type      string    `gorm:"column:type;type:varchar(64);not null;default:'';comment:类型:chatgpt chatroom" json:"type"`                 // 类型:chatgpt chatroom
		Content   string    `gorm:"column:content;type:varchar(4096);not null;default:'';comment:聊天内容" json:"content"`                        // 聊天内容
		Status    int64     `gorm:"column:status;type:tinyint;not null;default:0;comment:状态:0正常 1编辑 2撤回 3删除" json:"status"`                   // 状态:0正常 1编辑 2撤回 3删除
		CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`        // 创建时间
		UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`        // 更新时间
	}

	// 接口实现
	defaultTChatModel struct {
		baseModel[TChat]
	}
)

func NewTChatModel(db *gorm.DB) TChatModel {
	return &defaultTChatModel{newBaseModel[TChat](db, "`t_chat`")}
}

// 在事务中操作
func (m *defaultTChatModel) WithTx(tx *gorm.DB) (out TChatModel) {
	return &defaultTChatModel{newBaseModel[TChat](tx, m.table)}
}
