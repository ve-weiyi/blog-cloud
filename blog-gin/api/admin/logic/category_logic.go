package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type CategoryLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewCategoryLogic(svcCtx *svctx.ServiceContext) *CategoryLogic {
	return &CategoryLogic{
		svcCtx: svcCtx,
	}
}

// 创建文章分类
func (s *CategoryLogic) AddCategory(reqCtx *request.Context, in *types.NewCategoryReq) (out *types.CategoryBackVO, err error) {
	// todo

	return
}

// 删除文章分类
func (s *CategoryLogic) DeletesCategory(reqCtx *request.Context, in *types.IdsReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 分页获取文章分类列表
func (s *CategoryLogic) FindCategoryList(reqCtx *request.Context, in *types.QueryCategoryReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 更新文章分类
func (s *CategoryLogic) UpdateCategory(reqCtx *request.Context, in *types.NewCategoryReq) (out *types.CategoryBackVO, err error) {
	// todo

	return
}
