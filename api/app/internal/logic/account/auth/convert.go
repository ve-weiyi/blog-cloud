package auth

import (
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/tokenx"
)

// TokenTypeBearer 令牌类型，前端按 Bearer 方案携带 AccessToken。
const TokenTypeBearer = "Bearer"

// toToken 将令牌对转换为前端契约结构。
func toToken(tk *tokenx.TokenPair) *types.Token {
	return &types.Token{
		TokenType:        TokenTypeBearer,
		AccessToken:      tk.AccessToken,
		ExpiresIn:        tk.AccessExpiresIn,
		RefreshToken:     tk.RefreshToken,
		RefreshExpiresIn: tk.RefreshExpiresIn,
		RefreshExpiresAt: tk.RefreshExpiresAt.Unix(),
	}
}
