package oauth

type UserResult struct {
	OpenID   string `json:"open_id"`    // 用户在开放平台的唯一标识
	NickName string `json:"nick_name"`  // 用户昵称
	Name     string `json:"name"`       // 用户姓名
	EnName   string `json:"en_name"`    // 用户英文名
	Avatar   string `json:"avatar_url"` // 头像URL
	Email    string `json:"email"`      // 邮箱
	Mobile   string `json:"mobile"`     // 手机号码
	//Province string `json:"province"`   // 省份
	//City     string `json:"city"`       // 城市
}
