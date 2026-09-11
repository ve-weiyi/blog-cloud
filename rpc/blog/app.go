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

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/config"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/analyticsrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/configrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/socialrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/userauthrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/userrpc"
	analyticsserviceServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/analyticsservice"
	articleserviceServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/articleservice"
	discussionserviceServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/discussionservice"
	guestrpcServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/guestservice"
	permissionserviceServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/permissionservice"
	resourceserviceServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/resourceservice"
	socialserviceServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/socialservice"
	syslogserviceServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/syslogservice"
	userauthserviceServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/userauthservice"
	userserviceServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/userservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/adapter/nacosx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/mq/mqlogic"
	"github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/pb/guestrpc"
	configserviceServer "github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/server/configservice"
	notificationserviceServer "github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/server/notificationservice"

	"github.com/ve-weiyi/blog-cloud/infra/interceptorx"
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
		err := nacosx.LoadConfigFromNacos(
			&nacosx.NacosConfig{
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

	// 初始化消息队列并注册消费者
	mq.Init(ctx)
	mqlogic.RegisterMqConsumers(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		// 用户认证
		userauthrpc.RegisterUserAuthServiceServer(grpcServer, userauthserviceServer.NewUserAuthServiceServer(ctx))
		// 用户管理
		userrpc.RegisterUserServiceServer(grpcServer, userserviceServer.NewUserServiceServer(ctx))
		// 游客信息
		guestrpc.RegisterGuestServiceServer(grpcServer, guestrpcServer.NewGuestServiceServer(ctx))
		// 权限管理
		permissionrpc.RegisterPermissionServiceServer(grpcServer, permissionserviceServer.NewPermissionServiceServer(ctx))
		// 文章管理
		articlerpc.RegisterArticleServiceServer(grpcServer, articleserviceServer.NewArticleServiceServer(ctx))
		// 消息互动
		discussionrpc.RegisterDiscussionServiceServer(grpcServer, discussionserviceServer.NewDiscussionServiceServer(ctx))
		// 社交管理
		socialrpc.RegisterSocialServiceServer(grpcServer, socialserviceServer.NewSocialServiceServer(ctx))
		// 资源管理
		resourcerpc.RegisterResourceServiceServer(grpcServer, resourceserviceServer.NewResourceServiceServer(ctx))
		// 通知管理
		notificationrpc.RegisterNotificationServiceServer(grpcServer, notificationserviceServer.NewNotificationServiceServer(ctx))
		// 系统日志
		syslogrpc.RegisterSyslogServiceServer(grpcServer, syslogserviceServer.NewSyslogServiceServer(ctx))
		// 配置管理
		configrpc.RegisterConfigServiceServer(grpcServer, configserviceServer.NewConfigServiceServer(ctx))
		// 数据分析
		analyticsrpc.RegisterAnalyticsServiceServer(grpcServer, analyticsserviceServer.NewAnalyticsServiceServer(ctx))

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
