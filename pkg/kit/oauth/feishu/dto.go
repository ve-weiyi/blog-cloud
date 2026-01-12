package feishu

type AppTokenResp struct {
	AppAccessToken    string `json:"app_access_token"`
	Code              int    `json:"code"`
	Expire            int    `json:"expire"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
}

type TenantTokenResp struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}

type UserAccessTokenResp struct {
	Code int                 `json:"code"`
	Msg  string              `json:"msg"`
	Data UserAccessTokenData `json:"data"`
}

type RefreshAccessTokenResp struct {
	Code int                 `json:"code"`
	Msg  string              `json:"msg"`
	Data UserAccessTokenData `json:"data"`
}

type UserInfoResp struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data UserInfoData `json:"data"`
}

type UserAccessTokenData struct {
	AccessToken      string `json:"access_token"`       // 访问令牌
	TokenType        string `json:"token_type"`         // 令牌类型
	ExpiresIn        int    `json:"expires_in"`         // 令牌过期时间（秒）
	RefreshToken     string `json:"refresh_token"`      // 刷新令牌
	RefreshExpiresIn int    `json:"refresh_expires_in"` // 刷新令牌过期时间（秒）
	Scope            string `json:"scope"`              // 作用域
}

type UserInfoData struct {
	Name            string `json:"name"`             // 用户名
	EnName          string `json:"en_name"`          // 用户英文名
	AvatarUrl       string `json:"avatar_url"`       // 头像URL
	AvatarThumb     string `json:"avatar_thumb"`     // 头像缩略图URL
	AvatarMiddle    string `json:"avatar_middle"`    // 头像中等尺寸URL
	AvatarBig       string `json:"avatar_big"`       // 头像大尺寸URL
	OpenId          string `json:"open_id"`          // 用户在开放平台的唯一标识
	UnionId         string `json:"union_id"`         // 用户在开放平台的统一标识
	Email           string `json:"email"`            // 邮箱
	EnterpriseEmail string `json:"enterprise_email"` // 企业邮箱
	UserId          string `json:"user_id"`          // 用户ID
	Mobile          string `json:"mobile"`           // 手机号码
	TenantKey       string `json:"tenant_key"`       // 租户Key
	EmployeeNo      string `json:"employee_no"`      // 员工工号
}
