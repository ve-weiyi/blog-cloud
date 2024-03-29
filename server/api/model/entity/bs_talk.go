package entity

import "time"

// TableNameTalk return the table name of <talk>
const TableNameTalk = "talk"

// Talk mapped from table <talk>
type Talk struct {
	ID        int       `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:说说id" json:"id"`                  // 说说id
	UserID    int       `gorm:"column:user_id;type:int;not null;comment:用户id" json:"user_id"`                                      // 用户id
	Content   string    `gorm:"column:content;type:varchar(2048);not null;comment:说说内容" json:"content"`                            // 说说内容
	Images    string    `gorm:"column:images;type:varchar(2048);not null;comment:图片" json:"images"`                                // 图片
	IsTop     int       `gorm:"column:is_top;type:tinyint;not null;comment:是否置顶" json:"is_top"`                                    // 是否置顶
	Status    int       `gorm:"column:status;type:tinyint;not null;default:1;comment:状态 1.公开 2.私密" json:"status"`                  // 状态 1.公开 2.私密
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName Talk's table name
func (*Talk) TableName() string {
	return TableNameTalk
}
