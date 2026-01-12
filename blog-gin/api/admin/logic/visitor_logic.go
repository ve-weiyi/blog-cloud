package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type VisitorLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewVisitorLogic(svcCtx *svctx.ServiceContext) *VisitorLogic {
	return &VisitorLogic{
		svcCtx: svcCtx,
	}
}

// 分页获取游客列表
func (s *VisitorLogic) FindVisitorList(reqCtx *request.Context, in *types.QueryVisitorReq) (out *types.PageResp, err error) {
	// todo

	return
}
