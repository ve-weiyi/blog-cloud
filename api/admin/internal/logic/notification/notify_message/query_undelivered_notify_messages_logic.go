package notify_message

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type QueryUndeliveredNotifyMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询已发布但未投递的通知消息
func NewQueryUndeliveredNotifyMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryUndeliveredNotifyMessagesLogic {
	return &QueryUndeliveredNotifyMessagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// QueryUndeliveredNotifyMessages 列出已发布但尚未生成投递记录的消息。
//
// 这类消息对读者不可见且不产生任何报错，是"发布与投递之间中断"留下的静默失效。
// 查出来后用 RedeliverNotifyMessage 修补。
func (l *QueryUndeliveredNotifyMessagesLogic) QueryUndeliveredNotifyMessages(req *types.QueryUndeliveredNotifyMessagesReq) (resp *types.ListResult, err error) {
	// 不传与传 0 同义：都由 rpc 侧取默认值
	var beforeSeconds, limit int64
	if req.BeforeSeconds != nil {
		beforeSeconds = *req.BeforeSeconds
	}
	if req.Limit != nil {
		limit = *req.Limit
	}

	out, err := l.svcCtx.NotificationService.ListUndeliveredNotifyMessages(l.ctx, &notificationservice.ListUndeliveredNotifyMessagesRequest{
		BeforeSeconds: beforeSeconds,
		Limit:         limit,
	})
	if err != nil {
		return nil, err
	}

	list := make([]*types.NotifyMessageVO, 0, len(out.List))
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
		Page:     1,
		PageSize: int64(len(list)),
		Total:    int64(len(list)),
		List:     list,
	}, nil
}
