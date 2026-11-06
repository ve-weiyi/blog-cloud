package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TRoleModel = (*defaultTRoleModel)(nil)

type (
	// 接口定义
	TRoleModel interface {
		Table[TRole]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TRoleModel)
		// add extra method in here
		FindOneByRoleKey(ctx context.Context, roleKey string) (out *TRole, err error)
	}

	// 表字段定义
	TRole struct {
		Id          int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键id" json:"id"`                                           // 主键id
		ParentId    int64     `gorm:"column:parent_id;type:bigint;not null;default:0;comment:父角色id" json:"parent_id"`                                       // 父角色id
		RoleKey     string    `gorm:"column:role_key;type:varchar(64);not null;uniqueIndex:uk_role_key,priority:1;default:'';comment:角色标识" json:"role_key"` // 角色标识
		RoleLabel   string    `gorm:"column:role_label;type:varchar(64);not null;default:'';comment:角色标签" json:"role_label"`                                // 角色标签
		RoleComment string    `gorm:"column:role_comment;type:varchar(64);not null;default:'';comment:角色备注" json:"role_comment"`                            // 角色备注
		IsDefault   int64     `gorm:"column:is_default;type:tinyint(1);not null;default:0;comment:是否默认角色 0否 1是" json:"is_default"`                          // 是否默认角色 0否 1是
		Status      int64     `gorm:"column:status;type:tinyint;not null;default:0;comment:状态  0正常 1禁用" json:"status"`                                      // 状态  0正常 1禁用
		CreatedAt   time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                    // 创建时间
		UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                    // 更新时间
	}

	// 接口实现
	defaultTRoleModel struct {
		baseModel[TRole]
	}
)

func NewTRoleModel(db *gorm.DB) TRoleModel {
	return &defaultTRoleModel{newBaseModel[TRole](db, "`t_role`")}
}

// 在事务中操作
func (m *defaultTRoleModel) WithTx(tx *gorm.DB) (out TRoleModel) {
	return &defaultTRoleModel{newBaseModel[TRole](tx, m.table)}
}

// add extra method in here
func (m *defaultTRoleModel) FindOneByRoleKey(ctx context.Context, roleKey string) (out *TRole, err error) {
	db := m.scope(ctx)

	err = db.Where("`role_key` = ?", roleKey).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
