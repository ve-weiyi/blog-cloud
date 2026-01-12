package types

type QueryCommentReq struct {
	PageQuery
	UserId string `json:"user_id,optional"` // 用户ID
	Status int64  `json:"status,optional"`  // 状态
	Type   int64  `json:"type,optional"`    // 评论类型 1.文章 2.友链 3.说说
}

type UpdateCommentStatusReq struct {
	Ids    []int64 `json:"ids,optional"`
	Status int64   `json:"status"` // 状态
}
