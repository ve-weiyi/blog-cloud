package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TVisitDailyStatsModel = (*defaultTVisitDailyStatsModel)(nil)

type (
	// 接口定义
	TVisitDailyStatsModel interface {
		TableName() string
		// 在事务中操作
		WithTransaction(tx *gorm.DB) (out TVisitDailyStatsModel)
		// 插入
		Insert(ctx context.Context, in *TVisitDailyStats) (rows int64, err error)
		Inserts(ctx context.Context, in ...*TVisitDailyStats) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TVisitDailyStats) (rows int64, err error)
		Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存
		Save(ctx context.Context, in *TVisitDailyStats) (rows int64, err error)
		// 查询
		FindById(ctx context.Context, id int64) (out *TVisitDailyStats, err error)
		FindOne(ctx context.Context, conditions string, args ...interface{}) (out *TVisitDailyStats, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TVisitDailyStats, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TVisitDailyStats, total int64, err error)
		// add extra method in here
		FindOneByDateVisitType(ctx context.Context, date string, visit_type int64) (out *TVisitDailyStats, err error)
	}

	// 表字段定义
	TVisitDailyStats struct {
		Id        int64     `json:"id" gorm:"column:id"`                 // id
		Date      string    `json:"date" gorm:"column:date"`             // 日期
		ViewCount int64     `json:"view_count" gorm:"column:view_count"` // 访问量
		VisitType int64     `json:"visit_type" gorm:"column:visit_type"` // 1 访客数 2 浏览数
		CreatedAt time.Time `json:"created_at" gorm:"column:created_at"` // 创建时间
		UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"` // 更新时间
	}

	// 接口实现
	defaultTVisitDailyStatsModel struct {
		DbEngin *gorm.DB
		table   string
	}
)

func NewTVisitDailyStatsModel(db *gorm.DB) TVisitDailyStatsModel {
	return &defaultTVisitDailyStatsModel{
		DbEngin: db,
		table:   "`t_visit_daily_stats`",
	}
}

func (m *defaultTVisitDailyStatsModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTVisitDailyStatsModel) WithTransaction(tx *gorm.DB) (out TVisitDailyStatsModel) {
	return NewTVisitDailyStatsModel(tx)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTVisitDailyStatsModel) Insert(ctx context.Context, in *TVisitDailyStats) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTVisitDailyStatsModel) Inserts(ctx context.Context, in ...*TVisitDailyStats) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTVisitDailyStatsModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TVisitDailyStats{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTVisitDailyStatsModel) Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TVisitDailyStats{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTVisitDailyStatsModel) Update(ctx context.Context, in *TVisitDailyStats) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTVisitDailyStatsModel) Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(columns)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存记录（更新零值）
func (m *defaultTVisitDailyStatsModel) Save(ctx context.Context, in *TVisitDailyStats) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTVisitDailyStatsModel) FindById(ctx context.Context, id int64) (out *TVisitDailyStats, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTVisitDailyStatsModel) FindOne(ctx context.Context, conditions string, args ...interface{}) (out *TVisitDailyStats, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	err = db.First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询列表
func (m *defaultTVisitDailyStatsModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TVisitDailyStats, err error) {
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
func (m *defaultTVisitDailyStatsModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
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
func (m *defaultTVisitDailyStatsModel) FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TVisitDailyStats, total int64, err error) {
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
func (m *defaultTVisitDailyStatsModel) FindOneByDateVisitType(ctx context.Context, date string, visit_type int64) (out *TVisitDailyStats, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`date` = ? and `visit_type` = ?", date, visit_type).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
