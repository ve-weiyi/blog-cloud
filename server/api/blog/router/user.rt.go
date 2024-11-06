package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/controller"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type UserRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewUserRouter(svcCtx *svctx.ServiceContext) *UserRouter {
	return &UserRouter{
		svcCtx: svcCtx,
	}
}

func (s *UserRouter) Register(r *gin.RouterGroup) {
	// User
	// [SignToken JwtToken]
	{
		group := r.Group("/api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)
		group.Use(s.svcCtx.MiddlewareJwtToken)

		handler := controller.NewUserController(s.svcCtx)
		// 获取用户信息
		group.GET("/user/get_user_info", handler.GetUserInfo)
		// 获取用户点赞列表
		group.GET("/user/get_user_like", handler.GetUserLike)
		// 修改用户头像
		group.POST("/user/update_user_avatar", handler.UpdateUserAvatar)
		// 修改用户信息
		group.POST("/user/update_user_info", handler.UpdateUserInfo)
	}
}
