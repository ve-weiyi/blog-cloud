package announcement

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type QueryAnnouncementListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取公告列表
func NewQueryAnnouncementListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryAnnouncementListLogic {
	return &QueryAnnouncementListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryAnnouncementListLogic) QueryAnnouncementList(req *types.QueryAnnouncementListReq) (resp *types.ListResult, err error) {
	out, err := l.svcCtx.NotificationService.ListAnnouncements(l.ctx, &notificationservice.ListAnnouncementsRequest{
		ListQuery:  &notificationservice.ListQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
		Category:   req.Category,
		Level:      req.Level,
		Status:     req.Status,
		TargetType: req.TargetType,
	})
	if err != nil {
		return nil, err
	}

	list := make([]*types.AnnouncementVO, 0, len(out.List))
	for _, item := range out.List {
		list = append(list, convertAnnouncement(item))
	}

	return &types.ListResult{
		Page:     out.GetListResult().GetPage(),
		PageSize: out.GetListResult().GetPageSize(),
		Total:    out.GetListResult().GetTotal(),
		List:     list,
	}, nil
}
