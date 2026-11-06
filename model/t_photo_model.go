package model

import (
	"time"

	"gorm.io/gorm"
)

var _ TPhotoModel = (*defaultTPhotoModel)(nil)

type (
	// 接口定义
	TPhotoModel interface {
		Table[TPhoto]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TPhotoModel)
	}

	// 表字段定义
	TPhoto struct {
		Id        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键" json:"id"`                               // 主键
		AlbumId   int64     `gorm:"column:album_id;type:bigint;not null;index:idx_album,priority:1;default:0;comment:相册id" json:"album_id"` // 相册id
		PhotoName string    `gorm:"column:photo_name;type:varchar(64);not null;default:'';comment:照片名" json:"photo_name"`                   // 照片名
		PhotoDesc string    `gorm:"column:photo_desc;type:varchar(128);not null;default:'';comment:照片描述" json:"photo_desc"`                 // 照片描述
		PhotoSrc  string    `gorm:"column:photo_src;type:varchar(255);not null;default:'';comment:照片地址" json:"photo_src"`                   // 照片地址
		IsDelete  int64     `gorm:"column:is_delete;type:tinyint(1);not null;default:0;comment:是否删除" json:"is_delete"`                      // 是否删除
		CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`      // 创建时间
		UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`      // 更新时间
	}

	// 接口实现
	defaultTPhotoModel struct {
		baseModel[TPhoto]
	}
)

func NewTPhotoModel(db *gorm.DB) TPhotoModel {
	return &defaultTPhotoModel{newBaseModel[TPhoto](db, "`t_photo`")}
}

// 在事务中操作
func (m *defaultTPhotoModel) WithTx(tx *gorm.DB) (out TPhotoModel) {
	return &defaultTPhotoModel{newBaseModel[TPhoto](tx, m.table)}
}
