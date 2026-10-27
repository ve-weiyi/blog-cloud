package model

import (
	"time"

	"gorm.io/gorm"
)

var _ TAlbumModel = (*defaultTAlbumModel)(nil)

type (
	// 接口定义
	TAlbumModel interface {
		Table[TAlbum]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TAlbumModel)
	}

	// 表字段定义
	TAlbum struct {
		Id         int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键" json:"id"`                          // 主键
		AlbumName  string    `gorm:"column:album_name;type:varchar(64);not null;default:'';comment:相册名" json:"album_name"`              // 相册名
		AlbumDesc  string    `gorm:"column:album_desc;type:varchar(128);not null;default:'';comment:相册描述" json:"album_desc"`            // 相册描述
		AlbumCover string    `gorm:"column:album_cover;type:varchar(255);not null;default:'';comment:相册封面" json:"album_cover"`          // 相册封面
		IsDelete   int64     `gorm:"column:is_delete;type:tinyint(1);not null;default:0;comment:是否删除" json:"is_delete"`                 // 是否删除
		Status     int64     `gorm:"column:status;type:tinyint;not null;default:1;comment:状态值 1公开 2私密" json:"status"`                   // 状态值 1公开 2私密
		CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
		UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
	}

	// 接口实现
	defaultTAlbumModel struct {
		baseModel[TAlbum]
	}
)

func NewTAlbumModel(db *gorm.DB) TAlbumModel {
	return &defaultTAlbumModel{newBaseModel[TAlbum](db, "`t_album`")}
}

// 在事务中操作
func (m *defaultTAlbumModel) WithTx(tx *gorm.DB) (out TAlbumModel) {
	return &defaultTAlbumModel{newBaseModel[TAlbum](tx, m.table)}
}
