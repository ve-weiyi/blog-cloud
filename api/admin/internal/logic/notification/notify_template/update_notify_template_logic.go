package notify_template

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type UpdateNotifyTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新通知模板
func NewUpdateNotifyTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNotifyTemplateLogic {
	return &UpdateNotifyTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateNotifyTemplateLogic) UpdateNotifyTemplate(req *types.UpdateNotifyTemplateReq) (resp *types.NotifyTemplateVO, err error) {
	_, err = l.svcCtx.NotificationService.UpdateNotifyTemplate(l.ctx, &notificationservice.UpdateNotifyTemplateRequest{
		Id:      req.Id,
		Code:    req.Code,
		Channel: req.Channel,
		Scene:   req.Scene,
		Title:   req.Title,
		Content: req.Content,
		Enabled: req.Enabled,
	})
	if err != nil {
		return nil, err
	}

	detail, err := l.svcCtx.NotificationService.GetNotifyTemplate(l.ctx, &notificationservice.GetNotifyTemplateRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &types.NotifyTemplateVO{
		Id:        detail.Template.Id,
		Code:      detail.Template.Code,
		Channel:   detail.Template.Channel,
		Scene:     detail.Template.Scene,
		Title:     detail.Template.Title,
		Content:   detail.Template.Content,
		Enabled:   detail.Template.Enabled,
		CreatedAt: detail.Template.CreatedAt,
		UpdatedAt: detail.Template.UpdatedAt,
	}, nil
}
