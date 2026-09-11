package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TLoginLogModel = (*defaultTLoginLogModel)(nil)

type (
	// 接口定义
	TLoginLogModel interface {
		TableName() string
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TLoginLogModel)
		// 插入
		Insert(ctx context.Context, in *TLoginLog) (rows int64, err error)
		InsertBatch(ctx context.Context, in ...*TLoginLog) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TLoginLog) (rows int64, err error)
		UpdateFields(ctx context.Context, fields map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存或更新（更新零值）
		Save(ctx context.Context, in *TLoginLog) (rows int64, err error)
		// 查询
		FindById(ctx context.Context, id int64) (out *TLoginLog, err error)
		FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TLoginLog, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TLoginLog, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TLoginLog, total int64, err error)
		// add extra method in here
	}

	// 表字段定义
	TLoginLog struct {
		Id         int64      `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:日志ID" json:"id"`                                                        // 日志ID
		UserId     string     `gorm:"column:user_id;type:varchar(64);not null;index:idx_user_id,priority:1;default:'';comment:用户id" json:"user_id"`                      // 用户id
		DeviceId   string     `gorm:"column:device_id;type:varchar(64);not null;default:'';comment:设备id" json:"device_id"`                                               // 设备id
		LoginType  string     `gorm:"column:login_type;type:varchar(20);not null;default:'';comment:登录类型：password-密码登录 sms-短信登录" json:"login_type"`                      // 登录类型：password-密码登录 sms-短信登录
		Status     int64      `gorm:"column:status;type:bigint;not null;default:1;comment:登录状态：0-失败 1-成功" json:"status"`                                                 // 登录状态：0-失败 1-成功
		FailReason string     `gorm:"column:fail_reason;type:varchar(500);not null;default:'';comment:失败原因" json:"fail_reason"`                                          // 失败原因
		LogoutAt   *time.Time `gorm:"column:logout_at;type:datetime;comment:登出时间" json:"logout_at"`                                                                      // 登出时间
		CreatedAt  time.Time  `gorm:"column:created_at;type:datetime;not null;index:idx_created_at,priority:1;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
		UpdatedAt  time.Time  `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                                 // 更新时间
	}

	// 接口实现
	defaultTLoginLogModel struct {
		DbEngin *gorm.DB
		table   string
	}
)

func NewTLoginLogModel(db *gorm.DB) TLoginLogModel {
	return &defaultTLoginLogModel{
		DbEngin: db,
		table:   "`t_login_log`",
	}
}

func (m *defaultTLoginLogModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTLoginLogModel) WithTx(tx *gorm.DB) (out TLoginLogModel) {
	return NewTLoginLogModel(tx)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTLoginLogModel) Insert(ctx context.Context, in *TLoginLog) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTLoginLogModel) InsertBatch(ctx context.Context, in ...*TLoginLog) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTLoginLogModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TLoginLog{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTLoginLogModel) DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TLoginLog{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTLoginLogModel) Update(ctx context.Context, in *TLoginLog) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTLoginLogModel) UpdateFields(ctx context.Context, feilds map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(feilds)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存或更新（更新零值）
func (m *defaultTLoginLogModel) Save(ctx context.Context, in *TLoginLog) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTLoginLogModel) FindById(ctx context.Context, id int64) (out *TLoginLog, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTLoginLogModel) FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TLoginLog, err error) {
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
func (m *defaultTLoginLogModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TLoginLog, err error) {
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
func (m *defaultTLoginLogModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
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
func (m *defaultTLoginLogModel) FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TLoginLog, total int64, err error) {
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
