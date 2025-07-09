package types

type TagQueryReq struct {
	PageQuery
	TagName string `json:"tag_name,optional"` // 标签名
}
