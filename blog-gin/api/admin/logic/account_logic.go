package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type AccountLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewAccountLogic(svcCtx *svctx.ServiceContext) *AccountLogic {
	return &AccountLogic{
		svcCtx: svcCtx,
	}
}

// 查询用户列表
func (s *AccountLogic) FindAccountList(reqCtx *request.Context, in *types.AccountQuery) (out *types.PageResp, err error) {
	// todo

	return
}

// 查询在线用户列表
func (s *AccountLogic) FindAccountOnlineList(reqCtx *request.Context, in *types.AccountQuery) (out *types.PageResp, err error) {
	// todo

	return
}

// 修改用户密码
func (s *AccountLogic) UpdateAccountPassword(reqCtx *request.Context, in *types.UpdateAccountPasswordReq) (out *types.EmptyResp, err error) {
	// todo

	return
}

// 修改用户角色
func (s *AccountLogic) UpdateAccountRoles(reqCtx *request.Context, in *types.UpdateAccountRolesReq) (out *types.EmptyResp, err error) {
	// todo

	return
}

// 修改用户状态
func (s *AccountLogic) UpdateAccountStatus(reqCtx *request.Context, in *types.UpdateAccountStatusReq) (out *types.EmptyResp, err error) {
	// todo

	return
}
