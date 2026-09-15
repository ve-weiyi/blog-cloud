package contentservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/contentrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/adapter/gormx/queryx"
)

type ListTagsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListTagsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTagsLogic {
	return &ListTagsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListTagsLogic) ListTags(in *contentrpc.ListTagsRequest) (*contentrpc.ListTagsResponse, error) {
	helper := NewArticleHelper(l.ctx, l.svcCtx)

	var opts []queryx.Option
	if in.ListQuery != nil {
		opts = append(opts, queryx.WithPage(int(in.ListQuery.Page)))
		opts = append(opts, queryx.WithSize(int(in.ListQuery.PageSize)))
		opts = append(opts, queryx.WithSorts(in.ListQuery.Sorts...))
	}
	if in.TagName != nil {
		opts = append(opts, queryx.WithCondition("tag_name like ?", "%"+*in.TagName+"%"))
	}

	page, size, sorts, conditions, params := queryx.NewQueryBuilder(opts...).Build()
	records, total, err := l.svcCtx.TTagModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	acm, _ := helper.findArticleCountGroupTag(records)

	var list []*contentrpc.Tag
	for _, entity := range records {
		m := &contentrpc.Tag{
			Id:        entity.Id,
			TagName:   entity.TagName,
			CreatedAt: entity.CreatedAt.UnixMilli(),
			UpdatedAt: entity.UpdatedAt.UnixMilli(),
		}
		if acm != nil {
			m.ArticleCount = acm[entity.Id]
		}
		list = append(list, m)
	}

	return &contentrpc.ListTagsResponse{
		ListResult: &contentrpc.ListResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: list,
	}, nil
}
