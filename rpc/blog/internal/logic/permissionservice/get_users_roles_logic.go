package permissionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type GetUsersRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUsersRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersRolesLogic {
	return &GetUsersRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量查询用户角色
func (l *GetUsersRolesLogic) GetUsersRoles(in *permissionrpc.GetUsersRolesRequest) (*permissionrpc.GetUsersRolesResponse, error) {
	if len(in.UserIds) == 0 {
		return &permissionrpc.GetUsersRolesResponse{}, nil
	}

	userRoles, err := l.svcCtx.TUserRoleModel.FindALL(l.ctx, "user_id IN (?)", in.UserIds)
	if err != nil {
		return nil, err
	}

	if len(userRoles) == 0 {
		return &permissionrpc.GetUsersRolesResponse{}, nil
	}

	var roleIds []int64
	roleIdSet := make(map[int64]bool)
	for _, ur := range userRoles {
		if !roleIdSet[ur.RoleId] {
			roleIdSet[ur.RoleId] = true
			roleIds = append(roleIds, ur.RoleId)
		}
	}

	roles, err := l.svcCtx.TRoleModel.FindALL(l.ctx, "id IN (?)", roleIds)
	if err != nil {
		return nil, err
	}

	roleMap := make(map[int64]*permissionrpc.Role)
	for _, r := range roles {
		roleMap[r.Id] = convertRoleOut(r)
	}

	userRoleMap := make(map[string][]*permissionrpc.Role)
	for _, ur := range userRoles {
		if role, ok := roleMap[ur.RoleId]; ok {
			userRoleMap[ur.UserId] = append(userRoleMap[ur.UserId], role)
		}
	}

	var list []*permissionrpc.UserRoles
	for _, uid := range in.UserIds {
		list = append(list, &permissionrpc.UserRoles{
			UserId: uid,
			List:   userRoleMap[uid],
		})
	}

	return &permissionrpc.GetUsersRolesResponse{
		List: list,
	}, nil
}
