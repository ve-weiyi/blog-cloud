package types

type LoginLogQuery struct {
	PageQuery
	UserId string `json:"user_id,optional"` // 用户id
}
