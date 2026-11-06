package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
)

type ListNotificationEventsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListNotificationEventsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNotificationEventsLogic {
	return &ListNotificationEventsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页查询通知事件与扇出进度
func (l *ListNotificationEventsLogic) ListNotificationEvents(in *notificationrpc.ListNotificationEventsRequest) (*notificationrpc.ListNotificationEventsResponse, error) {
	page, size := pageOf(in.ListQuery)

	var w where
	w.eq("status", in.Status)
	w.eq("scene", in.Scene)
	conditions, args := w.build()

	list, total, err := l.svcCtx.TNotificationEventModel.FindListAndTotal(
		l.ctx, page, size, "id DESC", conditions, args...)
	if err != nil {
		return nil, err
	}

	out := make([]*notificationrpc.NotificationEvent, 0, len(list))
	for _, item := range list {
		out = append(out, convertEvent(item))
	}

	return &notificationrpc.ListNotificationEventsResponse{
		ListResult: &notificationrpc.ListResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: out,
	}, nil
}
