package user

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserBindThirdPartyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改用户绑定第三方平台账号
func NewUpdateUserBindThirdPartyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserBindThirdPartyLogic {
	return &UpdateUserBindThirdPartyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserBindThirdPartyLogic) UpdateUserBindThirdParty(req *types.UpdateUserBindThirdPartyReq) (resp *types.EmptyResp, err error) {
	in := &accountrpc.BindUserOauthReq{
		Platform: req.Platform,
		Code:     req.Code,
	}

	_, err = l.svcCtx.AccountRpc.BindUserOauth(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
