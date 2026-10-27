package userservicelogic

import (
	"context"
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/userrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/auth/oauthstate"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
	"github.com/ve-weiyi/blog-cloud/model"
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

	nonce := ""
	if in.State != nil {
		nonce = *in.State
	}
	state, err := oauthstate.Consume(l.ctx, l.svcCtx.Redis, nonce)
	if err != nil || state.Platform != in.Platform {
		return nil, bizerr.NewBizError(bizcode.CodeRequestSignInvalid, "OAuth state 无效或已过期")
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
	oa, err := l.svcCtx.TUserOauthModel.FindOneByOpenIdPlatform(l.ctx, info.OpenId, in.Platform)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if oa != nil {
		return nil, fmt.Errorf("open_id %s is already exist", info.OpenId)
	}

	// 绑定第三方账号
	if _, err = l.svcCtx.TUserOauthModel.Insert(l.ctx, &model.TUserOauth{
		Id:       0,
		UserId:   userId,
		Platform: in.Platform,
		OpenId:   info.OpenId,
		Nickname: info.NickName,
		Avatar:   info.Avatar,
	}); err != nil {
		return nil, err
	}

	return &userrpc.BindMeThirdPartyResponse{
		Success: true,
	}, nil
}
