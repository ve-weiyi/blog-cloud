package delivery

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type QueryDeliveryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取外部渠道投递记录
func NewQueryDeliveryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryDeliveryListLogic {
	return &QueryDeliveryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryDeliveryListLogic) QueryDeliveryList(req *types.QueryDeliveryListReq) (resp *types.ListResult, err error) {
	out, err := l.svcCtx.NotificationService.ListDeliveries(l.ctx, &notificationservice.ListDeliveriesRequest{
		ListQuery:      &notificationservice.ListQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
		Channel:        req.Channel,
		Status:         req.Status,
		Recipient:      req.Recipient,
		NotificationId: req.NotificationId,
	})
	if err != nil {
		return nil, err
	}

	list := make([]*types.DeliveryVO, 0, len(out.List))
	for _, item := range out.List {
		list = append(list, convertDelivery(item))
	}

	return &types.ListResult{
		Page:     out.GetListResult().GetPage(),
		PageSize: out.GetListResult().GetPageSize(),
		Total:    out.GetListResult().GetTotal(),
		List:     list,
	}, nil
}
