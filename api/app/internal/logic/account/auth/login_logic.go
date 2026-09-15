package auth

import (
	"context"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/authservice"
)

func onLogin(ctx context.Context, svcCtx *svc.ServiceContext, login *authservice.LoginResponse) (resp *types.LoginResp, err error) {
	tk, err := svcCtx.TokenStore.GenerateToken(login.UserId)
	if err != nil {
		return nil, err
	}

	return &types.LoginResp{
		UserId: login.UserId,
		Scope:  svcCtx.Config.Name,
		Token: &types.Token{
			TokenType:        tk.TokenType,
			AccessToken:      tk.AccessToken,
			ExpiresIn:        tk.ExpiresIn,
			RefreshToken:     tk.RefreshToken,
			RefreshExpiresIn: tk.RefreshExpiresIn,
			RefreshExpiresAt: tk.RefreshExpiresAt,
		},
	}, nil
}
