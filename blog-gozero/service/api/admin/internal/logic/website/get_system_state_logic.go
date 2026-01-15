package website

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/systemx"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSystemStateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取服务器信息
func NewGetSystemStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSystemStateLogic {
	return &GetSystemStateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSystemStateLogic) GetSystemState(req *types.EmptyReq) (resp *types.Server, err error) {
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
