package socialrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/query"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/socialrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindFriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindFriendListLogic {
	return &FindFriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询友链列表
func (l *FindFriendListLogic) FindFriendList(in *socialrpc.FindFriendListReq) (*socialrpc.FindFriendListResp, error) {
	page, size, sorts, conditions, params := convertFriendQuery(in)

	records, total, err := l.svcCtx.TFriendModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*socialrpc.Friend
	for _, v := range records {
		list = append(list, convertFriendOut(v))
	}

	return &socialrpc.FindFriendListResp{
		List: list,
		Pagination: &socialrpc.PageResp{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
	}, nil
}

func convertFriendQuery(in *socialrpc.FindFriendListReq) (page int, size int, sorts string, conditions string, params []any) {
	var opts []query.Option
	if in.Paginate != nil {
		opts = append(opts, query.WithPage(int(in.Paginate.Page)))
		opts = append(opts, query.WithSize(int(in.Paginate.PageSize)))
		opts = append(opts, query.WithSorts(in.Paginate.Sorts...))
	}

	if in.LinkName != "" {
		opts = append(opts, query.WithCondition("link_name like ?", "%"+in.LinkName+"%"))
	}

	return query.NewQueryBuilder(opts...).Build()
}
