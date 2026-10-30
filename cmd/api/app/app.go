package main

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/rest"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/config"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/handler"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/job"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/plugins"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/vkit/adapter/nacosx"

	"github.com/ve-weiyi/blog-cloud/infra/http/headerx"
	"github.com/ve-weiyi/blog-cloud/infra/http/middlewarex"
)

var (
	nacosHost      = flag.String("nacos-host", "veweiyi.cn", "Input Your Nacos Host")
	nacosPort      = flag.Uint64("nacos-port", 8848, "Input Your Nacos Port")
	nacosUsername  = flag.String("nacos-username", "nacos", "Input Your Nacos Username")
	nacosPassword  = flag.String("nacos-password", "nacos", "Input Your Nacos Password")
	nacosNamespace = flag.String("nacos-namespace", "test", "Input Your Nacos NameSpaceId")
	nacosGroup     = flag.String("nacos-group", "veweiyi.cn", "nacos group")
	nacosDataID    = flag.String("nacos-data-id", "blog-api", "Input Your Nacos DataId")
)

var configFile = flag.String("f", "", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	if *configFile != "" {
		conf.MustLoad(*configFile, &c)
	} else {
		fmt.Println("load config from nacos:", *nacosHost, *nacosPort, *nacosNamespace, *nacosGroup, *nacosDataID)
		nacosCloser, err := nacosx.LoadConfigFromNacos(
			&nacosx.NacosConfig{
				NacosHost:       *nacosHost,
				NacosPort:       *nacosPort,
				NacosNamespace:  *nacosNamespace,
				NacosUsername:   *nacosUsername,
				NacosPassword:   *nacosPassword,
				NacosDataID:     *nacosDataID,
				NacosGroup:      *nacosGroup,
				NacosRuntimeDir: "runtime/blog-api/nacos",
			},
			func(content string) {
				err := conf.LoadFromYamlBytes([]byte(content), &c)
				if err != nil {
					fmt.Printf("nacos config content changed, but failed to load: %v\n", err)
					return
				}
			})
		if err != nil {
			panic(err)
		}
		proc.AddShutdownListener(func() { _ = nacosCloser.Close() })
	}

	server := rest.MustNewServer(
		c.RestConf,
		rest.WithCors("*"),
	)
	server.Use(middlewarex.NewCtxMetaMiddleware(headerx.AppNameBlogWeb).Handle)
	server.Use(middlewarex.NewDeviceCookieMiddleware().Handle)
	server.Use(middlewarex.NewAntiReplyMiddleware().Handle)
	server.Use(middlewarex.NewDeviceTokenMiddleware().Handle)

	ctx := svc.NewServiceContext(c)
	server.Use(ctx.VisitLog)

	// 关闭必须注册为关闭监听：go-zero 收到信号后会跑监听并强制退出，main 不会返回，
	// 写在 main 里的 defer 不可靠。
	job.Init(ctx)
	proc.AddShutdownListener(func() {
		ctx.Shutdown()
		job.Stop()
	})

	handler.RegisterHandlers(server, ctx)
	plugins.RegisterPluginHandlers(server, ctx)

	server.PrintRoutes()
	//httpx.SetErrorHandler(func(err error) (int, interface{}) {
	//	return http.StatusInternalServerError, err
	//})
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	fmt.Printf(`
	默认接口文档地址: http://localhost:%d/api/v1/swagger/index.html
`, c.Port)

	server.Start()
}
