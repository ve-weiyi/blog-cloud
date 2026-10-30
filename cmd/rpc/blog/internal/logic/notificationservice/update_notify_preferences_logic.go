package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/notification"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
)

type UpdateNotifyPreferencesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateNotifyPreferencesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNotifyPreferencesLogic {
	return &UpdateNotifyPreferencesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新通知偏好
//
// 请求里的 mandatory 一律忽略：它由场景目录决定，客户端传什么都不会改变强制分组。
func (l *UpdateNotifyPreferencesLogic) UpdateNotifyPreferences(in *notificationrpc.UpdateNotifyPreferencesRequest) (*notificationrpc.UpdateNotifyPreferencesResponse, error) {
	uid, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, "未登录")
	}

	items := make([]notification.Preference, 0, len(in.List))
	for _, p := range in.List {
		items = append(items, notification.Preference{
			Group:   notify.Group(p.Group),
			Channel: notify.Channel(p.Channel),
			Enabled: p.Enabled,
		})
	}

	list, err := l.svcCtx.Notify.UpdatePreferences(l.ctx, uid, items)
	if err != nil {
		return nil, err
	}

	out := make([]*notificationrpc.NotificationPreference, 0, len(list))
	for _, p := range list {
		out = append(out, convertPreference(p))
	}
	return &notificationrpc.UpdateNotifyPreferencesResponse{List: out}, nil
}
