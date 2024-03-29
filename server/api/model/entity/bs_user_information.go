package entity

import "time"

// TableNameUserInformation return the table name of <user_information>
const TableNameUserInformation = "user_information"

// UserInformation mapped from table <user_information>
type UserInformation struct {
	ID        int       `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:id" json:"id"`                    // id
	UserID    int       `gorm:"column:user_id;type:int;not null;uniqueIndex:uk_uuid,priority:1;comment:用户id" json:"user_id"`       // 用户id
	Email     string    `gorm:"column:email;type:varchar(128);not null;comment:用户邮箱" json:"email"`                                 // 用户邮箱
	Nickname  string    `gorm:"column:nickname;type:varchar(128);not null;comment:用户昵称" json:"nickname"`                           // 用户昵称
	Avatar    string    `gorm:"column:avatar;type:varchar(1024);not null;comment:用户头像" json:"avatar"`                              // 用户头像
	Phone     string    `gorm:"column:phone;type:varchar(32);not null;comment:用户手机号" json:"phone"`                                 // 用户手机号
	Intro     string    `gorm:"column:intro;type:varchar(255);not null;comment:个人简介" json:"intro"`                                 // 个人简介
	Website   string    `gorm:"column:website;type:varchar(255);not null;comment:个人网站" json:"website"`                             // 个人网站
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName UserInformation's table name
func (*UserInformation) TableName() string {
	return TableNameUserInformation
}
