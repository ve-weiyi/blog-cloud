package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TOperationLogModel = (*defaultTOperationLogModel)(nil)

type (
	// 接口定义
	TOperationLogModel interface {
		TableName() string
		// 在事务中操作
		WithTransaction(tx *gorm.DB) (out TOperationLogModel)
		// 插入
		Insert(ctx context.Context, in *TOperationLog) (rows int64, err error)
		Inserts(ctx context.Context, in ...*TOperationLog) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TOperationLog) (rows int64, err error)
		Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存
		Save(ctx context.Context, in *TOperationLog) (rows int64, err error)
		// 查询
		FindById(ctx context.Context, id int64) (out *TOperationLog, err error)
		FindOne(ctx context.Context, conditions string, args ...interface{}) (out *TOperationLog, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TOperationLog, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TOperationLog, total int64, err error)
		// add extra method in here
	}

	// 表字段定义
	TOperationLog struct {
		Id             int64     `json:"id" gorm:"column:id"`                           // 主键id
		UserId         string    `json:"user_id" gorm:"column:user_id"`                 // 用户id
		TerminalId     string    `json:"terminal_id" gorm:"column:terminal_id"`         // 设备id
		IpAddress      string    `json:"ip_address" gorm:"column:ip_address"`           // 操作ip
		IpSource       string    `json:"ip_source" gorm:"column:ip_source"`             // 操作地址
		OptModule      string    `json:"opt_module" gorm:"column:opt_module"`           // 操作模块
		OptDesc        string    `json:"opt_desc" gorm:"column:opt_desc"`               // 操作描述
		RequestUri     string    `json:"request_uri" gorm:"column:request_uri"`         // 请求地址
		RequestMethod  string    `json:"request_method" gorm:"column:request_method"`   // 请求方式
		RequestData    string    `json:"request_data" gorm:"column:request_data"`       // 请求参数
		ResponseData   string    `json:"response_data" gorm:"column:response_data"`     // 返回数据
		ResponseStatus int64     `json:"response_status" gorm:"column:response_status"` // 响应状态码
		Cost           string    `json:"cost" gorm:"column:cost"`                       // 耗时（ms）
		CreatedAt      time.Time `json:"created_at" gorm:"column:created_at"`           // 创建时间
		UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at"`           // 更新时间
	}

	// 接口实现
	defaultTOperationLogModel struct {
		DbEngin *gorm.DB
		table   string
	}
)

func NewTOperationLogModel(db *gorm.DB) TOperationLogModel {
	return &defaultTOperationLogModel{
		DbEngin: db,
		table:   "`t_operation_log`",
	}
}

func (m *defaultTOperationLogModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTOperationLogModel) WithTransaction(tx *gorm.DB) (out TOperationLogModel) {
	return NewTOperationLogModel(tx)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTOperationLogModel) Insert(ctx context.Context, in *TOperationLog) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTOperationLogModel) Inserts(ctx context.Context, in ...*TOperationLog) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTOperationLogModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TOperationLog{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTOperationLogModel) Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TOperationLog{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTOperationLogModel) Update(ctx context.Context, in *TOperationLog) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTOperationLogModel) Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(columns)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存记录（更新零值）
func (m *defaultTOperationLogModel) Save(ctx context.Context, in *TOperationLog) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTOperationLogModel) FindById(ctx context.Context, id int64) (out *TOperationLog, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTOperationLogModel) FindOne(ctx context.Context, conditions string, args ...interface{}) (out *TOperationLog, err error) {
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
func (m *defaultTOperationLogModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TOperationLog, err error) {
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
func (m *defaultTOperationLogModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
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
func (m *defaultTOperationLogModel) FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TOperationLog, total int64, err error) {
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
