package accountrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/global/constant"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/biz/bizerr"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/crypto"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/patternx"

	"github.com/zeromicro/go-zero/core/logx"
)

type EmailLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEmailLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailLoginLogic {
	return &EmailLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 邮箱登录
func (l *EmailLoginLogic) EmailLogin(in *accountrpc.EmailLoginReq) (*accountrpc.LoginResp, error) {
	// 校验邮箱格式
	if !patternx.IsValidEmail(in.Email) {
		return nil, bizerr.NewBizError(bizerr.CodeInvalidParam, "邮箱格式不正确")
	}

	// 验证用户是否存在
	account, err := l.svcCtx.TUserModel.FindOneByUsername(l.ctx, in.Email)
	if err != nil {
		return nil, bizerr.NewBizError(bizerr.CodeUserNotExist, "用户不存在")
	}

	// 验证密码是否正确
	if !crypto.BcryptCheck(in.Password, account.Password) {
		return nil, bizerr.NewBizError(bizerr.CodeUserPasswordError, "密码不正确")
	}

	return onLogin(l.ctx, l.svcCtx, account, constant.LoginTypeEmail)
}
