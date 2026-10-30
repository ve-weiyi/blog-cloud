package announcement

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type UpdateAnnouncementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新公告
func NewUpdateAnnouncementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAnnouncementLogic {
	return &UpdateAnnouncementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAnnouncementLogic) UpdateAnnouncement(req *types.UpdateAnnouncementReq) (resp *types.AnnouncementVO, err error) {
	out, err := l.svcCtx.NotificationService.UpdateAnnouncement(l.ctx, &notificationservice.UpdateAnnouncementRequest{
		Id:         req.Id,
		Title:      req.Title,
		Content:    req.Content,
		Category:   req.Category,
		Level:      req.Level,
		TargetType: req.TargetType,
		TargetIds:  req.TargetIds,
	})
	if err != nil {
		return nil, err
	}
	return convertAnnouncement(out.Announcement), nil
}
