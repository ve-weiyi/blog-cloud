package permissionrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateRoleLogic) UpdateRole(in *permissionrpc.UpdateRoleReq) (*permissionrpc.UpdateRoleResp, error) {
	entity, err := l.svcCtx.TRoleModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	entity.ParentId = in.ParentId
	entity.RoleKey = in.RoleKey
	entity.RoleLabel = in.RoleLabel
	entity.RoleComment = in.RoleComment
	entity.Status = in.Status
	entity.IsDefault = in.IsDefault

	_, err = l.svcCtx.TRoleModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &permissionrpc.UpdateRoleResp{
		Role: convertRoleOut(entity),
	}, nil
}
