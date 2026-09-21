package tokenx

import (
	"time"

	"github.com/ve-weiyi/blog-cloud/infra/storex"
)

// LoginMode 登录模式。
type LoginMode int

const (
	MultiPoint  LoginMode = iota // 多点登录（零值默认）
	SinglePoint                  // 单点登录
)

// Config 令牌管理器配置。
type Config struct {
	Signer    Signer
	Store     storex.KVStore     // 令牌与会话的键值存储
	Devices   storex.MemberStore // userID → deviceID 索引，供 RevokeAll / ListSessions 使用
	LoginMode LoginMode
	KeyPrefix string // key 前缀，如 "myapp"，零值为空

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}
