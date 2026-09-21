package contentservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/contentrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchDeleteCategoriesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteCategoriesLogic {
	return &BatchDeleteCategoriesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除分类
func (l *BatchDeleteCategoriesLogic) BatchDeleteCategories(in *contentrpc.BatchDeleteCategoriesRequest) (*contentrpc.BatchDeleteCategoriesResponse, error) {
	rows, err := l.svcCtx.TCategoryModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &contentrpc.BatchDeleteCategoriesResponse{SuccessCount: rows}, nil
}
