package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TFriendModel = (*defaultTFriendModel)(nil)

type (
	// 接口定义
	TFriendModel interface {
		TableName() string
		// 在事务中操作
		WithTransaction(tx *gorm.DB) (out TFriendModel)
		// 插入
		Insert(ctx context.Context, in *TFriend) (rows int64, err error)
		Inserts(ctx context.Context, in ...*TFriend) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TFriend) (rows int64, err error)
		Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存
		Save(ctx context.Context, in *TFriend) (rows int64, err error)
		// 查询
		FindById(ctx context.Context, id int64) (out *TFriend, err error)
		FindOne(ctx context.Context, conditions string, args ...interface{}) (out *TFriend, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TFriend, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TFriend, total int64, err error)
		// add extra method in here
		FindOneByLinkName(ctx context.Context, link_name string) (out *TFriend, err error)
	}

	// 表字段定义
	TFriend struct {
		Id          int64     `json:"id" gorm:"column:id"`                     // id
		LinkName    string    `json:"link_name" gorm:"column:link_name"`       // 链接名
		LinkAvatar  string    `json:"link_avatar" gorm:"column:link_avatar"`   // 链接头像
		LinkAddress string    `json:"link_address" gorm:"column:link_address"` // 链接地址
		LinkIntro   string    `json:"link_intro" gorm:"column:link_intro"`     // 链接介绍
		CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`     // 创建时间
		UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`     // 更新时间
	}

	// 接口实现
	defaultTFriendModel struct {
		DbEngin *gorm.DB
		table   string
	}
)

func NewTFriendModel(db *gorm.DB) TFriendModel {
	return &defaultTFriendModel{
		DbEngin: db,
		table:   "`t_friend`",
	}
}

func (m *defaultTFriendModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTFriendModel) WithTransaction(tx *gorm.DB) (out TFriendModel) {
	return NewTFriendModel(tx)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTFriendModel) Insert(ctx context.Context, in *TFriend) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTFriendModel) Inserts(ctx context.Context, in ...*TFriend) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTFriendModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TFriend{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTFriendModel) Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TFriend{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTFriendModel) Update(ctx context.Context, in *TFriend) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTFriendModel) Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(columns)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存记录（更新零值）
func (m *defaultTFriendModel) Save(ctx context.Context, in *TFriend) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTFriendModel) FindById(ctx context.Context, id int64) (out *TFriend, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTFriendModel) FindOne(ctx context.Context, conditions string, args ...interface{}) (out *TFriend, err error) {
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
func (m *defaultTFriendModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TFriend, err error) {
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
func (m *defaultTFriendModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
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
func (m *defaultTFriendModel) FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TFriend, total int64, err error) {
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
func (m *defaultTFriendModel) FindOneByLinkName(ctx context.Context, link_name string) (out *TFriend, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`link_name` = ?", link_name).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
