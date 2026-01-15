package accountrpclogic

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/constant"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rpcutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizcode"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizerr"
	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/oauth"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/cryptox"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/ipx"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/randomx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ThirdLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewThirdLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThirdLoginLogic {
	return &ThirdLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 第三方登录
func (l *ThirdLoginLogic) ThirdLogin(in *accountrpc.ThirdLoginReq) (*accountrpc.LoginResp, error) {
	app, err := rpcutils.GetAppNameFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	auth, err := GetPlatformOauth(l.ctx, l.svcCtx, app, in.Platform)
	if err != nil {
		return nil, err
	}

	// 获取第三方用户信息
	info, err := auth.GetAuthUserInfo(in.Code)
	if err != nil {
		return nil, err
	}

	if info.OpenId == "" {
		return nil, fmt.Errorf("open_id is empty")
	}

	// 查询用户是否存在
	userOauth, err := l.svcCtx.TUserOauthModel.FindOneByOpenIdPlatform(l.ctx, info.OpenId, in.Platform)
	if userOauth == nil {
		// 用户未注册,先注册用户
		err = l.svcCtx.Gorm.Transaction(func(tx *gorm.DB) error {
			userOauth, err = l.oauthRegister(tx, in.Platform, info)
			return err
		})
		if err != nil {
			return nil, err
		}
	}

	// 用户已经注册,查询用户信息
	user, err := l.svcCtx.TUserModel.FindOneByUserId(l.ctx, userOauth.UserId)
	if err != nil {
		return nil, bizerr.NewBizError(bizcode.CodeUserNotExist, err.Error())
	}

	return onLogin(l.ctx, l.svcCtx, user, userOauth.Platform)
}

func (l *ThirdLoginLogic) oauthRegister(tx *gorm.DB, platform string, info *oauth.UserResult) (out *model.TUserOauth, err error) {
	// 用户未注册,先注册用户
	ip, _ := rpcutils.GetRemoteIPFromCtx(l.ctx)
	is := ipx.GetIpSourceByBaidu(ip)

	// 用户账号
	user := &model.TUser{
		UserId:       randomx.GenerateUID(32),
		Username:     generateUID(l.ctx, l.svcCtx, tx),
		Password:     cryptox.BcryptHash(info.EnName),
		Nickname:     info.NickName,
		Avatar:       info.Avatar,
		Email:        info.Email,
		Phone:        info.Mobile,
		Info:         "",
		Status:       constant.UserStatusNormal,
		RegisterType: platform,
		IpAddress:    ip,
		IpSource:     is,
	}

	/** 创建用户 **/
	ua, err := onRegister(l.ctx, l.svcCtx, tx, user)
	if err != nil {
		return nil, err
	}

	// 绑定用户第三方信息
	userOauth := &model.TUserOauth{
		UserId:   ua.UserId,
		OpenId:   info.OpenId,
		Platform: platform,
		Nickname: info.NickName,
		Avatar:   info.Avatar,
	}

	/** 创建用户第三方信息 **/
	_, err = l.svcCtx.TUserOauthModel.WithTransaction(tx).Insert(l.ctx, userOauth)
	if err != nil {
		return nil, err
	}

	return userOauth, nil
}

// 生成唯一UID
func generateUID(ctx context.Context, svcCtx *svc.ServiceContext, tx *gorm.DB) string {
	uid := randomx.GenerateQQNumber()
	user, _ := svcCtx.TUserModel.FindOneByUsername(ctx, uid)
	if user != nil {
		return generateUID(ctx, svcCtx, tx) // 如果用户已存在，继续生成新的UID
	}

	return uid
}
