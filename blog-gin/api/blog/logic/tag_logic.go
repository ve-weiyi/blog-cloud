package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type TagLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewTagLogic(svcCtx *svctx.ServiceContext) *TagLogic {
	return &TagLogic{
		svcCtx: svcCtx,
	}
}

// 分页获取标签列表
func (s *TagLogic) FindTagList(reqCtx *request.Context, in *types.QueryTagReq) (out *types.PageResp, err error) {
	// todo

	return
}
