package notify_record

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type QueryNotifyRecordListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取统一投递记录列表
func NewQueryNotifyRecordListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryNotifyRecordListLogic {
	return &QueryNotifyRecordListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryNotifyRecordListLogic) QueryNotifyRecordList(req *types.QueryNotifyRecordListReq) (resp *types.ListResult, err error) {
	out, err := l.svcCtx.NotificationService.ListNotifyRecords(l.ctx, &notificationservice.ListNotifyRecordsRequest{
		ListQuery: &notificationservice.ListQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
		Channel:   req.Channel,
		Status:    req.Status,
		Recipient: req.Recipient,
	})
	if err != nil {
		return nil, err
	}

	var list []*types.NotifyRecordVO
	for _, v := range out.List {
		list = append(list, &types.NotifyRecordVO{
			Id:           v.Id,
			MessageId:    v.MessageId,
			Channel:      v.Channel,
			Recipient:    v.Recipient,
			TemplateCode: v.TemplateCode,
			Content:      v.Content,
			Status:       v.Status,
			BizId:        v.BizId,
			ErrorMsg:     v.ErrorMsg,
			ReadAt:       v.ReadAt,
			SentAt:       v.SentAt,
			CreatedAt:    v.CreatedAt,
		})
	}

	return &types.ListResult{
		Page:     out.ListResult.Page,
		PageSize: out.ListResult.PageSize,
		Total:    out.ListResult.Total,
		List:     list,
	}, nil
}
