package model

import (
	"time"
)

// json 推荐下划线
type MODEL struct {
	//gorm.Model
	ID        uint      `gorm:"column:id;primaryKey;autoIncrement:true;type:int;" json:"ID"` // 主键ID
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"CreatedAt"`  // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime" json:"UpdatedAt"`  // 更新时间
	//DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`                            // 删除时间
}
