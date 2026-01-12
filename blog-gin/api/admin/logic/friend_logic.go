package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type FriendLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewFriendLogic(svcCtx *svctx.ServiceContext) *FriendLogic {
	return &FriendLogic{
		svcCtx: svcCtx,
	}
}

// 创建友链
func (s *FriendLogic) AddFriend(reqCtx *request.Context, in *types.NewFriendReq) (out *types.FriendBackVO, err error) {
	// todo

	return
}

// 删除友链
func (s *FriendLogic) DeletesFriend(reqCtx *request.Context, in *types.IdsReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 分页获取友链列表
func (s *FriendLogic) FindFriendList(reqCtx *request.Context, in *types.QueryFriendReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 更新友链
func (s *FriendLogic) UpdateFriend(reqCtx *request.Context, in *types.NewFriendReq) (out *types.FriendBackVO, err error) {
	// todo

	return
}
