package preference

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type UpdateNotifyPreferenceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新通知偏好
func NewUpdateNotifyPreferenceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNotifyPreferenceLogic {
	return &UpdateNotifyPreferenceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateNotifyPreference 只提交开关；mandatory 由服务端按场景目录回填，传了也不生效。
func (l *UpdateNotifyPreferenceLogic) UpdateNotifyPreference(req *types.UpdateNotifyPreferenceReq) (resp *types.NotifyPreferenceListResp, err error) {
	items := make([]*notificationservice.NotificationPreference, 0, len(req.List))
	for _, p := range req.List {
		items = append(items, &notificationservice.NotificationPreference{
			Group:   p.Group,
			Channel: p.Channel,
			Enabled: p.Enabled,
		})
	}

	out, err := l.svcCtx.NotificationService.UpdateNotifyPreferences(l.ctx,
		&notificationservice.UpdateNotifyPreferencesRequest{List: items})
	if err != nil {
		return nil, err
	}

	list := make([]types.NotificationPreferenceVO, 0, len(out.List))
	for _, p := range out.List {
		list = append(list, convertPreference(p))
	}
	return &types.NotifyPreferenceListResp{List: list}, nil
}
