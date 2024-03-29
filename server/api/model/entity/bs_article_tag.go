package entity

// TableNameArticleTag return the table name of <article_tag>
const TableNameArticleTag = "article_tag"

// ArticleTag mapped from table <article_tag>
type ArticleTag struct {
	ID        int `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:id" json:"id"`                       // id
	ArticleID int `gorm:"column:article_id;type:int;not null;index:fk_article_tag_1,priority:1;comment:文章id" json:"article_id"` // 文章id
	TagID     int `gorm:"column:tag_id;type:int;not null;index:fk_article_tag_2,priority:1;comment:标签id" json:"tag_id"`         // 标签id
}

// TableName ArticleTag's table name
func (*ArticleTag) TableName() string {
	return TableNameArticleTag
}
