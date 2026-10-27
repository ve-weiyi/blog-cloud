package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TNotifyTemplateModel = (*defaultTNotifyTemplateModel)(nil)

type (
	// 接口定义
	TNotifyTemplateModel interface {
		Table[TNotifyTemplate]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TNotifyTemplateModel)
		// add extra method in here
		FindOneBySceneChannel(ctx context.Context, scene string, channel string) (out *TNotifyTemplate, err error)
	}

	// 表字段定义
	TNotifyTemplate struct {
		Id        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:ID" json:"id"`                                                                     // ID
		Code      string    `gorm:"column:code;type:varchar(64);not null;default:'';comment:模板编码" json:"code"`                                                                    // 模板编码
		Channel   string    `gorm:"column:channel;type:varchar(16);not null;uniqueIndex:uk_scene_channel,priority:2;default:'';comment:渠道：sms | email | inbox" json:"channel"`    // 渠道：sms | email | inbox
		Scene     string    `gorm:"column:scene;type:varchar(32);not null;uniqueIndex:uk_scene_channel,priority:1;default:'';comment:业务场景（login/register/notify 等）" json:"scene"` // 业务场景（login/register/notify 等）
		Title     string    `gorm:"column:title;type:varchar(128);not null;default:'';comment:标题（邮件/站内通知）" json:"title"`                                                          // 标题（邮件/站内通知）
		Content   string    `gorm:"column:content;type:varchar(2048);not null;default:'';comment:模板内容" json:"content"`                                                            // 模板内容
		Enabled   int64     `gorm:"column:enabled;type:tinyint;not null;default:1;comment:是否启用" json:"enabled"`                                                                   // 是否启用
		CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                                            // 创建时间
		UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                                            // 更新时间
	}

	// 接口实现
	defaultTNotifyTemplateModel struct {
		baseModel[TNotifyTemplate]
	}
)

func NewTNotifyTemplateModel(db *gorm.DB) TNotifyTemplateModel {
	return &defaultTNotifyTemplateModel{newBaseModel[TNotifyTemplate](db, "`t_notify_template`")}
}

// 在事务中操作
func (m *defaultTNotifyTemplateModel) WithTx(tx *gorm.DB) (out TNotifyTemplateModel) {
	return &defaultTNotifyTemplateModel{newBaseModel[TNotifyTemplate](tx, m.table)}
}

// add extra method in here
func (m *defaultTNotifyTemplateModel) FindOneBySceneChannel(ctx context.Context, scene string, channel string) (out *TNotifyTemplate, err error) {
	db := m.scope(ctx)

	err = db.Where("`scene` = ? and `channel` = ?", scene, channel).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
