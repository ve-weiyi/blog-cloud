package user

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserBindEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改用户绑定邮箱
func NewUpdateUserBindEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserBindEmailLogic {
	return &UpdateUserBindEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserBindEmailLogic) UpdateUserBindEmail(req *types.UpdateUserBindEmailReq) (resp *types.EmptyResp, err error) {
	in := &accountrpc.BindUserEmailReq{
		Email:      req.Email,
		VerifyCode: req.VerifyCode,
	}

	_, err = l.svcCtx.AccountRpc.BindUserEmail(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
