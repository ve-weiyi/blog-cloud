package user

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserBindPhoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改用户绑定手机号
func NewUpdateUserBindPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserBindPhoneLogic {
	return &UpdateUserBindPhoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserBindPhoneLogic) UpdateUserBindPhone(req *types.UpdateUserBindPhoneReq) (resp *types.EmptyResp, err error) {
	in := &accountrpc.BindUserPhoneReq{
		Phone:      req.Phone,
		VerifyCode: req.VerifyCode,
	}

	_, err = l.svcCtx.AccountRpc.BindUserPhone(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
