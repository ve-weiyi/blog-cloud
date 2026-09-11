package page

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/resourceservice"
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

func (l *QueryPageListLogic) QueryPageList(req *types.QueryPageListReq) (resp *types.PageResult, err error) {
	in := &resourceservice.ListPagesRequest{
		PageQuery: &resourceservice.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
	}

	out, err := l.svcCtx.ResourceService.ListPages(l.ctx, in)
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

	resp = &types.PageResult{
		Page:     out.PageResult.Page,
		PageSize: out.PageResult.PageSize,
		Total:    out.PageResult.Total,
		List:     list,
	}
	return
}
