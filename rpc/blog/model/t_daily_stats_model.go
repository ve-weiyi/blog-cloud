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
		TableName() string
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TDailyStatsModel)
		// 插入
		Insert(ctx context.Context, in *TDailyStats) (rows int64, err error)
		InsertBatch(ctx context.Context, in ...*TDailyStats) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TDailyStats) (rows int64, err error)
		UpdateFields(ctx context.Context, fields map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存或更新（更新零值）
		Save(ctx context.Context, in *TDailyStats) (rows int64, err error)
		// 查询
		FindById(ctx context.Context, id int64) (out *TDailyStats, err error)
		FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TDailyStats, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TDailyStats, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TDailyStats, total int64, err error)
		// add extra method in here
		FindOneByDate(ctx context.Context, date string) (out *TDailyStats, err error)
	}

	// 表字段定义
	TDailyStats struct {
		Id           int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键ID" json:"id"`                                          // 主键ID
		Date         string    `gorm:"column:date;type:varchar(10);not null;uniqueIndex:uk_date,priority:1;default:'';comment:统计日期 YYYY-MM-DD" json:"date"` // 统计日期 YYYY-MM-DD
		NewUsers     int64     `gorm:"column:new_users;type:bigint;not null;default:0;comment:当日新增用户数" json:"new_users"`                                    // 当日新增用户数
		TotalUsers   int64     `gorm:"column:total_users;type:bigint;not null;default:0;comment:累计用户数" json:"total_users"`                                  // 累计用户数
		ActiveUsers  int64     `gorm:"column:active_users;type:bigint;not null;default:0;comment:当日活跃用户数" json:"active_users"`                              // 当日活跃用户数
		UvCount      int64     `gorm:"column:uv_count;type:bigint;not null;default:0;comment:当日独立访客数(UV)" json:"uv_count"`                                  // 当日独立访客数(UV)
		PvCount      int64     `gorm:"column:pv_count;type:bigint;not null;default:0;comment:当日页面浏览数(PV)" json:"pv_count"`                                  // 当日页面浏览数(PV)
		TotalUvCount int64     `gorm:"column:total_uv_count;type:bigint;not null;default:0;comment:累计访客数" json:"total_uv_count"`                            // 累计访客数
		TotalPvCount int64     `gorm:"column:total_pv_count;type:bigint;not null;default:0;comment:累计浏览量" json:"total_pv_count"`                            // 累计浏览量
		CreatedAt    time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                   // 创建时间
		UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                   // 更新时间
	}

	// 接口实现
	defaultTDailyStatsModel struct {
		DbEngin *gorm.DB
		table   string
	}
)

func NewTDailyStatsModel(db *gorm.DB) TDailyStatsModel {
	return &defaultTDailyStatsModel{
		DbEngin: db,
		table:   "`t_daily_stats`",
	}
}

func (m *defaultTDailyStatsModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTDailyStatsModel) WithTx(tx *gorm.DB) (out TDailyStatsModel) {
	return NewTDailyStatsModel(tx)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTDailyStatsModel) Insert(ctx context.Context, in *TDailyStats) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTDailyStatsModel) InsertBatch(ctx context.Context, in ...*TDailyStats) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTDailyStatsModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TDailyStats{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTDailyStatsModel) DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TDailyStats{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTDailyStatsModel) Update(ctx context.Context, in *TDailyStats) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTDailyStatsModel) UpdateFields(ctx context.Context, feilds map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(feilds)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存或更新（更新零值）
func (m *defaultTDailyStatsModel) Save(ctx context.Context, in *TDailyStats) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTDailyStatsModel) FindById(ctx context.Context, id int64) (out *TDailyStats, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTDailyStatsModel) FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TDailyStats, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	// 如果有排序参数
	if len(sorts) != 0 {
		db = db.Order(sorts)
	}

	err = db.First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询列表
func (m *defaultTDailyStatsModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TDailyStats, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	err = db.Find(&out).Error
	if err != nil {
		return nil, err
	}
	return out, err
}

// 查询总数
func (m *defaultTDailyStatsModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	err = db.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 分页查询记录
func (m *defaultTDailyStatsModel) FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TDailyStats, total int64, err error) {
	// 插入db
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有搜索条件
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	// 如果有排序参数
	if len(sorts) != 0 {
		db = db.Order(sorts)
	}

	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 如果有分页参数
	if page > 0 && size > 0 {
		limit := size
		offset := (page - 1) * limit
		db = db.Limit(limit).Offset(offset)
	}

	// 查询数据
	err = db.Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// add extra method in here
func (m *defaultTDailyStatsModel) FindOneByDate(ctx context.Context, date string) (out *TDailyStats, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`date` = ?", date).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
