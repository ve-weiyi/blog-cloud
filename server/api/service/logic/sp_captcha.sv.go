package logic

import (
	"fmt"

	"github.com/ve-weiyi/ve-blog-golang/server/api/model/request"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/response"
	"github.com/ve-weiyi/ve-blog-golang/server/api/service/svc"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/apierr"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/constant"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/mail"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/jsonconv"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/temputil"
)

type CaptchaService struct {
	svcCtx *svc.ServiceContext
}

func NewCaptchaService(svcCtx *svc.ServiceContext) *CaptchaService {
	return &CaptchaService{
		svcCtx: svcCtx,
	}
}

// 发送验证码
func (s *CaptchaService) SendCaptchaEmail(reqCtx *request.Context, req *request.CaptchaEmailReq) (result interface{}, err error) {
	// 验证用户是否存在
	account, err := s.svcCtx.UserAccountRepository.LoadUserByUsername(reqCtx, req.Email)
	if err != nil {
		return nil, apierr.ErrorUserNotExist
	}

	userinfo, err := s.svcCtx.UserAccountRepository.FindUserInfo(reqCtx, account.ID)
	if err != nil {
		return nil, apierr.ErrorUserNotExist
	}

	// 设置key
	key := fmt.Sprintf("%s:%s", req.Service, req.Email)
	code := s.svcCtx.Captcha.GetCodeCaptcha(key)
	data := mail.NewEmailContent()

	data.Title = fmt.Sprintf("重置密码邮件提醒")
	data.DearUser = fmt.Sprintf("你好，%s", userinfo.Nickname)
	data.Content = fmt.Sprintf("你的账号 %s 正在尝试重置密码。验证码为 %s，有效期15分钟！如果您没有尝试重置密码，请忽略此邮件。", req.Email, code)
	data.ButtonTips = fmt.Sprintf("点击重置密码")
	data.ButtonLink = "https://veweiyi.cn/blog/reset_password"

	var temp string
	switch req.Service {
	case constant.ForgetPassword:
		temp = mail.TempForgetPassword
	case constant.Register:
		temp = mail.TempRegister
	default:
		temp = mail.TempSimpleCode
	}

	content, err := temputil.TempParseString(temp, data)
	if err != nil {
		return nil, err
	}

	msg := &mail.EmailMessage{
		To:      []string{req.Email},
		Subject: data.Title,
		Content: content,
		Type:    0,
	}

	err = s.svcCtx.EmailPublisher.PublishMessage([]byte(jsonconv.ObjectToJson(msg)))
	if err != nil {
		return nil, err
	}
	return true, nil
}

// 获取图片验证码
func (s *CaptchaService) GetCaptchaImage(reqCtx *request.Context, req *request.CaptchaReq) (resp *response.CaptchaDTO, err error) {
	id, b64s, err := s.svcCtx.Captcha.GetImageCaptcha(req.CaptchaType, req.Height, req.Width, req.Length)
	if err != nil {
		return nil, err
	}

	resp = &response.CaptchaDTO{
		ID:         id,
		EncodeData: b64s,
		Length:     req.Length,
	}
	return resp, nil
}

func (s *CaptchaService) VerifyImageCaptcha(reqCtx *request.Context, req *request.CaptchaVerifyReq) (resp interface{}, err error) {
	if !s.svcCtx.Captcha.VerifyCaptcha(req.ID, req.Code) {
		return nil, apierr.ErrorCaptchaVerify
	}

	return resp, nil
}
