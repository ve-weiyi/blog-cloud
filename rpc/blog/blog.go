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

	"github.com/ve-weiyi/blog-cloud/infra/interceptorx"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/config"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq/mqlogic"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/accessrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/authrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/chatrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/contentrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/guestrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/mediarpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/siterpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/statsrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/userrpc"
	accessserviceServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/accessservice"
	authserviceServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/authservice"
	chatserviceServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/chatservice"
	contentserviceServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/contentservice"
	discussionserviceServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/discussionservice"
	guestserviceServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/guestservice"
	mediaserviceServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/mediaservice"
	notificationserviceServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/notificationservice"
	siteserviceServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/siteservice"
	statsserviceServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/statsservice"
	syslogserviceServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/syslogservice"
	userserviceServer "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/server/userservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/adapter/nacosx"
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
		// 用户
		userrpc.RegisterUserServiceServer(grpcServer, userserviceServer.NewUserServiceServer(ctx))
		// 访客
		guestrpc.RegisterGuestServiceServer(grpcServer, guestserviceServer.NewGuestServiceServer(ctx))
		// 认证
		authrpc.RegisterAuthServiceServer(grpcServer, authserviceServer.NewAuthServiceServer(ctx))
		// 访问控制
		accessrpc.RegisterAccessServiceServer(grpcServer, accessserviceServer.NewAccessServiceServer(ctx))
		// 内容
		contentrpc.RegisterContentServiceServer(grpcServer, contentserviceServer.NewContentServiceServer(ctx))
		// 媒体
		mediarpc.RegisterMediaServiceServer(grpcServer, mediaserviceServer.NewMediaServiceServer(ctx))
		// 讨论
		discussionrpc.RegisterDiscussionServiceServer(grpcServer, discussionserviceServer.NewDiscussionServiceServer(ctx))
		// 聊天
		chatrpc.RegisterChatServiceServer(grpcServer, chatserviceServer.NewChatServiceServer(ctx))
		// 站点
		siterpc.RegisterSiteServiceServer(grpcServer, siteserviceServer.NewSiteServiceServer(ctx))
		// 通知
		notificationrpc.RegisterNotificationServiceServer(grpcServer, notificationserviceServer.NewNotificationServiceServer(ctx))
		// 日志
		syslogrpc.RegisterSyslogServiceServer(grpcServer, syslogserviceServer.NewSyslogServiceServer(ctx))
		// 统计
		statsrpc.RegisterStatsServiceServer(grpcServer, statsserviceServer.NewStatsServiceServer(ctx))

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
