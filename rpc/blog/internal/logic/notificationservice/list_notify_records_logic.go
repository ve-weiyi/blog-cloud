package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/queryx"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type ListNotifyRecordsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListNotifyRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNotifyRecordsLogic {
	return &ListNotifyRecordsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListNotifyRecordsLogic) ListNotifyRecords(in *notificationrpc.ListNotifyRecordsRequest) (*notificationrpc.ListNotifyRecordsResponse, error) {
	var opts []queryx.Option
	if in.ListQuery != nil {
		opts = append(opts, queryx.WithPage(int(in.ListQuery.Page)))
		opts = append(opts, queryx.WithSize(int(in.ListQuery.PageSize)))
		opts = append(opts, queryx.WithSorts(in.ListQuery.Sorts...))
	}

	if in.Channel != nil && *in.Channel != "" {
		opts = append(opts, queryx.WithCondition("channel = ?", *in.Channel))
	}
	if in.Status != nil && *in.Status != "" {
		opts = append(opts, queryx.WithCondition("status = ?", *in.Status))
	}
	if in.Recipient != nil && *in.Recipient != "" {
		opts = append(opts, queryx.WithCondition("recipient = ?", *in.Recipient))
	}

	page, size, sorts, conditions, params := queryx.NewQueryBuilder(opts...).Build()
	records, total, err := l.svcCtx.TNotifyRecordModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*notificationrpc.NotifyRecord
	for _, v := range records {
		list = append(list, convertTNotifyRecordToProto(v))
	}

	return &notificationrpc.ListNotifyRecordsResponse{
		ListResult: &notificationrpc.ListResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: list,
	}, nil
}
