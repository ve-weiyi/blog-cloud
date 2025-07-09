package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type WebsiteLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewWebsiteLogic(svcCtx *svctx.ServiceContext) *WebsiteLogic {
	return &WebsiteLogic{
		svcCtx: svcCtx,
	}
}

// 获取博客前台首页信息
func (s *WebsiteLogic) GetBlogHomeInfo(reqCtx *request.Context, in *types.GetBlogHomeInfoReq) (out *types.GetBlogHomeInfoResp, err error) {
	// todo

	return
}

// 获取关于我的信息
func (s *WebsiteLogic) GetAboutMe(reqCtx *request.Context, in *types.GetAboutMeReq) (out *types.GetAboutMeResp, err error) {
	// todo

	return
}
