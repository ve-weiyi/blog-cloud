package accessservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/accessrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
)

type CleanApisLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCleanApisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CleanApisLogic {
	return &CleanApisLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 清空 API
func (l *CleanApisLogic) CleanApis(in *accessrpc.CleanApisRequest) (*accessrpc.CleanApisResponse, error) {
	var rows int64
	err := l.svcCtx.GormDB.Transaction(func(tx *gorm.DB) error {
		if _, err := l.svcCtx.TRoleApiModel.WithTx(tx).DeleteBatch(l.ctx, "1 = 1"); err != nil {
			return err
		}

		var err error
		rows, err = l.svcCtx.TApiModel.WithTx(tx).DeleteBatch(l.ctx, "1 = 1")
		return err
	})
	if err != nil {
		return nil, err
	}

	return &accessrpc.CleanApisResponse{SuccessCount: rows}, nil
}
