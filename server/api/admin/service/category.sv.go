package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type CategoryService struct {
	svcCtx *svctx.ServiceContext
}

func NewCategoryService(svcCtx *svctx.ServiceContext) *CategoryService {
	return &CategoryService{
		svcCtx: svcCtx,
	}
}

// 分页获取文章分类列表
func (s *CategoryService) FindCategoryList(reqCtx *request.Context, in *dto.CategoryQuery) (out *dto.PageResp, err error) {
	// todo

	return
}

// 创建文章分类
func (s *CategoryService) AddCategory(reqCtx *request.Context, in *dto.CategoryNewReq) (out *dto.CategoryBackDTO, err error) {
	// todo

	return
}

// 批量删除文章分类
func (s *CategoryService) BatchDeleteCategory(reqCtx *request.Context, in *dto.IdsReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 删除文章分类
func (s *CategoryService) DeleteCategory(reqCtx *request.Context, in *dto.IdReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 更新文章分类
func (s *CategoryService) UpdateCategory(reqCtx *request.Context, in *dto.CategoryNewReq) (out *dto.CategoryBackDTO, err error) {
	// todo

	return
}
