package accessservicelogic

import (
	"context"
	"slices"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/accessrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type ListUserRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserRolesLogic {
	return &ListUserRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询用户角色
func (l *ListUserRolesLogic) ListUserRoles(in *accessrpc.ListUserRolesRequest) (*accessrpc.ListUserRolesResponse, error) {
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
		return &accessrpc.ListUserRolesResponse{}, nil
	}

	roles, err := l.svcCtx.TRoleModel.FindALL(l.ctx, "id in (?)", roleIds)
	if err != nil {
		return nil, err
	}

	out := &accessrpc.ListUserRolesResponse{}
	for _, r := range roles {
		out.List = append(out.List, convertRoleOut(r))
	}

	return out, nil
}
