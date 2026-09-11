package friend

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/socialservice"
)

type QueryFriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取友链列表
func NewQueryFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryFriendListLogic {
	return &QueryFriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryFriendListLogic) QueryFriendList(req *types.QueryFriendListReq) (resp *types.PageResult, err error) {
	in := &socialservice.ListFriendsRequest{
		PageQuery: &socialservice.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
	}

	out, err := l.svcCtx.SocialService.ListFriends(l.ctx, in)
	if err != nil {
		return nil, err
	}

	list := make([]*types.Friend, 0)
	for _, v := range out.List {
		list = append(list, &types.Friend{
			Id:          v.Id,
			LinkName:    v.LinkName,
			LinkAvatar:  v.LinkAvatar,
			LinkAddress: v.LinkAddress,
			LinkIntro:   v.LinkIntro,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}

	resp = &types.PageResult{
		Page:     out.PageResult.Page,
		PageSize: out.PageResult.PageSize,
		Total:    out.PageResult.Total,
		List:     list,
	}
	return
}
