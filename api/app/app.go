package main

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/config"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/handler"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/job"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/plugins"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/vkit/adapter/nacosx"

	"github.com/ve-weiyi/blog-cloud/infra/middlewarex"
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
		defer func() { _ = nacosCloser.Close() }()
	}

	server := rest.MustNewServer(
		c.RestConf,
		rest.WithCors("*"),
	)
	defer server.Stop()

	server.Use(middlewarex.NewCtxMetaMiddleware().Handle)
	server.Use(middlewarex.NewAntiReplyMiddleware().Handle)
	server.Use(middlewarex.NewDeviceTokenMiddleware().Handle)

	ctx := svc.NewServiceContext(c)
	server.Use(ctx.VisitLog)

	// 启动定时任务
	job.Init(ctx)
	defer job.Stop()

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
