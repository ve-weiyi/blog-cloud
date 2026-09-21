package auth

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/metax"
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
	deviceId, _ := metax.GetApiDeviceIdFromCtx(l.ctx)
	tk, err := l.svcCtx.TokenManager.Refresh(l.ctx, req.UserId, deviceId, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &types.LoginResp{
		UserId: req.UserId,
		Scope:  l.svcCtx.Config.Name,
		Token:  toToken(tk),
	}, nil
}
