package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TArticleChunkModel = (*defaultTArticleChunkModel)(nil)

// ArticleChunkHit 是检索命中的分块及其所属文章元信息。
type ArticleChunkHit struct {
	ArticleId    int64   `gorm:"column:article_id" json:"article_id"`
	ArticleTitle string  `gorm:"column:article_title" json:"article_title"`
	Heading      string  `gorm:"column:heading" json:"heading"`
	Content      string  `gorm:"column:content" json:"content"`
	Summary      string  `gorm:"column:summary" json:"summary"`
	Score        float64 `gorm:"column:score" json:"score"`
}

type (
	// 接口定义
	TArticleChunkModel interface {
		Table[TArticleChunk]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TArticleChunkModel)
		// Search 按关键词检索指定状态文章的分块；关联文章状态与删除标记在 SQL 内过滤
		Search(ctx context.Context, keyword string, status int64, limit int) (list []*ArticleChunkHit, err error)
	}

	// 表字段定义
	TArticleChunk struct {
		Id           int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键" json:"id"`                            // 主键
		ArticleId    int64     `gorm:"column:article_id;type:bigint;not null;index:idx_article;default:0;comment:所属文章id" json:"article_id"` // 所属文章id
		ArticleTitle string    `gorm:"column:article_title;type:varchar(64);not null;default:'';comment:文章标题" json:"article_title"`         // 文章标题
		Ordinal      int64     `gorm:"column:ordinal;type:int;not null;default:0;comment:分块序号" json:"ordinal"`                              // 分块序号
		Heading      string    `gorm:"column:heading;type:varchar(255);not null;default:'';comment:标题路径" json:"heading"`                    // 标题路径
		Content      string    `gorm:"column:content;type:text;not null;comment:片段正文" json:"content"`                                       // 片段正文
		Summary      string    `gorm:"column:summary;type:varchar(512);not null;default:'';comment:截断摘要" json:"summary"`                    // 截断摘要
		ContentHash  string    `gorm:"column:content_hash;type:char(64);not null;default:'';comment:内容摘要，用于判断是否变化" json:"content_hash"`     // 内容摘要，用于判断是否变化
		UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`   // 更新时间
	}

	// 接口实现
	defaultTArticleChunkModel struct {
		baseModel[TArticleChunk]
	}
)

func NewTArticleChunkModel(db *gorm.DB) TArticleChunkModel {
	return &defaultTArticleChunkModel{newBaseModel[TArticleChunk](db, "`t_article_chunk`")}
}

// 在事务中操作
func (m *defaultTArticleChunkModel) WithTx(tx *gorm.DB) (out TArticleChunkModel) {
	return &defaultTArticleChunkModel{newBaseModel[TArticleChunk](tx, m.table)}
}

// Search 用中文全文索引检索分块。
//
// 文章可见性在 SQL 里再过滤一次：分块表理论上只存公开文章，但状态变更与重建
// 之间有窗口，不能把权限完全交给上一步的清理。
func (m *defaultTArticleChunkModel) Search(ctx context.Context, keyword string, status int64, limit int) (list []*ArticleChunkHit, err error) {
	const query = `
SELECT c.article_id,
       a.article_title,
       c.heading,
       c.content,
       c.summary,
       MATCH(c.article_title, c.heading, c.content) AGAINST(? IN NATURAL LANGUAGE MODE) AS score
FROM t_article_chunk c
JOIN t_article a ON a.id = c.article_id
WHERE a.status = ?
  AND a.is_delete = 0
  AND MATCH(c.article_title, c.heading, c.content) AGAINST(? IN NATURAL LANGUAGE MODE)
ORDER BY score DESC, c.id ASC
LIMIT ?`

	if err = m.db.WithContext(ctx).Raw(query, keyword, status, keyword, limit).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
