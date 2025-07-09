package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type BlogApiLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewBlogApiLogic(svcCtx *svctx.ServiceContext) *BlogApiLogic {
	return &BlogApiLogic{
		svcCtx: svcCtx,
	}
}

// ping
func (s *BlogApiLogic) Ping(reqCtx *request.Context, in *types.PingReq) (out *types.PingResp, err error) {
	// todo

	return
}
