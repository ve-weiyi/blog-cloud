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
	"github.com/ve-weiyi/blog-cloud/infra/middlewarex/limitx"
	"github.com/ve-weiyi/blog-cloud/infra/storex/captchax"
	tokenx2 "github.com/ve-weiyi/blog-cloud/infra/storex/tokenx"
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
	"github.com/ve-weiyi/stompws/logws"
	"github.com/ve-weiyi/stompws/server/client"
	"github.com/ve-weiyi/vkit/adapter/storagex"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/userservice"

	"github.com/ve-weiyi/blog-cloud/infra/constants/cachekey"
	"github.com/ve-weiyi/blog-cloud/infra/interceptorx"
	"github.com/ve-weiyi/blog-cloud/infra/middlewarex"
	"github.com/ve-weiyi/blog-cloud/infra/storex"
)

type ServiceContext struct {
	Config    config.Config
	UserAuth  rest.Middleware
	RateLimit rest.Middleware
	AgentLog  rest.Middleware
	VisitLog  rest.Middleware

	RedisClient     *redis.Client
	TokenManager    tokenx2.Manager
	CaptchaStore    *captchax.Store
	StorageProvider storagex.StorageProvider
	StompHubServer  *client.StompHubServer

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
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds, err := ConnectRedis(c.RedisConf)
	if err != nil {
		panic(err)
	}

	// 图形验证码 15 分钟过期
	captchaStore := captchax.New(storex.NewRedisStore(rds), cachekey.CaptchaStorePrefixApp, 15*time.Minute)

	tokenManager, err := tokenx2.New(tokenx2.Config{
		Signer:          tokenx2.NewJWTSigner([]byte(c.Name), c.Name),
		Store:           storex.NewRedisStore(rds),
		Devices:         storex.NewRedisMemberStore(rds),
		LoginMode:       tokenx2.SinglePoint,
		KeyPrefix:       cachekey.TokenStorePrefixApp,
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

	tracker := stomphook.NewRedisOnlineTracker(rds, "")

	hub := client.NewStompHubServer(
		client.WithOnlineTracker(tracker),
		client.WithEventHooks(
			stomphook.NewChatRoomEventHook(userService, chatService),
			stomphook.NewOnlineCatchupHook(tracker),
		),
		client.WithAuthenticator(stomphook.NewSignAuthenticator(tokenManager)),
		client.WithLogger(logws.NewDefaultLogger()),
	)

	return &ServiceContext{
		Config:              c,
		UserAuth:            middleware.NewUserAuthMiddleware(tokenManager).Handle,
		RateLimit:           middlewarex.NewRateLimitMiddleware(limitx.NewPeriodLimit(60, 5, rds, cachekey.RateLimitStrictPrefix)).Handle,
		VisitLog:            middleware.NewVisitLogMiddleware(visitx.NewVisitEnforcer(), syslogService).Handle,
		RedisClient:         rds,
		TokenManager:        tokenManager,
		CaptchaStore:        captchaStore,
		StorageProvider:     storageProvider,
		StompHubServer:      hub,
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
