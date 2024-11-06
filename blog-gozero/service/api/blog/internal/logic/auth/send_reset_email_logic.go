package auth

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendResetEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发送重置密码邮件
func NewSendResetEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendResetEmailLogic {
	return &SendResetEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendResetEmailLogic) SendResetEmail(req *types.UserEmailReq) (resp *types.EmptyResp, err error) {
	in := &accountrpc.UserEmailReq{
		Username: req.Username,
	}

	_, err = l.svcCtx.AccountRpc.SendResetPasswordEmail(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
