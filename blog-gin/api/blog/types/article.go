package types

type ArticleArchivesQueryReq struct {
	PageQuery
}

type ArticleClassifyQueryReq struct {
	PageQuery
	ClassifyName string `json:"classify_name,optional"` // 分类名
}

type ArticleDetails struct {
	ArticleHome
	Author               *UserInfoVO       `json:"author"`                 // 作者
	LastArticle          *ArticlePreview   `json:"last_article"`           // 上一篇文章
	NextArticle          *ArticlePreview   `json:"next_article"`           // 下一篇文章
	RecommendArticleList []*ArticlePreview `json:"recommend_article_list"` // 推荐文章列表
	NewestArticleList    []*ArticlePreview `json:"newest_article_list"`    // 最新文章列表
}

type ArticleHomeQueryReq struct {
	PageQuery
	ArticleTitle string `json:"article_title,optional"` // 标题
}
