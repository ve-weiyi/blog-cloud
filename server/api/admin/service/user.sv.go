package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
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

// 获取用户接口权限
func (s *UserService) GetUserApis(reqCtx *request.Context, in *dto.EmptyReq) (out *dto.UserApisResp, err error) {
	// todo

	return
}

// 获取用户信息
func (s *UserService) GetUserInfo(reqCtx *request.Context, in *dto.EmptyReq) (out *dto.UserInfoResp, err error) {
	// todo

	return
}

// 查询用户登录历史
func (s *UserService) GetUserLoginHistoryList(reqCtx *request.Context, in *dto.UserLoginHistoryQuery) (out *dto.PageResp, err error) {
	// todo

	return
}

// 获取用户菜单权限
func (s *UserService) GetUserMenus(reqCtx *request.Context, in *dto.EmptyReq) (out *dto.UserMenusResp, err error) {
	// todo

	return
}

// 获取用户角色
func (s *UserService) GetUserRoles(reqCtx *request.Context, in *dto.EmptyReq) (out *dto.UserRolesResp, err error) {
	// todo

	return
}

// 修改用户信息
func (s *UserService) UpdateUserInfo(reqCtx *request.Context, in *dto.UserInfoReq) (out *dto.EmptyResp, err error) {
	// todo

	return
}
