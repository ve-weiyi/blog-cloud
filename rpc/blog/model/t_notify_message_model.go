package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TNotifyMessageModel = (*defaultTNotifyMessageModel)(nil)

type (
	// 接口定义
	TNotifyMessageModel interface {
		TableName() string
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TNotifyMessageModel)
		// 插入
		Insert(ctx context.Context, in *TNotifyMessage) (rows int64, err error)
		InsertBatch(ctx context.Context, in ...*TNotifyMessage) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TNotifyMessage) (rows int64, err error)
		UpdateFields(ctx context.Context, fields map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存或更新（更新零值）
		Save(ctx context.Context, in *TNotifyMessage) (rows int64, err error)
		// 查询
		FindById(ctx context.Context, id int64) (out *TNotifyMessage, err error)
		FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TNotifyMessage, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TNotifyMessage, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TNotifyMessage, total int64, err error)
		// add extra method in here
	}

	// 表字段定义
	TNotifyMessage struct {
		Id          int64      `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:ID" json:"id"`                                                                                 // ID
		Title       string     `gorm:"column:title;type:varchar(256);not null;default:'';comment:消息标题" json:"title"`                                                                             // 消息标题
		Content     *string    `gorm:"column:content;type:text;comment:消息内容" json:"content"`                                                                                                     // 消息内容
		Category    string     `gorm:"column:category;type:varchar(50);not null;index:idx_category,priority:1;default:system;comment:分类：system | maintenance | update | remind" json:"category"` // 分类：system | maintenance | update | remind
		Level       string     `gorm:"column:level;type:varchar(20);not null;default:info;comment:等级：info | warning | error" json:"level"`                                                       // 等级：info | warning | error
		TargetType  string     `gorm:"column:target_type;type:varchar(20);not null;default:all;comment:目标类型：all | user_ids" json:"target_type"`                                                  // 目标类型：all | user_ids
		TargetIds   *string    `gorm:"column:target_ids;type:text;comment:目标用户ID列表，逗号分隔" json:"target_ids"`                                                                                      // 目标用户ID列表，逗号分隔
		Status      string     `gorm:"column:status;type:varchar(20);not null;index:idx_status,priority:1;default:draft;comment:状态：draft | published | revoked" json:"status"`                   // 状态：draft | published | revoked
		PublishedAt *time.Time `gorm:"column:published_at;type:datetime;index:idx_published_at,priority:1;comment:发布时间" json:"published_at"`                                                     // 发布时间
		PublishedBy string     `gorm:"column:published_by;type:varchar(64);not null;default:'';comment:发布人" json:"published_by"`                                                                 // 发布人
		CreatedAt   time.Time  `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                                                        // 创建时间
		UpdatedAt   time.Time  `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                                                        // 更新时间
	}

	// 接口实现
	defaultTNotifyMessageModel struct {
		DbEngin *gorm.DB
		table   string
	}
)

func NewTNotifyMessageModel(db *gorm.DB) TNotifyMessageModel {
	return &defaultTNotifyMessageModel{
		DbEngin: db,
		table:   "`t_notify_message`",
	}
}

func (m *defaultTNotifyMessageModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTNotifyMessageModel) WithTx(tx *gorm.DB) (out TNotifyMessageModel) {
	return NewTNotifyMessageModel(tx)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTNotifyMessageModel) Insert(ctx context.Context, in *TNotifyMessage) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTNotifyMessageModel) InsertBatch(ctx context.Context, in ...*TNotifyMessage) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTNotifyMessageModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TNotifyMessage{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTNotifyMessageModel) DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TNotifyMessage{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTNotifyMessageModel) Update(ctx context.Context, in *TNotifyMessage) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTNotifyMessageModel) UpdateFields(ctx context.Context, feilds map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(feilds)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存或更新（更新零值）
func (m *defaultTNotifyMessageModel) Save(ctx context.Context, in *TNotifyMessage) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTNotifyMessageModel) FindById(ctx context.Context, id int64) (out *TNotifyMessage, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTNotifyMessageModel) FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TNotifyMessage, err error) {
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
func (m *defaultTNotifyMessageModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TNotifyMessage, err error) {
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
func (m *defaultTNotifyMessageModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
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
func (m *defaultTNotifyMessageModel) FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TNotifyMessage, total int64, err error) {
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
