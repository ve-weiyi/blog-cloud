package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TArticleModel = (*defaultTArticleModel)(nil)

type (
	// 接口定义
	TArticleModel interface {
		TableName() string
		// 在事务中操作
		WithTransaction(tx *gorm.DB) (out TArticleModel)
		// 插入
		Insert(ctx context.Context, in *TArticle) (rows int64, err error)
		Inserts(ctx context.Context, in ...*TArticle) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TArticle) (rows int64, err error)
		Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存
		Save(ctx context.Context, in *TArticle) (rows int64, err error)
		// 查询
		FindById(ctx context.Context, id int64) (out *TArticle, err error)
		FindOne(ctx context.Context, conditions string, args ...interface{}) (out *TArticle, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TArticle, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TArticle, total int64, err error)
		// add extra method in here
	}

	// 表字段定义
	TArticle struct {
		Id             int64     `json:"id" gorm:"column:id"`                           // id
		UserId         string    `json:"user_id" gorm:"column:user_id"`                 // 作者
		CategoryId     int64     `json:"category_id" gorm:"column:category_id"`         // 文章分类
		ArticleCover   string    `json:"article_cover" gorm:"column:article_cover"`     // 文章缩略图
		ArticleTitle   string    `json:"article_title" gorm:"column:article_title"`     // 标题
		ArticleContent string    `json:"article_content" gorm:"column:article_content"` // 内容
		ArticleType    int64     `json:"article_type" gorm:"column:article_type"`       // 文章类型 1原创 2转载 3翻译
		OriginalUrl    string    `json:"original_url" gorm:"column:original_url"`       // 原文链接
		IsTop          int64     `json:"is_top" gorm:"column:is_top"`                   // 是否置顶 0否 1是
		IsDelete       int64     `json:"is_delete" gorm:"column:is_delete"`             // 是否删除  0否 1是
		Status         int64     `json:"status" gorm:"column:status"`                   // 状态值 1公开 2私密 3草稿 4评论可见
		LikeCount      int64     `json:"like_count" gorm:"column:like_count"`           // 点赞数
		ViewCount      int64     `json:"view_count" gorm:"column:view_count"`           // 查看数
		CreatedAt      time.Time `json:"created_at" gorm:"column:created_at"`           // 发表时间
		UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at"`           // 更新时间
	}

	// 接口实现
	defaultTArticleModel struct {
		DbEngin *gorm.DB
		table   string
	}
)

func NewTArticleModel(db *gorm.DB) TArticleModel {
	return &defaultTArticleModel{
		DbEngin: db,
		table:   "`t_article`",
	}
}

func (m *defaultTArticleModel) TableName() string {
	return m.table
}

// 在事务中操作
func (m *defaultTArticleModel) WithTransaction(tx *gorm.DB) (out TArticleModel) {
	return NewTArticleModel(tx)
}

// 插入记录 (返回的是受影响行数，如需获取自增id，请通过data参数获取)
func (m *defaultTArticleModel) Insert(ctx context.Context, in *TArticle) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Create(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 插入记录（批量操作）
func (m *defaultTArticleModel) Inserts(ctx context.Context, in ...*TArticle) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.CreateInBatches(&in, len(in))
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录
func (m *defaultTArticleModel) Delete(ctx context.Context, id int64) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	db = db.Where("id = ?", id)

	result := db.Delete(&TArticle{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 删除记录（批量操作）
func (m *defaultTArticleModel) Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	result := db.Delete(&TArticle{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（不更新零值）
func (m *defaultTArticleModel) Update(ctx context.Context, in *TArticle) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Updates(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 更新记录（批量操作）
func (m *defaultTArticleModel) Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(columns)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存记录（更新零值）
func (m *defaultTArticleModel) Save(ctx context.Context, in *TArticle) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Omit("created_at").Save(&in)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 查询记录
func (m *defaultTArticleModel) FindById(ctx context.Context, id int64) (out *TArticle, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	err = db.Where("`id` = ?", id).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, err
}

// 查询记录
func (m *defaultTArticleModel) FindOne(ctx context.Context, conditions string, args ...interface{}) (out *TArticle, err error) {
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
func (m *defaultTArticleModel) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*TArticle, err error) {
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
func (m *defaultTArticleModel) FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
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
func (m *defaultTArticleModel) FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TArticle, total int64, err error) {
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
