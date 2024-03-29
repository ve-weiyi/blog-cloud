package entity

import "time"

// TableNameTag return the table name of <tag>
const TableNameTag = "tag"

// Tag mapped from table <tag>
type Tag struct {
	ID        int       `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:id" json:"id"`                    // id
	TagName   string    `gorm:"column:tag_name;type:varchar(32);not null;comment:标签名" json:"tag_name"`                             // 标签名
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName Tag's table name
func (*Tag) TableName() string {
	return TableNameTag
}
