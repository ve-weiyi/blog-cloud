package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TTalkModel = (*defaultTTalkModel)(nil)

type (
	// 接口定义
	TTalkModel interface {
		TableName() string
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TTalkModel)
		// 插入
		Insert(ctx context.Context, in *TTalk) (rows int64, err error)
		InsertBatch(ctx context.Context, in ...*TTalk) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TTalk) (rows int64, err error)
		UpdateFields(ctx context.Context, fields map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存或更新（更新零值）
		Save(ctx context.Context, in *TTalk) (rows int64, err error)
		// 查询
		FindById(ctx context.Context, id int64) (out *TTalk, err error)
		FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TTalk, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TTalk, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TTalk, total int64, err error)
		// add extra method in here
	}

	// 表字段定义
	TTalk struct {
		Id        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:说说id" json:"id"`                        // 说说id
		UserId    string    `gorm:"column:user_id;type:varchar(64);not null;default:'';comment:用户id" json:"user_id"`                   // 用户id
		Content   string    `gorm:"column:content;type:varchar(2048);not null;default:'';comment:说说内容" json:"content"`                 // 说说内容
		Images    string    `gorm:"column:images;type:varchar(2048);not null;default:'';comment:图片" json:"images"`                     // 图片
		IsTop     int64     `gorm:"column:is_top;type:tinyint;not null;default:0;comment:是否置顶" json:"is_top"`                          // 是否置顶
		Status    int64     `gorm:"column:status;type:tinyint;not null;default:1;comment:状态 1.公开 2.私密" json:"status"`                  // 状态 1.公开 2.私密
		LikeCount int64     `gorm:"column:like_count;type:bigint;not null;default:0;comment:点赞数" json:"like_count"`                    // 点赞数
		CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
		UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
	}

	// 接口实现
	defaultTTalkModel struct {
		DbEngin *gorm.DB
		table   string
	}
)

func NewTTalkModel(db *gorm.DB) TTalkModel {
	return &defaultTTalkModel{
		DbEngin: db,
		table:   "`t_talk`",
	}
}

func (m *defaultTTalkModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTTalkModel) WithTx(tx *gorm.DB) (out TTalkModel) {
	return NewTTalkModel(tx)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTTalkModel) Insert(ctx context.Context, in *TTalk) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTTalkModel) InsertBatch(ctx context.Context, in ...*TTalk) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTTalkModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TTalk{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTTalkModel) DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TTalk{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTTalkModel) Update(ctx context.Context, in *TTalk) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTTalkModel) UpdateFields(ctx context.Context, feilds map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(feilds)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存或更新（更新零值）
func (m *defaultTTalkModel) Save(ctx context.Context, in *TTalk) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTTalkModel) FindById(ctx context.Context, id int64) (out *TTalk, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTTalkModel) FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TTalk, err error) {
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
func (m *defaultTTalkModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TTalk, err error) {
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
func (m *defaultTTalkModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
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
func (m *defaultTTalkModel) FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TTalk, total int64, err error) {
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
