package notificationservicelogic

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/x/patternx"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
)

type VerifyEmailCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyEmailCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyEmailCodeLogic {
	return &VerifyEmailCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 校验邮件验证码
func (l *VerifyEmailCodeLogic) VerifyEmailCode(in *notificationrpc.VerifyEmailCodeRequest) (*notificationrpc.VerifyEmailCodeResponse, error) {
	// 校验邮箱格式
	if !patternx.IsValidEmail(in.Email) {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "邮箱格式不正确")
	}

	// 生成验证码存储的key
	var key string
	if in.BizId != "" {
		key = fmt.Sprintf("email:code:%s", in.BizId)
	} else {
		key = fmt.Sprintf("email:code:%s:%s", in.Scene, in.Email)
	}

	// 验证验证码
	success, err := l.svcCtx.OTPStore.Verify(l.ctx, key, in.Code)
	if err != nil {
		return &notificationrpc.VerifyEmailCodeResponse{
			Success: false,
			Message: "验证码验证失败",
		}, nil
	}

	if !success {
		return &notificationrpc.VerifyEmailCodeResponse{
			Success: false,
			Message: "验证码错误或已过期",
		}, nil
	}

	return &notificationrpc.VerifyEmailCodeResponse{
		Success: true,
		Message: "验证成功",
	}, nil
}
