package notify_template

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type CreateNotifyTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建通知模板
func NewCreateNotifyTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNotifyTemplateLogic {
	return &CreateNotifyTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateNotifyTemplateLogic) CreateNotifyTemplate(req *types.CreateNotifyTemplateReq) (resp *types.NotifyTemplateVO, err error) {
	out, err := l.svcCtx.NotificationService.CreateNotifyTemplate(l.ctx, &notificationservice.CreateNotifyTemplateRequest{
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

	return &types.NotifyTemplateVO{
		Id:        out.Template.Id,
		Code:      out.Template.Code,
		Channel:   out.Template.Channel,
		Scene:     out.Template.Scene,
		Title:     out.Template.Title,
		Content:   out.Template.Content,
		Enabled:   out.Template.Enabled,
		CreatedAt: out.Template.CreatedAt,
		UpdatedAt: out.Template.UpdatedAt,
	}, nil
}
