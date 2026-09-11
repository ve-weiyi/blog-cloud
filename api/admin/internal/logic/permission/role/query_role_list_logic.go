package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/permissionservice"
)

type QueryRoleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取角色列表
func NewQueryRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRoleListLogic {
	return &QueryRoleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryRoleListLogic) QueryRoleList(req *types.QueryRoleListReq) (resp *types.PageResult, err error) {
	out, err := l.svcCtx.PermissionService.ListRoles(l.ctx, &permissionservice.ListRolesRequest{})
	if err != nil {
		return nil, err
	}

	var list []*types.RoleVO
	for _, v := range out.List {
		list = append(list, &types.RoleVO{
			Id:          v.Id,
			ParentId:    v.ParentId,
			RoleKey:     v.RoleKey,
			RoleLabel:   v.RoleLabel,
			RoleComment: v.RoleComment,
			IsDefault:   v.IsDefault,
			Status:      v.Status,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}

	return &types.PageResult{
		Page:     1,
		PageSize: int64(len(list)),
		Total:    int64(len(list)),
		List:     list,
	}, nil
}
