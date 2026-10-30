package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
)

type ListUserNotificationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUserNotificationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserNotificationsLogic {
	return &ListUserNotificationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 按用户查询收件箱条目
//
// 只读：管理员不代用户改读状态。这里带 user_id 是因为它是排障入口，
// 授权由管理端的权限点负责。
func (l *ListUserNotificationsLogic) ListUserNotifications(in *notificationrpc.ListUserNotificationsRequest) (*notificationrpc.ListUserNotificationsResponse, error) {
	if in.UserId == "" {
		return nil, bizerr.NewBizError(bizcode.CodeParamMissing, "缺少用户标识")
	}

	page, size := pageOf(in.ListQuery)
	items, total, err := l.svcCtx.Notify.ListItems(l.ctx, in.UserId, page, size, "")
	if err != nil {
		return nil, err
	}

	out := make([]*notificationrpc.NotificationItem, 0, len(items))
	for _, item := range items {
		out = append(out, convertItem(item))
	}

	return &notificationrpc.ListUserNotificationsResponse{
		ListResult: &notificationrpc.ListResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: out,
	}, nil
}
