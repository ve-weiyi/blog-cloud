package userservicelogic

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/metax"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/userrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

type BindMeThirdPartyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindMeThirdPartyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindMeThirdPartyLogic {
	return &BindMeThirdPartyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 绑定当前用户第三方平台
func (l *BindMeThirdPartyLogic) BindMeThirdParty(in *userrpc.BindMeThirdPartyRequest) (*userrpc.BindMeThirdPartyResponse, error) {
	// 查找当前用户是否存在
	userId, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	auth, ok := l.svcCtx.OAuthProviders[in.Platform]
	if !ok {
		return nil, fmt.Errorf("platform %s is not support", in.Platform)
	}

	// 获取第三方用户信息
	info, err := auth.GetAuthUserInfo(l.ctx, in.Code)
	if err != nil {
		return nil, err
	}

	if info.OpenId == "" {
		return nil, fmt.Errorf("open_id is empty")
	}

	// 查找这个第三方账号是否已绑定用户
	oa, _ := l.svcCtx.TUserOauthModel.FindOneByPlatformOpenId(l.ctx, in.Platform, info.OpenId)
	if oa != nil {
		return nil, fmt.Errorf("open_id %s is already exist", info.OpenId)
	}

	// 绑定第三方账号
	_, err = l.svcCtx.TUserOauthModel.Insert(l.ctx, &model.TUserOauth{
		Id:       0,
		UserId:   userId,
		Platform: in.Platform,
		OpenId:   info.OpenId,
		Nickname: info.NickName,
		Avatar:   info.Avatar,
	})

	return &userrpc.BindMeThirdPartyResponse{
		Success: true,
	}, nil
}
