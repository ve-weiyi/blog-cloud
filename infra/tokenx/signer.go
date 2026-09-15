package tokenx

import "context"

// Signer 令牌签名与验证接口。
// HMAC 实现：Verify 返回 (nil, nil) 表示不透明令牌不可自解析。
// JWT 实现：Verify 返回解析后的 *TokenClaims。
type Signer interface {
	Sign(ctx context.Context, claims TokenClaims) (string, error)

	// Verify 返回 (*TokenClaims, error)
	//   claims != nil, err == nil  → 自含令牌（JWT），已验签并解析
	//   claims == nil, err == nil  → 不透明令牌（HMAC），不可逆，请查 Store
	//   claims == nil, err != nil  → 验签失败（ErrTokenExpired 或 ErrTokenInvalid）
	Verify(ctx context.Context, token string) (*TokenClaims, error)
}
