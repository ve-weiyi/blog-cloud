package contentservicelogic

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/likes"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/contentrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
	"github.com/ve-weiyi/blog-cloud/infra/store/queryx"
	"github.com/ve-weiyi/blog-cloud/model"

	"github.com/ve-weiyi/blog-cloud/infra/domain/cachekey"
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

// toArticleQuery 把预览查询请求适配成文章查询请求。
// 两者字段完全一致，分开定义是为满足「每个 RPC 持有独立请求类型」的约束；
// 用显式拷贝而非取值方法，是为了保住 optional 字段的 nil 语义
// （is_delete = 0 是有效筛选值，不能用零值判断有无）。
func toArticleQuery(in *contentrpc.ListArticlePreviewsRequest) *contentrpc.ListArticlesRequest {
	return &contentrpc.ListArticlesRequest{
		ListQuery:    in.ListQuery,
		Ids:          in.Ids,
		ArticleTitle: in.ArticleTitle,
		ArticleType:  in.ArticleType,
		CategoryName: in.CategoryName,
		TagName:      in.TagName,
		IsTop:        in.IsTop,
		IsDelete:     in.IsDelete,
		Status:       in.Status,
	}
}

func (h *ArticleHelper) convertArticleQuery(in *contentrpc.ListArticlesRequest) (page int, size int, sorts string, conditions string, params []any, err error) {
	var opts []queryx.Option
	if in.ListQuery != nil {
		opts = append(opts, queryx.WithPage(int(in.ListQuery.Page)))
		opts = append(opts, queryx.WithSize(int(in.ListQuery.PageSize)))
		opts = append(opts, queryx.WithSorts(in.ListQuery.Sorts...))
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
		category, findErr := h.svcCtx.TCategoryModel.FindOneByCategoryName(h.ctx, *in.CategoryName)
		switch {
		case findErr == nil:
			opts = append(opts, queryx.WithCondition("category_id = ?", category.Id))
		case errors.Is(findErr, gorm.ErrRecordNotFound):
			// 条件不存在时返回空结果，而不是悄悄丢掉筛选条件
			opts = append(opts, queryx.WithCondition("1 = 0"))
		default:
			return 0, 0, "", "", nil, findErr
		}
	}
	if in.TagName != nil {
		tag, findErr := h.svcCtx.TTagModel.FindOneByTagName(h.ctx, *in.TagName)
		switch {
		case findErr == nil:
			ats, err := h.svcCtx.TArticleTagModel.FindAll(h.ctx, "tag_id = ?", tag.Id)
			if err != nil {
				return 0, 0, "", "", nil, err
			}
			if len(ats) == 0 {
				opts = append(opts, queryx.WithCondition("1 = 0"))
				break
			}
			articleIds := make([]int64, 0, len(ats))
			for _, v := range ats {
				articleIds = append(articleIds, v.ArticleId)
			}
			opts = append(opts, queryx.WithCondition("id in (?)", articleIds))
		case errors.Is(findErr, gorm.ErrRecordNotFound):
			opts = append(opts, queryx.WithCondition("1 = 0"))
		default:
			return 0, 0, "", "", nil, findErr
		}
	}
	page, size, sorts, conditions, params = queryx.NewQueryBuilder(opts...).Build()
	return page, size, sorts, conditions, params, nil
}

func (h *ArticleHelper) convertArticleOut(records []*model.TArticle) ([]*contentrpc.Article, error) {
	acm, err := h.findCategoryGroupArticle(records)
	if err != nil {
		return nil, err
	}
	atm, err := h.findTagGroupArticle(records)
	if err != nil {
		return nil, err
	}
	vcm, err := h.getArticleViewCounts(records)
	if err != nil {
		return nil, err
	}
	lkm, lcm, err := h.getLikeStats(articleIDs(records))
	if err != nil {
		return nil, err
	}

	var list []*contentrpc.Article
	for _, entity := range records {
		m := &contentrpc.Article{
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
			LikeCount:      lcm[entity.Id],
			ViewCount:      vcm[entity.Id],
			Liked:          lkm[entity.Id],
		}
		if v, ok := acm[entity.Id]; ok {
			m.Category = &contentrpc.ArticleCategory{
				Id:           v.Id,
				CategoryName: v.CategoryName,
			}
		}
		if v, ok := atm[entity.Id]; ok {
			for _, tag := range v {
				m.Tags = append(m.Tags, &contentrpc.ArticleTag{
					Id:      tag.Id,
					TagName: tag.TagName,
				})
			}
		}
		list = append(list, m)
	}
	return list, nil
}

// convertArticlePreviewOuts 批量转换预览：浏览量一次 pipeline 取回，避免逐条 Redis 往返。
func (h *ArticleHelper) convertArticlePreviewOuts(records []*model.TArticle) ([]*contentrpc.ArticlePreview, error) {
	vcm, err := h.getArticleViewCounts(records)
	if err != nil {
		return nil, err
	}
	lkm, lcm, err := h.getLikeStats(articleIDs(records))
	if err != nil {
		return nil, err
	}

	list := make([]*contentrpc.ArticlePreview, 0, len(records))
	for _, record := range records {
		list = append(list, &contentrpc.ArticlePreview{
			Id:           record.Id,
			ArticleCover: record.ArticleCover,
			ArticleTitle: record.ArticleTitle,
			CreatedAt:    record.CreatedAt.UnixMilli(),
			LikeCount:    lcm[record.Id],
			ViewCount:    vcm[record.Id],
			Liked:        lkm[record.Id],
		})
	}
	return list, nil
}

// articleIDs 提取记录 ID，供批量查询使用。
func articleIDs(records []*model.TArticle) []int64 {
	ids := make([]int64, 0, len(records))
	for _, record := range records {
		ids = append(ids, record.Id)
	}
	return ids
}

// getLikeStats 批量取当前用户的点赞状态与各文章的点赞数。
func (h *ArticleHelper) getLikeStats(ids []int64) (liked map[int64]bool, counts map[int64]int64, err error) {
	uid, _ := metax.GetUserIdFromCtx(h.ctx)
	deps := h.svcCtx.LikeDeps()

	liked, err = likes.Liked(h.ctx, deps, enums.LikeTargetArticle, uid, ids)
	if err != nil {
		return nil, nil, err
	}
	counts, err = likes.Counts(h.ctx, deps, enums.LikeTargetArticle, ids)
	if err != nil {
		return nil, nil, err
	}
	return liked, counts, nil
}

func (h *ArticleHelper) findCategoryGroupArticle(list []*model.TArticle) (map[int64]*model.TCategory, error) {
	var categoryIds []int64
	for _, v := range list {
		categoryIds = append(categoryIds, v.CategoryId)
	}
	cs, err := h.svcCtx.TCategoryModel.FindAll(h.ctx, "id IN ?", categoryIds)
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
	ats, err := h.svcCtx.TArticleTagModel.FindAll(h.ctx, "article_id in (?)", articleIds)
	if err != nil {
		return nil, err
	}
	var tagIds []int64
	for _, v := range ats {
		tagIds = append(tagIds, v.TagId)
	}
	ts, err := h.svcCtx.TTagModel.FindAll(h.ctx, "id in (?)", tagIds)
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

func (h *ArticleHelper) findOrAddCategory(tx *gorm.DB, name string) (int64, error) {
	if name == "" {
		return 0, nil
	}

	categories := h.svcCtx.TCategoryModel
	if tx != nil {
		categories = categories.WithTx(tx)
	}

	category, err := categories.FindOneByCategoryName(h.ctx, name)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
		insert := &model.TCategory{CategoryName: name}
		if _, err := categories.Insert(h.ctx, insert); err != nil {
			return 0, err
		}
		return insert.Id, nil
	}
	return category.Id, nil
}

func (h *ArticleHelper) findOrAddTag(tx *gorm.DB, name string) (int64, error) {
	if name == "" {
		return 0, nil
	}

	tags := h.svcCtx.TTagModel
	if tx != nil {
		tags = tags.WithTx(tx)
	}

	tag, err := tags.FindOneByTagName(h.ctx, name)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
		insert := &model.TTag{TagName: name}
		if _, err := tags.Insert(h.ctx, insert); err != nil {
			return 0, err
		}
		return insert.Id, nil
	}
	return tag.Id, nil
}

// findOrAddTags 去重并过滤空标签名，返回标签 ID 列表。
func (h *ArticleHelper) findOrAddTags(tx *gorm.DB, names []string) ([]int64, error) {
	seen := make(map[string]struct{}, len(names))
	ids := make([]int64, 0, len(names))
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}

		id, err := h.findOrAddTag(tx, name)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// getArticleViewCounts 批量读取浏览量：Redis pipeline 一次取回，缺失的成员用 DB 值回种。
func (h *ArticleHelper) getArticleViewCounts(records []*model.TArticle) (map[int64]int64, error) {
	counts := make(map[int64]int64, len(records))
	if len(records) == 0 {
		return counts, nil
	}

	key := cachekey.ArticleViewCountKey
	ids := make([]string, 0, len(records))
	for _, record := range records {
		ids = append(ids, strconv.FormatInt(record.Id, 10))
	}

	pipe := h.svcCtx.Redis.Pipeline()
	cmds := make([]*redis.FloatCmd, len(ids))
	for i, id := range ids {
		cmds[i] = pipe.ZScore(h.ctx, key, id)
	}
	if _, err := pipe.Exec(h.ctx); err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	var missing []int64
	for i, record := range records {
		score, err := cmds[i].Result()
		if err != nil {
			missing = append(missing, record.Id)
			continue
		}
		counts[record.Id] = int64(score)
	}
	if len(missing) == 0 {
		return counts, nil
	}

	articles, err := h.svcCtx.TArticleModel.FindAll(h.ctx, "id IN ?", missing)
	if err != nil {
		return nil, err
	}
	dbCounts := make(map[int64]int64, len(articles))
	for _, article := range articles {
		dbCounts[article.Id] = article.ViewCount
	}

	pipe = h.svcCtx.Redis.Pipeline()
	for _, id := range missing {
		count := dbCounts[id]
		counts[id] = count
		pipe.ZAddNX(h.ctx, key, redis.Z{
			Score:  float64(count),
			Member: strconv.FormatInt(id, 10),
		})
	}
	if _, err := pipe.Exec(h.ctx); err != nil {
		h.Errorf("getArticleViewCounts ZAddNX error: %v", err)
	}

	return counts, nil
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
	err := h.svcCtx.GormDB.Table("`t_article`").
		Select("category_id, COUNT(*) as article_count").
		Where("category_id IN ? AND is_delete = ? AND status = ?", ids, 0, enums.ArticleStatusPublic).
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
	err := h.svcCtx.GormDB.Table("`t_article_tag`").
		Joins("JOIN `t_article` ON `t_article`.`id` = `t_article_tag`.`article_id`").
		Select("tag_id, COUNT(*) as article_count").
		Where("tag_id IN ? AND `t_article`.`is_delete` = ? AND `t_article`.`status` = ?", ids, 0, enums.ArticleStatusPublic).
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
