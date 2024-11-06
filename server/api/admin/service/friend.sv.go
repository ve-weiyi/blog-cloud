package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
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
func (s *FriendService) FindFriendList(reqCtx *request.Context, in *dto.FriendQuery) (out *dto.PageResp, err error) {
	// todo

	return
}

// 创建友链
func (s *FriendService) AddFriend(reqCtx *request.Context, in *dto.FriendNewReq) (out *dto.FriendBackDTO, err error) {
	// todo

	return
}

// 批量删除友链
func (s *FriendService) BatchDeleteFriend(reqCtx *request.Context, in *dto.IdsReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 删除友链
func (s *FriendService) DeleteFriend(reqCtx *request.Context, in *dto.IdReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 更新友链
func (s *FriendService) UpdateFriend(reqCtx *request.Context, in *dto.FriendNewReq) (out *dto.FriendBackDTO, err error) {
	// todo

	return
}
