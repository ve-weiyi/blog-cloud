package entity

import "time"

// TableNameUserOauth return the table name of <user_oauth>
const TableNameUserOauth = "user_oauth"

// UserOauth mapped from table <user_oauth>
type UserOauth struct {
	ID        int       `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:id" json:"id"`                                       // id
	UserID    int       `gorm:"column:user_id;type:int;not null;index:idx_uuid,priority:1;comment:用户id" json:"user_id"`                               // 用户id
	OpenID    string    `gorm:"column:open_id;type:varchar(128);not null;uniqueIndex:uk_oid_plat,priority:1;comment:开发平台id，标识唯一用户" json:"open_id"`    // 开发平台id，标识唯一用户
	Platform  string    `gorm:"column:platform;type:varchar(64);not null;uniqueIndex:uk_oid_plat,priority:2;comment:平台:手机号、邮箱、微信、飞书" json:"platform"` // 平台:手机号、邮箱、微信、飞书
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                    // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                    // 更新时间
}

// TableName UserOauth's table name
func (*UserOauth) TableName() string {
	return TableNameUserOauth
}
