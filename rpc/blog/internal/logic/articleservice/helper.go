package articleservicelogic

import (
	"context"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/adapter/gormx/queryx"

	"github.com/ve-weiyi/blog-cloud/infra/constants/cachekey"
)

type ArticleHelper struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleHelper(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleHelper {
	return &ArticleHelper{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (h *ArticleHelper) convertArticleQuery(in *articlerpc.ListArticlesRequest) (page int, size int, sorts string, conditions string, params []any) {
	var opts []queryx.Option
	if in.PageQuery != nil {
		opts = append(opts, queryx.WithPage(int(in.PageQuery.Page)))
		opts = append(opts, queryx.WithSize(int(in.PageQuery.PageSize)))
		opts = append(opts, queryx.WithSorts(in.PageQuery.Sorts...))
	}
	if len(in.Ids) > 0 {
		opts = append(opts, queryx.WithCondition("id in (?)", in.Ids))
	}
	if in.ArticleTitle != nil {
		opts = append(opts, queryx.WithCondition("article_title like ?", "%"+*in.ArticleTitle+"%"))
	}
	if in.ArticleType != nil {
		opts = append(opts, queryx.WithCondition("article_type = ?", *in.ArticleType))
	}
	if in.IsTop != nil {
		opts = append(opts, queryx.WithCondition("is_top = ?", *in.IsTop))
	}
	if in.IsDelete != nil {
		opts = append(opts, queryx.WithCondition("is_delete = ?", *in.IsDelete))
	}
	if in.Status != nil {
		opts = append(opts, queryx.WithCondition("status = ?", *in.Status))
	}
	if in.CategoryName != nil {
		category, err := h.svcCtx.TCategoryModel.FindOneByCategoryName(h.ctx, *in.CategoryName)
		if err == nil {
			opts = append(opts, queryx.WithCondition("category_id = ?", category.Id))
		}
	}
	if in.TagName != nil {
		tag, err := h.svcCtx.TTagModel.FindOneByTagName(h.ctx, *in.TagName)
		if err == nil {
			ats, err := h.svcCtx.TArticleTagModel.FindALL(h.ctx, "tag_id = ?", tag.Id)
			if err == nil {
				var articleIds []int64
				for _, v := range ats {
					articleIds = append(articleIds, v.ArticleId)
				}
				opts = append(opts, queryx.WithCondition("id in (?)", articleIds))
			}
		}
	}
	return queryx.NewQueryBuilder(opts...).Build()
}

func (h *ArticleHelper) convertArticleOut(records []*model.TArticle) ([]*articlerpc.Article, error) {
	acm, err := h.findCategoryGroupArticle(records)
	if err != nil {
		return nil, err
	}
	atm, err := h.findTagGroupArticle(records)
	if err != nil {
		return nil, err
	}

	var list []*articlerpc.Article
	for _, entity := range records {
		m := &articlerpc.Article{
			Id:             entity.Id,
			UserId:         entity.UserId,
			CategoryId:     entity.CategoryId,
			ArticleCover:   entity.ArticleCover,
			ArticleTitle:   entity.ArticleTitle,
			ArticleContent: entity.ArticleContent,
			ArticleType:    entity.ArticleType,
			OriginalUrl:    entity.OriginalUrl,
			IsTop:          entity.IsTop,
			IsDelete:       entity.IsDelete,
			Status:         entity.Status,
			CreatedAt:      entity.CreatedAt.UnixMilli(),
			UpdatedAt:      entity.UpdatedAt.UnixMilli(),
			LikeCount:      entity.LikeCount,
			ViewCount:      h.getArticleViewCount(entity.Id),
		}
		if v, ok := acm[entity.Id]; ok {
			m.Category = &articlerpc.ArticleCategory{
				Id:           v.Id,
				CategoryName: v.CategoryName,
			}
		}
		if v, ok := atm[entity.Id]; ok {
			for _, tag := range v {
				m.Tags = append(m.Tags, &articlerpc.ArticleTag{
					Id:      tag.Id,
					TagName: tag.TagName,
				})
			}
		}
		list = append(list, m)
	}
	return list, nil
}

func (h *ArticleHelper) convertArticlePreviewOut(record *model.TArticle) *articlerpc.ArticlePreview {
	return &articlerpc.ArticlePreview{
		Id:           record.Id,
		ArticleCover: record.ArticleCover,
		ArticleTitle: record.ArticleTitle,
		CreatedAt:    record.CreatedAt.UnixMilli(),
		LikeCount:    record.LikeCount,
		ViewCount:    h.getArticleViewCount(record.Id),
	}
}

func (h *ArticleHelper) findCategoryGroupArticle(list []*model.TArticle) (map[int64]*model.TCategory, error) {
	var categoryIds []int64
	for _, v := range list {
		categoryIds = append(categoryIds, v.CategoryId)
	}
	cs, err := h.svcCtx.TCategoryModel.FindALL(h.ctx, "id IN ?", categoryIds)
	if err != nil {
		return nil, err
	}
	acm := make(map[int64]*model.TCategory)
	for _, v := range list {
		for _, category := range cs {
			if category.Id == v.CategoryId {
				acm[v.Id] = category
			}
		}
	}
	return acm, nil
}

func (h *ArticleHelper) findTagGroupArticle(list []*model.TArticle) (map[int64][]*model.TTag, error) {
	var articleIds []int64
	for _, v := range list {
		articleIds = append(articleIds, v.Id)
	}
	ats, err := h.svcCtx.TArticleTagModel.FindALL(h.ctx, "article_id in (?)", articleIds)
	if err != nil {
		return nil, err
	}
	var tagIds []int64
	for _, v := range ats {
		tagIds = append(tagIds, v.TagId)
	}
	ts, err := h.svcCtx.TTagModel.FindALL(h.ctx, "id in (?)", tagIds)
	if err != nil {
		return nil, err
	}
	atm := make(map[int64][]*model.TTag)
	for _, v := range ats {
		for _, tag := range ts {
			if tag.Id == v.TagId {
				atm[v.ArticleId] = append(atm[v.ArticleId], tag)
			}
		}
	}
	return atm, nil
}

func (h *ArticleHelper) findOrAddCategory(name string) (int64, error) {
	if name == "" {
		return 0, nil
	}
	category, err := h.svcCtx.TCategoryModel.FindOneByCategoryName(h.ctx, name)
	if err != nil {
		insert := &model.TCategory{CategoryName: name}
		_, err := h.svcCtx.TCategoryModel.Insert(h.ctx, insert)
		if err != nil {
			return 0, err
		}
		return insert.Id, nil
	}
	return category.Id, nil
}

func (h *ArticleHelper) findOrAddTag(name string) (int64, error) {
	if name == "" {
		return 0, nil
	}
	tag, err := h.svcCtx.TTagModel.FindOneByTagName(h.ctx, name)
	if err != nil {
		insert := &model.TTag{TagName: name}
		_, err := h.svcCtx.TTagModel.Insert(h.ctx, insert)
		if err != nil {
			return 0, err
		}
		return insert.Id, nil
	}
	return tag.Id, nil
}

func (h *ArticleHelper) getArticleViewCount(articleId int64) int64 {
	id := cast.ToString(articleId)
	key := cachekey.ArticleViewCountKey
	result, err := h.svcCtx.Redis.ZScore(h.ctx, key, id).Result()
	if err != nil {
		article, err := h.svcCtx.TArticleModel.FindById(h.ctx, articleId)
		if err != nil {
			return 0
		}
		if err := h.svcCtx.Redis.ZIncrBy(h.ctx, key, float64(article.ViewCount), id).Err(); err != nil {
			h.Errorf("getArticleViewCount ZIncrBy error: %v", err)
		}
		return article.ViewCount
	}
	return int64(result)
}

func (h *ArticleHelper) findArticleCountGroupCategory(list []*model.TCategory) (map[int64]int64, error) {
	var ids []int64
	for _, v := range list {
		ids = append(ids, v.Id)
	}
	var results []struct {
		CategoryID   int64 `gorm:"column:category_id"`
		ArticleCount int64 `gorm:"column:article_count"`
	}
	err := h.svcCtx.GormDB.Model(&model.TArticle{}).
		Select("category_id, COUNT(*) as article_count").
		Where("category_id IN ?", ids).
		Group("category_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	acm := make(map[int64]int64)
	for _, r := range results {
		acm[r.CategoryID] = r.ArticleCount
	}
	return acm, nil
}

func (h *ArticleHelper) findArticleCountGroupTag(list []*model.TTag) (map[int64]int64, error) {
	var ids []int64
	for _, v := range list {
		ids = append(ids, v.Id)
	}
	var results []struct {
		TagID        int64 `gorm:"column:tag_id"`
		ArticleCount int64 `gorm:"column:article_count"`
	}
	err := h.svcCtx.GormDB.Model(&model.TArticleTag{}).
		Select("tag_id, COUNT(*) as article_count").
		Where("tag_id IN ?", ids).
		Group("tag_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	acm := make(map[int64]int64)
	for _, r := range results {
		acm[r.TagID] = r.ArticleCount
	}
	return acm, nil
}
