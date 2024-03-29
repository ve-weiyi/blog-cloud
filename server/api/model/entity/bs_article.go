package entity

import "time"

// TableNameArticle return the table name of <article>
const TableNameArticle = "article"

// Article mapped from table <article>
type Article struct {
	ID             int       `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:id" json:"id"`                    // id
	UserID         int       `gorm:"column:user_id;type:int;not null;comment:作者" json:"user_id"`                                        // 作者
	CategoryID     int       `gorm:"column:category_id;type:int;not null;comment:文章分类" json:"category_id"`                              // 文章分类
	ArticleCover   string    `gorm:"column:article_cover;type:varchar(1024);not null;comment:文章缩略图" json:"article_cover"`               // 文章缩略图
	ArticleTitle   string    `gorm:"column:article_title;type:varchar(64);not null;comment:标题" json:"article_title"`                    // 标题
	ArticleContent string    `gorm:"column:article_content;type:longtext;not null;comment:内容" json:"article_content"`                   // 内容
	Type           int       `gorm:"column:type;type:tinyint;not null;comment:文章类型 1原创 2转载 3翻译" json:"type"`                            // 文章类型 1原创 2转载 3翻译
	OriginalURL    string    `gorm:"column:original_url;type:varchar(255);not null;comment:原文链接" json:"original_url"`                   // 原文链接
	IsTop          int       `gorm:"column:is_top;type:tinyint(1);not null;comment:是否置顶 0否 1是" json:"is_top"`                           // 是否置顶 0否 1是
	IsDelete       int       `gorm:"column:is_delete;type:tinyint(1);not null;comment:是否删除  0否 1是" json:"is_delete"`                    // 是否删除  0否 1是
	Status         int       `gorm:"column:status;type:tinyint;not null;default:1;comment:状态值 1公开 2私密 3评论可见" json:"status"`             // 状态值 1公开 2私密 3评论可见
	CreatedAt      time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:发表时间" json:"created_at"` // 发表时间
	UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName Article's table name
func (*Article) TableName() string {
	return TableNameArticle
}
