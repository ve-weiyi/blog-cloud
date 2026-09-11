package auth

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/metax"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/userauthservice"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 登出
func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.LogoutReq) (resp *types.LogoutResp, err error) {
	uid, _ := metax.GetApiUserIdFromCtx(l.ctx)

	in := userauthservice.LogoutRequest{}
	_, err = l.svcCtx.UserAuthService.Logout(l.ctx, &in)
	if err != nil {
		return
	}

	err = l.svcCtx.TokenStore.RevokeToken(uid, false)
	if err != nil {
		return nil, err
	}

	return &types.LogoutResp{}, nil
}
