package notify_record

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type BatchMarkRecordsReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量标记投递记录为已读
func NewBatchMarkRecordsReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchMarkRecordsReadLogic {
	return &BatchMarkRecordsReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchMarkRecordsReadLogic) BatchMarkRecordsRead(req *types.BatchMarkRecordsReadReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.NotificationService.BatchMarkRecordsRead(l.ctx, &notificationservice.BatchMarkRecordsReadRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
