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

type MarkNotificationsReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkNotificationsReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkNotificationsReadLogic {
	return &MarkNotificationsReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 标记通知已读
func (l *MarkNotificationsReadLogic) MarkNotificationsRead(in *notificationrpc.MarkNotificationsReadRequest) (*notificationrpc.MarkNotificationsReadResponse, error) {
	uid, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, "未登录")
	}
	if in.GroupKey == "" {
		return nil, bizerr.NewBizError(bizcode.CodeParamMissing, "缺少分组标识")
	}

	rows, err := l.svcCtx.Notify.MarkRead(l.ctx, uid, in.GroupKey)
	if err != nil {
		return nil, err
	}
	return &notificationrpc.MarkNotificationsReadResponse{SuccessCount: rows}, nil
}
