package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchDeleteNotifyTemplatesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteNotifyTemplatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteNotifyTemplatesLogic {
	return &BatchDeleteNotifyTemplatesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除通知模板
func (l *BatchDeleteNotifyTemplatesLogic) BatchDeleteNotifyTemplates(in *notificationrpc.BatchDeleteNotifyTemplatesRequest) (*notificationrpc.BatchDeleteNotifyTemplatesResponse, error) {
	var successCount int64
	for _, id := range in.Ids {
		rows, err := l.svcCtx.TNotifyTemplateModel.Delete(l.ctx, id)
		if err != nil {
			return nil, err
		}
		successCount += rows
	}

	return &notificationrpc.BatchDeleteNotifyTemplatesResponse{
		SuccessCount: successCount,
	}, nil
}
