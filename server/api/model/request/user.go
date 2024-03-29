package request

type UserInfoReq struct {
	Nickname string `json:"nickname" example:"nickname"` // 昵称
	Website  string `json:"website" example:"website"`   // 网站
	Intro    string `json:"intro" example:"intro"`       // 简介
	Avatar   string `json:"avatar" example:"avatar"`     // 头像
}
