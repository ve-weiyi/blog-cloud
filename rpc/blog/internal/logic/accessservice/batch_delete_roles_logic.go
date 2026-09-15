package accessservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/accessrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchDeleteRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteRolesLogic {
	return &BatchDeleteRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除角色
func (l *BatchDeleteRolesLogic) BatchDeleteRoles(in *accessrpc.BatchDeleteRolesRequest) (*accessrpc.BatchDeleteRolesResponse, error) {
	if len(in.Ids) == 0 {
		return &accessrpc.BatchDeleteRolesResponse{SuccessCount: 0}, nil
	}

	_, err := l.svcCtx.TRoleApiModel.DeleteBatch(l.ctx, "role_id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.TRoleMenuModel.DeleteBatch(l.ctx, "role_id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.TUserRoleModel.DeleteBatch(l.ctx, "role_id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	rows, err := l.svcCtx.TRoleModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &accessrpc.BatchDeleteRolesResponse{SuccessCount: rows}, nil
}
