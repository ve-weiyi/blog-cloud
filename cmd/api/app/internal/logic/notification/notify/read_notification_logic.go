package notify

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type ReadNotificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 标记通知已读
func NewReadNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadNotificationLogic {
	return &ReadNotificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadNotificationLogic) ReadNotification(req *types.ReadNotificationReq) (resp *types.AffectedResp, err error) {
	out, err := l.svcCtx.NotificationService.MarkNotificationsRead(l.ctx,
		&notificationservice.MarkNotificationsReadRequest{GroupKey: req.GroupKey})
	if err != nil {
		return nil, err
	}
	return &types.AffectedResp{SuccessCount: out.SuccessCount}, nil
}
