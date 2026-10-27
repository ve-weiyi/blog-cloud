package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TFriendModel = (*defaultTFriendModel)(nil)

type (
	// 接口定义
	TFriendModel interface {
		Table[TFriend]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TFriendModel)
		// add extra method in here
		FindOneByLinkName(ctx context.Context, linkName string) (out *TFriend, err error)
	}

	// 表字段定义
	TFriend struct {
		Id          int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:id" json:"id"`                                          // id
		LinkName    string    `gorm:"column:link_name;type:varchar(32);not null;uniqueIndex:uk_name,priority:1;default:'';comment:链接名" json:"link_name"` // 链接名
		LinkAvatar  string    `gorm:"column:link_avatar;type:varchar(255);not null;default:'';comment:链接头像" json:"link_avatar"`                          // 链接头像
		LinkAddress string    `gorm:"column:link_address;type:varchar(64);not null;default:'';comment:链接地址" json:"link_address"`                         // 链接地址
		LinkIntro   string    `gorm:"column:link_intro;type:varchar(100);not null;default:'';comment:链接介绍" json:"link_intro"`                            // 链接介绍
		CreatedAt   time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                 // 创建时间
		UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                 // 更新时间
	}

	// 接口实现
	defaultTFriendModel struct {
		baseModel[TFriend]
	}
)

func NewTFriendModel(db *gorm.DB) TFriendModel {
	return &defaultTFriendModel{newBaseModel[TFriend](db, "`t_friend`")}
}

// 在事务中操作
func (m *defaultTFriendModel) WithTx(tx *gorm.DB) (out TFriendModel) {
	return &defaultTFriendModel{newBaseModel[TFriend](tx, m.table)}
}

// add extra method in here
func (m *defaultTFriendModel) FindOneByLinkName(ctx context.Context, linkName string) (out *TFriend, err error) {
	db := m.scope(ctx)

	err = db.Where("`link_name` = ?", linkName).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
