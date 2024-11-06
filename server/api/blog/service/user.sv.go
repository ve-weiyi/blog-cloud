package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type UserService struct {
	svcCtx *svctx.ServiceContext
}

func NewUserService(svcCtx *svctx.ServiceContext) *UserService {
	return &UserService{
		svcCtx: svcCtx,
	}
}

// 获取用户信息
func (s *UserService) GetUserInfo(reqCtx *request.Context, in *dto.EmptyReq) (out *dto.UserInfoResp, err error) {
	// todo

	return
}

// 获取用户点赞列表
func (s *UserService) GetUserLike(reqCtx *request.Context, in *dto.EmptyReq) (out *dto.UserLikeResp, err error) {
	// todo

	return
}

// 修改用户头像
func (s *UserService) UpdateUserAvatar(reqCtx *request.Context, in *dto.UpdateUserAvatarReq) (out *dto.EmptyResp, err error) {
	// todo

	return
}

// 修改用户信息
func (s *UserService) UpdateUserInfo(reqCtx *request.Context, in *dto.UpdateUserInfoReq) (out *dto.EmptyResp, err error) {
	// todo

	return
}
