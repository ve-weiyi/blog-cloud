package friend

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/friendrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindFriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页获取友链列表
func NewFindFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindFriendListLogic {
	return &FindFriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindFriendListLogic) FindFriendList(req *types.FriendQuery) (resp *types.PageResp, err error) {
	in := &friendrpc.FindFriendListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Sorts:    req.Sorts,
		LinkName: req.LinkName,
	}

	out, err := l.svcCtx.FriendRpc.FindFriendList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var list []*types.FriendBackDTO
	for _, v := range out.List {
		m := ConvertFriendTypes(v)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = in.Page
	resp.PageSize = in.PageSize
	resp.Total = out.Total
	resp.List = list
	return resp, nil
}
