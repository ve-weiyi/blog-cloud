package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/accessservice"
)

type UpdateRoleApiPermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新角色接口权限
func NewUpdateRoleApiPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleApiPermissionsLogic {
	return &UpdateRoleApiPermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRoleApiPermissionsLogic) UpdateRoleApiPermissions(req *types.UpdateRoleApiPermissionsReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.AccessService.UpdateRoleApi(l.ctx, &accessservice.UpdateRoleApiRequest{
		RoleId: req.RoleId,
		ApiIds: req.ApiIds,
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
