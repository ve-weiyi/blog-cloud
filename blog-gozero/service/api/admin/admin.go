package main

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"

	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/nacos"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/infra/middlewarex"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/config"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/plugins"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
)

var (
	nacosHost      = flag.String("nacos-host", "veweiyi.cn", "Input Your Nacos Host")
	nacosPort      = flag.Uint64("nacos-port", 8848, "Input Your Nacos Port")
	nacosUsername  = flag.String("nacos-username", "nacos", "Input Your Nacos Username")
	nacosPassword  = flag.String("nacos-password", "nacos", "Input Your Nacos Password")
	nacosNamespace = flag.String("nacos-namespace", "test", "Input Your Nacos NameSpaceId")
	nacosGroup     = flag.String("nacos-group", "veweiyi.cn", "nacos group")
	nacosDataID    = flag.String("nacos-data-id", "admin-api", "Input Your Nacos DataId")
)

var configFile = flag.String("f", "", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	if *configFile != "" {
		conf.MustLoad(*configFile, &c)
	} else {
		fmt.Println("load config from nacos:", *nacosHost, *nacosPort, *nacosNamespace, *nacosGroup, *nacosDataID)
		err := nacos.LoadConfigFromNacos(
			&nacos.NacosConfig{
				NacosHost:       *nacosHost,
				NacosPort:       *nacosPort,
				NacosNamespace:  *nacosNamespace,
				NacosUsername:   *nacosUsername,
				NacosPassword:   *nacosPassword,
				NacosDataID:     *nacosDataID,
				NacosGroup:      *nacosGroup,
				NacosRuntimeDir: "runtime/admin-api/nacos",
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
	}

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	server.Use(middlewarex.NewCorsMiddleware().Handle)
	server.Use(middlewarex.NewCtxMetaMiddleware().Handle)
	server.Use(middlewarex.NewAntiReplyMiddleware().Handle)
	server.Use(middlewarex.NewTerminalTokenMiddleware().Handle)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	plugins.RegisterPluginHandlers(server, ctx)

	server.PrintRoutes()
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	fmt.Printf(`
	默认接口文档地址: http://localhost:%d/admin-api/v1/swagger/index.html
`, c.Port)
	server.Start()
}
