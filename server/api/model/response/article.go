package response

import (
	"time"
)

// 首页文章列表
type ArticleHome struct {
	ArticleDTO
	ArticleCategory *CategoryDTO `json:"article_category"` // 文章分类
	ArticleTagList  []*TagDTO    `json:"article_tag_list"` // 文章标签列表
}

// 后台文章列表
type ArticleBack struct {
	ArticleDTO
	CategoryName string   `json:"category_name"` // 文章分类名
	TagNameList  []string `json:"tag_name_list"` // 文章标签列表
}

type ArticleDTO struct {
	ID             int       `json:"id"`              // 文章ID
	ArticleCover   string    `json:"article_cover"`   // 文章缩略图
	ArticleTitle   string    `json:"article_title"`   // 标题
	ArticleContent string    `json:"article_content"` // 内容
	LikeCount      int       `json:"like_count"`      // 点赞量
	ViewsCount     int       `json:"views_count"`     // 浏览量
	Type           int       `json:"type"`            // 文章类型
	OriginalURL    string    `json:"original_url"`    // 原文链接
	IsTop          int       `json:"is_top"`          // 是否置顶
	IsDelete       int       `json:"is_delete"`       // 是否删除
	Status         int       `json:"status"`          // 状态值 1 公开 2 私密 3 评论可见
	CreatedAt      time.Time `json:"created_at"`      // 发表时间
	UpdatedAt      time.Time `json:"updated_at"`      // 更新时间
}

// 文章系别列表
type ArticleConditionDTO struct {
	ArticleDTOList []*ArticleHome `json:"article_dto_list"` // 文章列表
	ConditionName  string         `json:"condition_name"`   // 条件名
}

// 文章详情，包含文章内容，上一篇文章，下一篇文章，推荐文章列表，最新文章列表
type ArticlePageDetailsDTO struct {
	ArticleHome
	LastArticle          *ArticlePreviewDTO   `json:"last_article"`           // 上一篇文章
	NextArticle          *ArticlePreviewDTO   `json:"next_article"`           // 下一篇文章
	RecommendArticleList []*ArticlePreviewDTO `json:"recommend_article_list"` // 推荐文章列表
	NewestArticleList    []*ArticlePreviewDTO `json:"newest_article_list"`    // 最新文章列表
}

// 文章预览
type ArticlePreviewDTO struct {
	ID           int       `json:"id"`            // 文章ID
	ArticleCover string    `json:"article_cover"` // 文章缩略图
	ArticleTitle string    `json:"article_title"` // 标题
	CreatedAt    time.Time `json:"created_at"`    // 创建时间
}
