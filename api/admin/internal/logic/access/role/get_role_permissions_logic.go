package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/accessservice"
)

type GetRolePermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询角色权限配置
func NewGetRolePermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRolePermissionsLogic {
	return &GetRolePermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRolePermissionsLogic) GetRolePermissions(req *types.GetRolePermissionsReq) (resp *types.RolePermissionsResp, err error) {
	out, err := l.svcCtx.AccessService.GetRoleResource(l.ctx, &accessservice.GetRoleResourceRequest{
		RoleId: req.RoleId,
	})
	if err != nil {
		return nil, err
	}

	var apiIds []int64
	for _, v := range out.ApiIds {
		apiIds = append(apiIds, v)
	}

	var menuIds []int64
	for _, v := range out.MenuIds {
		menuIds = append(menuIds, v)
	}

	return &types.RolePermissionsResp{
		ApiIds:  apiIds,
		MenuIds: menuIds,
	}, nil
}
