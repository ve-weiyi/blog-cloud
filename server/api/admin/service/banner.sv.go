package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
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

// 创建页面
func (s *BannerService) AddBanner(reqCtx *request.Context, in *dto.BannerNewReq) (out *dto.BannerBackDTO, err error) {
	// todo

	return
}

// 删除页面
func (s *BannerService) DeleteBanner(reqCtx *request.Context, in *dto.IdReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 分页获取页面列表
func (s *BannerService) FindBannerList(reqCtx *request.Context, in *dto.BannerQuery) (out *dto.PageResp, err error) {
	// todo

	return
}

// 更新页面
func (s *BannerService) UpdateBanner(reqCtx *request.Context, in *dto.BannerNewReq) (out *dto.BannerBackDTO, err error) {
	// todo

	return
}
