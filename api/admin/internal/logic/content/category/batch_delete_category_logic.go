package category

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/contentservice"
)

type BatchDeleteCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除分类
func NewBatchDeleteCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteCategoryLogic {
	return &BatchDeleteCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteCategoryLogic) BatchDeleteCategory(req *types.DeleteCategoryReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.ContentService.BatchDeleteCategories(l.ctx, &contentservice.BatchDeleteCategoriesRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
