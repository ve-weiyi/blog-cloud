package accountrpclogic

import (
	"context"
	"time"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 登出
func (l *LogoutLogic) Logout(in *accountrpc.LogoutReq) (*accountrpc.LogoutResp, error) {
	err := l.svcCtx.OnlineUserService.Logout(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	return &accountrpc.LogoutResp{
		LogoutAt: time.Now().UnixMilli(),
	}, nil
}
