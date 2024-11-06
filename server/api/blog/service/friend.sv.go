package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type FriendService struct {
	svcCtx *svctx.ServiceContext
}

func NewFriendService(svcCtx *svctx.ServiceContext) *FriendService {
	return &FriendService{
		svcCtx: svcCtx,
	}
}

// 分页获取友链列表
func (s *FriendService) FindFriendList(reqCtx *request.Context, in *dto.FriendQueryReq) (out *dto.PageResp, err error) {
	// todo

	return
}
