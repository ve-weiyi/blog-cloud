package service

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/dto"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
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
func (s *CategoryService) FindCategoryList(reqCtx *request.Context, in *dto.CategoryQueryReq) (out *dto.PageResp, err error) {
	// todo

	return
}
