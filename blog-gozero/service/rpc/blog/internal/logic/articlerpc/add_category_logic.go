package articlerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCategoryLogic {
	return &AddCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建文章分类
func (l *AddCategoryLogic) AddCategory(in *articlerpc.CategoryNewReq) (*articlerpc.CategoryDetails, error) {
	entity := convertCategoryIn(in)
	_, err := l.svcCtx.TCategoryModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &articlerpc.CategoryDetails{}, nil
}
