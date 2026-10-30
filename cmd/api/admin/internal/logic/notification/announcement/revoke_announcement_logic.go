package announcement

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type RevokeAnnouncementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 撤回公告
func NewRevokeAnnouncementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeAnnouncementLogic {
	return &RevokeAnnouncementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RevokeAnnouncementLogic) RevokeAnnouncement(req *types.RevokeAnnouncementReq) (resp *types.AnnouncementVO, err error) {
	out, err := l.svcCtx.NotificationService.RevokeAnnouncement(l.ctx,
		&notificationservice.RevokeAnnouncementRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return convertAnnouncement(out.Announcement), nil
}
