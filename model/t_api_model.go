package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TApiModel = (*defaultTApiModel)(nil)

type (
	// 接口定义
	TApiModel interface {
		Table[TApi]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TApiModel)
		// add extra method in here
		FindOneByPathMethodName(ctx context.Context, path string, method string, name string) (out *TApi, err error)
	}

	// 表字段定义
	TApi struct {
		Id        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键id" json:"id"`                                                   // 主键id
		ParentId  int64     `gorm:"column:parent_id;type:bigint;not null;default:0;comment:分组id" json:"parent_id"`                                                // 分组id
		Name      string    `gorm:"column:name;type:varchar(128);not null;uniqueIndex:idx_path_method_name,priority:3;default:'';comment:api名称" json:"name"`      // api名称
		Path      string    `gorm:"column:path;type:varchar(128);not null;uniqueIndex:idx_path_method_name,priority:1;default:'';comment:api路径" json:"path"`      // api路径
		Method    string    `gorm:"column:method;type:varchar(16);not null;uniqueIndex:idx_path_method_name,priority:2;default:'';comment:api请求方法" json:"method"` // api请求方法
		Traceable int64     `gorm:"column:traceable;type:tinyint;not null;default:0;comment:是否追溯操作记录 0需要，1是" json:"traceable"`                                    // 是否追溯操作记录 0需要，1是
		Status    int64     `gorm:"column:status;type:tinyint;not null;default:0;comment:是否禁用 0否 1是" json:"status"`                                               // 是否禁用 0否 1是
		CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                            // 创建时间
		UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                            // 更新时间
	}

	// 接口实现
	defaultTApiModel struct {
		baseModel[TApi]
	}
)

func NewTApiModel(db *gorm.DB) TApiModel {
	return &defaultTApiModel{newBaseModel[TApi](db, "`t_api`")}
}

// 在事务中操作
func (m *defaultTApiModel) WithTx(tx *gorm.DB) (out TApiModel) {
	return &defaultTApiModel{newBaseModel[TApi](tx, m.table)}
}

// add extra method in here
func (m *defaultTApiModel) FindOneByPathMethodName(ctx context.Context, path string, method string, name string) (out *TApi, err error) {
	db := m.scope(ctx)

	err = db.Where("`path` = ? and `method` = ? and `name` = ?", path, method, name).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
