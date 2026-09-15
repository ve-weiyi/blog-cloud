package notify_message

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type BatchDeleteNotifyMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除统一通知消息
func NewBatchDeleteNotifyMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteNotifyMessageLogic {
	return &BatchDeleteNotifyMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteNotifyMessageLogic) BatchDeleteNotifyMessage(req *types.DeleteNotifyMessageReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.NotificationService.BatchDeleteNotifyMessages(l.ctx, &notificationservice.BatchDeleteNotifyMessagesRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
