package model

import (
	"context"

	"gorm.io/gorm"
)

var _ TArticleTagModel = (*defaultTArticleTagModel)(nil)

type (
	// 接口定义
	TArticleTagModel interface {
		TableName() string
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TArticleTagModel)
		// 插入
		Insert(ctx context.Context, in *TArticleTag) (rows int64, err error)
		InsertBatch(ctx context.Context, in ...*TArticleTag) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TArticleTag) (rows int64, err error)
		UpdateFields(ctx context.Context, fields map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存或更新（更新零值）
		Save(ctx context.Context, in *TArticleTag) (rows int64, err error)
		// 查询
		FindById(ctx context.Context, id int64) (out *TArticleTag, err error)
		FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TArticleTag, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TArticleTag, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TArticleTag, total int64, err error)
		// add extra method in here
		FindOneByArticleIdTagId(ctx context.Context, article_id int64, tag_id int64) (out *TArticleTag, err error)
	}

	// 表字段定义
	TArticleTag struct {
		Id        int64 `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:id" json:"id"`                                              // id
		ArticleId int64 `gorm:"column:article_id;type:bigint;not null;uniqueIndex:uk_article_tag,priority:1;default:0;comment:文章id" json:"article_id"` // 文章id
		TagId     int64 `gorm:"column:tag_id;type:bigint;not null;uniqueIndex:uk_article_tag,priority:2;default:0;comment:标签id" json:"tag_id"`         // 标签id
	}

	// 接口实现
	defaultTArticleTagModel struct {
		DbEngin *gorm.DB
		table   string
	}
)

func NewTArticleTagModel(db *gorm.DB) TArticleTagModel {
	return &defaultTArticleTagModel{
		DbEngin: db,
		table:   "`t_article_tag`",
	}
}

func (m *defaultTArticleTagModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTArticleTagModel) WithTx(tx *gorm.DB) (out TArticleTagModel) {
	return NewTArticleTagModel(tx)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTArticleTagModel) Insert(ctx context.Context, in *TArticleTag) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTArticleTagModel) InsertBatch(ctx context.Context, in ...*TArticleTag) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTArticleTagModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TArticleTag{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTArticleTagModel) DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TArticleTag{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTArticleTagModel) Update(ctx context.Context, in *TArticleTag) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTArticleTagModel) UpdateFields(ctx context.Context, feilds map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(feilds)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存或更新（更新零值）
func (m *defaultTArticleTagModel) Save(ctx context.Context, in *TArticleTag) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTArticleTagModel) FindById(ctx context.Context, id int64) (out *TArticleTag, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTArticleTagModel) FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TArticleTag, err error) {
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
func (m *defaultTArticleTagModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TArticleTag, err error) {
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
func (m *defaultTArticleTagModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
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
func (m *defaultTArticleTagModel) FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TArticleTag, total int64, err error) {
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
func (m *defaultTArticleTagModel) FindOneByArticleIdTagId(ctx context.Context, article_id int64, tag_id int64) (out *TArticleTag, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`article_id` = ? and `tag_id` = ?", article_id, tag_id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}
