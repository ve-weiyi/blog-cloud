package entity

import "time"

// TableNamePage return the table name of <page>
const TableNamePage = "page"

// Page mapped from table <page>
type Page struct {
	ID        int       `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:页面id" json:"id"`                  // 页面id
	PageName  string    `gorm:"column:page_name;type:varchar(32);not null;comment:页面名" json:"page_name"`                           // 页面名
	PageLabel string    `gorm:"column:page_label;type:varchar(32);not null;comment:页面标签" json:"page_label"`                        // 页面标签
	PageCover string    `gorm:"column:page_cover;type:varchar(255);not null;comment:页面封面" json:"page_cover"`                       // 页面封面
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName Page's table name
func (*Page) TableName() string {
	return TableNamePage
}
