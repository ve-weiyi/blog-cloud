package role

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/permissionrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindRoleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页获取角色列表
func NewFindRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindRoleListLogic {
	return &FindRoleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindRoleListLogic) FindRoleList(req *types.RoleQuery) (resp *types.PageResp, err error) {
	in := &permissionrpc.FindRoleListReq{
		Page:      req.Page,
		PageSize:  req.PageSize,
		RoleName:  req.RoleName,
		RoleLabel: req.RoleLabel,
		IsDisable: req.IsDisable,
	}

	out, err := l.svcCtx.PermissionRpc.FindRoleList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var list []*types.RoleBackDTO
	for _, v := range out.List {
		m := convertRoleTypes(v)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = in.Page
	resp.PageSize = in.PageSize
	resp.Total = int64(len(list))
	resp.List = list
	return resp, nil
}
