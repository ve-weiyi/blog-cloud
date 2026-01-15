package articlerpclogic

import (
	"context"

	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
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
func (l *UpdateArticleLogic) UpdateArticle(in *articlerpc.UpdateArticleReq) (*articlerpc.UpdateArticleResp, error) {
	helper := NewArticleHelperLogic(l.ctx, l.svcCtx)

	entity, err := l.svcCtx.TArticleModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 插入文章分类
	categoryId, err := helper.findOrAddCategory(in.CategoryName)
	if err != nil {
		return nil, err
	}

	entity.ArticleTitle = in.ArticleTitle
	entity.ArticleContent = in.ArticleContent
	entity.ArticleCover = in.ArticleCover
	entity.ArticleType = in.ArticleType
	entity.IsTop = in.IsTop
	entity.Status = in.Status
	entity.CategoryId = categoryId

	_, err = l.svcCtx.TArticleModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	// 插入文章标签
	var ats []*model.TArticleTag
	for _, tagName := range in.TagNameList {
		tagId, err := helper.findOrAddTag(tagName)
		if err != nil {
			return nil, err
		}

		at := &model.TArticleTag{
			ArticleId: entity.Id,
			TagId:     tagId,
		}
		ats = append(ats, at)
	}

	err = l.svcCtx.Gorm.Transaction(func(tx *gorm.DB) error {
		// 删除旧的标签
		_, err = l.svcCtx.TArticleTagModel.WithTransaction(tx).Deletes(l.ctx, "article_id = ?", entity.Id)
		if err != nil {
			return err
		}

		// 插入新的标签
		_, err = l.svcCtx.TArticleTagModel.WithTransaction(tx).Inserts(l.ctx, ats...)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &articlerpc.UpdateArticleResp{
		Article: helper.convertArticlePreviewOut(entity),
	}, nil
}
