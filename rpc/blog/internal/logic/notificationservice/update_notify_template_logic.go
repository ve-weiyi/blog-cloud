package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/pb/notificationrpc"

	notificationrpc2 "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type UpdateNotifyTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateNotifyTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNotifyTemplateLogic {
	return &UpdateNotifyTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新通知模板
func (l *UpdateNotifyTemplateLogic) UpdateNotifyTemplate(in *notificationrpc2.UpdateNotifyTemplateRequest) (*notificationrpc2.UpdateNotifyTemplateResponse, error) {
	template, err := l.svcCtx.TNotifyTemplateModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	template.Code = in.Code
	template.Channel = in.Channel
	template.Scene = in.Scene
	template.Title = in.Title
	template.Content = in.Content
	template.Enabled = in.Enabled

	_, err = l.svcCtx.TNotifyTemplateModel.Update(l.ctx, template)
	if err != nil {
		return nil, err
	}

	return &notificationrpc.UpdateNotifyTemplateResponse{
		Success: true,
	}, nil
}
