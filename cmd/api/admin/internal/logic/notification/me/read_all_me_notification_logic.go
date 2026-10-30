package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type ReadAllMeNotificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 我的通知全部已读
func NewReadAllMeNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadAllMeNotificationLogic {
	return &ReadAllMeNotificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadAllMeNotificationLogic) ReadAllMeNotification(req *types.EmptyReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.NotificationService.MarkAllNotificationsRead(l.ctx,
		&notificationservice.MarkAllNotificationsReadRequest{})
	if err != nil {
		return nil, err
	}
	return &types.BatchResp{SuccessCount: out.SuccessCount}, nil
}
