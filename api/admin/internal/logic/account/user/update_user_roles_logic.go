package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/permissionservice"
)

type UpdateUserRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新用户角色
func NewUpdateUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserRolesLogic {
	return &UpdateUserRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserRolesLogic) UpdateUserRoles(req *types.UpdateUserRolesReq) (resp *types.UpdateUserRolesResp, err error) {
	out, err := l.svcCtx.PermissionService.UpdateUserRoles(l.ctx, &permissionservice.UpdateUserRolesRequest{
		UserId:  req.UserId,
		RoleIds: req.RoleIds,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateUserRolesResp{
		Success: out.Success,
	}, nil
}
