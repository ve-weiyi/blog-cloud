package role

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/permissionrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindRoleResourcesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取角色资源列表
func NewFindRoleResourcesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindRoleResourcesLogic {
	return &FindRoleResourcesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindRoleResourcesLogic) FindRoleResources(req *types.IdReq) (resp *types.RoleResourcesResp, err error) {
	in := &permissionrpc.IdReq{
		Id: req.Id,
	}
	out, err := l.svcCtx.PermissionRpc.FindRoleResources(l.ctx, in)
	if err != nil {
		return
	}

	resp = &types.RoleResourcesResp{}
	resp.RoleId = out.RoleId
	resp.ApiIds = out.ApiIds
	resp.MenuIds = out.MenuIds

	return
}
