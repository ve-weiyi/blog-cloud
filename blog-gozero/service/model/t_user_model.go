package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TUserModel = (*defaultTUserModel)(nil)

type (
	// 接口定义
	TUserModel interface {
		TableName() string
		// 在事务中操作
		WithTransaction(tx *gorm.DB) (out TUserModel)
		// 插入
		Insert(ctx context.Context, in *TUser) (rows int64, err error)
		Inserts(ctx context.Context, in ...*TUser) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TUser) (rows int64, err error)
		Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存
		Save(ctx context.Context, in *TUser) (rows int64, err error)
		// 查询
		FindById(ctx context.Context, id int64) (out *TUser, err error)
		FindOne(ctx context.Context, conditions string, args ...interface{}) (out *TUser, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TUser, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TUser, total int64, err error)
		// add extra method in here
		FindOneByUserId(ctx context.Context, user_id string) (out *TUser, err error)
		FindOneByUsername(ctx context.Context, username string) (out *TUser, err error)
	}

	// 表字段定义
	TUser struct {
		Id           int64     `json:"id" gorm:"column:id"`                       // id
		UserId       string    `json:"user_id" gorm:"column:user_id"`             // 用户id
		Username     string    `json:"username" gorm:"column:username"`           // 用户名
		Password     string    `json:"password" gorm:"column:password"`           // 用户密码
		Nickname     string    `json:"nickname" gorm:"column:nickname"`           // 用户昵称
		Avatar       string    `json:"avatar" gorm:"column:avatar"`               // 用户头像
		Email        string    `json:"email" gorm:"column:email"`                 // 邮箱
		Phone        string    `json:"phone" gorm:"column:phone"`                 // 手机号
		Info         string    `json:"info" gorm:"column:info"`                   // 用户信息
		Status       int64     `json:"status" gorm:"column:status"`               // 状态: -1删除 0正常 1禁用
		RegisterType string    `json:"register_type" gorm:"column:register_type"` // 注册方式
		IpAddress    string    `json:"ip_address" gorm:"column:ip_address"`       // 注册ip
		IpSource     string    `json:"ip_source" gorm:"column:ip_source"`         // 注册ip 源
		CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`       // 创建时间
		UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`       // 更新时间
	}

	// 接口实现
	defaultTUserModel struct {
		DbEngin *gorm.DB
		table   string
	}
)

func NewTUserModel(db *gorm.DB) TUserModel {
	return &defaultTUserModel{
		DbEngin: db,
		table:   "`t_user`",
	}
}

func (m *defaultTUserModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTUserModel) WithTransaction(tx *gorm.DB) (out TUserModel) {
	return NewTUserModel(tx)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTUserModel) Insert(ctx context.Context, in *TUser) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTUserModel) Inserts(ctx context.Context, in ...*TUser) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTUserModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TUser{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTUserModel) Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TUser{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTUserModel) Update(ctx context.Context, in *TUser) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTUserModel) Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(columns)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存记录（更新零值）
func (m *defaultTUserModel) Save(ctx context.Context, in *TUser) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTUserModel) FindById(ctx context.Context, id int64) (out *TUser, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTUserModel) FindOne(ctx context.Context, conditions string, args ...interface{}) (out *TUser, err error) {
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
func (m *defaultTUserModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TUser, err error) {
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
func (m *defaultTUserModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
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
func (m *defaultTUserModel) FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TUser, total int64, err error) {
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
func (m *defaultTUserModel) FindOneByUserId(ctx context.Context, user_id string) (out *TUser, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`user_id` = ?", user_id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
func (m *defaultTUserModel) FindOneByUsername(ctx context.Context, username string) (out *TUser, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`username` = ?", username).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
