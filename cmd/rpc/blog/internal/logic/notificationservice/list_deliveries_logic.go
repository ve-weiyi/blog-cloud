package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
)

type ListDeliveriesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListDeliveriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDeliveriesLogic {
	return &ListDeliveriesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页查询外部渠道投递记录
func (l *ListDeliveriesLogic) ListDeliveries(in *notificationrpc.ListDeliveriesRequest) (*notificationrpc.ListDeliveriesResponse, error) {
	page, size := pageOf(in.ListQuery)

	var w where
	w.eq("channel", in.Channel)
	w.eq("status", in.Status)
	w.eq("recipient", in.Recipient)
	w.eqInt("notification_id", in.NotificationId)
	conditions, args := w.build()

	list, total, err := l.svcCtx.TNotificationDeliveryModel.FindListAndTotal(
		l.ctx, page, size, "id DESC", conditions, args...)
	if err != nil {
		return nil, err
	}

	out := make([]*notificationrpc.Delivery, 0, len(list))
	for _, item := range list {
		out = append(out, convertDelivery(item))
	}

	return &notificationrpc.ListDeliveriesResponse{
		ListResult: &notificationrpc.ListResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: out,
	}, nil
}
