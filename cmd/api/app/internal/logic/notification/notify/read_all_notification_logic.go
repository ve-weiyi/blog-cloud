package notify

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type ReadAllNotificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 全部标记已读
func NewReadAllNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadAllNotificationLogic {
	return &ReadAllNotificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadAllNotificationLogic) ReadAllNotification(req *types.EmptyReq) (resp *types.AffectedResp, err error) {
	out, err := l.svcCtx.NotificationService.MarkAllNotificationsRead(l.ctx,
		&notificationservice.MarkAllNotificationsReadRequest{})
	if err != nil {
		return nil, err
	}
	return &types.AffectedResp{SuccessCount: out.SuccessCount}, nil
}
