package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-openapi/loads"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/ve-weiyi/blog-cloud/api/admin/docs"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/common/sse"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/common/stomphook"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/config"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/middleware"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/middleware/permissionx"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/middleware/tracelogx"
	"github.com/ve-weiyi/blog-cloud/infra/middlewarex/limitx"
	"github.com/ve-weiyi/blog-cloud/infra/storex/captchax"
	tokenx2 "github.com/ve-weiyi/blog-cloud/infra/storex/tokenx"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/accessservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/authservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/chatservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/contentservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/discussionservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/guestservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/mediaservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/siteservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/statsservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/syslogservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/userservice"
	"github.com/ve-weiyi/stompws/logws"
	"github.com/ve-weiyi/stompws/server/client"
	"github.com/ve-weiyi/vkit/adapter/storagex"

	"github.com/ve-weiyi/blog-cloud/infra/constants/cachekey"
	"github.com/ve-weiyi/blog-cloud/infra/interceptorx"
	"github.com/ve-weiyi/blog-cloud/infra/middlewarex"
	"github.com/ve-weiyi/blog-cloud/infra/notifyx"
	"github.com/ve-weiyi/blog-cloud/infra/storex"
)

type ServiceContext struct {
	Config config.Config

	AdminAuth         rest.Middleware
	Permission        rest.Middleware
	OperationLog      rest.Middleware
	RateLimit         rest.Middleware
	NotifyStreamLimit rest.Middleware

	RedisClient     *redis.Client
	TokenManager    tokenx2.Manager
	CaptchaStore    *captchax.Store
	StorageProvider storagex.StorageProvider // 存储提供者

	StompHubServer *client.StompHubServer

	// NotifyBroker 通知事件广播器（每频道一个 Broker，各 SSE 连接从它分食）
	NotifyBroker *sse.Broker[notifyx.Event]

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

	tokenManager, err := tokenx2.New(tokenx2.Config{
		Signer:          tokenx2.NewJWTSigner([]byte(c.Name), c.Name),
		Store:           storex.NewRedisStore(rds),
		Devices:         storex.NewRedisMemberStore(rds),
		LoginMode:       tokenx2.SinglePoint,
		KeyPrefix:       cachekey.TokenStorePrefixAdmin,
		AccessTokenTTL:  2 * time.Hour,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	})
	if err != nil {
		panic(err)
	}

	storageProvider, err := storagex.NewStorageProvider(&c.StorageConfig)
	if err != nil {
		panic(err)
	}

	tracker := stomphook.NewRedisOnlineTracker(rds, cachekey.OnlineAdminKey)

	hub := client.NewStompHubServer(
		client.WithOnlineTracker(tracker),
		client.WithEventHooks(
			stomphook.NewChatRoomEventHook(),
			stomphook.NewOnlineCatchupHook(tracker),
		),
		client.WithAuthenticator(stomphook.NewJwtAuthenticator(tokenManager)),
		client.WithLogger(logws.NewDefaultLogger()),
	)

	// 通知事件扇出：Redis 订阅是进程级一份，各 SSE 连接从广播器分食
	notifyCtx, stopNotify := context.WithCancel(context.Background())
	notifyBroker := sse.NewBroker[notifyx.Event]()
	sse.NewSubscriber(rds, notifyBroker).Start(notifyCtx)

	doc, err := loads.Analyzed(json.RawMessage(docs.Docs), "")
	if err != nil {
		panic(err)
	}

	var options []zrpc.ClientOption
	options = append(options,
		zrpc.WithUnaryClientInterceptor(interceptorx.ClientErrorInterceptor),
	)

	authService := authservice.NewAuthService(zrpc.MustNewClient(c.AppRpcConf, options...))
	chatService := chatservice.NewChatService(zrpc.MustNewClient(c.AppRpcConf, options...))
	userService := userservice.NewUserService(zrpc.MustNewClient(c.AppRpcConf, options...))
	statsService := statsservice.NewStatsService(zrpc.MustNewClient(c.AppRpcConf, options...))
	notificationService := notificationservice.NewNotificationService(zrpc.MustNewClient(c.AppRpcConf, options...))
	syslogService := syslogservice.NewSyslogService(zrpc.MustNewClient(c.AppRpcConf, options...))
	guestService := guestservice.NewGuestService(zrpc.MustNewClient(c.AppRpcConf, options...))
	contentService := contentservice.NewContentService(zrpc.MustNewClient(c.AppRpcConf, options...))
	discussionService := discussionservice.NewDiscussionService(zrpc.MustNewClient(c.AppRpcConf, options...))
	mediaService := mediaservice.NewMediaService(zrpc.MustNewClient(c.AppRpcConf, options...))
	siteService := siteservice.NewSiteService(zrpc.MustNewClient(c.AppRpcConf, options...))
	accessService := accessservice.NewAccessService(zrpc.MustNewClient(c.AppRpcConf, options...))

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

// Shutdown 停止后台推送：先停 Redis 订阅，再关闭各 SSE 连接的事件通道。
// 通道关闭后 handler 的循环随之结束，客户端会看到流正常收尾并按需重连。
//
// 注意：STOMP hub（StompHubServer.Shutdown）在本仓库仍未接线，属既有遗漏，不在此处扩大范围。
func (s *ServiceContext) Shutdown() {
	if s.stopNotify != nil {
		s.stopNotify()
	}
	if s.NotifyBroker != nil {
		s.NotifyBroker.Close()
	}
}

func ConnectRedis(c config.RedisConf) (*redis.Client, error) {
	address := c.Host + ":" + c.Port
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Username: "",
		Password: c.Password, // no password set
		DB:       c.DB,       // use default DB
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("redis 连接失败: %v", err)
	}

	return client, nil
}
