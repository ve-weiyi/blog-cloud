// Code scaffolded by goctl. Safe to edit.
package notify_stream

import (
	"context"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/sse"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AnnouncementStreamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 公告状态推送流（SSE）
func NewAnnouncementStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AnnouncementStreamLogic {
	return &AnnouncementStreamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// AnnouncementStream 把本实例收到的公告状态变化写进推送流。
//
// 只推「哪条公告变了」，供多管理员协作时刷新列表；读者端的「收件箱变了」是另一条流。
// 返回 nil 表示流正常结束（客户端断开或服务关停）。
func (l *AnnouncementStreamLogic) AnnouncementStream(client chan<- *types.AnnouncementStreamEvent) error {
	return sse.Bridge(l.ctx, l.svcCtx.NotifyBroker, client)
}
