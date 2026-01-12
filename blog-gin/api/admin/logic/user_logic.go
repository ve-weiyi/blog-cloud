package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type UserLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewUserLogic(svcCtx *svctx.ServiceContext) *UserLogic {
	return &UserLogic{
		svcCtx: svcCtx,
	}
}

// 删除用户绑定第三方平台账号
func (s *UserLogic) DeleteUserBindThirdParty(reqCtx *request.Context, in *types.DeleteUserBindThirdPartyReq) (out *types.EmptyResp, err error) {
	// todo

	return
}

// 获取用户接口权限
func (s *UserLogic) GetUserApis(reqCtx *request.Context, in *types.EmptyReq) (out *types.UserApisResp, err error) {
	// todo

	return
}

// 获取用户信息
func (s *UserLogic) GetUserInfo(reqCtx *request.Context, in *types.EmptyReq) (out *types.UserInfoResp, err error) {
	// todo

	return
}

// 查询用户登录历史
func (s *UserLogic) GetUserLoginHistoryList(reqCtx *request.Context, in *types.QueryUserLoginHistoryReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 获取用户菜单权限
func (s *UserLogic) GetUserMenus(reqCtx *request.Context, in *types.EmptyReq) (out *types.UserMenusResp, err error) {
	// todo

	return
}

// 获取用户角色
func (s *UserLogic) GetUserRoles(reqCtx *request.Context, in *types.EmptyReq) (out *types.UserRolesResp, err error) {
	// todo

	return
}

// 修改用户头像
func (s *UserLogic) UpdateUserAvatar(reqCtx *request.Context, in *types.UpdateUserAvatarReq) (out *types.EmptyResp, err error) {
	// todo

	return
}

// 修改用户绑定邮箱
func (s *UserLogic) UpdateUserBindEmail(reqCtx *request.Context, in *types.UpdateUserBindEmailReq) (out *types.EmptyResp, err error) {
	// todo

	return
}

// 修改用户绑定手机号
func (s *UserLogic) UpdateUserBindPhone(reqCtx *request.Context, in *types.UpdateUserBindPhoneReq) (out *types.EmptyResp, err error) {
	// todo

	return
}

// 修改用户绑定第三方平台账号
func (s *UserLogic) UpdateUserBindThirdParty(reqCtx *request.Context, in *types.UpdateUserBindThirdPartyReq) (out *types.EmptyResp, err error) {
	// todo

	return
}

// 修改用户信息
func (s *UserLogic) UpdateUserInfo(reqCtx *request.Context, in *types.UpdateUserInfoReq) (out *types.EmptyResp, err error) {
	// todo

	return
}

// 修改用户密码
func (s *UserLogic) UpdateUserPassword(reqCtx *request.Context, in *types.UpdateUserPasswordReq) (out *types.EmptyResp, err error) {
	// todo

	return
}
