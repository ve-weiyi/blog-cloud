package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/config"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/job"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq/consumer"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq/outbox"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/accessrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/authrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/chatrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/contentrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/guestrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/mediarpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/siterpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/statsrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/userrpc"
	accessserviceServer "github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/server/accessservice"
	authserviceServer "github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/server/authservice"
	chatserviceServer "github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/server/chatservice"
	contentserviceServer "github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/server/contentservice"
	discussionserviceServer "github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/server/discussionservice"
	guestserviceServer "github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/server/guestservice"
	mediaserviceServer "github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/server/mediaservice"
	notificationserviceServer "github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/server/notificationservice"
	siteserviceServer "github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/server/siteservice"
	statsserviceServer "github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/server/statsservice"
	syslogserviceServer "github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/server/syslogservice"
	userserviceServer "github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/server/userservice"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/interceptorx"
	"github.com/ve-weiyi/blog-cloud/infra/metrics"
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
		nacosCloser, err := nacosx.LoadConfigFromNacos(
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
		defer func() { _ = nacosCloser.Close() }()
	}

	svcCtx := svc.NewServiceContext(c)

	// 先声明消费拓扑，再启动 outbox relay：relay 可能立刻发布积压事件，
	// 队列尚未绑定时消息会被 exchange 直接丢弃。
	// 两者失败都属于「不修就不会好」的部署错误，启动即失败；
	// broker 抖动由懒连接与 relay 退避重试兜底，不受影响。
	if !startMQ(svcCtx) {
		panic("mq: notification subscriptions unavailable")
	}
	relay, err := startOutboxRelay(svcCtx)
	if err != nil {
		panic(fmt.Errorf("outbox relay: %w", err))
	}

	jobs := job.New()
	jobs.Add(job.Task{Name: "outbox-relay", Start: relay.Start, Stop: relay.Stop})
	jobs.Start()

	stopJobs := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		jobs.Stop(ctx)
	}

	// 关闭顺序：先停后台任务（relay 可能正在发布），再释放 MQ/Redis/DB。
	// go-zero 的多个 shutdown listener 是并发执行的，顺序只能靠合并成一个来保证。
	shutdown := func() {
		stopJobs()
		svcCtx.Close()
	}
	defer shutdown()
	proc.AddShutdownListener(shutdown)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		// 用户
		userrpc.RegisterUserServiceServer(grpcServer, userserviceServer.NewUserServiceServer(svcCtx))
		// 访客
		guestrpc.RegisterGuestServiceServer(grpcServer, guestserviceServer.NewGuestServiceServer(svcCtx))
		// 认证
		authrpc.RegisterAuthServiceServer(grpcServer, authserviceServer.NewAuthServiceServer(svcCtx))
		// 访问控制
		accessrpc.RegisterAccessServiceServer(grpcServer, accessserviceServer.NewAccessServiceServer(svcCtx))
		// 内容
		contentrpc.RegisterContentServiceServer(grpcServer, contentserviceServer.NewContentServiceServer(svcCtx))
		// 媒体
		mediarpc.RegisterMediaServiceServer(grpcServer, mediaserviceServer.NewMediaServiceServer(svcCtx))
		// 讨论
		discussionrpc.RegisterDiscussionServiceServer(grpcServer, discussionserviceServer.NewDiscussionServiceServer(svcCtx))
		// 聊天
		chatrpc.RegisterChatServiceServer(grpcServer, chatserviceServer.NewChatServiceServer(svcCtx))
		// 站点
		siterpc.RegisterSiteServiceServer(grpcServer, siteserviceServer.NewSiteServiceServer(svcCtx))
		// 通知
		notificationrpc.RegisterNotificationServiceServer(grpcServer, notificationserviceServer.NewNotificationServiceServer(svcCtx))
		// 日志
		syslogrpc.RegisterSyslogServiceServer(grpcServer, syslogserviceServer.NewSyslogServiceServer(svcCtx))
		// 统计
		statsrpc.RegisterStatsServiceServer(grpcServer, statsserviceServer.NewStatsServiceServer(svcCtx))

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

// startOutboxRelay 先确保 outbox 表存在，再构造 relay。
func startOutboxRelay(svcCtx *svc.ServiceContext) (*outbox.Relay, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := svcCtx.Outbox.Migrate(ctx); err != nil {
		return nil, err
	}
	if err := svcCtx.Outbox.Ping(ctx); err != nil {
		return nil, err
	}
	return outbox.NewRelay(svcCtx.Outbox, svcCtx.MQ.Publisher,
		outbox.WithMetrics(metrics.NewOutbox())), nil
}

// startMQ 注册事件处理器并启动通知域与审计域的订阅。
// 返回通知域订阅是否全部就绪；未就绪时不启动 outbox relay，避免消息被无绑定的 exchange 丢弃。
func startMQ(svcCtx *svc.ServiceContext) bool {
	router, err := svcCtx.MQ.Runtime.NewRouter()
	if err != nil {
		logx.Errorf("[mq] create router: %v", err)
		return false
	}
	if err := consumer.Register(router, svcCtx); err != nil {
		// 注册失败意味着处理器缺失：宁可不启动订阅，也不要让该事件在运行期
		// 全部按失败策略处理（死信或丢弃）。
		logx.Errorf("[mq] register handlers: %v", err)
		return false
	}

	// 任一订阅未就绪都不启动 relay：没有绑定的 exchange 会静默丢消息。
	allReady := true
	for _, spec := range mq.Subscriptions() {
		if err := svcCtx.MQ.Runtime.Start(router, spec); err != nil {
			logx.Errorf("[mq] start subscription %s: %v", spec.Group, err)
			allReady = false
		}
	}
	return allReady
}
