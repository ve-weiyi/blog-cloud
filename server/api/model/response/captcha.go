package response

type CaptchaDTO struct {
	ID         string `json:"id"`
	EncodeData string `json:"encode_data"` // 验证码内容，base64编码
	Length     int    `json:"length"`
}
