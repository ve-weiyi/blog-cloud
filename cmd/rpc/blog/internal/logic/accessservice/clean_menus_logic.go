package accessservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/accessrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
)

type CleanMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCleanMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CleanMenusLogic {
	return &CleanMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 清空菜单
func (l *CleanMenusLogic) CleanMenus(in *accessrpc.CleanMenusRequest) (*accessrpc.CleanMenusResponse, error) {
	var rows int64
	err := l.svcCtx.GormDB.Transaction(func(tx *gorm.DB) error {
		if _, err := l.svcCtx.TRoleMenuModel.WithTx(tx).DeleteBatch(l.ctx, "1 = 1"); err != nil {
			return err
		}

		var err error
		rows, err = l.svcCtx.TMenuModel.WithTx(tx).DeleteBatch(l.ctx, "1 = 1")
		return err
	})
	if err != nil {
		return nil, err
	}

	return &accessrpc.CleanMenusResponse{SuccessCount: rows}, nil
}
