package authservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/authrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/x/cryptox"
	"github.com/ve-weiyi/vkit/x/patternx"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
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
	if !patternx.IsValidEmail(in.Email) {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "邮箱格式不正确")
	}

	// 校验两次密码一致
	if in.Password != in.ConfirmPassword {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "两次密码输入不一致")
	}

	// 校验密码长度
	if len(in.Password) < 6 {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "密码不能少于6位")
	}

	// 验证用户是否存在
	exist, _ := l.svcCtx.TUserModel.FindOneByEmail(l.ctx, in.Email)
	if exist == nil {
		return nil, bizerr.NewBizError(bizcode.CodeResourceNotFound, "用户不存在")
	}

	// 更新密码
	exist.Password = cryptox.BcryptHash(in.Password)

	_, err := l.svcCtx.TUserModel.Save(l.ctx, exist)
	if err != nil {
		return nil, err
	}

	return &authrpc.ResetPasswordResponse{}, nil
}
