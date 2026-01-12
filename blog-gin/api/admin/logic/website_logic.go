package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
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

// 获取后台首页信息
func (s *WebsiteLogic) GetAdminHomeInfo(reqCtx *request.Context, in *types.EmptyReq) (out *types.AdminHomeInfo, err error) {
	// todo

	return
}

// 获取关于我的信息
func (s *WebsiteLogic) GetAboutMe(reqCtx *request.Context, in *types.EmptyReq) (out *types.AboutMeVO, err error) {
	// todo

	return
}

// 获取服务器信息
func (s *WebsiteLogic) GetSystemState(reqCtx *request.Context, in *types.EmptyReq) (out *types.Server, err error) {
	// todo

	return
}

// 获取用户分布地区
func (s *WebsiteLogic) GetUserAreaStats(reqCtx *request.Context, in *types.GetUserAreaStatsReq) (out *types.GetUserAreaStatsResp, err error) {
	// todo

	return
}

// 获取访客数据分析
func (s *WebsiteLogic) GetVisitStats(reqCtx *request.Context, in *types.EmptyReq) (out *types.GetVisitStatsResp, err error) {
	// todo

	return
}

// 获取访客数据趋势
func (s *WebsiteLogic) GetVisitTrend(reqCtx *request.Context, in *types.GetVisitTrendReq) (out *types.GetVisitTrendResp, err error) {
	// todo

	return
}

// 获取网站配置
func (s *WebsiteLogic) GetWebsiteConfig(reqCtx *request.Context, in *types.EmptyReq) (out *types.WebsiteConfigVO, err error) {
	// todo

	return
}

// 更新关于我的信息
func (s *WebsiteLogic) UpdateAboutMe(reqCtx *request.Context, in *types.AboutMeVO) (out *types.EmptyResp, err error) {
	// todo

	return
}

// 更新网站配置
func (s *WebsiteLogic) UpdateWebsiteConfig(reqCtx *request.Context, in *types.WebsiteConfigVO) (out *types.EmptyResp, err error) {
	// todo

	return
}
