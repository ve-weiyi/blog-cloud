package accountrpclogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/constant"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/rediskey"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rpcutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizcode"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizerr"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/cryptox"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/ipx"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/patternx"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/randomx"
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

// 注册
func (l *RegisterLogic) Register(in *accountrpc.RegisterReq) (*accountrpc.RegisterResp, error) {
	// 校验邮箱格式
	if !patternx.IsValidEmail(in.Email) {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "邮箱格式不正确")
	}

	// 获取用户
	exist, err := l.svcCtx.TUserModel.FindOneByUsername(l.ctx, in.Username)
	if err == nil && exist != nil {
		return nil, bizerr.NewBizError(bizcode.CodeUserAlreadyExist, "用户名已被注册")
	}

	// 获取用户
	exist, err = l.svcCtx.TUserModel.FindOne(l.ctx, "email = ?", in.Email)
	if err == nil && exist != nil {
		return nil, bizerr.NewBizError(bizcode.CodeUserAlreadyExist, "邮箱已被注册")
	}

	// 验证code是否正确
	key := rediskey.GetCaptchaKey(constant.CodeTypeRegister, in.Email)
	if !l.svcCtx.CaptchaHolder.VerifyCaptcha(key, in.VerifyCode) {
		return nil, bizerr.NewBizError(bizcode.CodeCaptchaVerify, "验证码错误")
	}

	var ua *model.TUser
	err = l.svcCtx.Gorm.Transaction(func(tx *gorm.DB) error {
		ua, err = l.register(tx, in)
		return err
	})
	if err != nil {
		return nil, err
	}

	user, err := l.svcCtx.TUserModel.FindOneByUserId(l.ctx, ua.UserId)
	if err != nil {
		return nil, err
	}

	// 查找用户角色
	rList, err := getUserRoles(l.ctx, l.svcCtx, user.UserId)
	if err != nil {
		return nil, err
	}

	resp := &accountrpc.RegisterResp{
		User: convertUserInfoOut(user, rList),
	}

	return resp, nil
}

func (l *RegisterLogic) register(tx *gorm.DB, in *accountrpc.RegisterReq) (out *model.TUser, err error) {
	ip, _ := rpcutils.GetRemoteIPFromCtx(l.ctx)
	is := ipx.GetIpSourceByBaidu(ip)

	// 邮箱注册
	user := &model.TUser{
		UserId:       randomx.GenerateUID(32),
		Username:     in.Username,
		Password:     cryptox.BcryptHash(in.Password),
		Nickname:     in.Email,
		Avatar:       "https://mms1.baidu.com/it/u=2815887849,1501151317&fm=253&app=138&f=JPEG",
		Email:        in.Email,
		Phone:        "",
		Info:         "",
		Status:       constant.UserStatusNormal,
		RegisterType: constant.LoginTypeEmail,
		IpAddress:    ip,
		IpSource:     is,
	}

	return onRegister(l.ctx, l.svcCtx, tx, user)
}

func onRegister(ctx context.Context, svcCtx *svc.ServiceContext, tx *gorm.DB, user *model.TUser) (out *model.TUser, err error) {
	/** 创建用户 **/
	_, err = svcCtx.TUserModel.WithTransaction(tx).Insert(ctx, user)
	if err != nil {
		return nil, err
	}

	// 查找默认用户角色
	roles, err := svcCtx.TRoleModel.WithTransaction(tx).FindALL(ctx, "is_default = ?", 1)
	if err != nil {
		return nil, err
	}

	var userRoles []*model.TUserRole
	for _, item := range roles {
		userRoles = append(userRoles, &model.TUserRole{
			UserId: user.UserId,
			RoleId: item.Id,
		})
	}

	/** 创建用户角色 **/
	_, err = svcCtx.TUserRoleModel.WithTransaction(tx).Inserts(ctx, userRoles...)
	if err != nil {
		return nil, err
	}

	return user, nil
}
