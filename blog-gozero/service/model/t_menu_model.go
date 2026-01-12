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
		TableName() string
		// 在事务中操作
		WithTransaction(tx *gorm.DB) (out TMenuModel)
		// 插入
		Insert(ctx context.Context, in *TMenu) (rows int64, err error)
		Inserts(ctx context.Context, in ...*TMenu) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TMenu) (rows int64, err error)
		Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存
		Save(ctx context.Context, in *TMenu) (rows int64, err error)
		// 查询
		FindById(ctx context.Context, id int64) (out *TMenu, err error)
		FindOne(ctx context.Context, conditions string, args ...interface{}) (out *TMenu, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TMenu, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TMenu, total int64, err error)
		// add extra method in here
		FindOneByPath(ctx context.Context, path string) (out *TMenu, err error)
	}

	// 表字段定义
	TMenu struct {
		Id         int64     `json:"id" gorm:"column:id"`                   // 主键
		ParentId   int64     `json:"parent_id" gorm:"column:parent_id"`     // 父id
		Path       string    `json:"path" gorm:"column:path"`               // 路由路径
		Name       string    `json:"name" gorm:"column:name"`               // 路由名称
		Component  string    `json:"component" gorm:"column:component"`     // 路由组件
		Redirect   string    `json:"redirect" gorm:"column:redirect"`       // 路由重定向
		Type       string    `json:"type" gorm:"column:type"`               // 菜单类型
		Title      string    `json:"title" gorm:"column:title"`             // 菜单标题
		Icon       string    `json:"icon" gorm:"column:icon"`               // 菜单图标
		Rank       int64     `json:"rank" gorm:"column:rank"`               // 排序
		Perm       string    `json:"perm" gorm:"column:perm"`               // 权限标识
		Params     string    `json:"params" gorm:"column:params"`           // 路由参数
		KeepAlive  int64     `json:"keep_alive" gorm:"column:keep_alive"`   // 是否缓存
		AlwaysShow int64     `json:"always_show" gorm:"column:always_show"` // 是否一直显示菜单
		Visible    int64     `json:"visible" gorm:"column:visible"`         // 菜单是否可见
		Status     int64     `json:"status" gorm:"column:status"`           // 是否禁用
		Extra      string    `json:"extra" gorm:"column:extra"`             // 菜单元数据
		CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`   // 创建时间
		UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`   // 更新时间
	}

	// 接口实现
	defaultTMenuModel struct {
		DbEngin *gorm.DB
		table   string
	}
)

func NewTMenuModel(db *gorm.DB) TMenuModel {
	return &defaultTMenuModel{
		DbEngin: db,
		table:   "`t_menu`",
	}
}

func (m *defaultTMenuModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTMenuModel) WithTransaction(tx *gorm.DB) (out TMenuModel) {
	return NewTMenuModel(tx)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTMenuModel) Insert(ctx context.Context, in *TMenu) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTMenuModel) Inserts(ctx context.Context, in ...*TMenu) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTMenuModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TMenu{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTMenuModel) Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TMenu{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTMenuModel) Update(ctx context.Context, in *TMenu) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTMenuModel) Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(columns)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存记录（更新零值）
func (m *defaultTMenuModel) Save(ctx context.Context, in *TMenu) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTMenuModel) FindById(ctx context.Context, id int64) (out *TMenu, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTMenuModel) FindOne(ctx context.Context, conditions string, args ...interface{}) (out *TMenu, err error) {
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
func (m *defaultTMenuModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TMenu, err error) {
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
func (m *defaultTMenuModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
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
func (m *defaultTMenuModel) FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TMenu, total int64, err error) {
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
func (m *defaultTMenuModel) FindOneByPath(ctx context.Context, path string) (out *TMenu, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`path` = ?", path).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
