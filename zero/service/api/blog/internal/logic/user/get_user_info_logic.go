package user

import (
	"context"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/accountrpc"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户信息
func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.EmptyReq) (resp *types.UserInfoResp, err error) {
	in := &accountrpc.UserIdReq{
		UserId: cast.ToString(l.ctx.Value("uid")),
	}

	info, err := l.svcCtx.AccountRpc.GetUserInfo(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return ConvertUserInfoTypes(info), nil
}

func ConvertUserInfoTypes(in *accountrpc.UserInfoResp) (out *types.UserInfoResp) {
	out = &types.UserInfoResp{
		UserId:      in.UserId,
		Username:    in.Username,
		Nickname:    in.Nickname,
		Avatar:      in.Avatar,
		Email:       in.Email,
		Phone:       in.Phone,
		UserInfoExt: types.UserInfoExt{},
	}

	jsonconv.JsonToAny(in.Info, &out.UserInfoExt)

	return out
}
