package socialservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/adapter/gormx/queryx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/pb/socialrpc"
)

type ListTalksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListTalksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTalksLogic {
	return &ListTalksLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListTalksLogic) ListTalks(in *socialrpc.ListTalksRequest) (*socialrpc.ListTalksResponse, error) {
	var opts []queryx.Option
	if in.PageQuery != nil {
		opts = append(opts, queryx.WithPage(int(in.PageQuery.Page)))
		opts = append(opts, queryx.WithSize(int(in.PageQuery.PageSize)))
		opts = append(opts, queryx.WithSorts(in.PageQuery.Sorts...))
	}
	if in.Status != nil {
		opts = append(opts, queryx.WithCondition("status = ?", *in.Status))
	}

	page, size, sorts, conditions, params := queryx.NewQueryBuilder(opts...).Build()
	records, total, err := l.svcCtx.TTalkModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*socialrpc.Talk
	for _, v := range records {
		list = append(list, convertTalkOut(v))
	}

	return &socialrpc.ListTalksResponse{
		PageResult: &socialrpc.PageResult{Page: int64(page), PageSize: int64(size), Total: total},
		List:       list,
	}, nil
}
