package tag

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/contentservice"
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

func (l *QueryTagListLogic) QueryTagList(req *types.QueryTagListReq) (resp *types.ListResult, err error) {
	out, err := l.svcCtx.ContentService.ListTags(l.ctx, &contentservice.ListTagsRequest{
		ListQuery: &contentservice.ListQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
		TagName:   req.TagName,
	})
	if err != nil {
		return nil, err
	}

	var list []*types.TagVO
	for _, v := range out.List {
		list = append(list, &types.TagVO{
			Id:           v.Id,
			TagName:      v.TagName,
			ArticleCount: v.ArticleCount,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
		})
	}

	return &types.ListResult{
		Page:     out.ListResult.Page,
		PageSize: out.ListResult.PageSize,
		Total:    out.ListResult.Total,
		List:     list,
	}, nil
}
