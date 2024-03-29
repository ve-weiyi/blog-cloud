package logic

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/router/svc"
)

type RemarkRouter struct {
	svcCtx *svc.RouterContext
}

func NewRemarkRouter(svcCtx *svc.RouterContext) *RemarkRouter {
	return &RemarkRouter{
		svcCtx: svcCtx,
	}
}

// 初始化 Remark 路由信息
// publicRouter 公开路由，不登录就可以访问
// loginRouter  登录路由，登录后才可以访问
func (s *RemarkRouter) InitRemarkRouter(publicRouter *gin.RouterGroup, loginRouter *gin.RouterGroup) {

	var handler = s.svcCtx.RemarkController
	{
		loginRouter.POST("remark", handler.CreateRemark)       // 新建Remark
		loginRouter.PUT("remark", handler.UpdateRemark)        // 更新Remark
		loginRouter.DELETE("remark/:id", handler.DeleteRemark) // 删除Remark
		loginRouter.GET("remark/:id", handler.FindRemark)      // 查询Remark

		loginRouter.DELETE("remark/batch_delete", handler.DeleteRemarkByIds) // 批量删除Remark列表
		loginRouter.POST("remark/list", handler.FindRemarkList)              // 分页查询Remark列表
	}
}
