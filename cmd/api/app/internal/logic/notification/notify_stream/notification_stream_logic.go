package notify_stream

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/sse"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
)

type NotificationStreamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 通知实时推送流（SSE）
func NewNotificationStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotificationStreamLogic {
	return &NotificationStreamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// NotificationStream 订阅当前登录用户的收件箱信号。
//
// 只推「收件箱变了」，不带正文；客户端收到后回查列表与未读数。
// 返回 nil 表示流正常结束（客户端断开或服务关停）。
func (l *NotificationStreamLogic) NotificationStream(client chan<- *types.NotificationStreamEvent) error {
	uid, err := metax.GetApiUserIdFromCtx(l.ctx)
	if err != nil || uid == "" {
		// 未登录不建立长连接：没有收件人就没有可订阅的频道
		return nil
	}
	return sse.Bridge(l.ctx, l.svcCtx.RedisClient, uid, client)
}
