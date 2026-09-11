package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	notificationrpc2 "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/adapter/gormx/queryx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/pb/notificationrpc"
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
func (l *ListNotifyTemplatesLogic) ListNotifyTemplates(in *notificationrpc2.ListNotifyTemplatesRequest) (*notificationrpc2.ListNotifyTemplatesResponse, error) {
	var opts []queryx.Option
	if in.PageQuery != nil {
		opts = append(opts, queryx.WithPage(int(in.PageQuery.Page)))
		opts = append(opts, queryx.WithSize(int(in.PageQuery.PageSize)))
		opts = append(opts, queryx.WithSorts(in.PageQuery.Sorts...))
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
		PageResult: &notificationrpc.PageResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		Templates: list,
	}, nil
}
