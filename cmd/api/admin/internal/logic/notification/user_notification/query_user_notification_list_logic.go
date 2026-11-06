package user_notification

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type QueryUserNotificationListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 按用户查看收件箱条目
func NewQueryUserNotificationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryUserNotificationListLogic {
	return &QueryUserNotificationListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// QueryUserNotificationList 是排障入口：条目级、只读，不下钻投递。
func (l *QueryUserNotificationListLogic) QueryUserNotificationList(req *types.QueryUserNotificationListReq) (resp *types.ListResult, err error) {
	out, err := l.svcCtx.NotificationService.ListUserNotifications(l.ctx, &notificationservice.ListUserNotificationsRequest{
		ListQuery: &notificationservice.ListQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
		UserId:    req.UserId,
	})
	if err != nil {
		return nil, err
	}

	list := make([]*types.NotificationItemVO, 0, len(out.List))
	for _, item := range out.List {
		list = append(list, convertItem(item))
	}

	return &types.ListResult{
		Page:     out.GetListResult().GetPage(),
		PageSize: out.GetListResult().GetPageSize(),
		Total:    out.GetListResult().GetTotal(),
		List:     list,
	}, nil
}
