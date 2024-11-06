package account

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/permissionrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAccountRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改用户角色
func NewUpdateAccountRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAccountRolesLogic {
	return &UpdateAccountRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAccountRolesLogic) UpdateAccountRoles(req *types.UpdateAccountRolesReq) (resp *types.EmptyResp, err error) {
	in := &permissionrpc.UpdateUserRoleReq{
		UserId:  req.UserId,
		RoleIds: req.RoleIds,
	}

	_, err = l.svcCtx.PermissionRpc.UpdateUserRole(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
