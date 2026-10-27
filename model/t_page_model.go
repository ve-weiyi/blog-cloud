package model

import (
	"time"

	"gorm.io/gorm"
)

var _ TPageModel = (*defaultTPageModel)(nil)

type (
	// 接口定义
	TPageModel interface {
		Table[TPage]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TPageModel)
	}

	// 表字段定义
	TPage struct {
		Id             int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:页面id" json:"id"`                        // 页面id
		PageName       string    `gorm:"column:page_name;type:varchar(32);not null;default:'';comment:页面名" json:"page_name"`                // 页面名
		PageLabel      string    `gorm:"column:page_label;type:varchar(32);not null;default:'';comment:页面标签" json:"page_label"`             // 页面标签
		PageCover      string    `gorm:"column:page_cover;type:varchar(255);not null;default:'';comment:页面封面" json:"page_cover"`            // 页面封面
		IsCarousel     int64     `gorm:"column:is_carousel;type:tinyint;not null;default:0;comment:是否轮播" json:"is_carousel"`                // 是否轮播
		CarouselCovers *string   `gorm:"column:carousel_covers;type:json;comment:轮播图片列表" json:"carousel_covers"`                            // 轮播图片列表
		CreatedAt      time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
		UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
	}

	// 接口实现
	defaultTPageModel struct {
		baseModel[TPage]
	}
)

func NewTPageModel(db *gorm.DB) TPageModel {
	return &defaultTPageModel{newBaseModel[TPage](db, "`t_page`")}
}

// 在事务中操作
func (m *defaultTPageModel) WithTx(tx *gorm.DB) (out TPageModel) {
	return &defaultTPageModel{newBaseModel[TPage](tx, m.table)}
}
