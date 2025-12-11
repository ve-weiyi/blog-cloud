package accountrpclogic

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
)

type SendPhoneVerifyCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendPhoneVerifyCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendPhoneVerifyCodeLogic {
	return &SendPhoneVerifyCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 发送手机号验证码
func (l *SendPhoneVerifyCodeLogic) SendPhoneVerifyCode(in *accountrpc.SendPhoneVerifyCodeReq) (*accountrpc.EmptyResp, error) {
	// todo
	if in.Phone != "" {
		return nil, fmt.Errorf("暂不支持手机号验证码")
	}

	return &accountrpc.EmptyResp{}, nil
}
