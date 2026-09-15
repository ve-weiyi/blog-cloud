package tokenx

import "time"

// LoginMode 登录模式。
type LoginMode int

const (
	MultiPoint  LoginMode = iota // 多点登录（零值默认）
	SinglePoint                  // 单点登录
)

// Config 令牌管理器配置。
type Config struct {
	Signer          Signer
	Store           Store
	LoginMode       LoginMode
	KeyPrefix       string // key 前缀，如 "myapp"，零值为空
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}
