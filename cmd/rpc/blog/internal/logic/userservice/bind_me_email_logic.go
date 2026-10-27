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

type BindMeEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindMeEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindMeEmailLogic {
	return &BindMeEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 绑定当前用户邮箱
func (l *BindMeEmailLogic) BindMeEmail(in *userrpc.BindMeEmailRequest) (*userrpc.BindMeEmailResponse, error) {
	userID, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil || userID == "" {
		return nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, "缺少用户身份")
	}
	if !validatex.IsEmail(in.Email) {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "邮箱格式不正确")
	}

	if err := l.verifyCode(in.Email, in.VerifyCode); err != nil {
		return nil, err
	}

	existing, err := l.svcCtx.TUserModel.FindOneByEmail(l.ctx, in.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil && existing.UserId != userID {
		return nil, bizerr.NewBizError(bizcode.CodeResourceAlreadyExist, "邮箱已被绑定")
	}

	if _, err := l.svcCtx.TUserModel.UpdateFields(l.ctx, map[string]interface{}{
		"email": in.Email,
	}, "user_id = ?", userID); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, bizerr.NewBizError(bizcode.CodeResourceAlreadyExist, "邮箱已被绑定")
		}
		return nil, err
	}

	return &userrpc.BindMeEmailResponse{Success: true}, nil
}

func (l *BindMeEmailLogic) verifyCode(email, code string) error {
	key := fmt.Sprintf("email:code:%s:%s", enums.CodeSceneBind, email)
	ok, err := l.svcCtx.OTPStore.Verify(l.ctx, key, code)
	if errors.Is(err, otpx.ErrTooManyAttempts) {
		return bizerr.NewBizError(bizcode.CodeVerifyCodeError, "验证码错误次数过多，请重新获取")
	}
	if err != nil {
		return err
	}
	if !ok {
		return bizerr.NewBizError(bizcode.CodeVerifyCodeError, "验证码错误或已过期")
	}
	return nil
}
