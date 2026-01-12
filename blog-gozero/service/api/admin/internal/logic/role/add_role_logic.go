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

func (l *AddRoleLogic) AddRole(req *types.NewRoleReq) (resp *types.RoleBackVO, err error) {
	in := &permissionrpc.NewRoleReq{
		Id:          req.Id,
		ParentId:    req.ParentId,
		RoleKey:     req.RoleKey,
		RoleLabel:   req.RoleLabel,
		RoleComment: req.RoleComment,
		IsDisable:   req.IsDisable,
		IsDefault:   req.IsDefault,
	}

	out, err := l.svcCtx.PermissionRpc.AddRole(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return &types.RoleBackVO{
		Id:          out.Id,
		ParentId:    out.ParentId,
		RoleKey:     out.RoleKey,
		RoleLabel:   out.RoleLabel,
		RoleComment: out.RoleComment,
		IsDisable:   out.IsDisable,
		IsDefault:   out.IsDefault,
		CreatedAt:   out.CreatedAt,
		UpdatedAt:   out.UpdatedAt,
	}, nil
}
