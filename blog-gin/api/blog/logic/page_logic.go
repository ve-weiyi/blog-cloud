package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
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

// 分页获取页面列表
func (s *PageLogic) FindPageList(reqCtx *request.Context, in *types.QueryPageReq) (out *types.PageResp, err error) {
	// todo

	return
}
