package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
)

type GetAnnouncementLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAnnouncementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAnnouncementLogic {
	return &GetAnnouncementLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询公告详情
func (l *GetAnnouncementLogic) GetAnnouncement(in *notificationrpc.GetAnnouncementRequest) (*notificationrpc.GetAnnouncementResponse, error) {
	row, err := l.svcCtx.TAnnouncementModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, bizerr.NewBizError(bizcode.CodeResourceNotFound, "公告不存在")
	}

	return &notificationrpc.GetAnnouncementResponse{Announcement: convertAnnouncement(row)}, nil
}
