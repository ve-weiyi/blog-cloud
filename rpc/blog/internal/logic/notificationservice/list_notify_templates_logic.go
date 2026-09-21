package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/queryx"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type ListNotifyTemplatesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListNotifyTemplatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNotifyTemplatesLogic {
	return &ListNotifyTemplatesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页查询通知模板列表
func (l *ListNotifyTemplatesLogic) ListNotifyTemplates(in *notificationrpc.ListNotifyTemplatesRequest) (*notificationrpc.ListNotifyTemplatesResponse, error) {
	var opts []queryx.Option
	if in.ListQuery != nil {
		opts = append(opts, queryx.WithPage(int(in.ListQuery.Page)))
		opts = append(opts, queryx.WithSize(int(in.ListQuery.PageSize)))
		opts = append(opts, queryx.WithSorts(in.ListQuery.Sorts...))
	}

	// 可选过滤条件
	if in.Channel != nil {
		opts = append(opts, queryx.WithCondition("channel = ?", in.Channel))
	}

	page, size, sorts, conditions, params := queryx.NewQueryBuilder(opts...).Build()
	records, total, err := l.svcCtx.TNotifyTemplateModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*notificationrpc.NotifyTemplate
	for _, v := range records {
		list = append(list, convertTNotifyTemplateToProto(v))
	}

	return &notificationrpc.ListNotifyTemplatesResponse{
		ListResult: &notificationrpc.ListResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: list,
	}, nil
}
