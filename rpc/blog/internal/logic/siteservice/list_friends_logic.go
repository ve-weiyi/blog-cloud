package siteservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/queryx"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/siterpc"
)

type ListFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFriendsLogic {
	return &ListFriendsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListFriendsLogic) ListFriends(in *siterpc.ListFriendsRequest) (*siterpc.ListFriendsResponse, error) {
	var opts []queryx.Option
	if in.ListQuery != nil {
		opts = append(opts, queryx.WithPage(int(in.ListQuery.Page)))
		opts = append(opts, queryx.WithSize(int(in.ListQuery.PageSize)))
		opts = append(opts, queryx.WithSorts(in.ListQuery.Sorts...))
	}
	if in.LinkName != nil {
		opts = append(opts, queryx.WithCondition("link_name like ?", "%"+*in.LinkName+"%"))
	}

	page, size, sorts, conditions, params := queryx.NewQueryBuilder(opts...).Build()
	records, total, err := l.svcCtx.TFriendModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*siterpc.Friend
	for _, v := range records {
		list = append(list, convertFriendOut(v))
	}

	return &siterpc.ListFriendsResponse{
		ListResult: &siterpc.ListResult{Page: int64(page), PageSize: int64(size), Total: total},
		List:       list,
	}, nil
}
