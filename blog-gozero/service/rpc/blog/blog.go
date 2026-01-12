package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/infra/interceptorx"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/config"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/configrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/socialrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/websiterpc"
	accountrpcServer "github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/server/accountrpc"
	articlerpcServer "github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/server/articlerpc"
	configrpcServer "github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/server/configrpc"
	messagerpcServer "github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/server/messagerpc"
	permissionrpcServer "github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/server/permissionrpc"
	resourcerpcServer "github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/server/resourcerpc"
	socialrpcServer "github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/server/socialrpc"
	syslogrpcServer "github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/server/syslogrpc"
	websiterpcServer "github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/server/websiterpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/nacos"
)

var (
	nacosHost      = flag.String("nacos-host", "veweiyi.cn", "Input Your Nacos Host")
	nacosPort      = flag.Uint64("nacos-port", 8848, "Input Your Nacos Port")
	nacosUsername  = flag.String("nacos-username", "nacos", "Input Your Nacos Username")
	nacosPassword  = flag.String("nacos-password", "nacos", "Input Your Nacos Password")
	nacosNamespace = flag.String("nacos-namespace", "test", "Input Your Nacos NameSpaceId")
	nacosGroup     = flag.String("nacos-group", "veweiyi.cn", "nacos group")
	nacosDataID    = flag.String("nacos-data-id", "blog-rpc", "Input Your Nacos DataId")
)

var configFile = flag.String("f", "", "the config file")

func main() {
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Llongfile)
	var c config.Config
	if *configFile != "" {
		fmt.Println("load config from file:", *configFile)
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
				NacosRuntimeDir: "runtime/blog-rpc/nacos",
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

	ctx := svc.NewServiceContext(c)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		accountrpc.RegisterAccountRpcServer(grpcServer, accountrpcServer.NewAccountRpcServer(ctx))
		articlerpc.RegisterArticleRpcServer(grpcServer, articlerpcServer.NewArticleRpcServer(ctx))
		permissionrpc.RegisterPermissionRpcServer(grpcServer, permissionrpcServer.NewPermissionRpcServer(ctx))
		messagerpc.RegisterMessageRpcServer(grpcServer, messagerpcServer.NewMessageRpcServer(ctx))
		socialrpc.RegisterSocialRpcServer(grpcServer, socialrpcServer.NewSocialRpcServer(ctx))
		syslogrpc.RegisterSyslogRpcServer(grpcServer, syslogrpcServer.NewSyslogRpcServer(ctx))
		websiterpc.RegisterWebsiteRpcServer(grpcServer, websiterpcServer.NewWebsiteRpcServer(ctx))
		configrpc.RegisterConfigRpcServer(grpcServer, configrpcServer.NewConfigRpcServer(ctx))
		resourcerpc.RegisterResourceRpcServer(grpcServer, resourcerpcServer.NewResourceRpcServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	s.AddUnaryInterceptors(interceptorx.ServerMetaInterceptor)
	s.AddUnaryInterceptors(interceptorx.ServerErrorInterceptor)
	s.AddUnaryInterceptors(interceptorx.ServerLogInterceptor)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
