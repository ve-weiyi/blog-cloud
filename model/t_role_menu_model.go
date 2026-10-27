package model

import (
	"gorm.io/gorm"
)

var _ TRoleMenuModel = (*defaultTRoleMenuModel)(nil)

type (
	// 接口定义
	TRoleMenuModel interface {
		Table[TRoleMenu]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TRoleMenuModel)
	}

	// 表字段定义
	TRoleMenu struct {
		Id     int64 `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键id" json:"id"`                          // 主键id
		RoleId int64 `gorm:"column:role_id;type:bigint;not null;index:idx_role,priority:1;default:0;comment:角色id" json:"role_id"` // 角色id
		MenuId int64 `gorm:"column:menu_id;type:bigint;not null;default:0;comment:菜单id" json:"menu_id"`                           // 菜单id
	}

	// 接口实现
	defaultTRoleMenuModel struct {
		baseModel[TRoleMenu]
	}
)

func NewTRoleMenuModel(db *gorm.DB) TRoleMenuModel {
	return &defaultTRoleMenuModel{newBaseModel[TRoleMenu](db, "`t_role_menu`")}
}

// 在事务中操作
func (m *defaultTRoleMenuModel) WithTx(tx *gorm.DB) (out TRoleMenuModel) {
	return &defaultTRoleMenuModel{newBaseModel[TRoleMenu](tx, m.table)}
}
