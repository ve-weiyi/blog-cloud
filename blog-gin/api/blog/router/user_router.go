package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
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
	// [TerminalToken UserToken]
	{
		group := r.Group("/blog-api/v1")
		group.Use(s.svcCtx.TerminalToken)
		group.Use(s.svcCtx.UserToken)

		h := handler.NewUserController(s.svcCtx)
		// 删除用户绑定第三方平台账号
		group.POST("/user/delete_user_bind_third_party", h.DeleteUserBindThirdParty)
		// 获取用户信息
		group.GET("/user/get_user_info", h.GetUserInfo)
		// 获取用户点赞列表
		group.GET("/user/get_user_like", h.GetUserLike)
		// 修改用户头像
		group.POST("/user/update_user_avatar", h.UpdateUserAvatar)
		// 修改用户绑定邮箱
		group.POST("/user/update_user_bind_email", h.UpdateUserBindEmail)
		// 修改用户绑定手机号
		group.POST("/user/update_user_bind_phone", h.UpdateUserBindPhone)
		// 修改用户绑定第三方平台账号
		group.POST("/user/update_user_bind_third_party", h.UpdateUserBindThirdParty)
		// 修改用户信息
		group.POST("/user/update_user_info", h.UpdateUserInfo)
		// 修改用户密码
		group.POST("/user/update_user_password", h.UpdateUserPassword)
	}
}
