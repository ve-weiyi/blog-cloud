package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TTagModel = (*defaultTTagModel)(nil)

type (
	// 接口定义
	TTagModel interface {
		Table[TTag]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TTagModel)
		// add extra method in here
		FindOneByTagName(ctx context.Context, tagName string) (out *TTag, err error)
	}

	// 表字段定义
	TTag struct {
		Id        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:id" json:"id"`                                        // id
		TagName   string    `gorm:"column:tag_name;type:varchar(32);not null;uniqueIndex:uk_name,priority:1;default:'';comment:标签名" json:"tag_name"` // 标签名
		CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`               // 创建时间
		UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`               // 更新时间
	}

	// 接口实现
	defaultTTagModel struct {
		baseModel[TTag]
	}
)

func NewTTagModel(db *gorm.DB) TTagModel {
	return &defaultTTagModel{newBaseModel[TTag](db, "`t_tag`")}
}

// 在事务中操作
func (m *defaultTTagModel) WithTx(tx *gorm.DB) (out TTagModel) {
	return &defaultTTagModel{newBaseModel[TTag](tx, m.table)}
}

// add extra method in here
func (m *defaultTTagModel) FindOneByTagName(ctx context.Context, tagName string) (out *TTag, err error) {
	db := m.scope(ctx)

	err = db.Where("`tag_name` = ?", tagName).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
