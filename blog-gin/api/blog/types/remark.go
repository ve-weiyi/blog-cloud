package types

type Remark struct {
	Id             int64       `json:"id,optional"`     // 主键id
	UserId         string      `json:"user_id"`         // 用户id
	TerminalId     string      `json:"terminal_id"`     // 终端id
	MessageContent string      `json:"message_content"` // 留言内容
	IpAddress      string      `json:"ip_address"`      // 用户ip
	IpSource       string      `json:"ip_source"`       // 用户地址
	IsReview       int64       `json:"is_review"`       // 是否审核
	CreatedAt      int64       `json:"created_at"`      // 发布时间
	UpdatedAt      int64       `json:"updated_at"`      // 更新时间
	User           *UserInfoVO `json:"user"`            // 用户信息
}

type RemarkNewReq struct {
	MessageContent string `json:"message_content"` // 留言内容
}

type RemarkQueryReq struct {
	PageQuery
}
