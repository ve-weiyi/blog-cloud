package articlerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/constant"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddArticleLogic {
	return &AddArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建文章
func (l *AddArticleLogic) AddArticle(in *articlerpc.NewArticleReq) (*articlerpc.ArticlePreview, error) {
	helper := NewArticleHelperLogic(l.ctx, l.svcCtx)

	entity := &model.TArticle{
		Id:             in.Id,
		UserId:         in.UserId,
		CategoryId:     0,
		ArticleCover:   in.ArticleCover,
		ArticleTitle:   in.ArticleTitle,
		ArticleContent: in.ArticleContent,
		ArticleType:    in.ArticleType,
		OriginalUrl:    in.OriginalUrl,
		IsTop:          constant.ArticleIsTopNo,
		IsDelete:       constant.ArticleIsDeleteNo,
		Status:         in.Status,
		LikeCount:      0,
	}

	// 插入文章分类
	categoryId, err := helper.findOrAddCategory(in.CategoryName)
	if err != nil {
		return nil, err
	}

	entity.CategoryId = categoryId
	_, err = l.svcCtx.TArticleModel.Insert(l.ctx, entity)
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

	return helper.convertArticlePreviewOut(entity), nil
}
