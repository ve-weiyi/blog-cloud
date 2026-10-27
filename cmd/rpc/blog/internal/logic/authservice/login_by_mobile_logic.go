package authservicelogic

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/authrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/model"
	"github.com/ve-weiyi/vkit/adapter/iplocx"
	"github.com/ve-weiyi/vkit/x/randomx"
	"github.com/ve-weiyi/vkit/x/validatex"

	"github.com/ve-weiyi/blog-cloud/infra/auth/identityx"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
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
	if !validatex.IsMobile(in.Mobile) {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "手机号格式不正确")
	}

	user, err := l.svcCtx.TUserModel.FindOneByMobile(l.ctx, in.Mobile)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if user == nil {
		// 归属地查询放在事务外，避免外部接口抖动拖住数据库事务
		ip, _ := metax.GetRemoteIPFromCtx(l.ctx)
		ipSource := iplocx.GetIpSourceByBaidu(ip)

		// 用户不存在，自动注册
		err = l.svcCtx.GormDB.Transaction(func(tx *gorm.DB) error {
			userId, err := identityx.NewUserID()
			if err != nil {
				return err
			}

			newUser := &model.TUser{
				UserId:       userId,
				Username:     randomx.GenerateQQNumber(),
				Password:     "",
				Nickname:     in.Mobile,
				Mobile:       &in.Mobile,
				Email:        nil,
				Status:       enums.UserStatusNormal,
				RegisterType: enums.LoginTypeMobile,
				IpAddress:    ip,
				IpSource:     ipSource,
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
