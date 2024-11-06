package permissionrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

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
func (l *AddRoleLogic) AddRole(in *permissionrpc.RoleNewReq) (*permissionrpc.RoleDetails, error) {
	entity := convertRoleIn(in)

	_, err := l.svcCtx.TRoleModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertRoleOut(entity), nil
}
