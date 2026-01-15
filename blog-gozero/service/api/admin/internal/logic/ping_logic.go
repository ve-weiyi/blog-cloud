package logic

import (
	"context"
	"time"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// ping
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
		Description: "",
		Runtime:     runtime.String(),
	}

	return
}
