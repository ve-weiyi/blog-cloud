package page

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/siteservice"
)

type QueryPageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取页面列表
func NewQueryPageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryPageListLogic {
	return &QueryPageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryPageListLogic) QueryPageList(req *types.QueryPageListReq) (resp *types.ListResult, err error) {
	in := &siteservice.ListPagesRequest{
		ListQuery: &siteservice.ListQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
	}

	out, err := l.svcCtx.SiteService.ListPages(l.ctx, in)
	if err != nil {
		return nil, err
	}

	list := make([]*types.Page, 0)
	for _, v := range out.List {
		list = append(list, &types.Page{
			Id:         v.Id,
			PageName:   v.PageName,
			PageLabel:  v.PageLabel,
			PageCover:  v.PageCover,
			IsCarousel: v.IsCarousel,
			CreatedAt:  v.CreatedAt,
			UpdatedAt:  v.UpdatedAt,
		})
	}

	resp = &types.ListResult{
		Page:     out.ListResult.Page,
		PageSize: out.ListResult.PageSize,
		Total:    out.ListResult.Total,
		List:     list,
	}
	return
}
