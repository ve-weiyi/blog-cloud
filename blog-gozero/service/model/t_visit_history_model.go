package model

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var _ TVisitHistoryModel = (*defaultTVisitHistoryModel)(nil)

type (
	// 接口定义
	TVisitHistoryModel interface {
		TableName() string
		// 在事务中操作
		WithTransaction(tx *gorm.DB) (out TVisitHistoryModel)
		// 插入
		Insert(ctx context.Context, in *TVisitHistory) (rows int64, err error)
		Inserts(ctx context.Context, in ...*TVisitHistory) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TVisitHistory) (rows int64, err error)
		Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存
		Save(ctx context.Context, in *TVisitHistory) (rows int64, err error)
		// 查询
		FindOne(ctx context.Context, id int64) (out *TVisitHistory, err error)
		First(ctx context.Context, conditions string, args ...interface{}) (out *TVisitHistory, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TVisitHistory, err error)
		FindList(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TVisitHistory, err error)
		// add extra method in here
		FindOneByDate(ctx context.Context, date string) (out *TVisitHistory, err error)
	}

	// 表字段定义
	TVisitHistory struct {
		Id         int64     `json:"id" gorm:"column:id" `                   // id
		Date       string    `json:"date" gorm:"column:date" `               // 日期
		ViewsCount int64     `json:"views_count" gorm:"column:views_count" ` // 访问量
		CreatedAt  time.Time `json:"created_at" gorm:"column:created_at" `   // 创建时间
		UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at" `   // 更新时间
	}

	// 接口实现
	defaultTVisitHistoryModel struct {
		DbEngin    *gorm.DB
		CacheEngin *redis.Client
		table      string
	}
)

func NewTVisitHistoryModel(db *gorm.DB, cache *redis.Client) TVisitHistoryModel {
	return &defaultTVisitHistoryModel{
		DbEngin:    db,
		CacheEngin: cache,
		table:      "`t_visit_history`",
	}
}

func (m *defaultTVisitHistoryModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTVisitHistoryModel) WithTransaction(tx *gorm.DB) (out TVisitHistoryModel) {
	return NewTVisitHistoryModel(tx, m.CacheEngin)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTVisitHistoryModel) Insert(ctx context.Context, in *TVisitHistory) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTVisitHistoryModel) Inserts(ctx context.Context, in ...*TVisitHistory) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTVisitHistoryModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TVisitHistory{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTVisitHistoryModel) Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TVisitHistory{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTVisitHistoryModel) Update(ctx context.Context, in *TVisitHistory) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTVisitHistoryModel) Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(columns)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存记录（更新零值）
func (m *defaultTVisitHistoryModel) Save(ctx context.Context, in *TVisitHistory) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTVisitHistoryModel) FindOne(ctx context.Context, id int64) (out *TVisitHistory, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTVisitHistoryModel) First(ctx context.Context, conditions string, args ...interface{}) (out *TVisitHistory, err error) {
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

// 查询总数
func (m *defaultTVisitHistoryModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	err = db.Model(&TVisitHistory{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 查询列表
func (m *defaultTVisitHistoryModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TVisitHistory, err error) {
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

// 分页查询记录
func (m *defaultTVisitHistoryModel) FindList(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TVisitHistory, err error) {
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

	// 如果有分页参数
	if page > 0 && size > 0 {
		limit := size
		offset := (page - 1) * limit
		db = db.Limit(limit).Offset(offset)
	}

	// 查询数据
	err = db.Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

// add extra method in here
func (m *defaultTVisitHistoryModel) FindOneByDate(ctx context.Context, date string) (out *TVisitHistory, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`date` = ?", date).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
