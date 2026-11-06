package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
)

type ListNotificationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListNotificationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNotificationsLogic {
	return &ListNotificationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页查询收件箱（按分组）
//
// 收件人一律取自会话主体：请求里没有 user_id，就没有越权构造的余地。
func (l *ListNotificationsLogic) ListNotifications(in *notificationrpc.ListNotificationsRequest) (*notificationrpc.ListNotificationsResponse, error) {
	uid, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, "未登录")
	}

	page, size := pageOf(in.ListQuery)
	unreadOnly := in.UnreadOnly != nil && *in.UnreadOnly

	groups, total, err := l.svcCtx.Notify.ListGroups(
		l.ctx, uid, groupsOf(in.Groups), unreadOnly, page, size, "")
	if err != nil {
		return nil, err
	}

	out := make([]*notificationrpc.NotificationGroup, 0, len(groups))
	for _, g := range groups {
		out = append(out, convertInboxGroup(g))
	}

	return &notificationrpc.ListNotificationsResponse{
		ListResult: &notificationrpc.ListResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: out,
	}, nil
}
