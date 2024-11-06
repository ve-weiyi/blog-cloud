package auth

import (
	"context"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindUserEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 绑定邮箱
func NewBindUserEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindUserEmailLogic {
	return &BindUserEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindUserEmailLogic) BindUserEmail(req *types.BindUserEmailReq) (resp *types.EmptyResp, err error) {

	in := &accountrpc.BindUserEmailReq{
		UserId:     cast.ToString(l.ctx.Value("uid")),
		Email:      req.Email,
		VerifyCode: req.VerifyCode,
	}

	_, err = l.svcCtx.AccountRpc.BindUserEmail(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
