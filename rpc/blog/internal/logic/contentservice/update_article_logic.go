package contentservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/contentrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

type UpdateArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateArticleLogic {
	return &UpdateArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新文章
func (l *UpdateArticleLogic) UpdateArticle(in *contentrpc.UpdateArticleRequest) (*contentrpc.UpdateArticleResponse, error) {
	helper := NewArticleHelper(l.ctx, l.svcCtx)

	entity, err := l.svcCtx.TArticleModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	categoryId, err := helper.findOrAddCategory(in.CategoryName)
	if err != nil {
		return nil, err
	}

	entity.ArticleTitle = in.ArticleTitle
	entity.ArticleContent = in.ArticleContent
	entity.ArticleCover = in.ArticleCover
	entity.ArticleType = in.ArticleType
	entity.OriginalUrl = in.OriginalUrl
	entity.IsTop = in.IsTop
	entity.Status = in.Status
	entity.CategoryId = categoryId

	_, err = l.svcCtx.TArticleModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	// 更新文章标签
	var ats []*model.TArticleTag
	for _, tagName := range in.TagNames {
		tagId, err := helper.findOrAddTag(tagName)
		if err != nil {
			return nil, err
		}
		ats = append(ats, &model.TArticleTag{ArticleId: entity.Id, TagId: tagId})
	}

	err = l.svcCtx.GormDB.Transaction(func(tx *gorm.DB) error {
		_, err = l.svcCtx.TArticleTagModel.WithTx(tx).DeleteBatch(l.ctx, "article_id = ?", entity.Id)
		if err != nil {
			return err
		}
		if len(ats) > 0 {
			_, err = l.svcCtx.TArticleTagModel.WithTx(tx).InsertBatch(l.ctx, ats...)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &contentrpc.UpdateArticleResponse{
		Success: true,
	}, nil
}
