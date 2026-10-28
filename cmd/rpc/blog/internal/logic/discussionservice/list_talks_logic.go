package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/store/queryx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/discussionrpc"
)

type ListTalksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListTalksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTalksLogic {
	return &ListTalksLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListTalksLogic) ListTalks(in *discussionrpc.ListTalksRequest) (*discussionrpc.ListTalksResponse, error) {
	var opts []queryx.Option
	if in.ListQuery != nil {
		opts = append(opts, queryx.WithPage(int(in.ListQuery.Page)))
		opts = append(opts, queryx.WithSize(int(in.ListQuery.PageSize)))
		opts = append(opts, queryx.WithSorts(in.ListQuery.Sorts...))
	}
	if in.Status != nil {
		opts = append(opts, queryx.WithCondition("status = ?", *in.Status))
	}

	page, size, sorts, conditions, params := queryx.NewQueryBuilder(opts...).Build()
	records, total, err := l.svcCtx.TTalkModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*discussionrpc.Talk
	ids := make([]int64, 0, len(records))
	for _, v := range records {
		ids = append(ids, v.Id)
	}
	lm, lcm, err := likeStats(l.ctx, l.svcCtx, enums.LikeTargetTalk, ids)
	if err != nil {
		return nil, err
	}

	for _, v := range records {
		list = append(list, convertTalkOut(v, lm, lcm))
	}

	return &discussionrpc.ListTalksResponse{
		ListResult: &discussionrpc.ListResult{Page: int64(page), PageSize: int64(size), Total: total},
		List:       list,
	}, nil
}
