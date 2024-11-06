package article

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/articlerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type TopArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 置顶文章
func NewTopArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TopArticleLogic {
	return &TopArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TopArticleLogic) TopArticle(req *types.ArticleTopReq) (resp *types.EmptyResp, err error) {
	in := &articlerpc.TopArticleReq{
		ArticleId: req.Id,
		IsTop:     req.IsTop,
	}

	_, err = l.svcCtx.ArticleRpc.TopArticle(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
