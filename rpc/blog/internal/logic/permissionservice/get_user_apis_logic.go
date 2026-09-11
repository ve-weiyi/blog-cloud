package permissionservicelogic

import (
	"context"
	"slices"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type GetUserApisLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserApisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserApisLogic {
	return &GetUserApisLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询用户 API 权限
func (l *GetUserApisLogic) GetUserApis(in *permissionrpc.GetUserApisRequest) (*permissionrpc.GetUserApisResponse, error) {
	urs, err := l.svcCtx.TUserRoleModel.FindALL(l.ctx, "user_id = ?", in.UserId)
	if err != nil {
		return nil, err
	}
	var roleIds []int64
	for _, v := range urs {
		roleIds = append(roleIds, v.RoleId)
	}
	if len(roleIds) == 0 {
		return &permissionrpc.GetUserApisResponse{}, nil
	}

	rLinks, err := l.svcCtx.TRoleApiModel.FindALL(l.ctx, "role_id in (?)", roleIds)
	if err != nil {
		return nil, err
	}
	var apiIds []int64
	for _, v := range rLinks {
		if !slices.Contains(apiIds, v.ApiId) {
			apiIds = append(apiIds, v.ApiId)
		}
	}
	if len(apiIds) == 0 {
		return &permissionrpc.GetUserApisResponse{}, nil
	}

	apis, err := l.svcCtx.TApiModel.FindALL(l.ctx, "id in (?)", apiIds)
	if err != nil {
		return nil, err
	}

	out := &permissionrpc.GetUserApisResponse{}
	for _, a := range apis {
		out.List = append(out.List, convertApiOut(a))
	}

	return out, nil
}
