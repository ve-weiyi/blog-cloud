package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
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

// 分页获取友链列表
func (s *FriendLogic) FindFriendList(reqCtx *request.Context, in *types.FriendQueryReq) (out *types.PageResp, err error) {
	// todo

	return
}
