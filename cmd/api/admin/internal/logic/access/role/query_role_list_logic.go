package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/accessservice"
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

func (l *QueryRoleListLogic) QueryRoleList(req *types.QueryRoleListReq) (resp *types.ListResult, err error) {
	out, err := l.svcCtx.AccessService.ListRoles(l.ctx, &accessservice.ListRolesRequest{})
	if err != nil {
		return nil, err
	}

	var list []*types.RoleVO
	for _, v := range out.List {
		list = append(list, convertRoleVO(v))
	}

	return &types.ListResult{
		Page:     1,
		PageSize: int64(len(list)),
		Total:    int64(len(list)),
		List:     list,
	}, nil
}
