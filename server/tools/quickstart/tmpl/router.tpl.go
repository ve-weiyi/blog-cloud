package tmpl

const AppRouter = `
package router

import (

)

type AppRouter struct {
	svcCtx *svc.RouterContext //持有的controller引用
}

func NewRouter(svcCtx *svc.RouterContext) *AppRouter {
	return &AppRouter{
		svcCtx: svcCtx,
	}
}
`
const RouterContext = `
package svc

import (

)

// 注册需要用到的api
type RouterContext struct {
	*controller.AppController
}

func NewRouterContext(cfg *config.Config) *RouterContext {
	ctx := svc.NewControllerContext(cfg)
	ctl := controller.NewController(ctx)
	if ctl == nil {
		panic("ctl cannot be null")
	}

	return &RouterContext{
		AppController: ctl,
	}
}

`

const Router = `
package logic

import (
	"github.com/gin-gonic/gin"

	{{range .ImportPkgPaths}}{{.}} ` + "\n" + `{{end}}
)

type {{.UpperStartCamelName}}Router struct {
	svcCtx *svc.RouterContext
}

func New{{.UpperStartCamelName}}Router(svcCtx *svc.RouterContext) *{{.UpperStartCamelName}}Router {
	return &{{.UpperStartCamelName}}Router{
		svcCtx: svcCtx,
	}
}

// 初始化 {{.CommentName}} 路由信息
// publicRouter 公开路由，不登录就可以访问
// loginRouter  登录路由，登录后才可以访问
func (s *{{.UpperStartCamelName}}Router) Init{{.UpperStartCamelName}}Router(publicRouter *gin.RouterGroup, loginRouter *gin.RouterGroup) {

	var handler = s.svcCtx.{{.UpperStartCamelName}}Controller
	{
		loginRouter.POST("{{.SnakeName}}", handler.Create{{.UpperStartCamelName}})             // 新建{{.CommentName}}
		loginRouter.PUT("{{.SnakeName}}", handler.Update{{.UpperStartCamelName}})              // 更新{{.CommentName}}
		loginRouter.DELETE("{{.SnakeName}}/:id", handler.Delete{{.UpperStartCamelName}})           // 删除{{.CommentName}}
		loginRouter.DELETE("{{.SnakeName}}/batch_delete", handler.Delete{{.UpperStartCamelName}}ByIds)	// 批量删除{{.CommentName}}列表

		publicRouter.GET("{{.SnakeName}}/:id", handler.Find{{.UpperStartCamelName}})				  // 查询{{.UpperStartCamelName}}
		publicRouter.POST("{{.SnakeName}}/list", handler.Find{{.UpperStartCamelName}}List)  				// 分页查询{{.CommentName}}列表
	}
}
`

const CommonRouter = `
package logic

import (

)

type {{.UpperStartCamelName}}Router struct {
	svcCtx *svc.RouterContext
}

func New{{.UpperStartCamelName}}Router(svcCtx *svc.RouterContext) *{{.UpperStartCamelName}}Router {
	return &{{.UpperStartCamelName}}Router{
		svcCtx: svcCtx,
	}
}
`
