package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TRoleModel = (*defaultTRoleModel)(nil)

type (
	// 接口定义
	TRoleModel interface {
		TableName() string
		// 在事务中操作
		WithTransaction(tx *gorm.DB) (out TRoleModel)
		// 插入
		Insert(ctx context.Context, in *TRole) (rows int64, err error)
		Inserts(ctx context.Context, in ...*TRole) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TRole) (rows int64, err error)
		Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存
		Save(ctx context.Context, in *TRole) (rows int64, err error)
		// 查询
		FindById(ctx context.Context, id int64) (out *TRole, err error)
		FindOne(ctx context.Context, conditions string, args ...interface{}) (out *TRole, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TRole, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TRole, total int64, err error)
		// add extra method in here
	}

	// 表字段定义
	TRole struct {
		Id          int64     `json:"id" gorm:"column:id"`                     // 主键id
		ParentId    int64     `json:"parent_id" gorm:"column:parent_id"`       // 父角色id
		RoleKey     string    `json:"role_key" gorm:"column:role_key"`         // 角色标识
		RoleLabel   string    `json:"role_label" gorm:"column:role_label"`     // 角色标签
		RoleComment string    `json:"role_comment" gorm:"column:role_comment"` // 角色备注
		IsDefault   int64     `json:"is_default" gorm:"column:is_default"`     // 是否默认角色 0否 1是
		Status      int64     `json:"status" gorm:"column:status"`             // 状态  0正常 1禁用
		CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`     // 创建时间
		UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`     // 更新时间
	}

	// 接口实现
	defaultTRoleModel struct {
		DbEngin *gorm.DB
		table   string
	}
)

func NewTRoleModel(db *gorm.DB) TRoleModel {
	return &defaultTRoleModel{
		DbEngin: db,
		table:   "`t_role`",
	}
}

func (m *defaultTRoleModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTRoleModel) WithTransaction(tx *gorm.DB) (out TRoleModel) {
	return NewTRoleModel(tx)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTRoleModel) Insert(ctx context.Context, in *TRole) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTRoleModel) Inserts(ctx context.Context, in ...*TRole) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTRoleModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TRole{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTRoleModel) Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TRole{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTRoleModel) Update(ctx context.Context, in *TRole) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTRoleModel) Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(columns)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存记录（更新零值）
func (m *defaultTRoleModel) Save(ctx context.Context, in *TRole) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTRoleModel) FindById(ctx context.Context, id int64) (out *TRole, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTRoleModel) FindOne(ctx context.Context, conditions string, args ...interface{}) (out *TRole, err error) {
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
func (m *defaultTRoleModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TRole, err error) {
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
func (m *defaultTRoleModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
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
func (m *defaultTRoleModel) FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TRole, total int64, err error) {
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
