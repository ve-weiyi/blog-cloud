package role

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/permissionrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新角色
func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRoleLogic) UpdateRole(req *types.RoleNewReq) (resp *types.RoleBackDTO, err error) {
	in := &permissionrpc.RoleNewReq{
		Id:          req.Id,
		ParentId:    req.ParentId,
		RoleName:    req.RoleName,
		RoleLabel:   req.RoleLabel,
		RoleComment: req.RoleComment,
		IsDisable:   req.IsDisable,
		IsDefault:   req.IsDefault,
	}

	out, err := l.svcCtx.PermissionRpc.UpdateRole(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return convertRoleTypes(out), nil
}
