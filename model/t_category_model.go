package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TCategoryModel = (*defaultTCategoryModel)(nil)

type (
	// 接口定义
	TCategoryModel interface {
		Table[TCategory]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TCategoryModel)
		// add extra method in here
		FindOneByCategoryName(ctx context.Context, categoryName string) (out *TCategory, err error)
	}

	// 表字段定义
	TCategory struct {
		Id           int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:id" json:"id"`                                                  // id
		CategoryName string    `gorm:"column:category_name;type:varchar(32);not null;uniqueIndex:uk_name,priority:1;default:'';comment:分类名" json:"category_name"` // 分类名
		CreatedAt    time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                         // 创建时间
		UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                         // 更新时间
	}

	// 接口实现
	defaultTCategoryModel struct {
		baseModel[TCategory]
	}
)

func NewTCategoryModel(db *gorm.DB) TCategoryModel {
	return &defaultTCategoryModel{newBaseModel[TCategory](db, "`t_category`")}
}

// 在事务中操作
func (m *defaultTCategoryModel) WithTx(tx *gorm.DB) (out TCategoryModel) {
	return &defaultTCategoryModel{newBaseModel[TCategory](tx, m.table)}
}

// add extra method in here
func (m *defaultTCategoryModel) FindOneByCategoryName(ctx context.Context, categoryName string) (out *TCategory, err error) {
	db := m.scope(ctx)

	err = db.Where("`category_name` = ?", categoryName).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
