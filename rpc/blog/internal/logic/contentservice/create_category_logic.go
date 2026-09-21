package contentservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/contentrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

type CreateCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCategoryLogic {
	return &CreateCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建分类
func (l *CreateCategoryLogic) CreateCategory(in *contentrpc.CreateCategoryRequest) (*contentrpc.CreateCategoryResponse, error) {
	entity := &model.TCategory{CategoryName: in.CategoryName}
	_, err := l.svcCtx.TCategoryModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &contentrpc.CreateCategoryResponse{Id: entity.Id}, nil
}
