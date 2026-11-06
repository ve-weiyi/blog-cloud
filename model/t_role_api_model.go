package model

import (
	"gorm.io/gorm"
)

var _ TRoleApiModel = (*defaultTRoleApiModel)(nil)

type (
	// 接口定义
	TRoleApiModel interface {
		Table[TRoleApi]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TRoleApiModel)
	}

	// 表字段定义
	TRoleApi struct {
		Id     int64 `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键id" json:"id"`                          // 主键id
		RoleId int64 `gorm:"column:role_id;type:bigint;not null;index:idx_role,priority:1;default:0;comment:角色id" json:"role_id"` // 角色id
		ApiId  int64 `gorm:"column:api_id;type:bigint;not null;default:0;comment:接口id" json:"api_id"`                             // 接口id
	}

	// 接口实现
	defaultTRoleApiModel struct {
		baseModel[TRoleApi]
	}
)

func NewTRoleApiModel(db *gorm.DB) TRoleApiModel {
	return &defaultTRoleApiModel{newBaseModel[TRoleApi](db, "`t_role_api`")}
}

// 在事务中操作
func (m *defaultTRoleApiModel) WithTx(tx *gorm.DB) (out TRoleApiModel) {
	return &defaultTRoleApiModel{newBaseModel[TRoleApi](tx, m.table)}
}
