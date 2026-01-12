package permissionrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRoleLogic {
	return &AddRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建角色
func (l *AddRoleLogic) AddRole(in *permissionrpc.AddRoleReq) (*permissionrpc.AddRoleResp, error) {
	entity := convertRoleIn(in)

	_, err := l.svcCtx.TRoleModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &permissionrpc.AddRoleResp{
		Role: convertRoleOut(entity),
	}, nil
}

func convertRoleIn(in *permissionrpc.AddRoleReq) (out *model.TRole) {
	out = &model.TRole{
		Id:          in.Id,
		ParentId:    in.ParentId,
		RoleKey:     in.RoleKey,
		RoleLabel:   in.RoleLabel,
		RoleComment: in.RoleComment,
		Status:      in.Status,
		IsDefault:   in.IsDefault,
	}

	return out
}

func convertRoleOut(in *model.TRole) (out *permissionrpc.Role) {
	out = &permissionrpc.Role{
		Id:          in.Id,
		ParentId:    in.ParentId,
		RoleKey:     in.RoleKey,
		RoleLabel:   in.RoleLabel,
		RoleComment: in.RoleComment,
		Status:      in.Status,
		IsDefault:   in.IsDefault,
		CreatedAt:   in.CreatedAt.UnixMilli(),
		UpdatedAt:   in.UpdatedAt.UnixMilli(),
	}

	return out
}
