package articleservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type ListArticlePreviewsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListArticlePreviewsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListArticlePreviewsLogic {
	return &ListArticlePreviewsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页查询文章预览列表
func (l *ListArticlePreviewsLogic) ListArticlePreviews(in *articlerpc.ListArticlesRequest) (*articlerpc.ListArticlePreviewsResponse, error) {
	helper := NewArticleHelper(l.ctx, l.svcCtx)
	page, size, sorts, conditions, params := helper.convertArticleQuery(in)

	records, total, err := l.svcCtx.TArticleModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*articlerpc.ArticlePreview
	for _, v := range records {
		list = append(list, helper.convertArticlePreviewOut(v))
	}

	return &articlerpc.ListArticlePreviewsResponse{
		PageResult: &articlerpc.PageResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: list,
	}, nil
}
