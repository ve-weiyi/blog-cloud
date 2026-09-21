package article

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/contentservice"
)

type UpdateArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新文章
func NewUpdateArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateArticleLogic {
	return &UpdateArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateArticleLogic) UpdateArticle(req *types.UpdateArticleReq) (resp *types.ArticleVO, err error) {
	_, err = l.svcCtx.ContentService.UpdateArticle(l.ctx, &contentservice.UpdateArticleRequest{
		Id:             req.Id,
		ArticleCover:   req.ArticleCover,
		ArticleTitle:   req.ArticleTitle,
		ArticleContent: req.ArticleContent,
		ArticleType:    req.ArticleType,
		OriginalUrl:    req.OriginalUrl,
		IsTop:          req.IsTop,
		Status:         req.Status,
		CategoryName:   req.CategoryName,
		TagNames:       req.TagNameList,
	})
	if err != nil {
		return nil, err
	}

	return &types.ArticleVO{
		Id: req.Id,
	}, nil
}
