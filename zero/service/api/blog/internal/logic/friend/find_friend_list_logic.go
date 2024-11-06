package friend

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/types"
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

func (l *FindFriendListLogic) FindFriendList(req *types.FriendQueryReq) (resp *types.PageResp, err error) {
	in := &friendrpc.FindFriendListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Sorts:    req.Sorts,
	}
	out, err := l.svcCtx.FriendRpc.FindFriendList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	list := make([]*types.Friend, 0)
	for _, v := range out.List {
		list = append(list, ConvertFriendTypes(v))
	}

	resp = &types.PageResp{}
	resp.Page = in.Page
	resp.PageSize = in.PageSize
	resp.Total = out.Total
	resp.List = list
	return resp, nil
}

func ConvertFriendTypes(req *friendrpc.FriendDetails) (out *types.Friend) {
	return &types.Friend{
		Id:          req.Id,
		LinkName:    req.LinkName,
		LinkAvatar:  req.LinkAvatar,
		LinkAddress: req.LinkAddress,
		LinkIntro:   req.LinkIntro,
		CreatedAt:   req.CreatedAt,
		UpdatedAt:   req.UpdatedAt,
	}
}
