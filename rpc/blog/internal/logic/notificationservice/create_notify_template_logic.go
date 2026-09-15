package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type CreateNotifyTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateNotifyTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNotifyTemplateLogic {
	return &CreateNotifyTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建通知模板
func (l *CreateNotifyTemplateLogic) CreateNotifyTemplate(in *notificationrpc.CreateNotifyTemplateRequest) (*notificationrpc.CreateNotifyTemplateResponse, error) {
	template := convertProtoToTNotifyTemplate(in)

	_, err := l.svcCtx.TNotifyTemplateModel.Insert(l.ctx, template)
	if err != nil {
		return nil, err
	}

	return &notificationrpc.CreateNotifyTemplateResponse{
		Id: template.Id,
	}, nil
}
