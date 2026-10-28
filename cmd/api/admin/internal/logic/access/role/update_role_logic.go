package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/accessservice"
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

func (l *UpdateRoleLogic) UpdateRole(req *types.UpdateRoleReq) (resp *types.RoleVO, err error) {
	out, err := l.svcCtx.AccessService.UpdateRole(l.ctx, &accessservice.UpdateRoleRequest{
		Id:          req.Id,
		ParentId:    req.ParentId,
		RoleKey:     req.RoleKey,
		RoleLabel:   req.RoleLabel,
		RoleComment: req.RoleComment,
		IsDefault:   req.IsDefault,
		Status:      req.Status,
	})
	if err != nil {
		return nil, err
	}

	return convertRoleVO(out.Role), nil
}
