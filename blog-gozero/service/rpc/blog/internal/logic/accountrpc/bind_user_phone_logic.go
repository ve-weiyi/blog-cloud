package accountrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/global/constant"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rediskey"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rpcutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/biz/bizerr"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/patternx"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindUserPhoneLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindUserPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindUserPhoneLogic {
	return &BindUserPhoneLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改用户登录手机号
func (l *BindUserPhoneLogic) BindUserPhone(in *accountrpc.BindUserPhoneReq) (*accountrpc.EmptyResp, error) {
	// 校验邮箱格式
	if !patternx.IsValidPhone(in.Phone) {
		return nil, bizerr.NewBizError(bizerr.CodeInvalidParam, "手机号格式不正确")
	}

	userId, err := rpcutils.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// 验证用户是否存在
	user, err := l.svcCtx.TUserModel.FindOneByUserId(l.ctx, userId)
	if err != nil {
		return nil, bizerr.NewBizError(bizerr.CodeUserNotExist, err.Error())
	}

	// 验证code是否正确
	key := rediskey.GetCaptchaKey(constant.CodeTypeBindPhone, in.Phone)
	if !l.svcCtx.CaptchaHolder.VerifyCaptcha(key, in.VerifyCode) {
		return nil, bizerr.NewBizError(bizerr.CodeCaptchaVerify, "验证码错误")
	}

	// 更新密码
	user.Phone = in.Phone

	_, err = l.svcCtx.TUserModel.Save(l.ctx, user)
	if err != nil {
		return nil, err
	}

	return &accountrpc.EmptyResp{}, nil
}
