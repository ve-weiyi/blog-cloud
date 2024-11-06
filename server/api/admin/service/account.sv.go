package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type AccountService struct {
	svcCtx *svctx.ServiceContext
}

func NewAccountService(svcCtx *svctx.ServiceContext) *AccountService {
	return &AccountService{
		svcCtx: svcCtx,
	}
}

// 获取用户分布地区
func (s *AccountService) FindAccountAreaAnalysis(reqCtx *request.Context, in *dto.AccountQuery) (out *dto.PageResp, err error) {
	// todo

	return
}

// 查询用户列表
func (s *AccountService) FindAccountList(reqCtx *request.Context, in *dto.AccountQuery) (out *dto.PageResp, err error) {
	// todo

	return
}

// 查询用户登录历史
func (s *AccountService) FindAccountLoginHistoryList(reqCtx *request.Context, in *dto.AccountQuery) (out *dto.PageResp, err error) {
	// todo

	return
}

// 查询在线用户列表
func (s *AccountService) FindAccountOnlineList(reqCtx *request.Context, in *dto.AccountQuery) (out *dto.PageResp, err error) {
	// todo

	return
}

// 修改用户角色
func (s *AccountService) UpdateAccountRoles(reqCtx *request.Context, in *dto.UpdateAccountRolesReq) (out *dto.EmptyResp, err error) {
	// todo

	return
}

// 修改用户状态
func (s *AccountService) UpdateAccountStatus(reqCtx *request.Context, in *dto.UpdateAccountStatusReq) (out *dto.EmptyResp, err error) {
	// todo

	return
}
