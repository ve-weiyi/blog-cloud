package guestservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/guestrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/adapter/gormx/queryx"
)

type ListGuestsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListGuestsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGuestsLogic {
	return &ListGuestsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页查询游客信息列表
func (l *ListGuestsLogic) ListGuests(in *guestrpc.ListGuestsRequest) (*guestrpc.ListGuestsResponse, error) {
	var opts []queryx.Option
	if in.PageQuery != nil {
		if in.PageQuery.Page > 0 {
			opts = append(opts, queryx.WithPage(int(in.PageQuery.Page)))
		}
		if in.PageQuery.PageSize > 0 {
			opts = append(opts, queryx.WithSize(int(in.PageQuery.PageSize)))
		}
		opts = append(opts, queryx.WithSorts(in.PageQuery.Sorts...))
	}
	if in.DeviceId != nil {
		opts = append(opts, queryx.WithCondition("device_id = ?", *in.DeviceId))
	}
	if in.IpAddress != nil {
		opts = append(opts, queryx.WithCondition("ip_address = ?", *in.IpAddress))
	}
	if len(in.DeviceIds) > 0 {
		opts = append(opts, queryx.WithCondition("device_id IN ?", in.DeviceIds))
	}

	page, size, sorts, conditions, params := queryx.NewQueryBuilder(opts...).Build()
	records, total, err := l.svcCtx.TGuestModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*guestrpc.Guest
	for _, v := range records {
		list = append(list, convertGuestOut(v))
	}

	return &guestrpc.ListGuestsResponse{
		PageResult: &guestrpc.PageResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: list,
	}, nil
}

func convertGuestOut(in *model.TGuest) *guestrpc.Guest {
	return &guestrpc.Guest{
		Id:        in.Id,
		DeviceId:  in.DeviceId,
		Os:        in.Os,
		Browser:   in.Browser,
		IpAddress: in.IpAddress,
		IpSource:  in.IpSource,
		CreatedAt: in.CreatedAt.UnixMilli(),
		UpdatedAt: in.UpdatedAt.UnixMilli(),
	}
}
