package accessservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/accessrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchGetUserRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetUserRolesLogic {
	return &BatchGetUserRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量查询用户角色
func (l *BatchGetUserRolesLogic) BatchGetUserRoles(in *accessrpc.BatchGetUserRolesRequest) (*accessrpc.BatchGetUserRolesResponse, error) {
	if len(in.UserIds) == 0 {
		return &accessrpc.BatchGetUserRolesResponse{}, nil
	}

	userRoles, err := l.svcCtx.TUserRoleModel.FindALL(l.ctx, "user_id IN (?)", in.UserIds)
	if err != nil {
		return nil, err
	}

	if len(userRoles) == 0 {
		return &accessrpc.BatchGetUserRolesResponse{}, nil
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

	roleMap := make(map[int64]*accessrpc.Role)
	for _, r := range roles {
		roleMap[r.Id] = convertRoleOut(r)
	}

	userRoleMap := make(map[string][]*accessrpc.Role)
	for _, ur := range userRoles {
		if role, ok := roleMap[ur.RoleId]; ok {
			userRoleMap[ur.UserId] = append(userRoleMap[ur.UserId], role)
		}
	}

	var list []*accessrpc.UserRoles
	for _, uid := range in.UserIds {
		list = append(list, &accessrpc.UserRoles{
			UserId: uid,
			List:   userRoleMap[uid],
		})
	}

	return &accessrpc.BatchGetUserRolesResponse{
		List: list,
	}, nil
}
