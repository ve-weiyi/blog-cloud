package category

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/articlerpc"

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

func (l *FindCategoryListLogic) FindCategoryList(req *types.CategoryQueryReq) (resp *types.PageResp, err error) {
	in := &articlerpc.FindCategoryListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	out, err := l.svcCtx.ArticleRpc.FindCategoryList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	list := make([]*types.Category, 0)
	for _, v := range out.List {
		m := ConvertCategoryTypes(v)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	resp.Total = out.Total
	resp.List = list
	return resp, nil
}

func ConvertCategoryTypes(in *articlerpc.CategoryDetails) (out *types.Category) {
	return &types.Category{
		Id:           in.Id,
		CategoryName: in.CategoryName,
		ArticleCount: in.ArticleCount,
		CreatedAt:    in.CreatedAt,
		UpdatedAt:    in.UpdatedAt,
	}
}
