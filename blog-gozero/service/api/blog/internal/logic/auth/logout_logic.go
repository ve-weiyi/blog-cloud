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

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 登出
func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.EmptyReq) (resp *types.EmptyResp, err error) {
	uid := cast.ToString(l.ctx.Value(bizheader.HeaderUid))
	in := &accountrpc.LogoutReq{
		UserId: uid,
	}

	_, err = l.svcCtx.AccountRpc.Logout(l.ctx, in)
	if err != nil {
		return nil, err
	}

	// 登录日志
	_, err = l.svcCtx.SyslogRpc.AddLogoutLog(l.ctx, &syslogrpc.AddLogoutLogReq{
		UserId:   uid,
		LogoutAt: time.Now().UnixMilli(),
	})

	l.svcCtx.TokenManager.RevokeToken(uid, false)
	return &types.EmptyResp{}, nil
}
