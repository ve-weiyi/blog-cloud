package entity

import "time"

// TableNameChatRecord return the table name of <chat_record>
const TableNameChatRecord = "chat_record"

// ChatRecord mapped from table <chat_record>
type ChatRecord struct {
	ID        int       `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:主键" json:"id"`                    // 主键
	UserID    int       `gorm:"column:user_id;type:int;not null;comment:用户id" json:"user_id"`                                      // 用户id
	Nickname  string    `gorm:"column:nickname;type:varchar(64);not null;comment:昵称" json:"nickname"`                              // 昵称
	Avatar    string    `gorm:"column:avatar;type:varchar(255);not null;comment:头像" json:"avatar"`                                 // 头像
	Content   string    `gorm:"column:content;type:varchar(1000);not null;comment:聊天内容" json:"content"`                            // 聊天内容
	IpAddress string    `gorm:"column:ip_address;type:varchar(64);not null;comment:ip地址" json:"ip_address"`                        // ip地址
	IpSource  string    `gorm:"column:ip_source;type:varchar(255);not null;comment:ip来源" json:"ip_source"`                         // ip来源
	Type      int       `gorm:"column:type;type:int;not null;comment:类型" json:"type"`                                              // 类型
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName ChatRecord's table name
func (*ChatRecord) TableName() string {
	return TableNameChatRecord
}
