package core

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin"
	"github.com/ve-weiyi/ve-blog-golang/server/api/blog"
	"github.com/ve-weiyi/ve-blog-golang/server/docs"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/middleware"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

// 初始化总路由
func RegisterRouters(engine *gin.Engine, svCtx *svctx.ServiceContext) {
	r := engine.Group("")

	RegisterStaticHandlers(r, svCtx)

	// r.Use(middleware.LoadTls())  // 如果需要使用https 请打开此中间件 然后前往 core/server.go 将启动模式 更变为 r.RunTLS("端口","你的cre/pem文件","你的key文件")
	// r.Use(middleware.CorsByRules()) // 按照配置的规则放行跨域请求
	// r.Use(middleware.GinLogger())  // 访问记录
	r.Use(middleware.Cors())                    // 直接放行全部跨域请求
	r.Use(middleware.TraceMiddleware())         // 打印请求的traceId
	r.Use(middleware.LimitIP(svCtx.RedisEngin)) // 限制IP
	now := time.Now()
	// 健康监测
	r.GET("/api/v1/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version":  "1.0.0",
			"runtime":  now.String(),
			"trace_id": c.Request.Context().Value("X-Trace-ID").(string),
		})
	})

	admin.RegisterHandlers(r, svCtx)
	blog.RegisterHandlers(r, svCtx)
}

func RegisterStaticHandlers(r *gin.RouterGroup, svCtx *svctx.ServiceContext) {
	staticRouter := r.Group("/api/v1")
	// 放行后端静态资源目录，为用户头像和文件提供静态地址
	staticRouter.StaticFS(svCtx.Config.System.RuntimePath, http.Dir(svCtx.Config.System.RuntimePath))

	// Generate Swagger JSON file
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", svCtx.Config.System.Port)
	docs.SwaggerInfo.Version = svCtx.Config.System.Version
	staticRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
