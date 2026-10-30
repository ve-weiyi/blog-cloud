package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
)

type RevokeAnnouncementLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRevokeAnnouncementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeAnnouncementLogic {
	return &RevokeAnnouncementLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 撤回公告
func (l *RevokeAnnouncementLogic) RevokeAnnouncement(in *notificationrpc.RevokeAnnouncementRequest) (*notificationrpc.RevokeAnnouncementResponse, error) {
	revoked, err := l.svcCtx.Notify.RevokeAnnouncement(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	entity, err := l.svcCtx.TAnnouncementModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, bizerr.NewBizError(bizcode.CodeResourceNotFound, "公告不存在")
	}

	return &notificationrpc.RevokeAnnouncementResponse{
		Success:      revoked,
		Announcement: convertAnnouncement(entity),
	}, nil
}
