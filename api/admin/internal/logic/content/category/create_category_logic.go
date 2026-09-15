package category

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/contentservice"
)

type CreateCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建分类
func NewCreateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCategoryLogic {
	return &CreateCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCategoryLogic) CreateCategory(req *types.CreateCategoryReq) (resp *types.CategoryVO, err error) {
	out, err := l.svcCtx.ContentService.CreateCategory(l.ctx, &contentservice.CreateCategoryRequest{
		CategoryName: req.CategoryName,
	})
	if err != nil {
		return nil, err
	}

	return &types.CategoryVO{
		Id:           out.Id,
		CategoryName: req.CategoryName,
		ArticleCount: 0,
	}, nil
}
