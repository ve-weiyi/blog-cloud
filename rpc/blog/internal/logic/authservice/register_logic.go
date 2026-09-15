package authservicelogic

import (
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/authrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/adapter/ipx"
	"github.com/ve-weiyi/vkit/x/cryptox"
	"github.com/ve-weiyi/vkit/x/patternx"
	"github.com/ve-weiyi/vkit/x/randomx"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/constants/enums"
	"github.com/ve-weiyi/blog-cloud/infra/metax"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户注册
func (l *RegisterLogic) Register(in *authrpc.RegisterRequest) (*authrpc.RegisterResponse, error) {
	// 校验邮箱格式
	if !patternx.IsValidEmail(in.Email) {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "邮箱格式不正确")
	}

	// 检查邮箱是否已注册
	user, err := l.svcCtx.TUserModel.FindOneByEmail(l.ctx, in.Email)
	if err == nil && user != nil {
		return nil, bizerr.NewBizError(bizcode.CodeResourceAlreadyExist, "邮箱已被注册")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 创建用户
	err = l.svcCtx.GormDB.Transaction(func(tx *gorm.DB) error {
		ip, _ := metax.GetRemoteIPFromCtx(l.ctx)
		newUser := &model.TUser{
			Id:           0,
			UserId:       randomx.GenerateRandomUUID(),
			Username:     randomx.GenerateQQNumber(),
			Password:     cryptox.BcryptHash(in.Password),
			Nickname:     *in.Nickname,
			Avatar:       "",
			Email:        &in.Email,
			Mobile:       nil,
			Status:       enums.UserStatusNormal,
			Info:         "",
			RegisterType: enums.LoginTypeEmail,
			IpAddress:    ip,
			IpSource:     ipx.GetIpSourceByBaidu(ip),
			CreatedAt:    time.Time{},
			UpdatedAt:    time.Time{},
			DeletedAt:    nil,
		}
		user, err = onRegister(l.ctx, l.svcCtx, tx, newUser)
		return err
	})
	if err != nil {
		return nil, err
	}

	// 注册与登录共用同一套凭据签发流程，但契约上两者是各自独立的响应类型
	login, err := onLogin(l.ctx, l.svcCtx, user, enums.LoginTypeRegister)
	if err != nil {
		return nil, err
	}

	return &authrpc.RegisterResponse{
		UserId:   login.UserId,
		Username: login.Username,
		Status:   login.Status,
	}, nil
}

func onRegister(ctx context.Context, svcCtx *svc.ServiceContext, tx *gorm.DB, user *model.TUser) (out *model.TUser, err error) {
	/** 创建用户 **/
	_, err = svcCtx.TUserModel.WithTx(tx).Insert(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
