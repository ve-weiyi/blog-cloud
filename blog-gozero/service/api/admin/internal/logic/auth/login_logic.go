package auth

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {

	in := accountrpc.LoginReq{
		Username: req.Username,
		Password: req.Password,
	}

	out, err := l.svcCtx.AccountRpc.Login(l.ctx, &in)
	if err != nil {
		return
	}

	// 登录日志
	_, err = l.svcCtx.SyslogRpc.AddLoginLog(l.ctx, &syslogrpc.AddLoginLogReq{
		UserId:    out.User.UserId,
		LoginType: out.LoginType,
	})

	return onLogin(l.ctx, l.svcCtx, out)
}

func onLogin(ctx context.Context, svcCtx *svc.ServiceContext, login *accountrpc.LoginResp) (resp *types.LoginResp, err error) {
	tk, err := svcCtx.TokenManager.GenerateToken(login.User.UserId)
	if err != nil {
		return nil, err
	}

	token := &types.Token{
		TokenType:        tk.TokenType,
		AccessToken:      tk.AccessToken,
		ExpiresIn:        tk.ExpiresIn,
		RefreshToken:     tk.RefreshToken,
		RefreshExpiresIn: tk.RefreshExpiresIn,
		RefreshExpiresAt: tk.RefreshExpiresAt,
	}

	return &types.LoginResp{
		UserId: login.User.UserId,
		Scope:  svcCtx.Config.Name,
		Token:  token,
	}, nil
}
