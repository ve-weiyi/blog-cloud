package friendrpclogic

import (
	"context"
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/friendrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

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
func (l *FindFriendListLogic) FindFriendList(in *friendrpc.FindFriendListReq) (*friendrpc.FindFriendListResp, error) {
	page, size, sorts, conditions, params := convertFriendQuery(in)

	result, err := l.svcCtx.TFriendModel.FindList(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	count, err := l.svcCtx.TFriendModel.FindCount(l.ctx, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*friendrpc.FriendDetails
	for _, v := range result {
		list = append(list, convertFriendOut(v))
	}

	return &friendrpc.FindFriendListResp{
		List:  list,
		Total: count,
	}, nil
}

func convertFriendQuery(in *friendrpc.FindFriendListReq) (page int, size int, sorts string, conditions string, params []any) {
	page = int(in.Page)
	size = int(in.PageSize)
	sorts = strings.Join(in.Sorts, ",")
	if sorts == "" {
		sorts = "id desc"
	}

	if in.LinkName != "" {
		conditions += " link_name like ?"
		params = append(params, "%"+in.LinkName+"%")
	}

	return
}
