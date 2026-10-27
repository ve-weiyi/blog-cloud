package authservicelogic

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/authrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/model"
	"github.com/ve-weiyi/vkit/adapter/iplocx"
	"github.com/ve-weiyi/vkit/x/cryptox"
	"github.com/ve-weiyi/vkit/x/randomx"
	"github.com/ve-weiyi/vkit/x/validatex"

	"github.com/ve-weiyi/blog-cloud/infra/auth/identityx"
	"github.com/ve-weiyi/blog-cloud/infra/auth/passwordx"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
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
	if !validatex.IsEmail(in.Email) {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "邮箱格式不正确")
	}
	if err := passwordx.ValidatePassword(in.Password); err != nil {
		return nil, err
	}

	// 检查邮箱是否已注册
	user, err := l.svcCtx.TUserModel.FindOneByEmail(l.ctx, in.Email)
	if err == nil && user != nil {
		return nil, bizerr.NewBizError(bizcode.CodeResourceAlreadyExist, "邮箱已被注册")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 归属地查询是外部调用，放在事务外，避免接口抖动拖住数据库事务
	ip, _ := metax.GetRemoteIPFromCtx(l.ctx)
	ipSource := iplocx.GetIpSourceByBaidu(ip)

	// 创建用户
	err = l.svcCtx.GormDB.Transaction(func(tx *gorm.DB) error {
		hashedPassword, err := cryptox.BcryptHash(in.Password)
		if err != nil {
			return err
		}
		nickname := ""
		if in.Nickname != nil {
			nickname = *in.Nickname
		}

		username := randomx.GenerateQQNumber()
		if in.Username != nil && strings.TrimSpace(*in.Username) != "" {
			username = strings.TrimSpace(*in.Username)
			count, err := l.svcCtx.TUserModel.FindCount(l.ctx, "username = ?", username)
			if err != nil {
				return err
			}
			if count > 0 {
				return bizerr.NewBizError(bizcode.CodeResourceAlreadyExist, "用户名已被占用")
			}
		}

		userId, err := identityx.NewUserID()
		if err != nil {
			return err
		}

		newUser := &model.TUser{
			Id:           0,
			UserId:       userId,
			Username:     username,
			Password:     hashedPassword,
			Nickname:     nickname,
			Avatar:       "",
			Email:        &in.Email,
			Mobile:       nil,
			Status:       enums.UserStatusNormal,
			Info:         "",
			RegisterType: enums.LoginTypeEmail,
			IpAddress:    ip,
			IpSource:     ipSource,
			CreatedAt:    time.Time{},
			UpdatedAt:    time.Time{},
			DeletedAt:    nil,
		}
		user, err = onRegister(l.ctx, l.svcCtx, tx, newUser)
		return err
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, bizerr.NewBizError(bizcode.CodeResourceAlreadyExist, "邮箱已被注册")
		}
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
