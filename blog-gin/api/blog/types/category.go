package types

type QueryCategoryReq struct {
	PageQuery
	CategoryName string `json:"category_name,optional"` // 分类名
}
