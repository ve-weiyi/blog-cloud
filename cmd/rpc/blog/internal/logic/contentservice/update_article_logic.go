package contentservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/contentrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/model"
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
	if !enums.IsValidArticleStatus(in.Status) {
		return nil, bizerr.NewBizError(bizcode.CodeParamValueInvalid, "文章状态不合法")
	}

	helper := NewArticleHelper(l.ctx, l.svcCtx)

	var entity *model.TArticle
	err := l.svcCtx.GormDB.Transaction(func(tx *gorm.DB) error {
		var err error
		entity, err = l.svcCtx.TArticleModel.WithTx(tx).FindById(l.ctx, in.Id)
		if err != nil {
			return err
		}

		categoryId, err := helper.findOrAddCategory(tx, in.CategoryName)
		if err != nil {
			return err
		}

		entity.ArticleTitle = in.ArticleTitle
		entity.ArticleContent = in.ArticleContent
		entity.ArticleCover = in.ArticleCover
		entity.ArticleType = in.ArticleType
		entity.OriginalUrl = in.OriginalUrl
		entity.IsTop = in.IsTop
		entity.Status = in.Status
		entity.CategoryId = categoryId

		if _, err := l.svcCtx.TArticleModel.WithTx(tx).Update(l.ctx, entity); err != nil {
			return err
		}

		tagIds, err := helper.findOrAddTags(tx, in.TagNames)
		if err != nil {
			return err
		}

		if _, err := l.svcCtx.TArticleTagModel.WithTx(tx).DeleteBatch(l.ctx, "article_id = ?", entity.Id); err != nil {
			return err
		}
		if len(tagIds) == 0 {
			return nil
		}

		ats := make([]*model.TArticleTag, 0, len(tagIds))
		for _, tagId := range tagIds {
			ats = append(ats, &model.TArticleTag{ArticleId: entity.Id, TagId: tagId})
		}
		_, err = l.svcCtx.TArticleTagModel.WithTx(tx).InsertBatch(l.ctx, ats...)
		return err
	})
	if err != nil {
		return nil, err
	}

	list, err := helper.convertArticleOut([]*model.TArticle{entity})
	if err != nil {
		return nil, err
	}

	return &contentrpc.UpdateArticleResponse{
		Success: true,
		Article: list[0],
	}, nil
}
