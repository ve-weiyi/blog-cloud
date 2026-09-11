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
		WithTx(tx *gorm.DB) (out TArticleModel)
		// 插入
		Insert(ctx context.Context, in *TArticle) (rows int64, err error)
		InsertBatch(ctx context.Context, in ...*TArticle) (rows int64, err error)
		// 删除
		Delete(ctx context.Context, id int64) (rows int64, err error)
		DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
		// 更新
		Update(ctx context.Context, in *TArticle) (rows int64, err error)
		UpdateFields(ctx context.Context, fields map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
		// 保存或更新（更新零值）
		Save(ctx context.Context, in *TArticle) (rows int64, err error)
		// 查询
		FindById(ctx context.Context, id int64) (out *TArticle, err error)
		FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TArticle, err error)
		FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*TArticle, err error)
		FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
		FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*TArticle, total int64, err error)
		// add extra method in here
	}

	// 表字段定义
	TArticle struct {
		Id             int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:id" json:"id"`                          // id
		UserId         string    `gorm:"column:user_id;type:varchar(64);not null;default:'';comment:作者" json:"user_id"`                     // 作者
		CategoryId     int64     `gorm:"column:category_id;type:bigint;not null;default:0;comment:文章分类" json:"category_id"`                 // 文章分类
		ArticleCover   string    `gorm:"column:article_cover;type:varchar(1024);not null;default:'';comment:文章缩略图" json:"article_cover"`    // 文章缩略图
		ArticleTitle   string    `gorm:"column:article_title;type:varchar(64);not null;default:'';comment:标题" json:"article_title"`         // 标题
		ArticleContent string    `gorm:"column:article_content;type:longtext;not null;default:'';comment:内容" json:"article_content"`        // 内容
		ArticleType    int64     `gorm:"column:article_type;type:tinyint;not null;default:0;comment:文章类型 1原创 2转载 3翻译" json:"article_type"`  // 文章类型 1原创 2转载 3翻译
		OriginalUrl    string    `gorm:"column:original_url;type:varchar(255);not null;default:'';comment:原文链接" json:"original_url"`        // 原文链接
		IsTop          int64     `gorm:"column:is_top;type:tinyint(1);not null;default:0;comment:是否置顶 0否 1是" json:"is_top"`                 // 是否置顶 0否 1是
		IsDelete       int64     `gorm:"column:is_delete;type:tinyint(1);not null;default:0;comment:是否删除  0否 1是" json:"is_delete"`          // 是否删除  0否 1是
		Status         int64     `gorm:"column:status;type:tinyint;not null;default:1;comment:状态值 1公开 2私密 3草稿 4评论可见" json:"status"`         // 状态值 1公开 2私密 3草稿 4评论可见
		LikeCount      int64     `gorm:"column:like_count;type:bigint;not null;default:0;comment:点赞数" json:"like_count"`                    // 点赞数
		ViewCount      int64     `gorm:"column:view_count;type:bigint;not null;default:0;comment:查看数" json:"view_count"`                    // 查看数
		CreatedAt      time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:发表时间" json:"created_at"` // 发表时间
		UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
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
func (m *defaultTArticleModel) WithTx(tx *gorm.DB) (out TArticleModel) {
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
func (m *defaultTArticleModel) InsertBatch(ctx context.Context, in ...*TArticle) (rows int64, err error) {
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
func (m *defaultTArticleModel) DeleteBatch(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
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
func (m *defaultTArticleModel) UpdateFields(ctx context.Context, feilds map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error) {
	db := m.DbEngin.WithContext(ctx).Table(m.table)

	result := db.Where(conditions, args...).Updates(feilds)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, err
}

// 保存或更新（更新零值）
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
func (m *defaultTArticleModel) FindOne(ctx context.Context, sorts string, conditions string, args ...interface{}) (out *TArticle, err error) {
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
