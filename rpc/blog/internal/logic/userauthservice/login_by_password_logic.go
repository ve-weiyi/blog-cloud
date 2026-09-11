package userauthservicelogic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/userauthrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/adapter/mqx"
	"github.com/ve-weiyi/vkit/x/cryptox"
	"github.com/ve-weiyi/vkit/x/jsonconv"
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
func (l *LoginByPasswordLogic) LoginByPassword(in *userauthrpc.LoginByPasswordRequest) (*userauthrpc.LoginResponse, error) {
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

func onLogin(ctx context.Context, svcCtx *svc.ServiceContext, user *model.TUser, loginType string) (resp *userauthrpc.LoginResponse, err error) {
	// 判断用户是否被禁用
	if user.Status == enums.UserStatusDisabled {
		return nil, bizerr.NewBizError(bizcode.CodeAccountDisabled, "用户已被禁用")
	}

	did, _ := metax.GetDeviceIdFromCtx(ctx)
	if mq.LoginProducer != nil {
		mq.LoginProducer.Send(ctx, &mqx.Message{
			Topic: mq.LoginQueue,
			Key:   mq.LoginRoutingKey,
			Body: []byte(jsonconv.AnyToJsonNE(mq.LoginEvent{
				UserId:    user.UserId,
				DeviceId:  did,
				LoginType: loginType,
			})),
			Timestamp: time.Now(),
		})
	}

	resp = &userauthrpc.LoginResponse{
		UserId:   user.UserId,
		Username: user.Username,
		Status:   user.Status,
	}

	return resp, nil
}
