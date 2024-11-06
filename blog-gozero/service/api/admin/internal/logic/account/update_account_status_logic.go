package account

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAccountStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改用户状态
func NewUpdateAccountStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAccountStatusLogic {
	return &UpdateAccountStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAccountStatusLogic) UpdateAccountStatus(req *types.UpdateAccountStatusReq) (resp *types.EmptyResp, err error) {
	in := &accountrpc.UpdateUserStatusReq{
		UserId: req.UserId,
		Status: req.Status,
	}

	_, err = l.svcCtx.AccountRpc.UpdateUserStatus(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
