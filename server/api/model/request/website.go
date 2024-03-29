package request

import (
	"mime/multipart"
)

type WebsiteConfigReq struct {
	Key   string `json:"key" from:"key" example:"about"`
	Value string `json:"value" from:"value" example:"about me"`
}

// VoiceVO 表示系统中的音频数据。
type VoiceVO struct {
	Type int                   `json:"type" validate:"required"` // 消息类型
	File *multipart.FileHeader `json:"file" validate:"required"` // 文件
	//UserID    int                   `json:"user_id" validate:"required"`    // 用户id
	//Nickname  string                `json:"nickname" validate:"required"`   // 用户昵称
	//Avatar    string                `json:"avatar" validate:"required"`     // 用户头像
	Content string `json:"content" validate:"required"` // 聊天内容
	//CreatedAt time.Time             `json:"created_at" validate:"required"` // 创建时间
	//IPAddress string                `json:"ip_address" validate:"required"` // 用户登录ip
	//IPSource  string                `json:"ip_source" validate:"required"`  // ip来源
}
