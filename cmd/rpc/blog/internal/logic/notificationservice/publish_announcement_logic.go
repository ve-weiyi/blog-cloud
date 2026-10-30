package notificationservicelogic

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/notification"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
)

type PublishAnnouncementLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishAnnouncementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishAnnouncementLogic {
	return &PublishAnnouncementLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 发布公告
func (l *PublishAnnouncementLogic) PublishAnnouncement(in *notificationrpc.PublishAnnouncementRequest) (*notificationrpc.PublishAnnouncementResponse, error) {
	// 发布人要落库并作为触发者：取不到身份说明网关没注入主体，继续写空值会留下无主公告
	publishedBy, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, "未登录")
	}

	published, err := l.svcCtx.Notify.PublishAnnouncement(l.ctx, in.Id, publishedBy)
	if errors.Is(err, notification.ErrAudienceEmpty) {
		return nil, bizerr.NewBizError(bizcode.CodePreconditionFailed, "公告没有定向受众，无法发布")
	}
	if err != nil {
		return nil, err
	}

	entity, err := l.svcCtx.TAnnouncementModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, bizerr.NewBizError(bizcode.CodeResourceNotFound, "公告不存在")
	}

	return &notificationrpc.PublishAnnouncementResponse{
		Success:      published,
		Announcement: convertAnnouncement(entity),
	}, nil
}
