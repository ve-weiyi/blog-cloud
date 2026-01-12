package accountrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserOauthInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserOauthInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserOauthInfoLogic {
	return &GetUserOauthInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户第三平台信息
func (l *GetUserOauthInfoLogic) GetUserOauthInfo(in *accountrpc.GetUserOauthInfoReq) (*accountrpc.GetUserOauthInfoResp, error) {
	oas, err := l.svcCtx.TUserOauthModel.FindALL(l.ctx, "user_id = ?", in.UserId)
	if err != nil {
		return nil, err
	}

	var list []*accountrpc.UserOauthInfo
	for _, v := range oas {
		list = append(list, &accountrpc.UserOauthInfo{
			Platform:  v.Platform,
			OpenId:    v.OpenId,
			Nickname:  v.Nickname,
			Avatar:    v.Avatar,
			CreatedAt: v.CreatedAt.UnixMilli(),
			UpdatedAt: v.UpdatedAt.UnixMilli(),
		})
	}

	return &accountrpc.GetUserOauthInfoResp{
		List: list,
	}, nil
}
