package entity

import "time"

// TableNameUserAccount return the table name of <user_account>
const TableNameUserAccount = "user_account"

// UserAccount mapped from table <user_account>
type UserAccount struct {
	ID           int       `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:id" json:"id"`                           // id
	Username     string    `gorm:"column:username;type:varchar(64);not null;uniqueIndex:uk_username,priority:1;comment:用户名" json:"username"` // 用户名
	Password     string    `gorm:"column:password;type:varchar(128);not null;comment:密码" json:"password"`                                    // 密码
	Status       int       `gorm:"column:status;type:tinyint;not null;comment:状态: -1删除 0正常 1禁用" json:"status"`                               // 状态: -1删除 0正常 1禁用
	RegisterType string    `gorm:"column:register_type;type:varchar(64);not null;comment:注册方式" json:"register_type"`                         // 注册方式
	IpAddress    string    `gorm:"column:ip_address;type:varchar(255);not null;comment:注册ip" json:"ip_address"`                              // 注册ip
	IpSource     string    `gorm:"column:ip_source;type:varchar(255);not null;comment:注册ip 源" json:"ip_source"`                              // 注册ip 源
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`        // 创建时间
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`        // 更新时间
}

// TableName UserAccount's table name
func (*UserAccount) TableName() string {
	return TableNameUserAccount
}
