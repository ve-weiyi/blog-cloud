package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TUserModel = (*defaultTUserModel)(nil)

type (
	// 接口定义
	TUserModel interface {
		Table[TUser]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TUserModel)
		// add extra method in here
		FindOneByEmail(ctx context.Context, email string) (out *TUser, err error)
		FindOneByMobile(ctx context.Context, mobile string) (out *TUser, err error)
		FindOneByUserId(ctx context.Context, userId string) (out *TUser, err error)
		FindOneByUsername(ctx context.Context, username string) (out *TUser, err error)
	}

	// 表字段定义
	TUser struct {
		Id           int64      `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:ID" json:"id"`                                            // ID
		UserId       string     `gorm:"column:user_id;type:varchar(64);not null;uniqueIndex:uk_uid,priority:1;default:'';comment:用户id" json:"user_id"`       // 用户id
		Username     string     `gorm:"column:username;type:varchar(50);not null;uniqueIndex:uk_username,priority:1;default:'';comment:用户名" json:"username"` // 用户名
		Password     string     `gorm:"column:password;type:varchar(255);not null;default:'';comment:加密后的密码（bcrypt）" json:"password"`                        // 加密后的密码（bcrypt）
		Nickname     string     `gorm:"column:nickname;type:varchar(100);not null;default:'';comment:用户昵称" json:"nickname"`                                  // 用户昵称
		Avatar       string     `gorm:"column:avatar;type:varchar(500);not null;default:'';comment:用户头像URL" json:"avatar"`                                   // 用户头像URL
		Email        *string    `gorm:"column:email;type:varchar(100);uniqueIndex:uk_email,priority:1;comment:邮箱" json:"email"`                              // 邮箱
		Mobile       *string    `gorm:"column:mobile;type:varchar(20);uniqueIndex:uk_mobile,priority:1;comment:手机号" json:"mobile"`                           // 手机号
		Info         *string    `gorm:"column:info;type:json;comment:用户信息" json:"info"`                                                                      // 用户信息
		Status       int64      `gorm:"column:status;type:tinyint;not null;index:idx_status,priority:1;default:0;comment:账号状态：0-正常 1-禁用 2-注销" json:"status"` // 账号状态：0-正常 1-禁用 2-注销
		RegisterType string     `gorm:"column:register_type;type:varchar(64);not null;default:'';comment:注册方式" json:"register_type"`                         // 注册方式
		IpAddress    string     `gorm:"column:ip_address;type:varchar(255);not null;default:'';comment:注册ip" json:"ip_address"`                              // 注册ip
		IpSource     string     `gorm:"column:ip_source;type:varchar(255);not null;default:'';comment:注册ip 源" json:"ip_source"`                              // 注册ip 源
		CreatedAt    time.Time  `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                   // 创建时间
		UpdatedAt    time.Time  `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                   // 更新时间
		DeletedAt    *time.Time `gorm:"column:deleted_at;type:datetime;comment:删除时间，软删除" json:"deleted_at"`                                                  // 删除时间，软删除
	}

	// 接口实现
	defaultTUserModel struct {
		baseModel[TUser]
	}
)

func NewTUserModel(db *gorm.DB) TUserModel {
	return &defaultTUserModel{newBaseModel[TUser](db, "`t_user`")}
}

// 在事务中操作
func (m *defaultTUserModel) WithTx(tx *gorm.DB) (out TUserModel) {
	return &defaultTUserModel{newBaseModel[TUser](tx, m.table)}
}

// add extra method in here
func (m *defaultTUserModel) FindOneByEmail(ctx context.Context, email string) (out *TUser, err error) {
	db := m.scope(ctx)

	err = db.Where("`email` = ?", email).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (m *defaultTUserModel) FindOneByMobile(ctx context.Context, mobile string) (out *TUser, err error) {
	db := m.scope(ctx)

	err = db.Where("`mobile` = ?", mobile).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (m *defaultTUserModel) FindOneByUserId(ctx context.Context, userId string) (out *TUser, err error) {
	db := m.scope(ctx)

	err = db.Where("`user_id` = ?", userId).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (m *defaultTUserModel) FindOneByUsername(ctx context.Context, username string) (out *TUser, err error) {
	db := m.scope(ctx)

	err = db.Where("`username` = ?", username).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
