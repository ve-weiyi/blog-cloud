package model

import (
	"time"

	"gorm.io/gorm"
)

var _ TTalkModel = (*defaultTTalkModel)(nil)

type (
	// 接口定义
	TTalkModel interface {
		Table[TTalk]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TTalkModel)
	}

	// 表字段定义
	TTalk struct {
		Id        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:说说id" json:"id"`                        // 说说id
		UserId    string    `gorm:"column:user_id;type:varchar(64);not null;default:'';comment:用户id" json:"user_id"`                   // 用户id
		Content   string    `gorm:"column:content;type:varchar(2048);not null;default:'';comment:说说内容" json:"content"`                 // 说说内容
		Images    *string   `gorm:"column:images;type:json;comment:图片" json:"images"`                                                  // 图片
		IsTop     int64     `gorm:"column:is_top;type:tinyint;not null;default:0;comment:是否置顶" json:"is_top"`                          // 是否置顶
		Status    int64     `gorm:"column:status;type:tinyint;not null;default:1;comment:状态 1.公开 2.私密" json:"status"`                  // 状态 1.公开 2.私密
		LikeCount int64     `gorm:"column:like_count;type:bigint;not null;default:0;comment:点赞数" json:"like_count"`                    // 点赞数
		CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
		UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
	}

	// 接口实现
	defaultTTalkModel struct {
		baseModel[TTalk]
	}
)

func NewTTalkModel(db *gorm.DB) TTalkModel {
	return &defaultTTalkModel{newBaseModel[TTalk](db, "`t_talk`")}
}

// 在事务中操作
func (m *defaultTTalkModel) WithTx(tx *gorm.DB) (out TTalkModel) {
	return &defaultTTalkModel{newBaseModel[TTalk](tx, m.table)}
}
