// Code scaffolded by goctl. Safe to edit.
package notify_stream

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/common/sse"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
)

type NotifyStreamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 通知消息推送流（SSE）
func NewNotifyStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyStreamLogic {
	return &NotifyStreamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// NotifyStream 把本实例收到的通知事件写进推送流。
//
// 只推「哪条消息变了」，不带收件人——客户端收到后按自身身份回查未读列表。
// 返回 nil 表示流正常结束（客户端断开或服务关停）。桥接细节见 sse.Bridge。
func (l *NotifyStreamLogic) NotifyStream(client chan<- *types.NotifyStreamEvent) error {
	return sse.Bridge(l.ctx, l.svcCtx.NotifyBroker, client)
}
