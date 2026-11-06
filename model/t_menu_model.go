package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TMenuModel = (*defaultTMenuModel)(nil)

type (
	// 接口定义
	TMenuModel interface {
		Table[TMenu]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TMenuModel)
		// add extra method in here
		FindOneByPathPerm(ctx context.Context, path string, perm string) (out *TMenu, err error)
	}

	// 表字段定义
	TMenu struct {
		Id         int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键" json:"id"`                                      // 主键
		ParentId   int64     `gorm:"column:parent_id;type:bigint;not null;default:0;comment:父id" json:"parent_id"`                                  // 父id
		Path       string    `gorm:"column:path;type:varchar(64);not null;uniqueIndex:uk_path_perm,priority:1;default:'';comment:路由路径" json:"path"` // 路由路径
		Name       string    `gorm:"column:name;type:varchar(64);not null;default:'';comment:路由名称" json:"name"`                                     // 路由名称
		Component  string    `gorm:"column:component;type:varchar(256);not null;default:'';comment:路由组件" json:"component"`                          // 路由组件
		Redirect   string    `gorm:"column:redirect;type:varchar(256);not null;default:'';comment:路由重定向" json:"redirect"`                           // 路由重定向
		Type       string    `gorm:"column:type;type:varchar(64);not null;default:'';comment:菜单类型" json:"type"`                                     // 菜单类型
		Title      string    `gorm:"column:title;type:varchar(64);not null;default:'';comment:菜单标题" json:"title"`                                   // 菜单标题
		Icon       string    `gorm:"column:icon;type:varchar(64);not null;default:'';comment:菜单图标" json:"icon"`                                     // 菜单图标
		Rank       int64     `gorm:"column:rank;type:int;not null;default:0;comment:排序" json:"rank"`                                                // 排序
		Perm       string    `gorm:"column:perm;type:varchar(64);not null;uniqueIndex:uk_path_perm,priority:2;default:'';comment:权限标识" json:"perm"` // 权限标识
		Params     string    `gorm:"column:params;type:varchar(256);not null;default:'';comment:路由参数" json:"params"`                                // 路由参数
		KeepAlive  int64     `gorm:"column:keep_alive;type:tinyint;not null;default:0;comment:是否缓存" json:"keep_alive"`                              // 是否缓存
		AlwaysShow int64     `gorm:"column:always_show;type:tinyint;not null;default:0;comment:是否一直显示菜单" json:"always_show"`                        // 是否一直显示菜单
		Visible    int64     `gorm:"column:visible;type:tinyint;not null;default:0;comment:菜单是否可见" json:"visible"`                                  // 菜单是否可见
		Status     int64     `gorm:"column:status;type:tinyint;not null;default:0;comment:是否禁用" json:"status"`                                      // 是否禁用
		Extra      string    `gorm:"column:extra;type:varchar(1024);not null;default:'';comment:菜单元数据" json:"extra"`                                // 菜单元数据
		CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`             // 创建时间
		UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`             // 更新时间
	}

	// 接口实现
	defaultTMenuModel struct {
		baseModel[TMenu]
	}
)

func NewTMenuModel(db *gorm.DB) TMenuModel {
	return &defaultTMenuModel{newBaseModel[TMenu](db, "`t_menu`")}
}

// 在事务中操作
func (m *defaultTMenuModel) WithTx(tx *gorm.DB) (out TMenuModel) {
	return &defaultTMenuModel{newBaseModel[TMenu](tx, m.table)}
}

// add extra method in here
func (m *defaultTMenuModel) FindOneByPathPerm(ctx context.Context, path string, perm string) (out *TMenu, err error) {
	db := m.scope(ctx)

	err = db.Where("`path` = ? and `perm` = ?", path, perm).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
