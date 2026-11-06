package notify

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type QueryNotificationListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取通知列表（按分组分页）
func NewQueryNotificationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryNotificationListLogic {
	return &QueryNotificationListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// QueryNotificationList 读取当前登录用户的收件箱。
func (l *QueryNotificationListLogic) QueryNotificationList(req *types.QueryNotificationListReq) (resp *types.ListResult, err error) {
	out, err := l.svcCtx.NotificationService.ListNotifications(l.ctx, &notificationservice.ListNotificationsRequest{
		ListQuery:  &notificationservice.ListQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
		Groups:     req.Groups,
		UnreadOnly: req.UnreadOnly,
	})
	if err != nil {
		return nil, err
	}

	list := make([]*types.NotificationGroupVO, 0, len(out.List))
	for _, item := range out.List {
		list = append(list, convertGroup(item))
	}

	return &types.ListResult{
		Page:     out.GetListResult().GetPage(),
		PageSize: out.GetListResult().GetPageSize(),
		Total:    out.GetListResult().GetTotal(),
		List:     list,
	}, nil
}
