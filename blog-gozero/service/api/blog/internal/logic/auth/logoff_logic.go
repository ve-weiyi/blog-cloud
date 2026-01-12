package auth

import (
	"context"
	"time"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizheader"

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
		UserId: cast.ToString(l.ctx.Value(bizheader.HeaderUid)),
	}

	_, err = l.svcCtx.AccountRpc.Logoff(l.ctx, &in)
	if err != nil {
		return
	}

	// 登录日志
	_, err = l.svcCtx.SyslogRpc.AddLogoutLog(l.ctx, &syslogrpc.AddLogoutLogReq{
		UserId:   cast.ToString(l.ctx.Value(bizheader.HeaderUid)),
		LogoutAt: time.Now().Unix(),
	})

	l.svcCtx.TokenHolder.RemoveToken(l.ctx, cast.ToString(l.ctx.Value(bizheader.HeaderUid)))
	return &types.EmptyResp{}, nil
}
