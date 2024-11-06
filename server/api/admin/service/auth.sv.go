package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type AuthService struct {
	svcCtx *svctx.ServiceContext
}

func NewAuthService(svcCtx *svctx.ServiceContext) *AuthService {
	return &AuthService{
		svcCtx: svcCtx,
	}
}

// 登录
func (s *AuthService) Login(reqCtx *request.Context, in *dto.LoginReq) (out *dto.LoginResp, err error) {
	// todo

	return
}

// 登出
func (s *AuthService) Logout(reqCtx *request.Context, in *dto.EmptyReq) (out *dto.EmptyResp, err error) {
	// todo

	return
}
