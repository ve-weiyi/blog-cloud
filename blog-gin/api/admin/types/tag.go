package types

type NewTagReq struct {
	Id      int64  `json:"id,optional"`
	TagName string `json:"tag_name"` // 标签名
}

type QueryTagReq struct {
	PageQuery
	TagName string `json:"tag_name,optional"` // 标签名
}

type TagBackVO struct {
	Id           int64  `json:"id,optional"`   // 标签ID
	TagName      string `json:"tag_name"`      // 标签名
	ArticleCount int64  `json:"article_count"` // 文章数量
	CreatedAt    int64  `json:"created_at"`    // 创建时间
	UpdatedAt    int64  `json:"updated_at"`    // 更新时间
}
