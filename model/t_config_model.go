package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TConfigModel = (*defaultTConfigModel)(nil)

type (
	// 接口定义
	TConfigModel interface {
		Table[TConfig]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TConfigModel)
		// add extra method in here
		FindOneByConfigKey(ctx context.Context, configKey string) (out *TConfig, err error)
	}

	// 表字段定义
	TConfig struct {
		Id        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:id" json:"id"`                                                  // id
		ConfigKey string    `gorm:"column:config_key;type:varchar(32);not null;uniqueIndex:uk_config_key,priority:1;default:'';comment:关键词" json:"config_key"` // 关键词
		Config    string    `gorm:"column:config;type:varchar(4096);not null;default:'';comment:配置信息" json:"config"`                                           // 配置信息
		CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                         // 创建时间
		UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                         // 更新时间
	}

	// 接口实现
	defaultTConfigModel struct {
		baseModel[TConfig]
	}
)

func NewTConfigModel(db *gorm.DB) TConfigModel {
	return &defaultTConfigModel{newBaseModel[TConfig](db, "`t_config`")}
}

// 在事务中操作
func (m *defaultTConfigModel) WithTx(tx *gorm.DB) (out TConfigModel) {
	return &defaultTConfigModel{newBaseModel[TConfig](tx, m.table)}
}

// add extra method in here
func (m *defaultTConfigModel) FindOneByConfigKey(ctx context.Context, configKey string) (out *TConfig, err error) {
	db := m.scope(ctx)

	err = db.Where("`config_key` = ?", configKey).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
