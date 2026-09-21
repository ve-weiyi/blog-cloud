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
		WithTx(tx *gorm.DB) (out TOperationLogModel)
		// 插入
		Insert(ctx context.Context, in *TOperationLog) (rows int64, err error)
		InsertBatch(ctx context.Context, in ...*TOperationLog) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TOperationLog) (rows int64, err error)
		UpdateFields(ctx context.Context, fields map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存或更新（更新零值）
		Save(ctx context.Context, in *TOperationLog) (rows int64, err error)
		// 查询
		FindById(ctx context.Context, id int64) (out *TOperationLog, err error)
		FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TOperationLog, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TOperationLog, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TOperationLog, total int64, err error)
		// add extra method in here
	}

	// 表字段定义
	TOperationLog struct {
		Id             int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键id" json:"id"`                        // 主键id
		UserId         string    `gorm:"column:user_id;type:varchar(64);not null;default:'';comment:用户id" json:"user_id"`                   // 用户id
		DeviceId       string    `gorm:"column:device_id;type:varchar(64);not null;default:'';comment:设备id" json:"device_id"`               // 设备id
		Module         string    `gorm:"column:module;type:varchar(50);not null;default:'';comment:操作模块" json:"module"`                     // 操作模块
		Description    string    `gorm:"column:description;type:varchar(500);not null;default:'';comment:操作描述" json:"description"`          // 操作描述
		RequestUri     string    `gorm:"column:request_uri;type:varchar(255);not null;default:'';comment:请求地址" json:"request_uri"`          // 请求地址
		RequestMethod  string    `gorm:"column:request_method;type:varchar(32);not null;default:'';comment:请求方式" json:"request_method"`     // 请求方式
		RequestData    string    `gorm:"column:request_data;type:varchar(4096);not null;default:'';comment:请求参数" json:"request_data"`       // 请求参数
		ResponseData   string    `gorm:"column:response_data;type:varchar(4096);not null;default:'';comment:返回数据" json:"response_data"`     // 返回数据
		ResponseStatus int64     `gorm:"column:response_status;type:bigint;not null;default:0;comment:响应状态码" json:"response_status"`        // 响应状态码
		Cost           string    `gorm:"column:cost;type:varchar(32);not null;default:'';comment:耗时（ms）" json:"cost"`                       // 耗时（ms）
		CreatedAt      time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
		UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
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
func (m *defaultTOperationLogModel) WithTx(tx *gorm.DB) (out TOperationLogModel) {
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
func (m *defaultTOperationLogModel) InsertBatch(ctx context.Context, in ...*TOperationLog) (rows int64, err error) {
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
func (m *defaultTOperationLogModel) DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
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
func (m *defaultTOperationLogModel) UpdateFields(ctx context.Context, feilds map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(feilds)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存或更新（更新零值）
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
func (m *defaultTOperationLogModel) FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TOperationLog, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

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
