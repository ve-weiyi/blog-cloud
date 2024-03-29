package entity

import "time"

// TableNameRemark return the table name of <remark>
const TableNameRemark = "remark"

// Remark mapped from table <remark>
type Remark struct {
	ID             int       `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:主键id" json:"id"`                  // 主键id
	Nickname       string    `gorm:"column:nickname;type:varchar(64);not null;comment:昵称" json:"nickname"`                              // 昵称
	Avatar         string    `gorm:"column:avatar;type:varchar(255);not null;comment:头像" json:"avatar"`                                 // 头像
	MessageContent string    `gorm:"column:message_content;type:varchar(255);not null;comment:留言内容" json:"message_content"`             // 留言内容
	IpAddress      string    `gorm:"column:ip_address;type:varchar(64);not null;comment:用户ip" json:"ip_address"`                        // 用户ip
	IpSource       string    `gorm:"column:ip_source;type:varchar(255);not null;comment:用户地址" json:"ip_source"`                         // 用户地址
	Time           int       `gorm:"column:time;type:int;not null;comment:弹幕速度" json:"time"`                                            // 弹幕速度
	IsReview       int       `gorm:"column:is_review;type:tinyint(1);not null;default:1;comment:是否审核" json:"is_review"`                 // 是否审核
	CreatedAt      time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:发布时间" json:"created_at"` // 发布时间
	UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName Remark's table name
func (*Remark) TableName() string {
	return TableNameRemark
}
