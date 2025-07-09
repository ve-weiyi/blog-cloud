package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type PageLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewPageLogic(svcCtx *svctx.ServiceContext) *PageLogic {
	return &PageLogic{
		svcCtx: svcCtx,
	}
}

// 创建页面
func (s *PageLogic) AddPage(reqCtx *request.Context, in *types.PageNewReq) (out *types.PageBackVO, err error) {
	// todo

	return
}

// 删除页面
func (s *PageLogic) DeletePage(reqCtx *request.Context, in *types.IdReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 分页获取页面列表
func (s *PageLogic) FindPageList(reqCtx *request.Context, in *types.PageQueryReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 更新页面
func (s *PageLogic) UpdatePage(reqCtx *request.Context, in *types.PageNewReq) (out *types.PageBackVO, err error) {
	// todo

	return
}
