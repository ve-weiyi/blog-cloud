package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type BannerService struct {
	svcCtx *svctx.ServiceContext
}

func NewBannerService(svcCtx *svctx.ServiceContext) *BannerService {
	return &BannerService{
		svcCtx: svcCtx,
	}
}

// 分页获取页面列表
func (s *BannerService) FindBannerList(reqCtx *request.Context, in *dto.BannerQueryReq) (out *dto.PageResp, err error) {
	// todo

	return
}
