package friend

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/siteservice"
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

func (l *QueryFriendListLogic) QueryFriendList(req *types.QueryFriendListReq) (resp *types.ListResult, err error) {
	in := &siteservice.ListFriendsRequest{
		ListQuery: &siteservice.ListQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
	}

	out, err := l.svcCtx.SiteService.ListFriends(l.ctx, in)
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

	resp = &types.ListResult{
		Page:     out.ListResult.Page,
		PageSize: out.ListResult.PageSize,
		Total:    out.ListResult.Total,
		List:     list,
	}
	return
}
