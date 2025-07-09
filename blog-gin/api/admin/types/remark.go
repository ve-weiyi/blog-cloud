package types

type RemarkQuery struct {
	PageQuery
	Nickname string `json:"nickname,optional"`  // 昵称
	IsReview int64  `json:"is_review,optional"` // 是否审核
}

type RemarkReviewReq struct {
	Ids      []int64 `json:"ids,optional"`
	IsReview int64   `json:"is_review"` // 是否审核
}
