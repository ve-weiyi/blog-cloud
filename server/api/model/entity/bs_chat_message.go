package entity

import "time"

// TableNameChatMessage return the table name of <chat_message>
const TableNameChatMessage = "chat_message"

// ChatMessage mapped from table <chat_message>
type ChatMessage struct {
	ID         int       `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:主键" json:"id"`                    // 主键
	ChatID     string    `gorm:"column:chat_id;type:varchar(128);not null;comment:群聊id" json:"chat_id"`                             // 群聊id
	UserID     int       `gorm:"column:user_id;type:int;not null;comment:用户id" json:"user_id"`                                      // 用户id
	ReplyMsgID int       `gorm:"column:reply_msg_id;type:int;not null;comment:回复消息id" json:"reply_msg_id"`                          // 回复消息id
	Content    string    `gorm:"column:content;type:varchar(1024);not null;comment:聊天内容" json:"content"`                            // 聊天内容
	IpAddress  string    `gorm:"column:ip_address;type:varchar(64);not null;comment:ip地址" json:"ip_address"`                        // ip地址
	IpSource   string    `gorm:"column:ip_source;type:varchar(255);not null;comment:ip来源" json:"ip_source"`                         // ip来源
	Type       int       `gorm:"column:type;type:int;not null;comment:类型" json:"type"`                                              // 类型
	Status     int       `gorm:"column:status;type:int;not null;comment:0正常 1撤回 2已编辑" json:"status"`                                // 0正常 1撤回 2已编辑
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName ChatMessage's table name
func (*ChatMessage) TableName() string {
	return TableNameChatMessage
}
