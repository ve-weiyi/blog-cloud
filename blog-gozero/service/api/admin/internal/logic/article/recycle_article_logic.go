package article

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/articlerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecycleArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 回收文章
func NewRecycleArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecycleArticleLogic {
	return &RecycleArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecycleArticleLogic) RecycleArticle(req *types.ArticleRecycleReq) (resp *types.EmptyResp, err error) {
	in := &articlerpc.UpdateArticleDeleteReq{
		ArticleId: req.Id,
		IsDelete:  req.IsDelete,
	}

	_, err = l.svcCtx.ArticleRpc.UpdateArticleDelete(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
