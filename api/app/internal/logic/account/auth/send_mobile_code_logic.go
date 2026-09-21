package auth

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type SendMobileCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发送手机验证码
func NewSendMobileCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMobileCodeLogic {
	return &SendMobileCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendMobileCodeLogic) SendMobileCode(req *types.SendMobileCodeReq) (resp *types.SendMobileCodeResp, err error) {
	_, err = l.svcCtx.NotificationService.SendMobileCode(l.ctx, &notificationservice.SendMobileCodeRequest{
		Mobile: req.Mobile,
		Scene:  req.Type,
		BizId:  "",
	})
	if err != nil {
		return nil, err
	}

	return &types.SendMobileCodeResp{}, nil
}
