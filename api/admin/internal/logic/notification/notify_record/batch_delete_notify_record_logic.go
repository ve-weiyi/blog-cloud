package notify_record

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type BatchDeleteNotifyRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除投递记录
func NewBatchDeleteNotifyRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteNotifyRecordLogic {
	return &BatchDeleteNotifyRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteNotifyRecordLogic) BatchDeleteNotifyRecord(req *types.DeleteNotifyRecordReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.NotificationService.BatchDeleteNotifyRecords(l.ctx, &notificationservice.BatchDeleteNotifyRecordsRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
