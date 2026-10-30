package announcement

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type CreateAnnouncementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建公告
func NewCreateAnnouncementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAnnouncementLogic {
	return &CreateAnnouncementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAnnouncementLogic) CreateAnnouncement(req *types.CreateAnnouncementReq) (resp *types.AnnouncementVO, err error) {
	out, err := l.svcCtx.NotificationService.CreateAnnouncement(l.ctx, &notificationservice.CreateAnnouncementRequest{
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
