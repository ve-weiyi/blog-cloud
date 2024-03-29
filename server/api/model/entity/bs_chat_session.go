package entity

import "time"

// TableNameChatSession return the table name of <chat_session>
const TableNameChatSession = "chat_session"

// ChatSession mapped from table <chat_session>
type ChatSession struct {
	ID        int       `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:主键" json:"id"`                    // 主键
	ChatID    string    `gorm:"column:chat_id;type:int;not null;comment:群聊id" json:"chat_id"`                                      // 群聊id
	ChatTitle string    `gorm:"column:chat_title;type:varchar(128);not null;comment:标题" json:"chat_title"`                         // 标题
	Type      string    `gorm:"column:type;type:varchar(128);not null;comment:类型" json:"type"`                                     // 类型
	Status    int       `gorm:"column:status;type:int;not null;comment:0正常 1删除" json:"status"`                                     // 0正常 1删除
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName ChatSession's table name
func (*ChatSession) TableName() string {
	return TableNameChatSession
}
