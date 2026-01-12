package types

type CategoryBackVO struct {
	Id           int64  `json:"id,optional"`
	CategoryName string `json:"category_name"` // 分类名
	ArticleCount int64  `json:"article_count"`
	CreatedAt    int64  `json:"created_at"` // 创建时间
	UpdatedAt    int64  `json:"updated_at"` // 更新时间
}

type NewCategoryReq struct {
	Id           int64  `json:"id,optional"`
	CategoryName string `json:"category_name"` // 分类名
}

type QueryCategoryReq struct {
	PageQuery
	CategoryName string `json:"category_name,optional"` // 分类名
}
