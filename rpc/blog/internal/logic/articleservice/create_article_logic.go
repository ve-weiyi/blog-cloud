package articleservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

type CreateArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateArticleLogic {
	return &CreateArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建文章
func (l *CreateArticleLogic) CreateArticle(in *articlerpc.CreateArticleRequest) (*articlerpc.CreateArticleResponse, error) {
	helper := NewArticleHelper(l.ctx, l.svcCtx)

	categoryId, err := helper.findOrAddCategory(in.CategoryName)
	if err != nil {
		return nil, err
	}

	entity := &model.TArticle{
		UserId:         in.UserId,
		CategoryId:     categoryId,
		ArticleCover:   in.ArticleCover,
		ArticleTitle:   in.ArticleTitle,
		ArticleContent: in.ArticleContent,
		ArticleType:    in.ArticleType,
		OriginalUrl:    in.OriginalUrl,
		IsTop:          in.IsTop,
		Status:         in.Status,
	}

	_, err = l.svcCtx.TArticleModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	// 插入文章标签
	var ats []*model.TArticleTag
	for _, tagName := range in.TagNames {
		tagId, err := helper.findOrAddTag(tagName)
		if err != nil {
			return nil, err
		}
		ats = append(ats, &model.TArticleTag{ArticleId: entity.Id, TagId: tagId})
	}
	if len(ats) > 0 {
		err = l.svcCtx.GormDB.Transaction(func(tx *gorm.DB) error {
			_, err = l.svcCtx.TArticleTagModel.WithTx(tx).DeleteBatch(l.ctx, "article_id = ?", entity.Id)
			if err != nil {
				return err
			}
			_, err = l.svcCtx.TArticleTagModel.WithTx(tx).InsertBatch(l.ctx, ats...)
			return err
		})
		if err != nil {
			return nil, err
		}
	}

	return &articlerpc.CreateArticleResponse{
		Id: entity.Id,
	}, nil
}
