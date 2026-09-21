package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TNotifyRecordModel = (*defaultTNotifyRecordModel)(nil)

type (
	// 接口定义
	TNotifyRecordModel interface {
		TableName() string
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TNotifyRecordModel)
		// 插入
		Insert(ctx context.Context, in *TNotifyRecord) (rows int64, err error)
		InsertBatch(ctx context.Context, in ...*TNotifyRecord) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TNotifyRecord) (rows int64, err error)
		UpdateFields(ctx context.Context, fields map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存或更新（更新零值）
		Save(ctx context.Context, in *TNotifyRecord) (rows int64, err error)
		// 查询
		FindById(ctx context.Context, id int64) (out *TNotifyRecord, err error)
		FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TNotifyRecord, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TNotifyRecord, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TNotifyRecord, total int64, err error)
		// add extra method in here
	}

	// 表字段定义
	TNotifyRecord struct {
		Id           int64      `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:ID" json:"id"`                                                                                        // ID
		MessageId    int64      `gorm:"column:message_id;type:bigint;not null;index:idx_message_id,priority:1;default:0;comment:关联消息ID" json:"message_id"`                                               // 关联消息ID
		Channel      string     `gorm:"column:channel;type:varchar(16);not null;index:idx_channel,priority:1;default:inbox;comment:渠道：inbox | sms | email" json:"channel"`                               // 渠道：inbox | sms | email
		Recipient    string     `gorm:"column:recipient;type:varchar(128);not null;index:idx_recipient,priority:1;default:'';comment:接收者（user_id | mobile | email）" json:"recipient"`                    // 接收者（user_id | mobile | email）
		TemplateCode string     `gorm:"column:template_code;type:varchar(64);not null;default:'';comment:模板编码" json:"template_code"`                                                                     // 模板编码
		Content      *string    `gorm:"column:content;type:text;comment:实际发送内容" json:"content"`                                                                                                          // 实际发送内容
		Status       string     `gorm:"column:status;type:varchar(16);not null;index:idx_status,priority:1;default:unread;comment:状态：inbox: unread|read / sms/email: pending|sent|failed" json:"status"` // 状态：inbox: unread|read / sms/email: pending|sent|failed
		BizId        string     `gorm:"column:biz_id;type:varchar(64);not null;index:idx_biz_id,priority:1;default:'';comment:业务幂等键" json:"biz_id"`                                                      // 业务幂等键
		ErrorMsg     string     `gorm:"column:error_msg;type:varchar(512);not null;default:'';comment:失败原因" json:"error_msg"`                                                                            // 失败原因
		ReadAt       *time.Time `gorm:"column:read_at;type:datetime;comment:阅读时间（inbox）" json:"read_at"`                                                                                                 // 阅读时间（inbox）
		SentAt       *time.Time `gorm:"column:sent_at;type:datetime;comment:发送时间（sms/email）" json:"sent_at"`                                                                                             // 发送时间（sms/email）
		CreatedAt    time.Time  `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                                                               // 创建时间
	}

	// 接口实现
	defaultTNotifyRecordModel struct {
		DbEngin *gorm.DB
		table   string
	}
)

func NewTNotifyRecordModel(db *gorm.DB) TNotifyRecordModel {
	return &defaultTNotifyRecordModel{
		DbEngin: db,
		table:   "`t_notify_record`",
	}
}

func (m *defaultTNotifyRecordModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTNotifyRecordModel) WithTx(tx *gorm.DB) (out TNotifyRecordModel) {
	return NewTNotifyRecordModel(tx)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTNotifyRecordModel) Insert(ctx context.Context, in *TNotifyRecord) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTNotifyRecordModel) InsertBatch(ctx context.Context, in ...*TNotifyRecord) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTNotifyRecordModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TNotifyRecord{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTNotifyRecordModel) DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TNotifyRecord{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTNotifyRecordModel) Update(ctx context.Context, in *TNotifyRecord) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTNotifyRecordModel) UpdateFields(ctx context.Context, feilds map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(feilds)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存或更新（更新零值）
func (m *defaultTNotifyRecordModel) Save(ctx context.Context, in *TNotifyRecord) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTNotifyRecordModel) FindById(ctx context.Context, id int64) (out *TNotifyRecord, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTNotifyRecordModel) FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TNotifyRecord, err error) {
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
func (m *defaultTNotifyRecordModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TNotifyRecord, err error) {
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
func (m *defaultTNotifyRecordModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
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
func (m *defaultTNotifyRecordModel) FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TNotifyRecord, total int64, err error) {
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
