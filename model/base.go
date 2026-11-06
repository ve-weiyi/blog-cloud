package model

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Table 是所有表共用的数据访问方法，由 baseModel 统一实现。
type Table[T any] interface {
	TableName() string
	// 插入
	Insert(ctx context.Context, in *T) (rows int64, err error)
	InsertBatch(ctx context.Context, in ...*T) (rows int64, err error)
	// 插入，唯一键冲突时不修改任何列、返回 0 行
	InsertIgnore(ctx context.Context, in *T) (rows int64, err error)
	// 批量插入，唯一键冲突的行跳过，其余照常写入
	InsertIgnoreBatch(ctx context.Context, in ...*T) (rows int64, err error)
	// 删除
	Delete(ctx context.Context, id int64) (rows int64, err error)
	DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
	// 整行更新（含零值，按主键）
	Update(ctx context.Context, in *T) (rows int64, err error)
	// 部分更新
	UpdateFields(ctx context.Context, fields map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
	// 插入或整行更新（按主键）
	Upsert(ctx context.Context, in *T) (rows int64, err error)
	// 查询
	FindById(ctx context.Context, id int64) (out *T, err error)
	FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *T, err error)
	FindAll(ctx context.Context, conditions string, args ...interface{}) (list []*T, err error)
	FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
	FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*T, total int64, err error)
}

// baseModel 是所有表共用的实现，各表的 defaultXxxModel 内嵌它。
type baseModel[T any] struct {
	db    *gorm.DB
	table string
}

func newBaseModel[T any](db *gorm.DB, table string) baseModel[T] {
	return baseModel[T]{db: db, table: table}
}

func (m *baseModel[T]) TableName() string {
	return m.table
}

func (m *baseModel[T]) scope(ctx context.Context) *gorm.DB {
	return m.db.WithContext(ctx).Table(m.table)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *baseModel[T]) Insert(ctx context.Context, in *T) (rows int64, err error) {
	result := m.scope(ctx).Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
// 插入记录，唯一键冲突时忽略冲突行
func (m *baseModel[T]) InsertIgnore(ctx context.Context, in *T) (rows int64, err error) {
	result := m.scope(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

// 批量插入，唯一键冲突的行跳过
func (m *baseModel[T]) InsertIgnoreBatch(ctx context.Context, in ...*T) (rows int64, err error) {
	if len(in) == 0 {
		return 0, nil
	}

	result := m.scope(ctx).Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func (m *baseModel[T]) InsertBatch(ctx context.Context, in ...*T) (rows int64, err error) {
	result := m.scope(ctx).CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *baseModel[T]) Delete(ctx context.Context, id int64) (rows int64, err error) {
	result := m.scope(ctx).Where("id = ?", id).Delete(new(T))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *baseModel[T]) DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.scope(ctx)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(new(T))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 整行更新（含零值，按主键）
func (m *baseModel[T]) Update(ctx context.Context, in *T) (rows int64, err error) {
	result := m.scope(ctx).Omit("id", "created_at", "updated_at").Select("*").Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 部分更新
func (m *baseModel[T]) UpdateFields(ctx context.Context, fields map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	result := m.scope(ctx).Where(conditions, args...).Updates(fields)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入或整行更新（按主键）
func (m *baseModel[T]) Upsert(ctx context.Context, in *T) (rows int64, err error) {
	result := m.scope(ctx).Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *baseModel[T]) FindById(ctx context.Context, id int64) (out *T, err error) {
	err = m.scope(ctx).Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *baseModel[T]) FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *T, err error) {
	db := m.scope(ctx)

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
func (m *baseModel[T]) FindAll(ctx context.Context, conditions string, args ...interface{}) (list []*T, err error) {
	db := m.scope(ctx)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	err = db.Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, err
}

// 查询总数
func (m *baseModel[T]) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
	db := m.scope(ctx)

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
func (m *baseModel[T]) FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*T, total int64, err error) {
	db := m.scope(ctx)

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
