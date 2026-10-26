package authservicelogic

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/authrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/model"
	"github.com/ve-weiyi/vkit/x/cryptox"
	"github.com/ve-weiyi/vkit/x/patternx"

	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
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
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerr.WrapBizError(bizcode.CodeDatabaseError, "查询用户失败", err)
		}

		// 账号不存在与密码错误返回同一标识与文案——分开返回等于给出一条可枚举账号的信道
		did, _ := metax.GetDeviceIdFromCtx(l.ctx)
		publishLoginEvent(l.ctx, l.svcCtx, "", did, enums.LoginTypeUsername, 0, "账号不存在")
		return nil, bizerr.NewBizError(bizcode.CodeCredentialsInvalid, "账号或密码不正确")
	}

	if !cryptox.BcryptCheck(in.Password, user.Password) {
		did, _ := metax.GetDeviceIdFromCtx(l.ctx)
		publishLoginEvent(l.ctx, l.svcCtx, user.UserId, did, enums.LoginTypeUsername, 0, "密码不正确")
		return nil, bizerr.NewBizError(bizcode.CodeCredentialsInvalid, "账号或密码不正确")
	}

	return onLogin(l.ctx, l.svcCtx, user, enums.LoginTypeUsername)
}

// publishLoginEvent 发布登录日志事件；失败只记录日志，不影响登录流程。
func publishLoginEvent(ctx context.Context, svcCtx *svc.ServiceContext, userId, deviceId, loginType string, status int64, reason string) {
	if _, err := mq.LoginMessage.Publish(ctx, svcCtx.MQ.Publisher, &mq.LoginEvent{
		UserId:     userId,
		DeviceId:   deviceId,
		LoginType:  loginType,
		Status:     status,
		FailReason: reason,
	}); err != nil {
		logx.WithContext(ctx).Errorf("发送登录日志消息失败: %v", err)
	}
}

func onLogin(ctx context.Context, svcCtx *svc.ServiceContext, user *model.TUser, loginType string) (resp *authrpc.LoginResponse, err error) {
	did, _ := metax.GetDeviceIdFromCtx(ctx)

	if user.DeletedAt != nil {
		publishLoginEvent(ctx, svcCtx, user.UserId, did, loginType, 0, "账号已注销")
		return nil, bizerr.NewBizError(bizcode.CodeAccountDisabled, "账号已注销")
	}

	switch user.Status {
	case enums.UserStatusNormal:
	case enums.UserStatusDisabled:
		publishLoginEvent(ctx, svcCtx, user.UserId, did, loginType, 0, "用户已被禁用")
		return nil, bizerr.NewBizError(bizcode.CodeAccountDisabled, "用户已被禁用")
	case enums.UserStatusLogoff:
		publishLoginEvent(ctx, svcCtx, user.UserId, did, loginType, 0, "账号已注销")
		return nil, bizerr.NewBizError(bizcode.CodeAccountDisabled, "账号已注销")
	default:
		publishLoginEvent(ctx, svcCtx, user.UserId, did, loginType, 0, "账号不可用")
		return nil, bizerr.NewBizError(bizcode.CodeAccountDisabled, "账号不可用")
	}

	// 推送登录日志消息，失败不影响登录流程
	publishLoginEvent(ctx, svcCtx, user.UserId, did, loginType, 1, "")

	resp = &authrpc.LoginResponse{
		UserId:   user.UserId,
		Username: user.Username,
		Status:   user.Status,
	}

	return resp, nil
}
