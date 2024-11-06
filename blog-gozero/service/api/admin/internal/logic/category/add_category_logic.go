package category

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/articlerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建文章分类
func NewAddCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCategoryLogic {
	return &AddCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddCategoryLogic) AddCategory(req *types.CategoryNewReq) (resp *types.CategoryBackDTO, err error) {
	in := ConvertCategoryPb(req)
	out, err := l.svcCtx.ArticleRpc.AddCategory(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = ConvertCategoryTypes(out)
	return resp, nil
}

func ConvertCategoryPb(in *types.CategoryNewReq) (out *articlerpc.CategoryNewReq) {
	out = &articlerpc.CategoryNewReq{
		Id:           in.Id,
		CategoryName: in.CategoryName,
	}

	return
}

func ConvertCategoryTypes(in *articlerpc.CategoryDetails) (out *types.CategoryBackDTO) {
	out = &types.CategoryBackDTO{
		Id:           in.Id,
		CategoryName: in.CategoryName,
		ArticleCount: in.ArticleCount,
		CreatedAt:    in.CreatedAt,
		UpdatedAt:    in.UpdatedAt,
	}

	return
}
