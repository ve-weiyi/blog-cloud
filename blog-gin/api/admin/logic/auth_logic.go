package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type AuthLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewAuthLogic(svcCtx *svctx.ServiceContext) *AuthLogic {
	return &AuthLogic{
		svcCtx: svcCtx,
	}
}

// 获取客户端信息
func (s *AuthLogic) GetClientInfo(reqCtx *request.Context, in *types.GetClientInfoReq) (out *types.GetClientInfoResp, err error) {
	// todo

	return
}

// 邮箱登录
func (s *AuthLogic) EmailLogin(reqCtx *request.Context, in *types.EmailLoginReq) (out *types.LoginResp, err error) {
	// todo

	return
}

// 获取验证码
func (s *AuthLogic) GetCaptchaCode(reqCtx *request.Context, in *types.GetCaptchaCodeReq) (out *types.GetCaptchaCodeResp, err error) {
	// todo

	return
}

// 第三方登录授权地址
func (s *AuthLogic) GetOauthAuthorizeUrl(reqCtx *request.Context, in *types.GetOauthAuthorizeUrlReq) (out *types.GetOauthAuthorizeUrlResp, err error) {
	// todo

	return
}

// 登录
func (s *AuthLogic) Login(reqCtx *request.Context, in *types.LoginReq) (out *types.LoginResp, err error) {
	// todo

	return
}

// 手机登录
func (s *AuthLogic) PhoneLogin(reqCtx *request.Context, in *types.PhoneLoginReq) (out *types.LoginResp, err error) {
	// todo

	return
}

// 刷新token
func (s *AuthLogic) RefreshToken(reqCtx *request.Context, in *types.RefreshTokenReq) (out *types.LoginResp, err error) {
	// todo

	return
}

// 注册
func (s *AuthLogic) Register(reqCtx *request.Context, in *types.RegisterReq) (out *types.EmptyResp, err error) {
	// todo

	return
}

// 重置密码
func (s *AuthLogic) ResetPassword(reqCtx *request.Context, in *types.ResetPasswordReq) (out *types.EmptyResp, err error) {
	// todo

	return
}

// 发送邮件验证码
func (s *AuthLogic) SendEmailVerifyCode(reqCtx *request.Context, in *types.SendEmailVerifyCodeReq) (out *types.EmptyResp, err error) {
	// todo

	return
}

// 发送手机验证码
func (s *AuthLogic) SendPhoneVerifyCode(reqCtx *request.Context, in *types.SendPhoneVerifyCodeReq) (out *types.EmptyResp, err error) {
	// todo

	return
}

// 第三方登录
func (s *AuthLogic) ThirdLogin(reqCtx *request.Context, in *types.ThirdLoginReq) (out *types.LoginResp, err error) {
	// todo

	return
}

// 注销
func (s *AuthLogic) Logoff(reqCtx *request.Context, in *types.EmptyReq) (out *types.EmptyResp, err error) {
	// todo

	return
}

// 登出
func (s *AuthLogic) Logout(reqCtx *request.Context, in *types.EmptyReq) (out *types.EmptyResp, err error) {
	// todo

	return
}
