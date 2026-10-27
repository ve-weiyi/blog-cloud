package authservicelogic

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/authrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/x/cryptox"
	"github.com/ve-weiyi/vkit/x/validatex"

	"github.com/ve-weiyi/blog-cloud/infra/auth/passwordx"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 重置密码
func (l *ResetPasswordLogic) ResetPassword(in *authrpc.ResetPasswordRequest) (*authrpc.ResetPasswordResponse, error) {
	// 校验邮箱格式
	if !validatex.IsEmail(in.Email) {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "邮箱格式不正确")
	}

	if err := passwordx.ValidatePasswordConfirm(in.Password, in.ConfirmPassword); err != nil {
		return nil, err
	}
	if err := passwordx.ValidatePassword(in.Password); err != nil {
		return nil, err
	}

	// 验证用户是否存在
	exist, err := l.svcCtx.TUserModel.FindOneByEmail(l.ctx, in.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if exist == nil {
		return nil, bizerr.NewBizError(bizcode.CodeResourceNotFound, "用户不存在")
	}

	// 更新密码
	hashedPassword, err := cryptox.BcryptHash(in.Password)
	if err != nil {
		return nil, err
	}
	exist.Password = hashedPassword

	_, err = l.svcCtx.TUserModel.Update(l.ctx, exist)
	if err != nil {
		return nil, err
	}

	return &authrpc.ResetPasswordResponse{Success: true, UserId: exist.UserId}, nil
}
