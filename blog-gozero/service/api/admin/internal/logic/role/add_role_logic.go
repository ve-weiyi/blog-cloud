package role

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/permissionrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建角色
func NewAddRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRoleLogic {
	return &AddRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddRoleLogic) AddRole(req *types.RoleNewReq) (resp *types.RoleBackDTO, err error) {
	in := &permissionrpc.RoleNewReq{
		Id:          req.Id,
		ParentId:    req.ParentId,
		RoleName:    req.RoleName,
		RoleLabel:   req.RoleLabel,
		RoleComment: req.RoleComment,
		IsDisable:   req.IsDisable,
		IsDefault:   req.IsDefault,
	}

	out, err := l.svcCtx.PermissionRpc.AddRole(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return convertRoleTypes(out), nil
}

func convertRoleTypes(req *permissionrpc.RoleDetails) *types.RoleBackDTO {
	out := &types.RoleBackDTO{
		Id:          req.Id,
		ParentId:    req.ParentId,
		RoleName:    req.RoleName,
		RoleLabel:   req.RoleLabel,
		RoleComment: req.RoleComment,
		IsDisable:   req.IsDisable,
		IsDefault:   req.IsDefault,
		CreatedAt:   req.CreatedAt,
		UpdatedAt:   req.UpdatedAt,
	}

	return out
}
