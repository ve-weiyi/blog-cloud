package auth

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendPhoneVerifyCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发送手机验证码
func NewSendPhoneVerifyCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendPhoneVerifyCodeLogic {
	return &SendPhoneVerifyCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendPhoneVerifyCodeLogic) SendPhoneVerifyCode(req *types.SendPhoneVerifyCodeReq) (resp *types.EmptyResp, err error) {
	in := &accountrpc.SendPhoneVerifyCodeReq{
		Phone: req.Phone,
		Type:  req.Type,
	}

	_, err = l.svcCtx.AccountRpc.SendPhoneVerifyCode(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
