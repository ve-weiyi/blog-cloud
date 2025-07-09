package types

type CategoryQueryReq struct {
	PageQuery
	CategoryName string `json:"category_name,optional"` // 分类名
}
