package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TDailyStatsModel = (*defaultTDailyStatsModel)(nil)

type (
	// 接口定义
	TDailyStatsModel interface {
		Table[TDailyStats]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TDailyStatsModel)
		// add extra method in here
		FindOneByDate(ctx context.Context, date time.Time) (out *TDailyStats, err error)
	}

	// 表字段定义
	TDailyStats struct {
		Id           int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键ID" json:"id"`                        // 主键ID
		Date         time.Time `gorm:"column:date;type:date;not null;uniqueIndex:uk_date,priority:1;comment:统计日期 YYYY-MM-DD" json:"date"` // 统计日期 YYYY-MM-DD
		NewUsers     int64     `gorm:"column:new_users;type:bigint;not null;default:0;comment:当日新增用户数" json:"new_users"`                  // 当日新增用户数
		TotalUsers   int64     `gorm:"column:total_users;type:bigint;not null;default:0;comment:累计用户数" json:"total_users"`                // 累计用户数
		ActiveUsers  int64     `gorm:"column:active_users;type:bigint;not null;default:0;comment:当日活跃用户数" json:"active_users"`            // 当日活跃用户数
		UvCount      int64     `gorm:"column:uv_count;type:bigint;not null;default:0;comment:当日独立访客数(UV)" json:"uv_count"`                // 当日独立访客数(UV)
		PvCount      int64     `gorm:"column:pv_count;type:bigint;not null;default:0;comment:当日页面浏览数(PV)" json:"pv_count"`                // 当日页面浏览数(PV)
		TotalUvCount int64     `gorm:"column:total_uv_count;type:bigint;not null;default:0;comment:累计访客数" json:"total_uv_count"`          // 累计访客数
		TotalPvCount int64     `gorm:"column:total_pv_count;type:bigint;not null;default:0;comment:累计浏览量" json:"total_pv_count"`          // 累计浏览量
		CreatedAt    time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
		UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
	}

	// 接口实现
	defaultTDailyStatsModel struct {
		baseModel[TDailyStats]
	}
)

func NewTDailyStatsModel(db *gorm.DB) TDailyStatsModel {
	return &defaultTDailyStatsModel{newBaseModel[TDailyStats](db, "`t_daily_stats`")}
}

// 在事务中操作
func (m *defaultTDailyStatsModel) WithTx(tx *gorm.DB) (out TDailyStatsModel) {
	return &defaultTDailyStatsModel{newBaseModel[TDailyStats](tx, m.table)}
}

// add extra method in here
func (m *defaultTDailyStatsModel) FindOneByDate(ctx context.Context, date time.Time) (out *TDailyStats, err error) {
	db := m.scope(ctx)

	err = db.Where("`date` = ?", date).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
