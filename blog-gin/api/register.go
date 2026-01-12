package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/otel/trace"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/middleware"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/docs"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

// 初始化总路由
func RegisterRouters(engine *gin.Engine, svCtx *svctx.ServiceContext) {
	v1 := engine.Group("/api/v1")
	RegisterPingHandlers(v1, svCtx)
	RegisterStaticHandlers(v1, svCtx)
	RegisterSwaggerHandlers(v1, svCtx)

	g := engine.Group("")

	// r.Use(middleware.GinLogger())  // 访问记录，gin.Default()自带了
	// r.Use(middleware.GinRecovery(true)) // 捕获panic，gin.Default()自带了
	// r.Use(middleware.Limit(svCtx.RedisEngin)) // 请求限频，一般在nginx层处理。
	g.Use(middleware.Cors())  // 直接放行全部跨域请求
	g.Use(middleware.Trace()) // 打印请求的traceId

	admin.RegisterHandlers(g, svCtx)
	blog.RegisterHandlers(g, svCtx)
}

// 健康检查
func RegisterPingHandlers(r *gin.RouterGroup, svCtx *svctx.ServiceContext) {
	now := time.Now()
	// 健康监测
	r.GET("/version", func(c *gin.Context) {
		spanCtx := trace.SpanContextFromContext(c.Request.Context())

		c.JSON(http.StatusOK, gin.H{
			"version":  svCtx.Config.System.Version,
			"runtime":  now.String(),
			"trace_id": spanCtx.TraceID().String(),
		})
	})
}

// 静态资源处理
func RegisterStaticHandlers(r *gin.RouterGroup, svCtx *svctx.ServiceContext) {
	// 放行后端静态资源目录，为用户头像和文件提供静态地址
	r.StaticFS(svCtx.Config.System.RuntimePath, http.Dir(svCtx.Config.System.RuntimePath))
}

func RegisterSwaggerHandlers(r *gin.RouterGroup, svCtx *svctx.ServiceContext) {
	// Generate Swagger JSON file
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", svCtx.Config.System.Port)
	docs.SwaggerInfo.Version = svCtx.Config.System.Version
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
