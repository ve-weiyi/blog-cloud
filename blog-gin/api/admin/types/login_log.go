package types

type QueryLoginLogReq struct {
	PageQuery
	UserId string `json:"user_id,optional"` // 用户id
}
