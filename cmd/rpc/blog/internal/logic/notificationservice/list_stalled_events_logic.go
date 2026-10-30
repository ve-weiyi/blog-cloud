package notificationservicelogic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/notification"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/model"
)

const (
	stalledDefaultSeconds = 300
	stalledDefaultLimit   = 100
)

type ListStalledEventsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListStalledEventsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListStalledEventsLogic {
	return &ListStalledEventsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询停滞的通知事件
//
// 停滞 = 停在扇出中超过阈值，或已失败。发布与投递之间中断不会自动恢复，也不报错，
// 这是它的排查入口。
func (l *ListStalledEventsLogic) ListStalledEvents(in *notificationrpc.ListStalledEventsRequest) (*notificationrpc.ListStalledEventsResponse, error) {
	seconds := in.BeforeSeconds
	if seconds <= 0 {
		seconds = stalledDefaultSeconds
	}
	limit := in.Limit
	if limit <= 0 || limit > stalledDefaultLimit {
		limit = stalledDefaultLimit
	}

	before := time.Now().Add(-time.Duration(seconds) * time.Second)

	rows := make([]*model.TNotificationEvent, 0, limit)
	err := l.svcCtx.GormDB.WithContext(l.ctx).
		Table(l.svcCtx.TNotificationEventModel.TableName()).
		Where("status IN ? AND updated_at < ?",
			[]string{notification.StatusFanning, notification.StatusFailed}, before).
		Order("id ASC").
		Limit(int(limit)).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	out := make([]*notificationrpc.NotificationEvent, 0, len(rows))
	for _, item := range rows {
		out = append(out, convertEvent(item))
	}
	return &notificationrpc.ListStalledEventsResponse{List: out}, nil
}
