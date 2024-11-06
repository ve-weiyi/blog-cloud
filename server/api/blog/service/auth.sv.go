package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/dto"
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

// 第三方登录授权地址
func (s *AuthService) OauthAuthorizeUrl(reqCtx *request.Context, in *dto.OauthLoginReq) (out *dto.OauthLoginUrlResp, err error) {
	// todo

	return
}

// 第三方登录
func (s *AuthService) OauthLogin(reqCtx *request.Context, in *dto.OauthLoginReq) (out *dto.LoginResp, err error) {
	// todo

	return
}

// 注册
func (s *AuthService) Register(reqCtx *request.Context, in *dto.RegisterReq) (out *dto.EmptyResp, err error) {
	// todo

	return
}

// 发送注册账号邮件
func (s *AuthService) SendRegisterEmail(reqCtx *request.Context, in *dto.UserEmailReq) (out *dto.EmptyResp, err error) {
	// todo

	return
}

// 重置密码
func (s *AuthService) ResetPassword(reqCtx *request.Context, in *dto.ResetPasswordReq) (out *dto.EmptyResp, err error) {
	// todo

	return
}

// 发送重置密码邮件
func (s *AuthService) SendResetEmail(reqCtx *request.Context, in *dto.UserEmailReq) (out *dto.EmptyResp, err error) {
	// todo

	return
}

// 绑定邮箱
func (s *AuthService) BindUserEmail(reqCtx *request.Context, in *dto.BindUserEmailReq) (out *dto.EmptyResp, err error) {
	// todo

	return
}

// 注销
func (s *AuthService) Logoff(reqCtx *request.Context, in *dto.EmptyReq) (out *dto.EmptyResp, err error) {
	// todo

	return
}

// 登出
func (s *AuthService) Logout(reqCtx *request.Context, in *dto.EmptyReq) (out *dto.EmptyResp, err error) {
	// todo

	return
}

// 发送绑定邮箱验证码
func (s *AuthService) SendBindEmail(reqCtx *request.Context, in *dto.UserEmailReq) (out *dto.EmptyResp, err error) {
	// todo

	return
}
