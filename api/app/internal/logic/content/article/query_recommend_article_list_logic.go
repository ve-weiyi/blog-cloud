package article

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/articleservice"
)

type QueryRecommendArticleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取推荐文章列表
func NewQueryRecommendArticleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRecommendArticleListLogic {
	return &QueryRecommendArticleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryRecommendArticleListLogic) QueryRecommendArticleList(req *types.QueryRecommendArticleListReq) (resp *types.PageResult, err error) {
	isTop := int64(1)
	isDelete := int64(0)
	status := int64(1)

	in := &articleservice.ListArticlesRequest{
		IsTop:    &isTop,
		IsDelete: &isDelete,
		Status:   &status,
	}

	out, err := l.svcCtx.ArticleService.ListArticles(l.ctx, in)
	if err != nil {
		return nil, err
	}

	list := make([]*types.ArticleHome, 0)
	for _, v := range out.List {
		m := convertArticleHomeTypes(v)
		list = append(list, m)
	}

	resp = &types.PageResult{
		Page:     out.PageResult.Page,
		PageSize: out.PageResult.PageSize,
		Total:    out.PageResult.Total,
		List:     list,
	}
	return
}
