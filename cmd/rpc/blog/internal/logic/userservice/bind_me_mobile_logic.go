package userservicelogic

import (
	"context"
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/userrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
	"github.com/ve-weiyi/blog-cloud/infra/store/otpx"
	"github.com/ve-weiyi/vkit/x/validatex"
)

type BindMeMobileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindMeMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindMeMobileLogic {
	return &BindMeMobileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 绑定当前用户手机号
func (l *BindMeMobileLogic) BindMeMobile(in *userrpc.BindMeMobileRequest) (*userrpc.BindMeMobileResponse, error) {
	userID, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil || userID == "" {
		return nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, "缺少用户身份")
	}
	if !validatex.IsMobile(in.Mobile) {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "手机号格式不正确")
	}

	key := fmt.Sprintf("sms:code:%s:%s", enums.CodeSceneBind, in.Mobile)
	ok, err := l.svcCtx.OTPStore.Verify(l.ctx, key, in.VerifyCode)
	if errors.Is(err, otpx.ErrTooManyAttempts) {
		return nil, bizerr.NewBizError(bizcode.CodeVerifyCodeError, "验证码错误次数过多，请重新获取")
	}
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, bizerr.NewBizError(bizcode.CodeVerifyCodeError, "验证码错误或已过期")
	}

	existing, err := l.svcCtx.TUserModel.FindOneByMobile(l.ctx, in.Mobile)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil && existing.UserId != userID {
		return nil, bizerr.NewBizError(bizcode.CodeResourceAlreadyExist, "手机号已被绑定")
	}

	if _, err := l.svcCtx.TUserModel.UpdateFields(l.ctx, map[string]interface{}{
		"mobile": in.Mobile,
	}, "user_id = ?", userID); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, bizerr.NewBizError(bizcode.CodeResourceAlreadyExist, "手机号已被绑定")
		}
		return nil, err
	}

	return &userrpc.BindMeMobileResponse{Success: true}, nil
}
