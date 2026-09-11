package svc

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/common/stomphook"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/config"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/middleware"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/middleware/visitx"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/analyticsservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/articleservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/configservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/discussionservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/guestservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/resourceservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/socialservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/syslogservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/userauthservice"
	"github.com/ve-weiyi/stompws/logws"
	"github.com/ve-weiyi/stompws/server/client"
	"github.com/ve-weiyi/vkit/adapter/storagex"
	"github.com/ve-weiyi/vkit/adapter/storex"
	"github.com/ve-weiyi/vkit/adapter/storex/captchastore"
	"github.com/ve-weiyi/vkit/adapter/storex/tokenstore"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/client/userservice"

	"github.com/ve-weiyi/blog-cloud/infra/constants/cachekey"
	"github.com/ve-weiyi/blog-cloud/infra/interceptorx"
	"github.com/ve-weiyi/blog-cloud/infra/limitx"
	"github.com/ve-weiyi/blog-cloud/infra/middlewarex"
)

type ServiceContext struct {
	Config    config.Config
	UserAuth  rest.Middleware
	RateLimit rest.Middleware
	AgentLog  rest.Middleware
	VisitLog  rest.Middleware

	RedisClient     *redis.Client
	TokenStore      tokenstore.TokenStore
	CaptchaStore    *captchastore.CaptchaStore
	StorageProvider storagex.StorageProvider
	StompHubServer  *client.StompHubServer

	UserAuthService     userauthservice.UserAuthService
	UserService         userservice.UserService
	AnalyticsService    analyticsservice.AnalyticsService
	NotificationService notificationservice.NotificationService
	SyslogService       syslogservice.SyslogService
	GuestService        guestservice.GuestService
	ArticleService      articleservice.ArticleService
	DiscussionService   discussionservice.DiscussionService
	ResourceService     resourceservice.ResourceService
	SocialService       socialservice.SocialService
	ConfigService       configservice.ConfigService
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds, err := ConnectRedis(c.RedisConf)
	if err != nil {
		panic(err)
	}

	tokenStore := tokenstore.NewJwtTokenStore(
		storex.NewRedisStore(rds, storex.WithPrefix(cachekey.TokenStorePrefixApp)),
		c.Name,
		c.Name,
		2*3600,
		7*24*3600,
	)

	captchaStore := captchastore.NewCaptchaStore(
		captchastore.WithStore(
			storex.NewRedisStore(rds, storex.WithPrefix(cachekey.CaptchaStorePrefixApp)),
			15*time.Minute,
		),
	)

	storageProvider := storagex.NewStorageProvider(&c.StorageConfig)

	var options []zrpc.ClientOption
	options = append(options,
		zrpc.WithUnaryClientInterceptor(interceptorx.ClientErrorInterceptor),
	)

	userauthService := userauthservice.NewUserAuthService(zrpc.MustNewClient(c.AppRpcConf, options...))
	userService := userservice.NewUserService(zrpc.MustNewClient(c.AppRpcConf, options...))
	analyticsService := analyticsservice.NewAnalyticsService(zrpc.MustNewClient(c.AppRpcConf, options...))
	notificationService := notificationservice.NewNotificationService(zrpc.MustNewClient(c.AppRpcConf, options...))
	syslogService := syslogservice.NewSyslogService(zrpc.MustNewClient(c.AppRpcConf, options...))
	clientService := guestservice.NewGuestService(zrpc.MustNewClient(c.AppRpcConf, options...))
	articleService := articleservice.NewArticleService(zrpc.MustNewClient(c.AppRpcConf, options...))
	discussionService := discussionservice.NewDiscussionService(zrpc.MustNewClient(c.AppRpcConf, options...))
	resourceService := resourceservice.NewResourceService(zrpc.MustNewClient(c.AppRpcConf, options...))
	socialService := socialservice.NewSocialService(zrpc.MustNewClient(c.AppRpcConf, options...))
	configService := configservice.NewConfigService(zrpc.MustNewClient(c.AppRpcConf, options...))

	tracker := stomphook.NewRedisOnlineTracker(rds, "")

	hub := client.NewStompHubServer(
		client.WithOnlineTracker(tracker),
		client.WithEventHooks(
			stomphook.NewChatRoomEventHook(userService, discussionService),
			stomphook.NewOnlineCatchupHook(tracker),
		),
		client.WithAuthenticator(stomphook.NewSignAuthenticator(tokenStore)),
		client.WithLogger(logws.NewDefaultLogger()),
	)

	return &ServiceContext{
		Config:              c,
		UserAuth:            middleware.NewUserAuthMiddleware(tokenStore).Handle,
		RateLimit:           middlewarex.NewRateLimitMiddleware(limitx.NewPeriodLimit(60, 5, rds, cachekey.RateLimitStrictPrefix)).Handle,
		VisitLog:            middleware.NewVisitLogMiddleware(visitx.NewVisitEnforcer(), syslogService).Handle,
		RedisClient:         rds,
		TokenStore:          tokenStore,
		CaptchaStore:        captchaStore,
		StorageProvider:     storageProvider,
		StompHubServer:      hub,
		UserAuthService:     userauthService,
		UserService:         userService,
		AnalyticsService:    analyticsService,
		NotificationService: notificationService,
		SyslogService:       syslogService,
		GuestService:        clientService,
		ArticleService:      articleService,
		DiscussionService:   discussionService,
		ResourceService:     resourceService,
		SocialService:       socialService,
		ConfigService:       configService,
	}
}

func ConnectRedis(c config.RedisConf) (*redis.Client, error) {
	address := c.Host + ":" + c.Port
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Username: "",
		Password: c.Password, // no password set
		DB:       c.DB,       // use default DB
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("redis 连接失败: %v", err)
	}

	return redisClient, nil
}
