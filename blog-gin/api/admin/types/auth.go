package types

type EmailLoginReq struct {
	Email       string `json:"email"`                 // 邮箱
	Password    string `json:"password"`              // 密码
	CaptchaKey  string `json:"captcha_key,optional"`  // 验证码key
	CaptchaCode string `json:"captcha_code,optional"` // 验证码
}

type GetCaptchaCodeReq struct {
	Width  int64 `json:"width,optional"`  // 宽度
	Height int64 `json:"height,optional"` // 高度
}

type GetCaptchaCodeResp struct {
	CaptchaKey    string `json:"captcha_key"`    // 验证码key
	CaptchaBase64 string `json:"captcha_base64"` // 验证码base64
	CaptchaCode   string `json:"captcha_code"`   // 验证码
}

type GetClientInfoReq struct {
}

type GetClientInfoResp struct {
	Id         int64  `json:"id"`          // 访客唯一ID
	TerminalId string `json:"terminal_id"` // 终端ID
	Os         string `json:"os"`          // 操作系统
	Browser    string `json:"browser"`     // 浏览器
	IpAddress  string `json:"ip_address"`  // IP地址
	IpSource   string `json:"ip_source"`   // IP归属地
}

type GetOauthAuthorizeUrlReq struct {
	Platform string `json:"platform"`       // 平台
	State    string `json:"state,optional"` // 状态
}

type GetOauthAuthorizeUrlResp struct {
	AuthorizeUrl string `json:"authorize_url"` // 授权地址
}

type LoginReq struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	CaptchaKey  string `json:"captcha_key,optional"`  // 验证码key
	CaptchaCode string `json:"captcha_code,optional"` // 验证码
}

type LoginResp struct {
	UserId string `json:"user_id"` // 用户id
	Scope  string `json:"scope"`   // 作用域
	Token  *Token `json:"token"`
}

type PhoneLoginReq struct {
	Phone      string `json:"phone"`       // 手机号
	VerifyCode string `json:"verify_code"` // 验证码
}

type RefreshTokenReq struct {
	UserId       string `json:"user_id"`       // 用户id
	GrantType    string `json:"grant_type"`    // 授权类型
	RefreshToken string `json:"refresh_token"` // 刷新令牌
}

type RegisterReq struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"` // 确认密码
	Email           string `json:"email"`            // 邮箱
	VerifyCode      string `json:"verify_code"`      // 验证码
}

type ResetPasswordReq struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"` // 确认密码
	Email           string `json:"email"`
	VerifyCode      string `json:"verify_code"` // 验证码
}

type SendEmailVerifyCodeReq struct {
	Email string `json:"email"` // 邮箱
	Type  string `json:"type"`  // 类型 register,reset_password,bind_email,bind_phone
}

type SendPhoneVerifyCodeReq struct {
	Phone string `json:"phone"` // 手机号
	Type  string `json:"type"`  // 类型 register,reset_password,bind_email,bind_phone
}

type ThirdLoginReq struct {
	Platform string `json:"platform"`      // 平台
	Code     string `json:"code,optional"` // 授权码
}
