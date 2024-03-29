package request

import (
	"time"
)

type ArticleConditionReq struct {
	TagID      int `json:"tag_id"`      // 文章标签ID
	CategoryID int `json:"category_id"` // 文章分类ID
}

type ArticleTopReq struct {
	ID    int `json:"id"`     // 文章ID
	IsTop int `json:"is_top"` // 是否置顶
}

type ArticleDeleteReq struct {
	ID       int `json:"id"`        // 文章ID
	IsDelete int `json:"is_delete"` // 是否删除
}

type ArticleDetailsDTOReq struct {
	ID             int       `json:"id"`              // 文章ID
	ArticleCover   string    `json:"article_cover"`   // 文章缩略图
	ArticleTitle   string    `json:"article_title"`   // 标题
	ArticleContent string    `json:"article_content"` // 内容
	LikeCount      int       `json:"like_count"`      // 点赞量
	ViewsCount     int       `json:"views_count"`     // 浏览量
	Type           int       `json:"type"`            // 文章类型
	OriginalURL    string    `json:"original_url"`    // 原文链接
	IsTop          int       `json:"is_top"`          // 是否置顶
	Status         int       `json:"status"`          // 状态值 1 公开 2 私密 3 评论可见
	CreatedAt      time.Time `json:"created_at"`      // 发表时间
	UpdatedAt      time.Time `json:"updated_at"`      // 更新时间
	CategoryName   string    `json:"category_name"`   // 文章分类名
	TagNameList    []string  `json:"tag_name_list"`   // 文章标签列表
}
