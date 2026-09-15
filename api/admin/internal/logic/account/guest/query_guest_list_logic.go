package guest

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/guestservice"
)

type QueryGuestListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取访客列表
func NewQueryGuestListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryGuestListLogic {
	return &QueryGuestListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryGuestListLogic) QueryGuestList(req *types.QueryGuestListReq) (resp *types.ListResult, err error) {
	out, err := l.svcCtx.GuestService.ListGuests(l.ctx, &guestservice.ListGuestsRequest{
		ListQuery: &guestservice.ListQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
		DeviceId:  req.DeviceId,
	})
	if err != nil {
		return nil, err
	}

	var list []*types.GuestItem
	for _, v := range out.List {
		list = append(list, &types.GuestItem{
			Id:        v.Id,
			DeviceId:  v.DeviceId,
			Os:        v.Os,
			Browser:   v.Browser,
			IpAddress: v.IpAddress,
			IpSource:  v.IpSource,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return &types.ListResult{
		Page:     out.ListResult.Page,
		PageSize: out.ListResult.PageSize,
		Total:    out.ListResult.Total,
		List:     list,
	}, nil
}
