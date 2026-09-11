package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	notificationrpc2 "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/adapter/gormx/queryx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/pb/notificationrpc"
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

func (l *ListNotifyRecordsLogic) ListNotifyRecords(in *notificationrpc2.ListNotifyRecordsRequest) (*notificationrpc2.ListNotifyRecordsResponse, error) {
	var opts []queryx.Option
	if in.PageQuery != nil {
		opts = append(opts, queryx.WithPage(int(in.PageQuery.Page)))
		opts = append(opts, queryx.WithSize(int(in.PageQuery.PageSize)))
		opts = append(opts, queryx.WithSorts(in.PageQuery.Sorts...))
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
		PageResult: &notificationrpc.PageResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		Records: list,
	}, nil
}
