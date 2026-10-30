package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TNotificationPreferenceModel = (*defaultTNotificationPreferenceModel)(nil)

type (
	// 接口定义
	TNotificationPreferenceModel interface {
		Table[TNotificationPreference]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TNotificationPreferenceModel)
		// add extra method in here
		FindOneByUserIdPrefGroupChannel(ctx context.Context, userId string, prefGroup string, channel string) (out *TNotificationPreference, err error)
	}

	// 表字段定义
	TNotificationPreference struct {
		Id        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:ID" json:"id"`                                                                                                             // ID
		UserId    string    `gorm:"column:user_id;type:varchar(64);not null;uniqueIndex:uk_user_group_channel,priority:1;default:'';comment:用户ID" json:"user_id"`                                                         // 用户ID
		PrefGroup string    `gorm:"column:pref_group;type:varchar(32);not null;uniqueIndex:uk_user_group_channel,priority:2;default:'';comment:偏好分组：announcement | comment | like | review | security" json:"pref_group"` // 偏好分组：announcement | comment | like | review | security
		Channel   string    `gorm:"column:channel;type:varchar(16);not null;uniqueIndex:uk_user_group_channel,priority:3;default:'';comment:渠道：inbox | email | sms | push" json:"channel"`                                // 渠道：inbox | email | sms | push
		Enabled   int64     `gorm:"column:enabled;type:tinyint;not null;default:1;comment:是否接收：1 接收 0 关闭" json:"enabled"`                                                                                                 // 是否接收：1 接收 0 关闭
		CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                                                                                    // 创建时间
		UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                                                                                    // 更新时间
	}

	// 接口实现
	defaultTNotificationPreferenceModel struct {
		baseModel[TNotificationPreference]
	}
)

func NewTNotificationPreferenceModel(db *gorm.DB) TNotificationPreferenceModel {
	return &defaultTNotificationPreferenceModel{newBaseModel[TNotificationPreference](db, "`t_notification_preference`")}
}

// 在事务中操作
func (m *defaultTNotificationPreferenceModel) WithTx(tx *gorm.DB) (out TNotificationPreferenceModel) {
	return &defaultTNotificationPreferenceModel{newBaseModel[TNotificationPreference](tx, m.table)}
}

// add extra method in here
func (m *defaultTNotificationPreferenceModel) FindOneByUserIdPrefGroupChannel(ctx context.Context, userId string, prefGroup string, channel string) (out *TNotificationPreference, err error) {
	db := m.scope(ctx)

	err = db.Where("`user_id` = ? and `pref_group` = ? and `channel` = ?", userId, prefGroup, channel).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
