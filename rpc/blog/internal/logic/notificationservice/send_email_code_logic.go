package notificationservicelogic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/otpx"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/x/patternx"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
)

type SendEmailCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendEmailCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendEmailCodeLogic {
	return &SendEmailCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 发送邮件验证码
func (l *SendEmailCodeLogic) SendEmailCode(in *notificationrpc.SendEmailCodeRequest) (*notificationrpc.SendEmailCodeResponse, error) {
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

	// 设置过期时间，默认15分钟
	expireSeconds := in.ExpireSeconds
	if expireSeconds <= 0 {
		expireSeconds = 15 * 60
	}
	expire := time.Duration(expireSeconds) * time.Second

	// 生成6位验证码并存储到Redis
	code, err := l.svcCtx.OTPStore.Generate(l.ctx, key, otpx.WithLength(6), otpx.WithExpire(expire))
	if err != nil {
		return nil, err
	}

	// 构造 Email 消息事件
	emailEvent := &mq.EmailMessageEvent{
		Email: in.Email,
		Scene: in.Scene,
		BizId: in.BizId,
		Params: map[string]string{
			"code": code,
			"time": fmt.Sprintf("%d", expireSeconds/60),
		},
	}
	// MQ 未就绪时跳过投递；投递失败则向上报错
	if err = mq.PublishEmailMessageEvent(l.ctx, emailEvent); err != nil && !errors.Is(err, mq.ErrUnavailable) {
		return nil, err
	}

	return &notificationrpc.SendEmailCodeResponse{
		Id: 0, // 消费者会创建记录
	}, nil
}
