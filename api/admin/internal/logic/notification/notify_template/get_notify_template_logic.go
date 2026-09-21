package notify_template

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type GetNotifyTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取通知模板详情
func NewGetNotifyTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotifyTemplateLogic {
	return &GetNotifyTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNotifyTemplateLogic) GetNotifyTemplate(req *types.GetNotifyTemplateReq) (resp *types.NotifyTemplateVO, err error) {
	out, err := l.svcCtx.NotificationService.GetNotifyTemplate(l.ctx, &notificationservice.GetNotifyTemplateRequest{
		Id: req.Id,
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
