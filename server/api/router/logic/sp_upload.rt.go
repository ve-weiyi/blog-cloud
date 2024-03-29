package logic

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/router/svc"
)

type UploadRouter struct {
	svcCtx *svc.RouterContext
}

func NewUploadRouter(svcCtx *svc.RouterContext) *UploadRouter {
	return &UploadRouter{
		svcCtx: svcCtx,
	}
}

// 初始化 Upload 路由信息
// publicRouter 公开路由，不登录就可以访问
// loginRouter  登录路由，登录后才可以访问
func (s *UploadRouter) InitUploadRouter(publicRouter *gin.RouterGroup, loginRouter *gin.RouterGroup) {

	var handler = s.svcCtx.UploadController
	{
		publicRouter.POST("upload/*label", handler.UploadFile) // 上传文件
	}
}
