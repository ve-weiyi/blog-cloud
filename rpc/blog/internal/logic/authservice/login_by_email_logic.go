package authservicelogic

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/authrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/x/patternx"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/constants/enums"
)

type LoginByEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginByEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginByEmailLogic {
	return &LoginByEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 邮箱验证码登录
func (l *LoginByEmailLogic) LoginByEmail(in *authrpc.LoginByEmailRequest) (*authrpc.LoginResponse, error) {
	if !patternx.IsValidEmail(in.Email) {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "邮箱格式不正确")
	}

	user, err := l.svcCtx.TUserModel.FindOneByEmail(l.ctx, in.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerr.WrapBizError(bizcode.CodeDatabaseError, "查询用户失败", err)
		}

		// 与密码登录保持同一标识与文案，避免账号枚举
		return nil, bizerr.NewBizError(bizcode.CodeCredentialsInvalid, "账号或密码不正确")
	}

	return onLogin(l.ctx, l.svcCtx, user, enums.LoginTypeEmail)
}
