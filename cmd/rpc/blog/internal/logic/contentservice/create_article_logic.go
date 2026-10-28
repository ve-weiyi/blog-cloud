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
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
	"github.com/ve-weiyi/blog-cloud/model"
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
func (l *CreateArticleLogic) CreateArticle(in *contentrpc.CreateArticleRequest) (*contentrpc.CreateArticleResponse, error) {
	if !enums.IsValidArticleStatus(in.Status) {
		return nil, bizerr.NewBizError(bizcode.CodeParamValueInvalid, "文章状态不合法")
	}

	helper := NewArticleHelper(l.ctx, l.svcCtx)

	userId := in.UserId
	if uid, err := metax.GetUserIdFromCtx(l.ctx); err == nil && uid != "" {
		userId = uid
	}

	var entity *model.TArticle
	err := l.svcCtx.GormDB.Transaction(func(tx *gorm.DB) error {
		categoryId, err := helper.findOrAddCategory(tx, in.CategoryName)
		if err != nil {
			return err
		}

		entity = &model.TArticle{
			UserId:         userId,
			CategoryId:     categoryId,
			ArticleCover:   in.ArticleCover,
			ArticleTitle:   in.ArticleTitle,
			ArticleContent: in.ArticleContent,
			ArticleType:    in.ArticleType,
			OriginalUrl:    in.OriginalUrl,
			IsTop:          in.IsTop,
			Status:         in.Status,
		}
		if _, err := l.svcCtx.TArticleModel.WithTx(tx).Insert(l.ctx, entity); err != nil {
			return err
		}

		tagIds, err := helper.findOrAddTags(tx, in.TagNames)
		if err != nil {
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

	return &contentrpc.CreateArticleResponse{
		Id:      entity.Id,
		Article: list[0],
	}, nil
}
