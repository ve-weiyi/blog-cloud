package tag

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/articleservice"
)

type QueryTagListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取标签列表
func NewQueryTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryTagListLogic {
	return &QueryTagListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryTagListLogic) QueryTagList(req *types.QueryTagListReq) (resp *types.PageResult, err error) {
	in := &articleservice.ListTagsRequest{
		PageQuery: &articleservice.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		TagName: req.TagName,
	}

	out, err := l.svcCtx.ArticleService.ListTags(l.ctx, in)
	if err != nil {
		return nil, err
	}

	list := make([]*types.Tag, 0)
	for _, v := range out.List {
		list = append(list, &types.Tag{
			Id:           v.Id,
			TagName:      v.TagName,
			ArticleCount: v.ArticleCount,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
		})
	}

	resp = &types.PageResult{
		Page:     out.PageResult.Page,
		PageSize: out.PageResult.PageSize,
		Total:    out.PageResult.Total,
		List:     list,
	}
	return
}
