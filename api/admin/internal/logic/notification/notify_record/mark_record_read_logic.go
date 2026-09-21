package notify_record

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type MarkRecordReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 标记单条投递记录为已读
func NewMarkRecordReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkRecordReadLogic {
	return &MarkRecordReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkRecordReadLogic) MarkRecordRead(req *types.MarkRecordReadReq) (resp *types.BatchResp, err error) {
	_, err = l.svcCtx.NotificationService.MarkRecordRead(l.ctx, &notificationservice.MarkRecordReadRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{}, nil
}
