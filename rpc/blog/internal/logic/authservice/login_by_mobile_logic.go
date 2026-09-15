package authservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/authrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/adapter/ipx"
	"github.com/ve-weiyi/vkit/x/patternx"
	"github.com/ve-weiyi/vkit/x/randomx"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/constants/enums"
	"github.com/ve-weiyi/blog-cloud/infra/metax"
)

type LoginByMobileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginByMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginByMobileLogic {
	return &LoginByMobileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 手机验证码登录
func (l *LoginByMobileLogic) LoginByMobile(in *authrpc.LoginByMobileRequest) (*authrpc.LoginResponse, error) {
	if !patternx.IsValidMobile(in.Mobile) {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "手机号格式不正确")
	}

	user, _ := l.svcCtx.TUserModel.FindOneByMobile(l.ctx, in.Mobile)
	if user == nil {
		// 用户不存在，自动注册
		var err error
		err = l.svcCtx.GormDB.Transaction(func(tx *gorm.DB) error {
			ip, _ := metax.GetRemoteIPFromCtx(l.ctx)
			newUser := &model.TUser{
				UserId:       randomx.GenerateRandomUUID(),
				Username:     randomx.GenerateQQNumber(),
				Password:     "",
				Nickname:     in.Mobile,
				Mobile:       &in.Mobile,
				Email:        nil,
				Status:       enums.UserStatusNormal,
				RegisterType: enums.LoginTypeMobile,
				IpAddress:    ip,
				IpSource:     ipx.GetIpSourceByBaidu(ip),
			}
			user, err = onRegister(l.ctx, l.svcCtx, tx, newUser)
			return err
		})
		if err != nil {
			return nil, err
		}
	}

	return onLogin(l.ctx, l.svcCtx, user, enums.LoginTypeMobile)
}
