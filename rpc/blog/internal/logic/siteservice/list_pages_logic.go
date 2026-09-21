package siteservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/queryx"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/siterpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type ListPagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPagesLogic {
	return &ListPagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPagesLogic) ListPages(in *siterpc.ListPagesRequest) (*siterpc.ListPagesResponse, error) {
	var opts []queryx.Option
	if in.ListQuery != nil {
		opts = append(opts, queryx.WithPage(int(in.ListQuery.Page)))
		opts = append(opts, queryx.WithSize(int(in.ListQuery.PageSize)))
		opts = append(opts, queryx.WithSorts(in.ListQuery.Sorts...))
	}
	if in.PageName != nil {
		opts = append(opts, queryx.WithCondition("page_name like ?", "%"+*in.PageName+"%"))
	}

	page, size, sorts, conditions, params := queryx.NewQueryBuilder(opts...).Build()
	records, total, err := l.svcCtx.TPageModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*siterpc.Page
	for _, v := range records {
		list = append(list, convertPageOut(v))
	}

	return &siterpc.ListPagesResponse{
		ListResult: &siterpc.ListResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: list,
	}, nil
}
