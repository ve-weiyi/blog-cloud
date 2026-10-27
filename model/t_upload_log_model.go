package model

import (
	"time"

	"gorm.io/gorm"
)

var _ TUploadLogModel = (*defaultTUploadLogModel)(nil)

type (
	// 接口定义
	TUploadLogModel interface {
		Table[TUploadLog]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TUploadLogModel)
	}

	// 表字段定义
	TUploadLog struct {
		Id        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:id" json:"id"`                                       // id
		UserId    string    `gorm:"column:user_id;type:varchar(64);not null;index:idx_uid,priority:1;default:'';comment:用户id" json:"user_id"`       // 用户id
		DeviceId  string    `gorm:"column:device_id;type:varchar(64);not null;default:'';comment:设备id" json:"device_id"`                            // 设备id
		FileBase  string    `gorm:"column:file_base;type:varchar(128);not null;index:idx_path,priority:1;default:'';comment:文件路径" json:"file_base"` // 文件路径
		FileName  string    `gorm:"column:file_name;type:varchar(128);not null;default:'';comment:文件名称" json:"file_name"`                           // 文件名称
		FileType  string    `gorm:"column:file_type;type:varchar(128);not null;default:'';comment:文件类型" json:"file_type"`                           // 文件类型
		FileSize  int64     `gorm:"column:file_size;type:bigint;not null;default:0;comment:文件大小" json:"file_size"`                                  // 文件大小
		FileMd5   string    `gorm:"column:file_md5;type:char(32);not null;default:'';comment:文件md5值" json:"file_md5"`                               // 文件md5值
		FileUrl   string    `gorm:"column:file_url;type:varchar(256);not null;default:'';comment:上传路径" json:"file_url"`                             // 上传路径
		CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`              // 创建时间
		UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`              // 更新时间
	}

	// 接口实现
	defaultTUploadLogModel struct {
		baseModel[TUploadLog]
	}
)

func NewTUploadLogModel(db *gorm.DB) TUploadLogModel {
	return &defaultTUploadLogModel{newBaseModel[TUploadLog](db, "`t_upload_log`")}
}

// 在事务中操作
func (m *defaultTUploadLogModel) WithTx(tx *gorm.DB) (out TUploadLogModel) {
	return &defaultTUploadLogModel{newBaseModel[TUploadLog](tx, m.table)}
}
