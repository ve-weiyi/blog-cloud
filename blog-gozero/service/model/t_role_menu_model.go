package model

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var _ TRoleMenuModel = (*defaultTRoleMenuModel)(nil)

type (
	// 接口定义
	TRoleMenuModel interface {
		TableName() string
		// 在事务中操作
		WithTransaction(tx *gorm.DB) (out TRoleMenuModel)
		// 插入
		Insert(ctx context.Context, in *TRoleMenu) (rows int64, err error)
		Inserts(ctx context.Context, in ...*TRoleMenu) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TRoleMenu) (rows int64, err error)
		Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存
		Save(ctx context.Context, in *TRoleMenu) (rows int64, err error)
		// 查询
		FindOne(ctx context.Context, id int64) (out *TRoleMenu, err error)
		First(ctx context.Context, conditions string, args ...interface{}) (out *TRoleMenu, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TRoleMenu, err error)
		FindList(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TRoleMenu, err error)
		// add extra method in here
	}

	// 表字段定义
	TRoleMenu struct {
		Id     int64 `json:"id" gorm:"column:id" `           // 主键id
		RoleId int64 `json:"role_id" gorm:"column:role_id" ` // 角色id
		MenuId int64 `json:"menu_id" gorm:"column:menu_id" ` // 菜单id
	}

	// 接口实现
	defaultTRoleMenuModel struct {
		DbEngin    *gorm.DB
		CacheEngin *redis.Client
		table      string
	}
)

func NewTRoleMenuModel(db *gorm.DB, cache *redis.Client) TRoleMenuModel {
	return &defaultTRoleMenuModel{
		DbEngin:    db,
		CacheEngin: cache,
		table:      "`t_role_menu`",
	}
}

func (m *defaultTRoleMenuModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTRoleMenuModel) WithTransaction(tx *gorm.DB) (out TRoleMenuModel) {
	return NewTRoleMenuModel(tx, m.CacheEngin)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTRoleMenuModel) Insert(ctx context.Context, in *TRoleMenu) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTRoleMenuModel) Inserts(ctx context.Context, in ...*TRoleMenu) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTRoleMenuModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TRoleMenu{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTRoleMenuModel) Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TRoleMenu{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTRoleMenuModel) Update(ctx context.Context, in *TRoleMenu) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTRoleMenuModel) Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(columns)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存记录（更新零值）
func (m *defaultTRoleMenuModel) Save(ctx context.Context, in *TRoleMenu) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTRoleMenuModel) FindOne(ctx context.Context, id int64) (out *TRoleMenu, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTRoleMenuModel) First(ctx context.Context, conditions string, args ...interface{}) (out *TRoleMenu, err error) {
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
func (m *defaultTRoleMenuModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	err = db.Model(&TRoleMenu{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 查询列表
func (m *defaultTRoleMenuModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TRoleMenu, err error) {
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
func (m *defaultTRoleMenuModel) FindList(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TRoleMenu, err error) {
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
