package overview

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/vkit/x/systemx"
)

type GetSystemInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取服务器信息
func NewGetSystemInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSystemInfoLogic {
	return &GetSystemInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSystemInfoLogic) GetSystemInfo(req *types.EmptyReq) (resp *types.Server, err error) {
	sv := types.Server{}

	sv.Os = systemx.InitOS()
	if sv.Cpu, err = systemx.InitCPU(); err != nil {
		return &sv, err
	}
	if sv.Ram, err = systemx.InitRAM(); err != nil {
		return &sv, err
	}
	if sv.Disk, err = systemx.InitDisk(); err != nil {
		return &sv, err
	}

	return &sv, err
}
