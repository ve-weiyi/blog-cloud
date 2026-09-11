package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TAlbumModel = (*defaultTAlbumModel)(nil)

type (
	// 接口定义
	TAlbumModel interface {
		TableName() string
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TAlbumModel)
		// 插入
		Insert(ctx context.Context, in *TAlbum) (rows int64, err error)
		InsertBatch(ctx context.Context, in ...*TAlbum) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TAlbum) (rows int64, err error)
		UpdateFields(ctx context.Context, fields map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存或更新（更新零值）
		Save(ctx context.Context, in *TAlbum) (rows int64, err error)
		// 查询
		FindById(ctx context.Context, id int64) (out *TAlbum, err error)
		FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TAlbum, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TAlbum, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TAlbum, total int64, err error)
		// add extra method in here
	}

	// 表字段定义
	TAlbum struct {
		Id         int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键" json:"id"`                          // 主键
		AlbumName  string    `gorm:"column:album_name;type:varchar(64);not null;default:'';comment:相册名" json:"album_name"`              // 相册名
		AlbumDesc  string    `gorm:"column:album_desc;type:varchar(128);not null;default:'';comment:相册描述" json:"album_desc"`            // 相册描述
		AlbumCover string    `gorm:"column:album_cover;type:varchar(255);not null;default:'';comment:相册封面" json:"album_cover"`          // 相册封面
		IsDelete   int64     `gorm:"column:is_delete;type:tinyint(1);not null;default:0;comment:是否删除" json:"is_delete"`                 // 是否删除
		Status     int64     `gorm:"column:status;type:tinyint;not null;default:1;comment:状态值 1公开 2私密" json:"status"`                   // 状态值 1公开 2私密
		CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
		UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
	}

	// 接口实现
	defaultTAlbumModel struct {
		DbEngin *gorm.DB
		table   string
	}
)

func NewTAlbumModel(db *gorm.DB) TAlbumModel {
	return &defaultTAlbumModel{
		DbEngin: db,
		table:   "`t_album`",
	}
}

func (m *defaultTAlbumModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTAlbumModel) WithTx(tx *gorm.DB) (out TAlbumModel) {
	return NewTAlbumModel(tx)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTAlbumModel) Insert(ctx context.Context, in *TAlbum) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTAlbumModel) InsertBatch(ctx context.Context, in ...*TAlbum) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTAlbumModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TAlbum{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTAlbumModel) DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TAlbum{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTAlbumModel) Update(ctx context.Context, in *TAlbum) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTAlbumModel) UpdateFields(ctx context.Context, feilds map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(feilds)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存或更新（更新零值）
func (m *defaultTAlbumModel) Save(ctx context.Context, in *TAlbum) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTAlbumModel) FindById(ctx context.Context, id int64) (out *TAlbum, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTAlbumModel) FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TAlbum, err error) {
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
func (m *defaultTAlbumModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TAlbum, err error) {
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
func (m *defaultTAlbumModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
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
func (m *defaultTAlbumModel) FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TAlbum, total int64, err error) {
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
