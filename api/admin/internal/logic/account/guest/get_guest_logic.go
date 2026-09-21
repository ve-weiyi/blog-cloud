package guest

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/guestservice"
)

type GetGuestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取游客信息
func NewGetGuestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGuestLogic {
	return &GetGuestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGuestLogic) GetGuest(req *types.GetGuestReq) (resp *types.GetGuestResp, err error) {
	in := guestservice.GetGuestRequest{}

	out, err := l.svcCtx.GuestService.GetGuest(l.ctx, &in)
	if err != nil {
		return
	}

	return &types.GetGuestResp{
		Id:        out.Guest.Id,
		DeviceId:  out.Guest.DeviceId,
		Os:        out.Guest.Os,
		Browser:   out.Guest.Browser,
		IpAddress: out.Guest.IpAddress,
		IpSource:  out.Guest.IpSource,
	}, nil
}
