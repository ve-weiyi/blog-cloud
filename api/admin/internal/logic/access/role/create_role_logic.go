package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/accessservice"
)

type CreateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建角色
func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRoleLogic) CreateRole(req *types.CreateRoleReq) (resp *types.RoleVO, err error) {
	out, err := l.svcCtx.AccessService.CreateRole(l.ctx, &accessservice.CreateRoleRequest{
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

	return &types.RoleVO{
		Id:          out.Id,
		ParentId:    req.ParentId,
		RoleKey:     req.RoleKey,
		RoleLabel:   req.RoleLabel,
		RoleComment: req.RoleComment,
		IsDefault:   req.IsDefault,
		Status:      req.Status,
	}, nil
}
