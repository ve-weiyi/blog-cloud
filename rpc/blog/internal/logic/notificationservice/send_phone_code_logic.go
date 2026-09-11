package notificationservicelogic

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/adapter/mqx"
	"github.com/ve-weiyi/vkit/x/jsonconv"
	"github.com/ve-weiyi/vkit/x/patternx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/pb/notificationrpc"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
)

type SendPhoneCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendPhoneCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendPhoneCodeLogic {
	return &SendPhoneCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 发送短信验证码
func (l *SendPhoneCodeLogic) SendPhoneCode(in *notificationrpc.SendPhoneCodeRequest) (*notificationrpc.SendPhoneCodeResponse, error) {
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

	// 设置过期时间，默认15分钟
	expireSeconds := in.ExpireSeconds
	if expireSeconds <= 0 {
		expireSeconds = 15 * 60
	}
	expire := time.Duration(expireSeconds) * time.Second

	// 生成6位验证码并存储到Redis
	code, err := l.svcCtx.CodeStore.Generate(key, 6, expire)
	if err != nil {
		return nil, err
	}

	// 构造 SMS 消息事件
	if mq.SmsProducer != nil {
		// 构造 SMS 消息事件
		smsEvent := &mq.SmsMessageEvent{
			Mobile: in.Phone,
			Scene:  in.Scene,
			BizId:  in.BizId,
			Params: map[string]string{
				"code": code,
				"time": fmt.Sprintf("%d", expireSeconds/60),
			},
		}
		err = mq.SmsProducer.Send(l.ctx, &mqx.Message{
			Topic:     mq.SmsQueue,
			Key:       mq.SmsRoutingKey,
			Body:      []byte(jsonconv.AnyToJsonNE(smsEvent)),
			Timestamp: time.Now(),
		})
		if err != nil {
			return nil, err
		}
	}

	return &notificationrpc.SendPhoneCodeResponse{
		Id:   0,    // 消费者会创建记录
		Code: code, // 返回验证码用于测试，生产环境可以移除
	}, nil
}
