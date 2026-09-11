package permissionservicelogic

import (
	"context"
	"slices"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type GetUserRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserRolesLogic {
	return &GetUserRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询用户角色
func (l *GetUserRolesLogic) GetUserRoles(in *permissionrpc.GetUserRolesRequest) (*permissionrpc.GetUserRolesResponse, error) {
	urs, err := l.svcCtx.TUserRoleModel.FindALL(l.ctx, "user_id = ?", in.UserId)
	if err != nil {
		return nil, err
	}
	var roleIds []int64
	for _, v := range urs {
		if !slices.Contains(roleIds, v.RoleId) {
			roleIds = append(roleIds, v.RoleId)
		}
	}
	if len(roleIds) == 0 {
		return &permissionrpc.GetUserRolesResponse{}, nil
	}

	roles, err := l.svcCtx.TRoleModel.FindALL(l.ctx, "id in (?)", roleIds)
	if err != nil {
		return nil, err
	}

	out := &permissionrpc.GetUserRolesResponse{}
	for _, r := range roles {
		out.List = append(out.List, convertRoleOut(r))
	}

	return out, nil
}
