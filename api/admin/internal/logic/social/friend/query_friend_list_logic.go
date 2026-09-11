package friend

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
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
	out, err := l.svcCtx.SocialService.ListFriends(l.ctx, &socialservice.ListFriendsRequest{
		PageQuery: &socialservice.PageQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
		LinkName:  req.LinkName,
	})
	if err != nil {
		return nil, err
	}

	var list []*types.FriendVO
	for _, v := range out.List {
		list = append(list, &types.FriendVO{
			Id:          v.Id,
			LinkName:    v.LinkName,
			LinkAvatar:  v.LinkAvatar,
			LinkAddress: v.LinkAddress,
			LinkIntro:   v.LinkIntro,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}

	return &types.PageResult{
		Page:     out.PageResult.Page,
		PageSize: out.PageResult.PageSize,
		Total:    out.PageResult.Total,
		List:     list,
	}, nil
}
