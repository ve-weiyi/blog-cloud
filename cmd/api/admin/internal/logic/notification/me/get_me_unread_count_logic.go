package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type GetMeUnreadCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取我的未读分组数
func NewGetMeUnreadCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMeUnreadCountLogic {
	return &GetMeUnreadCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMeUnreadCountLogic) GetMeUnreadCount(req *types.MeUnreadCountReq) (resp *types.MeUnreadCountResp, err error) {
	out, err := l.svcCtx.NotificationService.CountUnreadNotifications(l.ctx,
		&notificationservice.CountUnreadNotificationsRequest{Groups: req.Groups})
	if err != nil {
		return nil, err
	}
	return &types.MeUnreadCountResp{Total: out.Total, GroupCounts: out.GroupCounts}, nil
}
