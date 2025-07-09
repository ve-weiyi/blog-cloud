package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TRemarkModel = (*defaultTRemarkModel)(nil)

type (
	// 接口定义
	TRemarkModel interface {
		TableName() string
		// 在事务中操作
		WithTransaction(tx *gorm.DB) (out TRemarkModel)
		// 插入
		Insert(ctx context.Context, in *TRemark) (rows int64, err error)
		Inserts(ctx context.Context, in ...*TRemark) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TRemark) (rows int64, err error)
		Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存
		Save(ctx context.Context, in *TRemark) (rows int64, err error)
		// 查询
		FindById(ctx context.Context, id int64) (out *TRemark, err error)
		FindOne(ctx context.Context, conditions string, args ...interface{}) (out *TRemark, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TRemark, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TRemark, total int64, err error)
		// add extra method in here
	}

	// 表字段定义
	TRemark struct {
		Id             int64     `json:"id" gorm:"column:id"`                           // 主键id
		UserId         string    `json:"user_id" gorm:"column:user_id"`                 // 用户id
		TerminalId     string    `json:"terminal_id" gorm:"column:terminal_id"`         // 终端id
		MessageContent string    `json:"message_content" gorm:"column:message_content"` // 留言内容
		IpAddress      string    `json:"ip_address" gorm:"column:ip_address"`           // 用户ip 127.0.0.1
		IpSource       string    `json:"ip_source" gorm:"column:ip_source"`             // 用户地址 广东省深圳市
		Status         int64     `json:"status" gorm:"column:status"`                   // 状态:0正常 1编辑 2撤回 3删除
		IsReview       int64     `json:"is_review" gorm:"column:is_review"`             // 是否审核通过
		CreatedAt      time.Time `json:"created_at" gorm:"column:created_at"`           // 发布时间
		UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at"`           // 更新时间
	}

	// 接口实现
	defaultTRemarkModel struct {
		DbEngin *gorm.DB
		table   string
	}
)

func NewTRemarkModel(db *gorm.DB) TRemarkModel {
	return &defaultTRemarkModel{
		DbEngin: db,
		table:   "`t_remark`",
	}
}

func (m *defaultTRemarkModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTRemarkModel) WithTransaction(tx *gorm.DB) (out TRemarkModel) {
	return NewTRemarkModel(tx)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTRemarkModel) Insert(ctx context.Context, in *TRemark) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTRemarkModel) Inserts(ctx context.Context, in ...*TRemark) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTRemarkModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TRemark{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTRemarkModel) Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TRemark{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTRemarkModel) Update(ctx context.Context, in *TRemark) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTRemarkModel) Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(columns)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存记录（更新零值）
func (m *defaultTRemarkModel) Save(ctx context.Context, in *TRemark) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTRemarkModel) FindById(ctx context.Context, id int64) (out *TRemark, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTRemarkModel) FindOne(ctx context.Context, conditions string, args ...interface{}) (out *TRemark, err error) {
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
func (m *defaultTRemarkModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TRemark, err error) {
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
func (m *defaultTRemarkModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
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
func (m *defaultTRemarkModel) FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TRemark, total int64, err error) {
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
