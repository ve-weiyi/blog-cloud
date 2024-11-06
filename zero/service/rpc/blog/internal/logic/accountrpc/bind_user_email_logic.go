package accountrpclogic

import (
	"context"
	"fmt"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/apierr"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/apierr/codex"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/constant"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/valid"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindUserEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindUserEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindUserEmailLogic {
	return &BindUserEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改用户邮箱
func (l *BindUserEmailLogic) BindUserEmail(in *accountrpc.BindUserEmailReq) (*accountrpc.EmptyResp, error) {
	// 校验邮箱格式
	if !valid.IsEmailValid(in.Email) {
		return nil, apierr.NewApiError(codex.CodeInvalidParam, "邮箱格式不正确")
	}

	// 验证用户是否存在
	user, err := l.svcCtx.TUserModel.FindOneByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, apierr.NewApiError(codex.CodeUserNotExist, err.Error())
	}

	// 验证code是否正确
	key := fmt.Sprintf("%s:%s", constant.BindEmail, in.Email)
	if !l.svcCtx.CaptchaHolder.VerifyCaptcha(key, in.VerifyCode) {
		return nil, apierr.NewApiError(codex.CodeCaptchaVerify, "验证码错误")
	}

	// 更新密码
	user.Username = in.Email
	user.Email = in.Email

	_, err = l.svcCtx.TUserModel.Save(l.ctx, user)
	if err != nil {
		return nil, err
	}

	return &accountrpc.EmptyResp{}, nil
}
