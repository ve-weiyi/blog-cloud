package model

import (
	"gorm.io/gorm"
)

var _ TUserRoleModel = (*defaultTUserRoleModel)(nil)

type (
	// 接口定义
	TUserRoleModel interface {
		Table[TUserRole]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TUserRoleModel)
	}

	// 表字段定义
	TUserRole struct {
		Id     int64  `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键id" json:"id"`                                // 主键id
		UserId string `gorm:"column:user_id;type:varchar(64);not null;index:idx_user,priority:1;default:'';comment:用户id" json:"user_id"` // 用户id
		RoleId int64  `gorm:"column:role_id;type:bigint;not null;index:idx_role,priority:1;default:0;comment:角色id" json:"role_id"`       // 角色id
	}

	// 接口实现
	defaultTUserRoleModel struct {
		baseModel[TUserRole]
	}
)

func NewTUserRoleModel(db *gorm.DB) TUserRoleModel {
	return &defaultTUserRoleModel{newBaseModel[TUserRole](db, "`t_user_role`")}
}

// 在事务中操作
func (m *defaultTUserRoleModel) WithTx(tx *gorm.DB) (out TUserRoleModel) {
	return &defaultTUserRoleModel{newBaseModel[TUserRole](tx, m.table)}
}
