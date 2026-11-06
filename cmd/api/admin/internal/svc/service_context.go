package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-openapi/loads"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/docs"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/chat"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/config"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/middleware"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/middleware/permissionx"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/middleware/tracelogx"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/accessservice"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/authservice"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/chatservice"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/contentservice"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/discussionservice"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/guestservice"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/mediaservice"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/siteservice"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/statsservice"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/syslogservice"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/userservice"
	"github.com/ve-weiyi/blog-cloud/infra/auth/tokenx"
	"github.com/ve-weiyi/blog-cloud/infra/auth/tokenx/codec/jwtcodec"
	"github.com/ve-weiyi/blog-cloud/infra/auth/tokenx/codec/opaque"
	"github.com/ve-weiyi/blog-cloud/infra/auth/tokenx/store/redisstore"
	"github.com/ve-weiyi/blog-cloud/infra/domain/cachekey"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/interceptorx"
	"github.com/ve-weiyi/blog-cloud/infra/http/middlewarex"
	"github.com/ve-weiyi/blog-cloud/infra/http/middlewarex/limitx"
	"github.com/ve-weiyi/blog-cloud/infra/messaging/notifyx"
	"github.com/ve-weiyi/blog-cloud/infra/messaging/ssex"
	"github.com/ve-weiyi/blog-cloud/infra/messaging/stompwsx"
	"github.com/ve-weiyi/blog-cloud/infra/store/captchax"
	"github.com/ve-weiyi/blog-cloud/infra/store/presencex"
	"github.com/ve-weiyi/blog-cloud/infra/store/storex"
	"github.com/ve-weiyi/stompws/logws"
	"github.com/ve-weiyi/stompws/server"
	"github.com/ve-weiyi/vkit/adapter/storagex"
)

type ServiceContext struct {
	Config config.Config

	AdminAuth         rest.Middleware
	Permission        rest.Middleware
	OperationLog      rest.Middleware
	RateLimit         rest.Middleware
	NotifyStreamLimit rest.Middleware

	RedisClient     *redis.Client
	TokenManager    tokenx.Service
	CaptchaStore    *captchax.Store
	StorageProvider storagex.StorageProvider // 存储提供者

	StompHubServer *server.StompHubServer

	// SitePresence 是博客前台的在线名册，本服务只读 —— 名单页的数据源。
	SitePresence *presencex.Store

	// NotifyBroker 通知事件广播器（每频道一个 Broker，各 SSE 连接从它分食）
	NotifyBroker *ssex.Broker[notifyx.Event]

	// stopNotify 停止通知事件订阅，由 Shutdown 调用
	stopNotify context.CancelFunc

	AuthService         authservice.AuthService
	ChatService         chatservice.ChatService
	UserService         userservice.UserService
	StatsService        statsservice.StatsService
	NotificationService notificationservice.NotificationService
	SyslogService       syslogservice.SyslogService
	GuestService        guestservice.GuestService
	ContentService      contentservice.ContentService
	DiscussionService   discussionservice.DiscussionService
	MediaService        mediaservice.MediaService
	SiteService         siteservice.SiteService
	AccessService       accessservice.AccessService
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds, err := ConnectRedis(c.RedisConf)
	if err != nil {
		panic(err)
	}

	// 图形验证码 15 分钟过期
	captchaStore := captchax.New(storex.NewRedisStore(rds), cachekey.CaptchaStorePrefixAdmin, 15*time.Minute)

	accessCodec, err := jwtcodec.New(c.TokenConf.Secret)
	if err != nil {
		panic(fmt.Errorf("admin token codec: %w", err))
	}

	sessionStore, err := redisstore.New(rds, cachekey.TokenStorePrefixAdmin)
	if err != nil {
		panic(err)
	}

	tokenManager, err := tokenx.New(tokenx.Config{
		AccessCodec:        accessCodec,
		RefreshCodec:       opaque.New(),
		Store:              sessionStore,
		LoginMode:          tokenx.SinglePoint,
		AccessTTL:          c.TokenConf.AccessTokenTTL,
		RefreshIdleTTL:     c.TokenConf.RefreshTokenTTL,
		RefreshAbsoluteTTL: c.TokenConf.RefreshAbsoluteTTL,
	})
	if err != nil {
		panic(err)
	}

	storageProvider, err := storagex.NewStorageProvider(&c.StorageConfig)
	if err != nil {
		panic(err)
	}

	// 后台自己的在线名册：只有本服务（后台的 hub）往里写。
	presence := presencex.NewStore(storex.NewRedisMemberStore(rds), cachekey.OnlineAdminKey, presencex.DefaultTTL)
	// 博客前台的在线名册：本服务**只读**，供在线用户名单页用（名册归前台服务写）。
	sitePresence := presencex.NewStore(storex.NewRedisMemberStore(rds), cachekey.OnlineUserKey, presencex.DefaultTTL)

	hub := server.NewStompHubServer(
		server.WithPresenceManager(presence),
		server.WithPresenceTTL(presence.TTL()),
		server.WithEventHooks(
			chat.NewChatRoomEventHook(),
			stompwsx.NewOnlineBroadcastHook(),
			stompwsx.NewOnlineCatchupHook(),
		),
		server.WithAuthenticator(stompwsx.NewAuthenticator(tokenManager, false)),
		server.WithLogger(logws.NewDefaultLogger()),
	)

	// 通知事件扇出：Redis 订阅是进程级一份，各 SSE 连接从广播器分食
	notifyCtx, stopNotify := context.WithCancel(context.Background())
	notifyBroker := ssex.NewBroker[notifyx.Event]()
	ssex.NewSubscriber(rds, notifyBroker).Start(notifyCtx)

	doc, err := loads.Analyzed(json.RawMessage(docs.Docs), "")
	if err != nil {
		panic(err)
	}

	rpcClient := zrpc.MustNewClient(c.AppRpcConf,
		zrpc.WithUnaryClientInterceptor(interceptorx.ClientErrorInterceptor),
	)

	authService := authservice.NewAuthService(rpcClient)
	chatService := chatservice.NewChatService(rpcClient)
	userService := userservice.NewUserService(rpcClient)
	statsService := statsservice.NewStatsService(rpcClient)
	notificationService := notificationservice.NewNotificationService(rpcClient)
	syslogService := syslogservice.NewSyslogService(rpcClient)
	guestService := guestservice.NewGuestService(rpcClient)
	contentService := contentservice.NewContentService(rpcClient)
	discussionService := discussionservice.NewDiscussionService(rpcClient)
	mediaService := mediaservice.NewMediaService(rpcClient)
	siteService := siteservice.NewSiteService(rpcClient)
	accessService := accessservice.NewAccessService(rpcClient)

	return &ServiceContext{
		Config:              c,
		AdminAuth:           middleware.NewAdminAuthMiddleware(tokenManager).Handle,
		RateLimit:           middlewarex.NewRateLimitMiddleware(limitx.NewPeriodLimit(60, 10, rds, cachekey.RateLimitStrictPrefix)).Handle,
		NotifyStreamLimit:   middleware.NewNotifyStreamLimitMiddleware(notifyBroker, middleware.DefaultStreamLimit).Handle,
		Permission:          middleware.NewPermissionMiddleware(permissionx.NewRbacEnforcer(rds, accessService)).Handle,
		OperationLog:        middleware.NewOperationLogMiddleware(doc.Spec(), tracelogx.NewTraceEnforcer(rds, accessService), syslogService).Handle,
		RedisClient:         rds,
		TokenManager:        tokenManager,
		CaptchaStore:        captchaStore,
		StorageProvider:     storageProvider,
		StompHubServer:      hub,
		SitePresence:        sitePresence,
		NotifyBroker:        notifyBroker,
		stopNotify:          stopNotify,
		AuthService:         authService,
		ChatService:         chatService,
		UserService:         userService,
		StatsService:        statsService,
		NotificationService: notificationService,
		SyslogService:       syslogService,
		GuestService:        guestService,
		ContentService:      contentService,
		DiscussionService:   discussionService,
		MediaService:        mediaService,
		SiteService:         siteService,
		AccessService:       accessService,
	}
}

// Shutdown 停止后台推送与长连接：先停 Redis 订阅，再关闭各 SSE 连接的事件通道，
// 最后等待 STOMP 连接收尾。必须注册为关闭监听，defer 在信号路径上不会执行。
func (s *ServiceContext) Shutdown() {
	if s.stopNotify != nil {
		s.stopNotify()
	}
	if s.NotifyBroker != nil {
		s.NotifyBroker.Close()
	}
	if s.StompHubServer != nil {
		if err := s.StompHubServer.Shutdown(5 * time.Second); err != nil {
			logx.Errorf("stop stomp hub: %v", err)
		}
	}
}

func ConnectRedis(c config.RedisConf) (*redis.Client, error) {
	address := c.Host + ":" + c.Port
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Username: "",
		Password: c.Password,
		DB:       c.DB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("redis 连接失败: %v", err)
	}

	return client, nil
}
