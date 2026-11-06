package preference

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type GetNotifyPreferenceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取通知偏好
func NewGetNotifyPreferenceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotifyPreferenceLogic {
	return &GetNotifyPreferenceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNotifyPreferenceLogic) GetNotifyPreference(req *types.EmptyReq) (resp *types.NotifyPreferenceListResp, err error) {
	out, err := l.svcCtx.NotificationService.GetNotifyPreferences(l.ctx,
		&notificationservice.GetNotifyPreferencesRequest{})
	if err != nil {
		return nil, err
	}

	list := make([]types.NotificationPreferenceVO, 0, len(out.List))
	for _, p := range out.List {
		list = append(list, convertPreference(p))
	}
	return &types.NotifyPreferenceListResp{List: list}, nil
}
