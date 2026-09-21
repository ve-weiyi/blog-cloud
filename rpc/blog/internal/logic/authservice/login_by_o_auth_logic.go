package authservicelogic

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/authrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/adapter/iplocx"
	"github.com/ve-weiyi/vkit/adapter/oauthx"
	"github.com/ve-weiyi/vkit/x/randomx"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/constants/enums"
	"github.com/ve-weiyi/blog-cloud/infra/metax"
)

type LoginByOAuthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginByOAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginByOAuthLogic {
	return &LoginByOAuthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 第三方平台登录（自动注册）
func (l *LoginByOAuthLogic) LoginByOAuth(in *authrpc.LoginByOAuthRequest) (*authrpc.LoginResponse, error) {
	auth, ok := l.svcCtx.OAuthProviders[in.Platform]
	if !ok {
		return nil, fmt.Errorf("platform %s is not support", in.Platform)
	}

	info, err := auth.GetAuthUserInfo(l.ctx, in.Code)
	if err != nil {
		return nil, err
	}

	if info.OpenId == "" {
		return nil, fmt.Errorf("open_id is empty")
	}

	userOauth, err := l.svcCtx.TUserOauthModel.FindOneByPlatformOpenId(l.ctx, in.Platform, info.OpenId)
	if userOauth == nil {
		err = l.svcCtx.GormDB.Transaction(func(tx *gorm.DB) error {
			userOauth, err = l.oauthRegister(tx, in.Platform, info)
			return err
		})
		if err != nil {
			return nil, err
		}
	}

	user, err := l.svcCtx.TUserModel.FindOneByUserId(l.ctx, userOauth.UserId)
	if err != nil {
		return nil, bizerr.NewBizError(bizcode.CodeResourceNotFound, err.Error())
	}

	return onLogin(l.ctx, l.svcCtx, user, enums.LoginTypeOauth)
}

func (l *LoginByOAuthLogic) oauthRegister(tx *gorm.DB, platform string, info *oauthx.UserResult) (out *model.TUserOauth, err error) {
	ip, _ := metax.GetRemoteIPFromCtx(l.ctx)

	newUser := &model.TUser{
		UserId:       randomx.GenerateRandomUUID(),
		Username:     randomx.GenerateQQNumber(),
		Password:     "",
		Nickname:     info.NickName,
		Avatar:       info.Avatar,
		Email:        nil,
		Mobile:       nil,
		Status:       enums.UserStatusNormal,
		RegisterType: enums.LoginTypeOauth,
		IpAddress:    ip,
		IpSource:     iplocx.GetIpSourceByBaidu(ip),
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
		DeletedAt:    nil,
	}

	user, err := onRegister(l.ctx, l.svcCtx, tx, newUser)
	if err != nil {
		return nil, err
	}

	userOauth := &model.TUserOauth{
		UserId:   user.UserId,
		OpenId:   info.OpenId,
		Platform: platform,
		Nickname: info.NickName,
		Avatar:   info.Avatar,
	}

	_, err = l.svcCtx.TUserOauthModel.WithTx(tx).Insert(l.ctx, userOauth)
	if err != nil {
		return nil, err
	}

	return userOauth, nil
}
