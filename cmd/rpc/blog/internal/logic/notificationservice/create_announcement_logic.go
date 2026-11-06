package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/notification"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/model"
)

type CreateAnnouncementLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAnnouncementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAnnouncementLogic {
	return &CreateAnnouncementLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建公告
func (l *CreateAnnouncementLogic) CreateAnnouncement(in *notificationrpc.CreateAnnouncementRequest) (*notificationrpc.CreateAnnouncementResponse, error) {
	if err := validateAudience(in.TargetType, in.TargetIds); err != nil {
		return nil, err
	}
	if in.Title == "" {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "标题不能为空")
	}

	entity := &model.TAnnouncement{
		Title:      in.Title,
		Content:    strPtr(in.Content),
		Category:   in.Category,
		Level:      in.Level,
		TargetType: in.TargetType,
		TargetIds:  strPtr(in.TargetIds),
		Status:     notification.AnnouncementDraft,
	}
	if _, err := l.svcCtx.TAnnouncementModel.Insert(l.ctx, entity); err != nil {
		return nil, err
	}

	return &notificationrpc.CreateAnnouncementResponse{
		Id:           entity.Id,
		Announcement: convertAnnouncement(entity),
	}, nil
}
