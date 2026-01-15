package permissionrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletesRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesRoleLogic {
	return &DeletesRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除角色
func (l *DeletesRoleLogic) DeletesRole(in *permissionrpc.DeletesRoleReq) (*permissionrpc.DeletesRoleResp, error) {
	rows, err := l.svcCtx.TRoleModel.Deletes(l.ctx, "id in (?) or parent_id in (?) ", in.Ids, in.Ids)
	if err != nil {
		return nil, err
	}

	return &permissionrpc.DeletesRoleResp{
		SuccessCount: rows,
	}, nil
}
