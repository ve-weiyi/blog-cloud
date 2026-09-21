package userservicelogic

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/metax"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/userrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type UnbindMeThirdPartyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnbindMeThirdPartyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbindMeThirdPartyLogic {
	return &UnbindMeThirdPartyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 解绑当前用户第三方平台
func (l *UnbindMeThirdPartyLogic) UnbindMeThirdParty(in *userrpc.UnbindMeThirdPartyRequest) (*userrpc.UnbindMeThirdPartyResponse, error) {
	userId, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	oa, err := l.svcCtx.TUserOauthModel.FindOneByUserIdPlatform(l.ctx, userId, in.Platform)
	if err != nil {
		return nil, fmt.Errorf("未找到平台 %s 的绑定记录", in.Platform)
	}

	_, err = l.svcCtx.TUserOauthModel.Delete(l.ctx, oa.Id)
	if err != nil {
		return nil, err
	}

	return &userrpc.UnbindMeThirdPartyResponse{
		Success: true,
	}, nil
}
