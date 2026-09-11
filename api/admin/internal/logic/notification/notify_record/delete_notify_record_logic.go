package notify_record

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type DeleteNotifyRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除投递记录
func NewDeleteNotifyRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNotifyRecordLogic {
	return &DeleteNotifyRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteNotifyRecordLogic) DeleteNotifyRecord(req *types.DeleteNotifyRecordReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.NotificationService.DeleteNotifyRecords(l.ctx, &notificationservice.DeleteNotifyRecordsRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
