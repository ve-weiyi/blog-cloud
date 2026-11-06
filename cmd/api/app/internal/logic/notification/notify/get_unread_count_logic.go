package notify

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type GetUnreadCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取未读分组数
func NewGetUnreadCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUnreadCountLogic {
	return &GetUnreadCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUnreadCountLogic) GetUnreadCount(req *types.UnreadCountReq) (resp *types.UnreadCountResp, err error) {
	out, err := l.svcCtx.NotificationService.CountUnreadNotifications(l.ctx,
		&notificationservice.CountUnreadNotificationsRequest{Groups: req.Groups})
	if err != nil {
		return nil, err
	}
	return &types.UnreadCountResp{Total: out.Total, GroupCounts: out.GroupCounts}, nil
}
