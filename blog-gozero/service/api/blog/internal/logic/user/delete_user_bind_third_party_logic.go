package user

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserBindThirdPartyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除用户绑定第三方平台账号
func NewDeleteUserBindThirdPartyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserBindThirdPartyLogic {
	return &DeleteUserBindThirdPartyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserBindThirdPartyLogic) DeleteUserBindThirdParty(req *types.DeleteUserBindThirdPartyReq) (resp *types.EmptyResp, err error) {
	in := &accountrpc.UnbindUserOauthReq{
		Platform: req.Platform,
	}

	_, err = l.svcCtx.AccountRpc.UnbindUserOauth(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
