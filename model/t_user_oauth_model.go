package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TUserOauthModel = (*defaultTUserOauthModel)(nil)

type (
	// 接口定义
	TUserOauthModel interface {
		Table[TUserOauth]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TUserOauthModel)
		// add extra method in here
		FindOneByOpenIdPlatform(ctx context.Context, openId string, platform string) (out *TUserOauth, err error)
		FindOneByUserIdPlatform(ctx context.Context, userId string, platform string) (out *TUserOauth, err error)
	}

	// 表字段定义
	TUserOauth struct {
		Id        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:id" json:"id"`                                                                                           // id
		UserId    string    `gorm:"column:user_id;type:varchar(64);not null;uniqueIndex:uk_uid_plat,priority:1;default:'';comment:用户id" json:"user_id"`                                                 // 用户id
		Platform  string    `gorm:"column:platform;type:varchar(64);not null;uniqueIndex:uk_uid_plat,priority:2;uniqueIndex:uk_oid_plat,priority:2;default:'';comment:平台:手机号、邮箱、微信、飞书" json:"platform"` // 平台:手机号、邮箱、微信、飞书
		OpenId    string    `gorm:"column:open_id;type:varchar(128);not null;uniqueIndex:uk_oid_plat,priority:1;default:'';comment:第三方平台id，标识唯一用户" json:"open_id"`                                      // 第三方平台id，标识唯一用户
		Nickname  string    `gorm:"column:nickname;type:varchar(128);not null;default:'';comment:第三方平台昵称" json:"nickname"`                                                                              // 第三方平台昵称
		Avatar    string    `gorm:"column:avatar;type:varchar(256);not null;default:'';comment:第三方平台头像" json:"avatar"`                                                                                  // 第三方平台头像
		CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                                                                  // 创建时间
		UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                                                                  // 更新时间
	}

	// 接口实现
	defaultTUserOauthModel struct {
		baseModel[TUserOauth]
	}
)

func NewTUserOauthModel(db *gorm.DB) TUserOauthModel {
	return &defaultTUserOauthModel{newBaseModel[TUserOauth](db, "`t_user_oauth`")}
}

// 在事务中操作
func (m *defaultTUserOauthModel) WithTx(tx *gorm.DB) (out TUserOauthModel) {
	return &defaultTUserOauthModel{newBaseModel[TUserOauth](tx, m.table)}
}

// add extra method in here
func (m *defaultTUserOauthModel) FindOneByOpenIdPlatform(ctx context.Context, openId string, platform string) (out *TUserOauth, err error) {
	db := m.scope(ctx)

	err = db.Where("`open_id` = ? and `platform` = ?", openId, platform).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (m *defaultTUserOauthModel) FindOneByUserIdPlatform(ctx context.Context, userId string, platform string) (out *TUserOauth, err error) {
	db := m.scope(ctx)

	err = db.Where("`user_id` = ? and `platform` = ?", userId, platform).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
