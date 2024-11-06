package model

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var _ TUserOauthModel = (*defaultTUserOauthModel)(nil)

type (
	// 接口定义
	TUserOauthModel interface {
		TableName() string
		// 在事务中操作
		WithTransaction(tx *gorm.DB) (out TUserOauthModel)
		// 插入
		Insert(ctx context.Context, in *TUserOauth) (rows int64, err error)
		Inserts(ctx context.Context, in ...*TUserOauth) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TUserOauth) (rows int64, err error)
		Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存
		Save(ctx context.Context, in *TUserOauth) (rows int64, err error)
		// 查询
		FindOne(ctx context.Context, id int64) (out *TUserOauth, err error)
		First(ctx context.Context, conditions string, args ...interface{}) (out *TUserOauth, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TUserOauth, err error)
		FindList(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TUserOauth, err error)
		// add extra method in here
		FindOneByOpenIdPlatform(ctx context.Context, open_id string, platform string) (out *TUserOauth, err error)
	}

	// 表字段定义
	TUserOauth struct {
		Id        int64     `json:"id" gorm:"column:id" `                 // id
		UserId    string    `json:"user_id" gorm:"column:user_id" `       // 用户id
		OpenId    string    `json:"open_id" gorm:"column:open_id" `       // 开发平台id，标识唯一用户
		Platform  string    `json:"platform" gorm:"column:platform" `     // 平台:手机号、邮箱、微信、飞书
		CreatedAt time.Time `json:"created_at" gorm:"column:created_at" ` // 创建时间
		UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at" ` // 更新时间
	}

	// 接口实现
	defaultTUserOauthModel struct {
		DbEngin    *gorm.DB
		CacheEngin *redis.Client
		table      string
	}
)

func NewTUserOauthModel(db *gorm.DB, cache *redis.Client) TUserOauthModel {
	return &defaultTUserOauthModel{
		DbEngin:    db,
		CacheEngin: cache,
		table:      "`t_user_oauth`",
	}
}

func (m *defaultTUserOauthModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTUserOauthModel) WithTransaction(tx *gorm.DB) (out TUserOauthModel) {
	return NewTUserOauthModel(tx, m.CacheEngin)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTUserOauthModel) Insert(ctx context.Context, in *TUserOauth) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTUserOauthModel) Inserts(ctx context.Context, in ...*TUserOauth) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTUserOauthModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TUserOauth{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTUserOauthModel) Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TUserOauth{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTUserOauthModel) Update(ctx context.Context, in *TUserOauth) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTUserOauthModel) Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(columns)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存记录（更新零值）
func (m *defaultTUserOauthModel) Save(ctx context.Context, in *TUserOauth) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTUserOauthModel) FindOne(ctx context.Context, id int64) (out *TUserOauth, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTUserOauthModel) First(ctx context.Context, conditions string, args ...interface{}) (out *TUserOauth, err error) {
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
func (m *defaultTUserOauthModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	err = db.Model(&TUserOauth{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 查询列表
func (m *defaultTUserOauthModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TUserOauth, err error) {
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
func (m *defaultTUserOauthModel) FindList(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TUserOauth, err error) {
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
func (m *defaultTUserOauthModel) FindOneByOpenIdPlatform(ctx context.Context, open_id string, platform string) (out *TUserOauth, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`open_id` = ? and `platform` = ?", open_id, platform).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
