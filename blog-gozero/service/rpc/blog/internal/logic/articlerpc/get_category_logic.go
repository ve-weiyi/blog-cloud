package articlerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoryLogic {
	return &GetCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询文章分类
func (l *GetCategoryLogic) GetCategory(in *articlerpc.GetCategoryReq) (*articlerpc.GetCategoryResp, error) {
	entity, err := l.svcCtx.TCategoryModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &articlerpc.GetCategoryResp{
		Category: &articlerpc.CategoryDetails{
			Id:           entity.Id,
			CategoryName: entity.CategoryName,
			CreatedAt:    entity.CreatedAt.UnixMilli(),
			UpdatedAt:    entity.UpdatedAt.UnixMilli(),
			ArticleCount: 0,
		},
	}, nil
}
