package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type ReadMeNotificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 标记我的通知已读
func NewReadMeNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadMeNotificationLogic {
	return &ReadMeNotificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadMeNotificationLogic) ReadMeNotification(req *types.ReadMeNotificationReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.NotificationService.MarkNotificationsRead(l.ctx,
		&notificationservice.MarkNotificationsReadRequest{GroupKey: req.GroupKey})
	if err != nil {
		return nil, err
	}
	return &types.BatchResp{SuccessCount: out.SuccessCount}, nil
}
