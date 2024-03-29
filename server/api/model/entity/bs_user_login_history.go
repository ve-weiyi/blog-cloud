package entity

import "time"

// TableNameUserLoginHistory return the table name of <user_login_history>
const TableNameUserLoginHistory = "user_login_history"

// UserLoginHistory mapped from table <user_login_history>
type UserLoginHistory struct {
	ID        int       `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:id" json:"id"`                    // id
	UserID    int       `gorm:"column:user_id;type:int;not null;index:uk_uuid,priority:1;comment:用户id" json:"user_id"`             // 用户id
	LoginType string    `gorm:"column:login_type;type:varchar(64);not null;comment:登录类型" json:"login_type"`                        // 登录类型
	Agent     string    `gorm:"column:agent;type:varchar(255);not null;comment:代理" json:"agent"`                                   // 代理
	IpAddress string    `gorm:"column:ip_address;type:varchar(255);not null;comment:ip host" json:"ip_address"`                    // ip host
	IpSource  string    `gorm:"column:ip_source;type:varchar(255);not null;comment:ip 源" json:"ip_source"`                         // ip 源
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName UserLoginHistory's table name
func (*UserLoginHistory) TableName() string {
	return TableNameUserLoginHistory
}
