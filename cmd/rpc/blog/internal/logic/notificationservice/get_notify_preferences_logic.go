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

type GetNotifyPreferencesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNotifyPreferencesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotifyPreferencesLogic {
	return &GetNotifyPreferencesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询通知偏好
func (l *GetNotifyPreferencesLogic) GetNotifyPreferences(in *notificationrpc.GetNotifyPreferencesRequest) (*notificationrpc.GetNotifyPreferencesResponse, error) {
	uid, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, "未登录")
	}

	list, err := l.svcCtx.Notify.ListPreferences(l.ctx, uid)
	if err != nil {
		return nil, err
	}

	out := make([]*notificationrpc.NotificationPreference, 0, len(list))
	for _, p := range list {
		out = append(out, convertPreference(p))
	}
	return &notificationrpc.GetNotifyPreferencesResponse{List: out}, nil
}
