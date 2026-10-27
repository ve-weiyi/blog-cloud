package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TGuestModel = (*defaultTGuestModel)(nil)

type (
	// 接口定义
	TGuestModel interface {
		Table[TGuest]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TGuestModel)
		// add extra method in here
		FindOneByDeviceId(ctx context.Context, deviceId string) (out *TGuest, err error)
	}

	// 表字段定义
	TGuest struct {
		Id        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:id" json:"id"`                                          // id
		DeviceId  string    `gorm:"column:device_id;type:varchar(64);not null;uniqueIndex:uk_did,priority:1;default:'';comment:设备id" json:"device_id"` // 设备id
		Os        string    `gorm:"column:os;type:varchar(50);not null;default:'';comment:操作系统" json:"os"`                                             // 操作系统
		Browser   string    `gorm:"column:browser;type:varchar(50);not null;default:'';comment:浏览器" json:"browser"`                                    // 浏览器
		IpAddress string    `gorm:"column:ip_address;type:varchar(255);not null;default:'';comment:操作ip" json:"ip_address"`                            // 操作ip
		IpSource  string    `gorm:"column:ip_source;type:varchar(255);not null;default:'';comment:操作地址" json:"ip_source"`                              // 操作地址
		CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                 // 创建时间
		UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                 // 更新时间
	}

	// 接口实现
	defaultTGuestModel struct {
		baseModel[TGuest]
	}
)

func NewTGuestModel(db *gorm.DB) TGuestModel {
	return &defaultTGuestModel{newBaseModel[TGuest](db, "`t_guest`")}
}

// 在事务中操作
func (m *defaultTGuestModel) WithTx(tx *gorm.DB) (out TGuestModel) {
	return &defaultTGuestModel{newBaseModel[TGuest](tx, m.table)}
}

// add extra method in here
func (m *defaultTGuestModel) FindOneByDeviceId(ctx context.Context, deviceId string) (out *TGuest, err error) {
	db := m.scope(ctx)

	err = db.Where("`device_id` = ?", deviceId).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
