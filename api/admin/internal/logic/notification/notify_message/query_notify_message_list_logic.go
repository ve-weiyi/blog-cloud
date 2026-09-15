package notify_message

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type QueryNotifyMessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取统一通知消息列表
func NewQueryNotifyMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryNotifyMessageListLogic {
	return &QueryNotifyMessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryNotifyMessageListLogic) QueryNotifyMessageList(req *types.QueryNotifyMessageListReq) (resp *types.ListResult, err error) {
	out, err := l.svcCtx.NotificationService.ListNotifyMessages(l.ctx, &notificationservice.ListNotifyMessagesRequest{
		ListQuery:  &notificationservice.ListQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
		Category:   req.Category,
		Level:      req.Level,
		Status:     req.Status,
		TargetType: req.TargetType,
	})
	if err != nil {
		return nil, err
	}

	var list []*types.NotifyMessageVO
	for _, v := range out.List {
		list = append(list, &types.NotifyMessageVO{
			Id:          v.Id,
			Title:       v.Title,
			Content:     v.Content,
			Category:    v.Category,
			Level:       v.Level,
			TargetType:  v.TargetType,
			TargetIds:   v.TargetIds,
			Status:      v.Status,
			PublishedAt: v.PublishedAt,
			PublishedBy: v.PublishedBy,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}

	return &types.ListResult{
		Page:     out.ListResult.Page,
		PageSize: out.ListResult.PageSize,
		Total:    out.ListResult.Total,
		List:     list,
	}, nil
}
