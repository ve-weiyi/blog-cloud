package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type DeleteMeNotificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除我的通知
func NewDeleteMeNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMeNotificationLogic {
	return &DeleteMeNotificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMeNotificationLogic) DeleteMeNotification(req *types.DeleteMeNotificationReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.NotificationService.DeleteNotifications(l.ctx,
		&notificationservice.DeleteNotificationsRequest{GroupKey: req.GroupKey})
	if err != nil {
		return nil, err
	}
	return &types.BatchResp{SuccessCount: out.SuccessCount}, nil
}
