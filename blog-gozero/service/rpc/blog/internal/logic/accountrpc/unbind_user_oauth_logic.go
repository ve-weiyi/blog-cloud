package accountrpclogic

import (
	"context"
	"fmt"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rpcutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnbindUserOauthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnbindUserOauthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbindUserOauthLogic {
	return &UnbindUserOauthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 解绑第三方账号
func (l *UnbindUserOauthLogic) UnbindUserOauth(in *accountrpc.UnbindUserOauthReq) (*accountrpc.UnbindUserOauthResp, error) {
	// 查找当前用户是否存在
	userId, err := rpcutils.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// 查找这个第三方账号是否已绑定用户
	oa, _ := l.svcCtx.TUserOauthModel.FindOneByUserIdPlatform(l.ctx, userId, in.Platform)
	if oa == nil {
		return nil, fmt.Errorf("用户未绑定该第三方账号")
	}

	// 绑定第三方账号
	_, err = l.svcCtx.TUserOauthModel.Delete(l.ctx, oa.Id)
	if err != nil {
		return nil, err
	}

	return &accountrpc.UnbindUserOauthResp{}, nil
}
