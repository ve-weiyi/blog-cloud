package articleservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/articlerpc"
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

func (l *ListTagsLogic) ListTags(in *articlerpc.ListTagsRequest) (*articlerpc.ListTagsResponse, error) {
	helper := NewArticleHelper(l.ctx, l.svcCtx)

	var opts []queryx.Option
	if in.PageQuery != nil {
		opts = append(opts, queryx.WithPage(int(in.PageQuery.Page)))
		opts = append(opts, queryx.WithSize(int(in.PageQuery.PageSize)))
		opts = append(opts, queryx.WithSorts(in.PageQuery.Sorts...))
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

	var list []*articlerpc.Tag
	for _, entity := range records {
		m := &articlerpc.Tag{
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

	return &articlerpc.ListTagsResponse{
		PageResult: &articlerpc.PageResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: list,
	}, nil
}
