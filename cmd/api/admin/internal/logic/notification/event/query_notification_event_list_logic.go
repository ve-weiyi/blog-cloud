package event

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type QueryNotificationEventListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取通知事件与扇出进度
func NewQueryNotificationEventListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryNotificationEventListLogic {
	return &QueryNotificationEventListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryNotificationEventListLogic) QueryNotificationEventList(req *types.QueryNotificationEventListReq) (resp *types.ListResult, err error) {
	out, err := l.svcCtx.NotificationService.ListNotificationEvents(l.ctx, &notificationservice.ListNotificationEventsRequest{
		ListQuery: &notificationservice.ListQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
		Status:    req.Status,
		Scene:     req.Scene,
	})
	if err != nil {
		return nil, err
	}

	list := make([]*types.NotificationEventVO, 0, len(out.List))
	for _, item := range out.List {
		list = append(list, convertEvent(item))
	}

	return &types.ListResult{
		Page:     out.GetListResult().GetPage(),
		PageSize: out.GetListResult().GetPageSize(),
		Total:    out.GetListResult().GetTotal(),
		List:     list,
	}, nil
}
