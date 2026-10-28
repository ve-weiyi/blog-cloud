package model

import (
	"time"

	"gorm.io/gorm"
)

var _ TArticleModel = (*defaultTArticleModel)(nil)

type (
	// 接口定义
	TArticleModel interface {
		Table[TArticle]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TArticleModel)
	}

	// 表字段定义
	TArticle struct {
		Id             int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:id" json:"id"`                                                         // id
		UserId         string    `gorm:"column:user_id;type:varchar(64);not null;default:'';comment:作者" json:"user_id"`                                                    // 作者
		CategoryId     int64     `gorm:"column:category_id;type:bigint;not null;index:idx_category_delete_status_id,priority:1;default:0;comment:文章分类" json:"category_id"` // 文章分类
		ArticleCover   string    `gorm:"column:article_cover;type:varchar(1024);not null;default:'';comment:文章缩略图" json:"article_cover"`                                   // 文章缩略图
		ArticleTitle   string    `gorm:"column:article_title;type:varchar(64);not null;default:'';comment:标题" json:"article_title"`                                        // 标题
		ArticleContent string    `gorm:"column:article_content;type:longtext;not null;default:'';comment:内容" json:"article_content"`                                       // 内容
		ArticleType    int64     `gorm:"column:article_type;type:tinyint;not null;default:0;comment:文章类型 1原创 2转载 3翻译" json:"article_type"`                                 // 文章类型 1原创 2转载 3翻译
		OriginalUrl    string    `gorm:"column:original_url;type:varchar(255);not null;default:'';comment:原文链接" json:"original_url"`                                       // 原文链接
		IsTop          int64     `gorm:"column:is_top;type:tinyint(1);not null;default:0;comment:是否置顶 0否 1是" json:"is_top"`                                                // 是否置顶 0否 1是
		IsDelete       int64     `gorm:"column:is_delete;type:tinyint(1);not null;index:idx_delete_status_id,priority:1;default:0;comment:是否删除  0否 1是" json:"is_delete"`   // 是否删除  0否 1是
		Status         int64     `gorm:"column:status;type:tinyint;not null;index:idx_delete_status_id,priority:2;default:1;comment:状态值 1公开 2私密 3草稿 4评论可见" json:"status"`  // 状态值 1公开 2私密 3草稿 4评论可见
		ViewCount      int64     `gorm:"column:view_count;type:bigint;not null;default:0;comment:查看数" json:"view_count"`                                                   // 查看数
		CreatedAt      time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:发表时间" json:"created_at"`                                // 发表时间
		UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                                // 更新时间
	}

	// 接口实现
	defaultTArticleModel struct {
		baseModel[TArticle]
	}
)

func NewTArticleModel(db *gorm.DB) TArticleModel {
	return &defaultTArticleModel{newBaseModel[TArticle](db, "`t_article`")}
}

// 在事务中操作
func (m *defaultTArticleModel) WithTx(tx *gorm.DB) (out TArticleModel) {
	return &defaultTArticleModel{newBaseModel[TArticle](tx, m.table)}
}
