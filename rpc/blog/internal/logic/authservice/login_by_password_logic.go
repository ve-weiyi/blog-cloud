package authservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/authrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/x/cryptox"
	"github.com/ve-weiyi/vkit/x/patternx"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/constants/enums"
	"github.com/ve-weiyi/blog-cloud/infra/metax"
)

type LoginByPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginByPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginByPasswordLogic {
	return &LoginByPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 密码登录
func (l *LoginByPasswordLogic) LoginByPassword(in *authrpc.LoginByPasswordRequest) (*authrpc.LoginResponse, error) {
	var user *model.TUser
	var err error
	switch {
	case patternx.IsValidEmail(in.Account):
		user, err = l.svcCtx.TUserModel.FindOneByEmail(l.ctx, in.Account)
	case patternx.IsValidMobile(in.Account):
		user, err = l.svcCtx.TUserModel.FindOneByMobile(l.ctx, in.Account)
	default:
		user, err = l.svcCtx.TUserModel.FindOneByUsername(l.ctx, in.Account)
	}
	if err != nil {
		return nil, bizerr.NewBizError(bizcode.CodeResourceNotFound, "用户不存在")
	}

	if !cryptox.BcryptCheck(in.Password, user.Password) {
		return nil, bizerr.NewBizError(bizcode.CodePasswordError, "密码不正确")
	}

	return onLogin(l.ctx, l.svcCtx, user, enums.LoginTypeUsername)
}

func onLogin(ctx context.Context, svcCtx *svc.ServiceContext, user *model.TUser, loginType string) (resp *authrpc.LoginResponse, err error) {
	// 判断用户是否被禁用
	if user.Status == enums.UserStatusDisabled {
		return nil, bizerr.NewBizError(bizcode.CodeAccountDisabled, "用户已被禁用")
	}

	did, _ := metax.GetDeviceIdFromCtx(ctx)
	// 推送登录日志消息，失败不影响登录流程
	if err := mq.PublishLoginEvent(ctx, &mq.LoginEvent{
		UserId:    user.UserId,
		DeviceId:  did,
		LoginType: loginType,
	}); err != nil {
		logx.WithContext(ctx).Errorf("发送登录日志消息失败: %v", err)
	}

	resp = &authrpc.LoginResponse{
		UserId:   user.UserId,
		Username: user.Username,
		Status:   user.Status,
	}

	return resp, nil
}
