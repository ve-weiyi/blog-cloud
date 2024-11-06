package category

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/articlerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除文章分类
func NewBatchDeleteCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteCategoryLogic {
	return &BatchDeleteCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteCategoryLogic) BatchDeleteCategory(req *types.IdsReq) (resp *types.BatchResp, err error) {
	in := &articlerpc.IdsReq{
		Ids: req.Ids,
	}

	out, err := l.svcCtx.ArticleRpc.DeleteCategory(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
