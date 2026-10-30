package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
)

type ListAnnouncementsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAnnouncementsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAnnouncementsLogic {
	return &ListAnnouncementsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页查询公告列表
func (l *ListAnnouncementsLogic) ListAnnouncements(in *notificationrpc.ListAnnouncementsRequest) (*notificationrpc.ListAnnouncementsResponse, error) {
	page, size := pageOf(in.ListQuery)

	var w where
	w.eq("category", in.Category)
	w.eq("level", in.Level)
	w.eq("status", in.Status)
	w.eq("target_type", in.TargetType)
	conditions, args := w.build()

	list, total, err := l.svcCtx.TAnnouncementModel.FindListAndTotal(
		l.ctx, page, size, "id DESC", conditions, args...)
	if err != nil {
		return nil, err
	}

	out := make([]*notificationrpc.Announcement, 0, len(list))
	for _, item := range list {
		out = append(out, convertAnnouncement(item))
	}

	return &notificationrpc.ListAnnouncementsResponse{
		ListResult: &notificationrpc.ListResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: out,
	}, nil
}
