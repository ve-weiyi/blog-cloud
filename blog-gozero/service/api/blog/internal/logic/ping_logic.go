package logic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
)

type PingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var runtime = time.Now()

func (l *PingLogic) Ping(req *types.PingReq) (resp *types.PingResp, err error) {
	resp = &types.PingResp{
		Env:         l.svcCtx.Config.Mode,
		Name:        l.svcCtx.Config.Name,
		Version:     "v1.0.0",
		Runtime:     runtime.String(),
		Description: "",
		RpcStatus:   nil,
	}

	return
}
