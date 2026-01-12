package accountrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/constant"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/mailtemplate"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/rediskey"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizcode"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizerr"
	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/mail"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/patternx"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/tempx"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendEmailVerifyCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendEmailVerifyCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendEmailVerifyCodeLogic {
	return &SendEmailVerifyCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 发送邮件验证码
func (l *SendEmailVerifyCodeLogic) SendEmailVerifyCode(in *accountrpc.SendEmailVerifyCodeReq) (*accountrpc.EmptyResp, error) {
	// 校验邮箱格式
	if !patternx.IsValidEmail(in.Email) {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "邮箱格式不正确")
	}

	// 验证用户是否存在
	switch in.Type {
	case constant.CodeTypeRegister, constant.CodeTypeBindEmail:
		// 用户已存在
		exist, _ := l.svcCtx.TUserModel.FindOne(l.ctx, "email = ?", in.Email)
		if exist != nil {
			return nil, bizerr.NewBizError(bizcode.CodeUserAlreadyExist, "用户已存在")
		}
	case constant.CodeTypeResetPwd:
		// 用户不存在
		exist, _ := l.svcCtx.TUserModel.FindOne(l.ctx, "email = ?", in.Email)
		if exist == nil {
			return nil, bizerr.NewBizError(bizcode.CodeUserNotExist, "用户不存在")
		}
	default:
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "验证码类型不正确")
	}

	// 发送验证码邮件
	key := rediskey.GetCaptchaKey(in.Type, in.Email)
	code, _ := l.svcCtx.CaptchaHolder.GetCodeCaptcha(key)
	data := mailtemplate.CaptchaEmail{
		Username: in.Email,
		Content:  emailContent[in.Type],
		Code:     code,
	}

	// 组装邮件内容
	content, err := tempx.TempParseString(mailtemplate.TempCaptchaCode, data)
	if err != nil {
		return nil, err
	}

	msg := &mail.EmailMessage{
		To:      []string{in.Email},
		Subject: emailSubject[in.Type],
		Content: content,
	}
	// 发送邮件
	err = l.svcCtx.EmailDeliver.DeliveryEmail(msg)
	if err != nil {
		return nil, err
	}

	return &accountrpc.EmptyResp{}, nil
}

var emailContent = map[string]string{
	constant.CodeTypeRegister:  "欢迎注册我的博客平台。",
	constant.CodeTypeResetPwd:  "您的账号正在尝试修改密码。",
	constant.CodeTypeBindEmail: "您的账号正在尝试修改绑定邮箱。",
}

var emailSubject = map[string]string{
	constant.CodeTypeRegister:  "blog|注册邮件提醒",
	constant.CodeTypeResetPwd:  "blog|重置密码邮件提醒",
	constant.CodeTypeBindEmail: "blog|修改绑定邮箱邮件提醒",
}
