package entity

import "time"

// TableNameFriendLink return the table name of <friend_link>
const TableNameFriendLink = "friend_link"

// FriendLink mapped from table <friend_link>
type FriendLink struct {
	ID          int       `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:id" json:"id"`                               // id
	LinkName    string    `gorm:"column:link_name;type:varchar(32);not null;index:fk_friend_link_user,priority:1;comment:链接名" json:"link_name"` // 链接名
	LinkAvatar  string    `gorm:"column:link_avatar;type:varchar(255);not null;comment:链接头像" json:"link_avatar"`                                // 链接头像
	LinkAddress string    `gorm:"column:link_address;type:varchar(64);not null;comment:链接地址" json:"link_address"`                               // 链接地址
	LinkIntro   string    `gorm:"column:link_intro;type:varchar(100);not null;comment:链接介绍" json:"link_intro"`                                  // 链接介绍
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`            // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`            // 更新时间
}

// TableName FriendLink's table name
func (*FriendLink) TableName() string {
	return TableNameFriendLink
}
