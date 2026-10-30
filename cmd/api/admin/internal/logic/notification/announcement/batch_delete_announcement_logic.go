package announcement

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type BatchDeleteAnnouncementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除公告
func NewBatchDeleteAnnouncementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteAnnouncementLogic {
	return &BatchDeleteAnnouncementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteAnnouncementLogic) BatchDeleteAnnouncement(req *types.IdsReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.NotificationService.BatchDeleteAnnouncements(l.ctx,
		&notificationservice.BatchDeleteAnnouncementsRequest{Ids: req.Ids})
	if err != nil {
		return nil, err
	}
	return &types.BatchResp{SuccessCount: out.SuccessCount}, nil
}
