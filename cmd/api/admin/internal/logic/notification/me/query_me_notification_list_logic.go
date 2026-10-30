package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type QueryMeNotificationListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取我的通知列表（按分组）
func NewQueryMeNotificationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryMeNotificationListLogic {
	return &QueryMeNotificationListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// QueryMeNotificationList 读取管理员自己收到的通知。
//
// 收件人由 AdminAuth 写入 ctx，请求里没有 user_id，无法构造越权。
func (l *QueryMeNotificationListLogic) QueryMeNotificationList(req *types.QueryMeNotificationListReq) (resp *types.ListResult, err error) {
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
