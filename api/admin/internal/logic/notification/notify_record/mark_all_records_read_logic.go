package notify_record

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type MarkAllRecordsReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 全部标记投递记录为已读
func NewMarkAllRecordsReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkAllRecordsReadLogic {
	return &MarkAllRecordsReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkAllRecordsReadLogic) MarkAllRecordsRead(req *types.MarkAllRecordsReadReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.NotificationService.MarkAllRecordsRead(l.ctx, &notificationservice.MarkAllRecordsReadRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
