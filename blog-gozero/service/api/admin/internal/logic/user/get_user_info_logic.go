package user

import (
	"context"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizheader"

	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/jsonconv"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"

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
	userId := cast.ToString(l.ctx.Value(bizheader.HeaderUid))
	in := &accountrpc.UserIdReq{
		UserId: userId,
	}

	info, err := l.svcCtx.AccountRpc.GetUserInfo(l.ctx, in)
	if err != nil {
		return nil, err
	}

	thp, err := l.svcCtx.AccountRpc.GetUserOauthInfo(l.ctx, in)
	if err != nil {
		return nil, err
	}

	ur, err := l.svcCtx.PermissionRpc.FindUserRoles(l.ctx, &permissionrpc.UserIdReq{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	up, err := l.svcCtx.PermissionRpc.FindUserApis(l.ctx, &permissionrpc.UserIdReq{
		UserId: userId,
	})

	return ConvertUserInfoTypes(info, thp, ur, up), nil
}

func ConvertUserInfoTypes(in *accountrpc.UserInfoResp, thp *accountrpc.GetUserOauthInfoResp, ur *permissionrpc.FindRoleListResp, up *permissionrpc.FindApiListResp) (out *types.UserInfoResp) {
	var info types.UserInfoExt
	jsonconv.JsonToAny(in.Info, &info)

	thirdParty := make([]*types.UserThirdPartyInfo, 0)
	for _, v := range thp.List {
		thirdParty = append(thirdParty, &types.UserThirdPartyInfo{
			Platform:  v.Platform,
			OpenId:    v.OpenId,
			Nickname:  v.Nickname,
			Avatar:    v.Avatar,
			CreatedAt: v.CreatedAt,
		})
	}

	roles := make([]string, 0)
	for _, v := range in.Roles {
		roles = append(roles, v.RoleKey)
	}

	perms := make([]string, 0)
	for _, v := range up.List {
		perms = append(perms, v.Path)
	}

	out = &types.UserInfoResp{
		UserId:       in.UserId,
		Username:     in.Username,
		Nickname:     in.Nickname,
		Avatar:       in.Avatar,
		Email:        in.Email,
		Phone:        in.Phone,
		RegisterType: in.RegisterType,
		CreatedAt:    in.CreatedAt,
		UserInfoExt:  info,
		ThirdParty:   thirdParty,
		Roles:        roles,
		Perms:        perms,
	}

	return out
}
