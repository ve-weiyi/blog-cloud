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

type CountUnreadNotificationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCountUnreadNotificationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CountUnreadNotificationsLogic {
	return &CountUnreadNotificationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询未读分组数
func (l *CountUnreadNotificationsLogic) CountUnreadNotifications(in *notificationrpc.CountUnreadNotificationsRequest) (*notificationrpc.CountUnreadNotificationsResponse, error) {
	uid, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, "未登录")
	}

	total, counts, err := l.svcCtx.Notify.CountUnread(l.ctx, uid, groupsOf(in.Groups))
	if err != nil {
		return nil, err
	}

	out := make(map[string]int64, len(counts))
	for group, n := range counts {
		out[string(group)] = n
	}
	return &notificationrpc.CountUnreadNotificationsResponse{
		Total:       total,
		GroupCounts: out,
	}, nil
}
