package tokenx

import (
	"fmt"
)

const (
	TokenTypeBearer = "Bearer" // 使用jwt
	TokenTypeSign   = "Sign"   // 使用md5

	// Token存储前缀
	TokenPrefixAccess  = "access"  // AccessToken前缀
	TokenPrefixRefresh = "refresh" // RefreshToken前缀
)

var (
	ErrTokenEmpty   = fmt.Errorf("token is empty")
	ErrTokenInvalid = fmt.Errorf("token is invalid")
	ErrTokenExpired = fmt.Errorf("token is expired")
)

// Token 结构体定义 Token 相关的核心字段，适配不同 Token 实现的统一返回格式
type Token struct {
	TokenType         string `json:"token_type"`          // Token 类型（如 "Bearer"）
	AccessToken       string `json:"access_token"`        // 访问令牌：用于接口访问，有效期短
	ExpiresIn         int64  `json:"expires_in"`          // AccessToken 有效期（秒），如 3600（1小时）
	RefreshToken      string `json:"refresh_token"`       // 刷新令牌：仅用于刷新 AccessToken，有效期长
	RefreshExpiresIn  int64  `json:"refresh_expires_in"`  // RefreshToken 有效期（秒），如 604800（7天）
	RefreshExpiresAt  int64  `json:"refresh_expires_at"`  // RefreshToken 过期时间戳（秒）
}

// TokenManager 定义 Token 全生命周期管理的核心接口，适配多实现（JWT/Sign）
// 所有 Token 实现（如 JWTokenManager、SignTokenManager）都需遵循此接口规范
type TokenManager interface {
	// GenerateToken 基于用户唯一标识 uid 生成 Token 结构体
	// 入参：
	//   uid - 用户唯一标识（string 类型，可根据业务改为 int/int64）
	// 返回：
	//   *Token - 包含 AccessToken、RefreshToken 等信息的结构体
	//   error - 生成失败时返回具体错误（如 uid 为空、加密失败等）
	GenerateToken(uid string) (*Token, error)

	// ValidateToken 验证 AccessToken 的有效性
	// 验证维度：合法性（签名/加密规则）、是否过期、是否已被作废
	// 入参：
	//   uid - 用户唯一标识
	//   accessToken - 待验证的访问令牌字符串
	// 返回：
	//   error - 验证成功返回 nil，验证失败返回具体错误（如 ErrTokenInvalid、ErrTokenExpired 等）
	ValidateToken(uid, accessToken string) error

	// RefreshToken 使用旧的 RefreshToken 刷新获取新的 Token 结构体
	// 核心逻辑：验证旧 RefreshToken 有效性 → 作废旧 RefreshToken → 生成新 Token
	// 入参：
	//   uid - 用户唯一标识
	//   refreshToken - 旧的刷新令牌字符串
	// 返回：
	//   *Token - 新生成的 Token 结构体
	//   error - 刷新失败时返回具体错误（如旧 RefreshToken 无效、作废失败等）
	RefreshToken(uid, refreshToken string) (*Token, error)

	// RevokeToken 主动作废 Token（AccessToken 或 RefreshToken）
	// 作废后 Token 立即失效，无法再使用
	// 入参：
	//   uid - 用户唯一标识
	//   isRefresh - true 表示作废 RefreshToken，false 表示作废 AccessToken
	// 返回：
	//   error - 作废失败时返回具体错误（如 Token 不存在、存储操作失败等）
	RevokeToken(uid string, isRefresh bool) error
}

// TokenStore 定义 Token 存储层的通用接口，适配不同存储实现（Redis/本地缓存/数据库等）
// 核心职责：为 TokenManager 提供 Token 增删查、过期控制等存储能力，与具体存储介质解耦
// 存储结构：key=uid, value=token信息（JSON格式包含access_token和refresh_token）
type TokenStore interface {
	// Set 存储用户的 Token 信息，并设置过期时间
	// 入参：
	//   key - 用户唯一标识 uid
	//   value - Token 信息（JSON格式，包含 access_token 和 refresh_token）
	//   expireSeconds - 过期时间（秒），0 表示永不过期
	// 返回：error - 存储失败时返回具体错误
	Set(key string, value string, expireSeconds int) error

	// Get 根据 uid 获取存储的 Token 信息
	// 入参：key - 用户唯一标识 uid
	// 返回：
	//   string - Token 信息（JSON格式，未找到返回空字符串）
	//   error - 查询失败时返回具体错误
	Get(key string) (string, error)

	// Delete 根据 uid 删除存储的 Token 记录
	// 入参：key - 用户唯一标识 uid
	// 返回：error - 删除失败时返回具体错误
	Delete(key string) error

	// Exists 检查指定 uid 是否存在 Token
	// 入参：key - 用户唯一标识 uid
	// 返回：
	//   bool - true 存在，false 不存在
	//   error - 检查失败时返回具体错误
	Exists(key string) (bool, error)

	// SetExpire 为已存在的 uid 设置过期时间
	// 入参：
	//   key - 用户唯一标识 uid
	//   expireSeconds - 过期时间（秒）
	// 返回：error - 设置失败时返回具体错误
	SetExpire(key string, expireSeconds int) error
}
