package model

import (
	"gorm.io/gorm"
)

var _ TArticleTagModel = (*defaultTArticleTagModel)(nil)

type (
	// 接口定义
	TArticleTagModel interface {
		Table[TArticleTag]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TArticleTagModel)
	}

	// 表字段定义
	TArticleTag struct {
		Id        int64 `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:id" json:"id"`                                          // id
		ArticleId int64 `gorm:"column:article_id;type:bigint;not null;index:fk_article_tag_1,priority:1;default:0;comment:文章id" json:"article_id"` // 文章id
		TagId     int64 `gorm:"column:tag_id;type:bigint;not null;index:fk_article_tag_2,priority:1;default:0;comment:标签id" json:"tag_id"`         // 标签id
	}

	// 接口实现
	defaultTArticleTagModel struct {
		baseModel[TArticleTag]
	}
)

func NewTArticleTagModel(db *gorm.DB) TArticleTagModel {
	return &defaultTArticleTagModel{newBaseModel[TArticleTag](db, "`t_article_tag`")}
}

// 在事务中操作
func (m *defaultTArticleTagModel) WithTx(tx *gorm.DB) (out TArticleTagModel) {
	return &defaultTArticleTagModel{newBaseModel[TArticleTag](tx, m.table)}
}
