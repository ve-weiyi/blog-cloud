package response

// 用户登录信息
type Login struct {
	Token        *Token        `json:"token"`
	UserInfo     *UserInfo     `json:"user_info"`
	LoginHistory *LoginHistory `json:"login_history"`
}

type Token struct {
	TokenType        string `json:"token_type"`         // token类型,Bearer
	AccessToken      string `json:"access_token"`       // 访问token,过期时间较短。2h
	ExpiresIn        int64  `json:"expires_in"`         // 访问token过期时间
	RefreshToken     string `json:"refresh_token"`      // 刷新token,过期时间较长。30d
	RefreshExpiresIn int64  `json:"refresh_expires_in"` // 刷新token过期时间
	Scope            string `json:"scope"`              // 作用域
	UID              int    `json:"uid"`                // 用户id
}

type UserInfo struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Intro    string `json:"intro"`
	Website  string `json:"website"`
	Email    string `json:"email"`

	ArticleLikeSet []string `json:"article_like_set"` // 文章点赞集合
	CommentLikeSet []string `json:"comment_like_set"` // 评论点赞集合
	TalkLikeSet    []string `json:"talk_like_set"`    // 说说点赞集合

	Roles []*RoleDTO `json:"roles"`
}

type LoginHistory struct {
	ID        int    `json:"id"`
	LoginType string `json:"login_type"` // 登录类型
	Agent     string `json:"agent"`      // 代理
	IpAddress string `json:"ip_address"` // ip host
	IpSource  string `json:"ip_source"`  // ip 源
	LoginTime string `json:"login_time"` // 创建时间
}

type OauthLoginUrl struct {
	Url string `json:"url" example:""` // 授权地址
}
