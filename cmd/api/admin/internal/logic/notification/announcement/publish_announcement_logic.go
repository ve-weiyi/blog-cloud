package announcement

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type PublishAnnouncementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发布公告
func NewPublishAnnouncementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishAnnouncementLogic {
	return &PublishAnnouncementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishAnnouncementLogic) PublishAnnouncement(req *types.PublishAnnouncementReq) (resp *types.AnnouncementVO, err error) {
	out, err := l.svcCtx.NotificationService.PublishAnnouncement(l.ctx,
		&notificationservice.PublishAnnouncementRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return convertAnnouncement(out.Announcement), nil
}
