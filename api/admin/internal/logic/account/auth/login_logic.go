package auth

import (
	"context"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/metax"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/authservice"
)

func onLogin(ctx context.Context, svcCtx *svc.ServiceContext, login *authservice.LoginResponse) (resp *types.LoginResp, err error) {
	deviceId, _ := metax.GetApiDeviceIdFromCtx(ctx)
	tk, err := svcCtx.TokenManager.Generate(ctx, login.UserId, deviceId)
	if err != nil {
		return nil, err
	}

	return &types.LoginResp{
		UserId: login.UserId,
		Scope:  svcCtx.Config.Name,
		Token:  toToken(tk),
	}, nil
}
