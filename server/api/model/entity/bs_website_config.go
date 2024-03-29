package entity

import "time"

// TableNameWebsiteConfig return the table name of <website_config>
const TableNameWebsiteConfig = "website_config"

// WebsiteConfig mapped from table <website_config>
type WebsiteConfig struct {
	ID        int       `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:id" json:"id"`                    // id
	Key       string    `gorm:"column:key;type:varchar(32);not null;uniqueIndex:key,priority:1;comment:关键词" json:"key"`            // 关键词
	Config    string    `gorm:"column:config;type:varchar(2048);not null;comment:配置信息" json:"config"`                              // 配置信息
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName WebsiteConfig's table name
func (*WebsiteConfig) TableName() string {
	return TableNameWebsiteConfig
}
