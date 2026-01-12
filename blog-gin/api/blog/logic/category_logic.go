package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
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

// 分页获取文章分类列表
func (s *CategoryLogic) FindCategoryList(reqCtx *request.Context, in *types.QueryCategoryReq) (out *types.PageResp, err error) {
	// todo

	return
}
