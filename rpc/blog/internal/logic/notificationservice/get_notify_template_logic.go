package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/pb/notificationrpc"

	notificationrpc2 "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type GetNotifyTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNotifyTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotifyTemplateLogic {
	return &GetNotifyTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据ID查询通知模板
func (l *GetNotifyTemplateLogic) GetNotifyTemplate(in *notificationrpc2.GetNotifyTemplateRequest) (*notificationrpc2.GetNotifyTemplateResponse, error) {
	record, err := l.svcCtx.TNotifyTemplateModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &notificationrpc.GetNotifyTemplateResponse{
		Template: convertTNotifyTemplateToProto(record),
	}, nil
}
