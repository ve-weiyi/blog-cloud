package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TGuestModel = (*defaultTGuestModel)(nil)

type (
	// 接口定义
	TGuestModel interface {
		TableName() string
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TGuestModel)
		// 插入
		Insert(ctx context.Context, in *TGuest) (rows int64, err error)
		InsertBatch(ctx context.Context, in ...*TGuest) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TGuest) (rows int64, err error)
		UpdateFields(ctx context.Context, fields map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存或更新（更新零值）
		Save(ctx context.Context, in *TGuest) (rows int64, err error)
		// 查询
		FindById(ctx context.Context, id int64) (out *TGuest, err error)
		FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TGuest, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TGuest, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TGuest, total int64, err error)
		// add extra method in here
		FindOneByDeviceId(ctx context.Context, device_id string) (out *TGuest, err error)
	}

	// 表字段定义
	TGuest struct {
		Id        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:id" json:"id"`                                          // id
		DeviceId  string    `gorm:"column:device_id;type:varchar(64);not null;uniqueIndex:uk_did,priority:1;default:'';comment:设备id" json:"device_id"` // 设备id
		Os        string    `gorm:"column:os;type:varchar(50);not null;default:'';comment:操作系统" json:"os"`                                             // 操作系统
		Browser   string    `gorm:"column:browser;type:varchar(50);not null;default:'';comment:浏览器" json:"browser"`                                    // 浏览器
		IpAddress string    `gorm:"column:ip_address;type:varchar(255);not null;default:'';comment:操作ip" json:"ip_address"`                            // 操作ip
		IpSource  string    `gorm:"column:ip_source;type:varchar(255);not null;default:'';comment:操作地址" json:"ip_source"`                              // 操作地址
		CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                 // 创建时间
		UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                 // 更新时间
	}

	// 接口实现
	defaultTGuestModel struct {
		DbEngin *gorm.DB
		table   string
	}
)

func NewTGuestModel(db *gorm.DB) TGuestModel {
	return &defaultTGuestModel{
		DbEngin: db,
		table:   "`t_guest`",
	}
}

func (m *defaultTGuestModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTGuestModel) WithTx(tx *gorm.DB) (out TGuestModel) {
	return NewTGuestModel(tx)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTGuestModel) Insert(ctx context.Context, in *TGuest) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTGuestModel) InsertBatch(ctx context.Context, in ...*TGuest) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTGuestModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TGuest{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTGuestModel) DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TGuest{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTGuestModel) Update(ctx context.Context, in *TGuest) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTGuestModel) UpdateFields(ctx context.Context, feilds map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(feilds)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存或更新（更新零值）
func (m *defaultTGuestModel) Save(ctx context.Context, in *TGuest) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTGuestModel) FindById(ctx context.Context, id int64) (out *TGuest, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTGuestModel) FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TGuest, err error) {
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
func (m *defaultTGuestModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TGuest, err error) {
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
func (m *defaultTGuestModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
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
func (m *defaultTGuestModel) FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TGuest, total int64, err error) {
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
func (m *defaultTGuestModel) FindOneByDeviceId(ctx context.Context, device_id string) (out *TGuest, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`device_id` = ?", device_id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
