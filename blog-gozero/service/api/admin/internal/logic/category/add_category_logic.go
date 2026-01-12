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

func (l *AddCategoryLogic) AddCategory(req *types.NewCategoryReq) (resp *types.CategoryBackVO, err error) {
	in := &articlerpc.AddCategoryReq{
		Id:           req.Id,
		CategoryName: req.CategoryName,
	}
	out, err := l.svcCtx.ArticleRpc.AddCategory(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return convertCategoryTypes(out.Category), nil
}

func convertCategoryTypes(out *articlerpc.Category) *types.CategoryBackVO {
	return &types.CategoryBackVO{
		Id:           out.Id,
		CategoryName: out.CategoryName,
		ArticleCount: 0,
		CreatedAt:    out.CreatedAt,
		UpdatedAt:    out.UpdatedAt,
	}
}

func convertCategoryDetailsTypes(out *articlerpc.CategoryDetails) *types.CategoryBackVO {
	return &types.CategoryBackVO{
		Id:           out.Id,
		CategoryName: out.CategoryName,
		ArticleCount: out.ArticleCount,
		CreatedAt:    out.CreatedAt,
		UpdatedAt:    out.UpdatedAt,
	}
}
