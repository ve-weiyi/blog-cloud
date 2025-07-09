package accountrpclogic

import (
	"context"
	"fmt"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rpcutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindUserOauthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindUserOauthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindUserOauthLogic {
	return &BindUserOauthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改用户第三方账号
func (l *BindUserOauthLogic) BindUserOauth(in *accountrpc.BindUserOauthReq) (*accountrpc.EmptyResp, error) {
	// 查找当前用户是否存在
	userId, err := rpcutils.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	app, err := rpcutils.GetAppNameFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	auth, err := GetPlatformOauth(l.ctx, l.svcCtx, app, in.Platform)
	if err != nil {
		return nil, err
	}

	// 获取第三方用户信息
	info, err := auth.GetAuthUserInfo(in.Code)
	if err != nil {
		return nil, err
	}

	if info.OpenId == "" {
		return nil, fmt.Errorf("open_id is empty")
	}

	// 查找这个第三方账号是否已绑定用户
	oa, _ := l.svcCtx.TUserOauthModel.FindOneByOpenIdPlatform(l.ctx, info.OpenId, in.Platform)
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

	return &accountrpc.EmptyResp{}, nil
}
