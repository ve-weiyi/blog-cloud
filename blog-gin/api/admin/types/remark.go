package types

type QueryRemarkReq struct {
	PageQuery
	UserId string `json:"user_id,optional"` // 用户ID
	Status int64  `json:"status,optional"`  // 状态
}

type UpdateRemarkStatusReq struct {
	Ids    []int64 `json:"ids,optional"`
	Status int64   `json:"status"` // 状态
}
