package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type WebsiteService struct {
	svcCtx *svctx.ServiceContext
}

func NewWebsiteService(svcCtx *svctx.ServiceContext) *WebsiteService {
	return &WebsiteService{
		svcCtx: svcCtx,
	}
}

// 获取后台首页信息
func (s *WebsiteService) GetAdminHomeInfo(reqCtx *request.Context, in *dto.EmptyReq) (out *dto.AdminHomeInfo, err error) {
	// todo

	return
}

// 获取关于我的信息
func (s *WebsiteService) GetAboutMe(reqCtx *request.Context, in *dto.EmptyReq) (out *dto.AboutMe, err error) {
	// todo

	return
}

// 更新关于我的信息
func (s *WebsiteService) UpdateAboutMe(reqCtx *request.Context, in *dto.AboutMe) (out *dto.EmptyResp, err error) {
	// todo

	return
}

// 获取网站配置
func (s *WebsiteService) GetWebsiteConfig(reqCtx *request.Context, in *dto.EmptyReq) (out *dto.WebsiteConfig, err error) {
	// todo

	return
}

// 获取服务器信息
func (s *WebsiteService) GetSystemState(reqCtx *request.Context, in *dto.EmptyReq) (out *dto.Server, err error) {
	// todo

	return
}

// 更新网站配置
func (s *WebsiteService) UpdateWebsiteConfig(reqCtx *request.Context, in *dto.WebsiteConfig) (out *dto.EmptyResp, err error) {
	// todo

	return
}
