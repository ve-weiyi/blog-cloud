package role

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/permissionrpc"

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

func (l *FindRoleListLogic) FindRoleList(req *types.QueryRoleReq) (resp *types.PageResp, err error) {
	in := &permissionrpc.FindRoleListReq{
		Paginate: &permissionrpc.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		RoleKey:   req.RoleKey,
		RoleLabel: req.RoleLabel,
		Status:    req.Status,
	}

	out, err := l.svcCtx.PermissionRpc.FindRoleList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var list []*types.RoleBackVO
	for _, v := range out.List {
		list = append(list, convertRoleTypes(v))
	}

	resp = &types.PageResp{}
	resp.Page = 0
	resp.PageSize = int64(len(list))
	resp.Total = int64(len(list))
	resp.List = list
	return resp, nil
}
