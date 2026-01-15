package auth

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 刷新token
func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenReq) (resp *types.LoginResp, err error) {
	tk, err := l.svcCtx.TokenManager.RefreshToken(req.UserId, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &types.LoginResp{
		UserId: req.UserId,
		Scope:  l.svcCtx.Config.Name,
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
