package notificationservicelogic

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	notificationrpc2 "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/x/patternx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/pb/notificationrpc"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
)

type VerifyPhoneCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyPhoneCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyPhoneCodeLogic {
	return &VerifyPhoneCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 校验短信验证码
func (l *VerifyPhoneCodeLogic) VerifyPhoneCode(in *notificationrpc2.VerifyPhoneCodeRequest) (*notificationrpc.VerifyPhoneCodeResponse, error) {
	// 校验手机号格式
	if !patternx.IsValidMobile(in.Phone) {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "手机号格式不正确")
	}

	// 生成验证码存储的key
	var key string
	if in.BizId != "" {
		key = fmt.Sprintf("sms:code:%s", in.BizId)
	} else {
		key = fmt.Sprintf("sms:code:%s:%s", in.Scene, in.Phone)
	}

	// 验证验证码
	success, err := l.svcCtx.CodeStore.Verify(key, in.Code)
	if err != nil {
		return &notificationrpc.VerifyPhoneCodeResponse{
			Success: false,
			Message: "验证码验证失败",
		}, nil
	}

	if !success {
		return &notificationrpc.VerifyPhoneCodeResponse{
			Success: false,
			Message: "验证码错误或已过期",
		}, nil
	}

	return &notificationrpc.VerifyPhoneCodeResponse{
		Success: true,
		Message: "验证成功",
	}, nil
}
