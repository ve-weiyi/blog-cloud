package entity

import "time"

// TableNameUniqueView return the table name of <unique_view>
const TableNameUniqueView = "unique_view"

// UniqueView mapped from table <unique_view>
type UniqueView struct {
	ID         int       `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:id" json:"id"`                    // id
	ViewsCount int       `gorm:"column:views_count;type:int;not null;comment:访问量" json:"views_count"`                               // 访问量
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName UniqueView's table name
func (*UniqueView) TableName() string {
	return TableNameUniqueView
}
