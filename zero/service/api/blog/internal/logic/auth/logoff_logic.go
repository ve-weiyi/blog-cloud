package auth

import (
	"context"
	"time"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/accountrpc"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoffLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 注销
func NewLogoffLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoffLogic {
	return &LogoffLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoffLogic) Logoff(req *types.EmptyReq) (resp *types.EmptyResp, err error) {
	in := accountrpc.LogoffReq{
		UserId: cast.ToString(l.ctx.Value("uid")),
	}

	_, err = l.svcCtx.AccountRpc.Logoff(l.ctx, &in)
	if err != nil {
		return
	}

	l.svcCtx.TokenHolder.SetLogout(l.ctx, cast.ToString(l.ctx.Value("uid")), time.Now().Unix())
	return &types.EmptyResp{}, nil
}
