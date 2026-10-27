package model

import (
	"time"

	"gorm.io/gorm"
)

var _ TVisitLogModel = (*defaultTVisitLogModel)(nil)

type (
	// 接口定义
	TVisitLogModel interface {
		Table[TVisitLog]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TVisitLogModel)
	}

	// 表字段定义
	TVisitLog struct {
		Id        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:id" json:"id"`                                 // id
		UserId    string    `gorm:"column:user_id;type:varchar(64);not null;index:idx_uid,priority:1;default:'';comment:用户id" json:"user_id"` // 用户id
		DeviceId  string    `gorm:"column:device_id;type:varchar(64);not null;default:'';comment:设备id" json:"device_id"`                      // 设备id
		PageName  string    `gorm:"column:page_name;type:varchar(64);not null;default:'';comment:访问页面" json:"page_name"`                      // 访问页面
		CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`        // 创建时间
		UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`        // 更新时间
	}

	// 接口实现
	defaultTVisitLogModel struct {
		baseModel[TVisitLog]
	}
)

func NewTVisitLogModel(db *gorm.DB) TVisitLogModel {
	return &defaultTVisitLogModel{newBaseModel[TVisitLog](db, "`t_visit_log`")}
}

// 在事务中操作
func (m *defaultTVisitLogModel) WithTx(tx *gorm.DB) (out TVisitLogModel) {
	return &defaultTVisitLogModel{newBaseModel[TVisitLog](tx, m.table)}
}
