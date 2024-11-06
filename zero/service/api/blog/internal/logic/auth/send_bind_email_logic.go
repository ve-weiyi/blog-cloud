package auth

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/accountrpc"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendBindEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发送绑定邮箱验证码
func NewSendBindEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendBindEmailLogic {
	return &SendBindEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendBindEmailLogic) SendBindEmail(req *types.UserEmailReq) (resp *types.EmptyResp, err error) {
	in := &accountrpc.UserEmailReq{
		Username: req.Username,
	}

	_, err = l.svcCtx.AccountRpc.SendBindEmail(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
