package auth

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCaptchaCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取验证码
func NewGetCaptchaCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCaptchaCodeLogic {
	return &GetCaptchaCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCaptchaCodeLogic) GetCaptchaCode(req *types.GetCaptchaCodeReq) (resp *types.GetCaptchaCodeResp, err error) {
	in := &accountrpc.GenerateCaptchaCodeReq{
		Height: req.Height,
		Width:  req.Width,
	}

	out, err := l.svcCtx.AccountRpc.GenerateCaptchaCode(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return &types.GetCaptchaCodeResp{
		CaptchaKey:    out.CaptchaKey,
		CaptchaBase64: out.CaptchaBase64,
		CaptchaCode:   out.CaptchaCode,
	}, nil
}
