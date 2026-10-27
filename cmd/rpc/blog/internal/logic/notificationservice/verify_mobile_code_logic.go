package notificationservicelogic

import (
	"context"
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/x/validatex"

	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/store/otpx"
)

type VerifyMobileCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyMobileCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyMobileCodeLogic {
	return &VerifyMobileCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 校验短信验证码
func (l *VerifyMobileCodeLogic) VerifyMobileCode(in *notificationrpc.VerifyMobileCodeRequest) (*notificationrpc.VerifyMobileCodeResponse, error) {
	// 校验手机号格式
	if !validatex.IsMobile(in.Mobile) {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "手机号格式不正确")
	}

	// 生成验证码存储的key
	var key string
	if in.BizId != "" {
		key = fmt.Sprintf("sms:code:%s", in.BizId)
	} else {
		key = fmt.Sprintf("sms:code:%s:%s", in.Scene, in.Mobile)
	}

	// 验证失败必须返回错误，不能再用 Success:false 表达安全校验结果
	ok, err := l.svcCtx.OTPStore.Verify(l.ctx, key, in.Code)
	if errors.Is(err, otpx.ErrTooManyAttempts) {
		return nil, bizerr.NewBizError(bizcode.CodeVerifyCodeError, "验证码错误次数过多，请重新获取")
	}
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, bizerr.NewBizError(bizcode.CodeVerifyCodeError, "验证码错误或已过期")
	}

	return &notificationrpc.VerifyMobileCodeResponse{
		Success: true,
		Message: "验证成功",
	}, nil
}
