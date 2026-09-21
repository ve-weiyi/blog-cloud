package tokenx

import (
	"context"
	"errors"
	"time"
)

// Manager 令牌管理器接口，提供令牌全生命周期治理能力。
// 默认不保证并发安全。若需并发使用，调用方应自行同步。
type Manager interface {
	Generate(ctx context.Context, userID, deviceID string) (*TokenPair, error)
	Refresh(ctx context.Context, userID, deviceID, refreshToken string) (*TokenPair, error)
	Validate(ctx context.Context, userID, deviceID, accessToken string) error
	Revoke(ctx context.Context, userID, deviceID string) error
	RevokeAll(ctx context.Context, userID string) error
	ListSessions(ctx context.Context, userID string) ([]Session, error)
}

// TokenClaims 令牌载荷。Signer.Sign 的输入，自含令牌（JWT）Verify 的输出。
type TokenClaims struct {
	UserID    string
	DeviceID  string
	TokenType string // "access" 或 "refresh"
	IssuedAt  time.Time
	ExpiresAt time.Time
}

// TokenPair 颁发给客户端的令牌对。
type TokenPair struct {
	UserID           string // 令牌所属用户
	DeviceID         string // 令牌所属设备
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Time // AccessToken 过期时间
	AccessExpiresIn  int64     // AccessToken 有效秒数，客户端缓存友好
	RefreshExpiresAt time.Time // RefreshToken 过期时间
	RefreshExpiresIn int64     // RefreshToken 有效秒数
}

// Session 用户登录会话。
type Session struct {
	UserID    string
	DeviceID  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// 错误变量
var (
	ErrTokenEmpty   = errors.New("token is empty")
	ErrTokenInvalid = errors.New("token is invalid")
	ErrTokenExpired = errors.New("token is expired")
)
