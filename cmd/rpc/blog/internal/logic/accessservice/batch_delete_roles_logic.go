package accessservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/accessrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
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

	var rows int64
	err := l.svcCtx.GormDB.Transaction(func(tx *gorm.DB) error {
		if _, err := l.svcCtx.TRoleApiModel.WithTx(tx).DeleteBatch(l.ctx, "role_id in (?)", in.Ids); err != nil {
			return err
		}
		if _, err := l.svcCtx.TRoleMenuModel.WithTx(tx).DeleteBatch(l.ctx, "role_id in (?)", in.Ids); err != nil {
			return err
		}
		if _, err := l.svcCtx.TUserRoleModel.WithTx(tx).DeleteBatch(l.ctx, "role_id in (?)", in.Ids); err != nil {
			return err
		}

		var err error
		rows, err = l.svcCtx.TRoleModel.WithTx(tx).DeleteBatch(l.ctx, "id in (?)", in.Ids)
		return err
	})
	if err != nil {
		return nil, err
	}

	return &accessrpc.BatchDeleteRolesResponse{SuccessCount: rows}, nil
}
