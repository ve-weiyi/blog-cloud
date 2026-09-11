package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/pb/notificationrpc"

	notificationrpc2 "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type DeleteNotifyTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteNotifyTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNotifyTemplateLogic {
	return &DeleteNotifyTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除通知模板
func (l *DeleteNotifyTemplateLogic) DeleteNotifyTemplate(in *notificationrpc2.DeleteNotifyTemplateRequest) (*notificationrpc2.DeleteNotifyTemplateResponse, error) {
	_, err := l.svcCtx.TNotifyTemplateModel.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &notificationrpc.DeleteNotifyTemplateResponse{
		Success: true,
	}, nil
}
