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

type DeleteNotificationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteNotificationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNotificationsLogic {
	return &DeleteNotificationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除通知
func (l *DeleteNotificationsLogic) DeleteNotifications(in *notificationrpc.DeleteNotificationsRequest) (*notificationrpc.DeleteNotificationsResponse, error) {
	uid, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, "未登录")
	}
	if in.GroupKey == "" {
		return nil, bizerr.NewBizError(bizcode.CodeParamMissing, "缺少分组标识")
	}

	rows, err := l.svcCtx.Notify.DeleteGroup(l.ctx, uid, in.GroupKey)
	if err != nil {
		return nil, err
	}
	return &notificationrpc.DeleteNotificationsResponse{SuccessCount: rows}, nil
}
