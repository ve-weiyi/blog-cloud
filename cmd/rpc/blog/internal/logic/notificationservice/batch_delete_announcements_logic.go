package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
)

type BatchDeleteAnnouncementsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteAnnouncementsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteAnnouncementsLogic {
	return &BatchDeleteAnnouncementsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除公告（连同派生条目与投递记录）
func (l *BatchDeleteAnnouncementsLogic) BatchDeleteAnnouncements(in *notificationrpc.BatchDeleteAnnouncementsRequest) (*notificationrpc.BatchDeleteAnnouncementsResponse, error) {
	affected, err := l.svcCtx.Notify.DeleteAnnouncements(l.ctx, in.Ids)
	if err != nil {
		return nil, err
	}
	return &notificationrpc.BatchDeleteAnnouncementsResponse{SuccessCount: affected}, nil
}
