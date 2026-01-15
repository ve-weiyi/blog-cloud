package category

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindCategoryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页获取文章分类列表
func NewFindCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindCategoryListLogic {
	return &FindCategoryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindCategoryListLogic) FindCategoryList(req *types.QueryCategoryReq) (resp *types.PageResp, err error) {
	in := &articlerpc.FindCategoryListReq{
		Paginate: &articlerpc.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
	}
	out, err := l.svcCtx.ArticleRpc.FindCategoryList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	list := make([]*types.Category, 0)
	for _, v := range out.List {
		list = append(list, &types.Category{
			Id:           v.Id,
			CategoryName: v.CategoryName,
			ArticleCount: v.ArticleCount,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
		})
	}

	_, err = l.svcCtx.SyslogRpc.AddVisitLog(l.ctx, &syslogrpc.AddVisitLogReq{
		PageName: "分类",
	})
	if err != nil {
		return nil, err
	}

	resp = &types.PageResp{}
	resp.Page = out.Pagination.Page
	resp.PageSize = out.Pagination.PageSize
	resp.Total = out.Pagination.Total
	resp.List = list
	return resp, nil
}
