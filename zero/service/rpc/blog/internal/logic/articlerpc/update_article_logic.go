package articlerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
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
func (l *UpdateArticleLogic) UpdateArticle(in *articlerpc.ArticleNewReq) (*articlerpc.ArticleDetails, error) {
	helper := NewArticleHelperLogic(l.ctx, l.svcCtx)

	entity, err := l.svcCtx.TArticleModel.FindOne(l.ctx, in.Id)
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
	l.svcCtx.TArticleTagModel.Deletes(l.ctx, "article_id = ?", entity.Id)
	l.svcCtx.TArticleTagModel.Inserts(l.ctx, ats...)

	return &articlerpc.ArticleDetails{}, nil
}
