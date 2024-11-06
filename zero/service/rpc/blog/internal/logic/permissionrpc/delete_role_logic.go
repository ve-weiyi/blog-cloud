package permissionrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除角色
func (l *DeleteRoleLogic) DeleteRole(in *permissionrpc.IdsReq) (*permissionrpc.BatchResp, error) {
	rows, err := l.svcCtx.TRoleModel.Deletes(l.ctx, "id in (?) or parent_id in (?) ", in.Ids, in.Ids)
	if err != nil {
		return nil, err
	}

	return &permissionrpc.BatchResp{
		SuccessCount: rows,
	}, nil
}
