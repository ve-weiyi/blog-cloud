package svc

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/chat"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/config"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/middleware"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/middleware/visitx"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/aiservice"
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
	"github.com/ve-weiyi/blog-cloud/infra/ai"
	"github.com/ve-weiyi/blog-cloud/infra/auth/tokenx"
	"github.com/ve-weiyi/blog-cloud/infra/auth/tokenx/codec/opaque"
	"github.com/ve-weiyi/blog-cloud/infra/auth/tokenx/store/redisstore"
	"github.com/ve-weiyi/blog-cloud/infra/domain/cachekey"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/interceptorx"
	"github.com/ve-weiyi/blog-cloud/infra/http/middlewarex"
	"github.com/ve-weiyi/blog-cloud/infra/http/middlewarex/limitx"
	"github.com/ve-weiyi/blog-cloud/infra/messaging/stompwsx"
	"github.com/ve-weiyi/blog-cloud/infra/store/captchax"
	"github.com/ve-weiyi/blog-cloud/infra/store/presencex"
	"github.com/ve-weiyi/blog-cloud/infra/store/storex"
	"github.com/ve-weiyi/stompws/logws"
	"github.com/ve-weiyi/stompws/server"
	"github.com/ve-weiyi/vkit/adapter/storagex"
)

type ServiceContext struct {
	Config            config.Config
	UserAuth          rest.Middleware
	RateLimit         rest.Middleware
	MusicRateLimit    rest.Middleware
	VisitLog          rest.Middleware
	NotifyStreamLimit rest.Middleware
	AiSwitch          rest.Middleware
	OptionalUserAuth  rest.Middleware
	AiRateLimit       rest.Middleware

	AiClient        *ai.Client
	AiService       aiservice.AiService
	RedisClient     *redis.Client
	TokenManager    tokenx.Service
	CaptchaStore    *captchax.Store
	StorageProvider storagex.StorageProvider
	StompHubServer  *server.StompHubServer

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

	sessionStore, err := redisstore.New(rds, cachekey.TokenStorePrefixApp)
	if err != nil {
		panic(err)
	}
	codec := opaque.New()

	tokenManager, err := tokenx.New(tokenx.Config{
		AccessCodec:        codec,
		RefreshCodec:       codec,
		Store:              sessionStore,
		LoginMode:          tokenx.MultiPoint,
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
	aiService := aiservice.NewAiService(rpcClient)
	discussionService := discussionservice.NewDiscussionService(rpcClient)
	mediaService := mediaservice.NewMediaService(rpcClient)
	siteService := siteservice.NewSiteService(rpcClient)

	// 模型客户端惰性可用：没配端点时服务照常启动，请求期再报明确错误
	aiClient := newAiClient(c.AiConf)

	// 博客前台的在线名册：只有本服务（前台的 hub）往里写。
	presence := presencex.NewStore(storex.NewRedisMemberStore(rds), cachekey.OnlineUserKey, presencex.DefaultTTL)

	hub := server.NewStompHubServer(
		server.WithPresenceManager(presence),
		server.WithPresenceTTL(presence.TTL()),
		server.WithEventHooks(
			chat.NewChatRoomEventHook(userService, chatService),
			stompwsx.NewOnlineBroadcastHook(),
			stompwsx.NewOnlineCatchupHook(),
		),
		server.WithAuthenticator(stompwsx.NewAuthenticator(tokenManager, true)),
		server.WithLogger(logws.NewDefaultLogger()),
	)

	return &ServiceContext{
		Config:              c,
		UserAuth:            middleware.NewUserAuthMiddleware(tokenManager).Handle,
		RateLimit:           middlewarex.NewRateLimitMiddleware(limitx.NewPeriodLimit(60, 5, rds, cachekey.RateLimitStrictPrefix)).Handle,
		MusicRateLimit:      middlewarex.NewRateLimitMiddleware(limitx.NewPeriodLimit(60, 60, rds, cachekey.RateLimitMusicPrefix)).Handle,
		VisitLog:            middleware.NewVisitLogMiddleware(visitx.NewVisitEnforcer(), syslogService).Handle,
		NotifyStreamLimit:   middleware.NewNotifyStreamLimitMiddleware(middleware.DefaultStreamGlobalLimit, middleware.DefaultStreamUserLimit).Handle,
		AiClient:            aiClient,
		AiService:           aiService,
		AiSwitch:            middleware.NewAiSwitchMiddleware(siteService, c.AiConf.SwitchTTL).Handle,
		OptionalUserAuth:    middleware.NewOptionalUserAuthMiddleware(tokenManager).Handle,
		AiRateLimit:         middleware.NewAiRateLimitMiddleware(rds, c.AiConf.Quota).Handle,
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

// Shutdown 停止长连接：等待 STOMP 连接处理完毕，超时后返回。
// 必须注册为关闭监听，defer 在信号路径上不会执行。
func (s *ServiceContext) Shutdown() {
	if s.StompHubServer != nil {
		if err := s.StompHubServer.Shutdown(5 * time.Second); err != nil {
			logx.Errorf("stop stomp hub: %v", err)
		}
	}
}

// newAiClient 构造模型客户端；没配端点时返回 nil，由请求路径给出明确错误。
func newAiClient(conf config.AiConf) *ai.Client {
	if conf.BaseURL == "" {
		return nil
	}
	apiKey := conf.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("AI_API_KEY")
	}
	if apiKey == "" {
		return nil
	}
	return ai.NewClient(ai.Config{
		BaseURL:     conf.BaseURL,
		APIKey:      apiKey,
		Model:       conf.Model,
		Temperature: conf.Temperature,
		MaxTokens:   conf.MaxTokens,
		Timeout:     conf.Timeout,
	})
}

func ConnectRedis(c config.RedisConf) (*redis.Client, error) {
	address := c.Host + ":" + c.Port
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Username: "",
		Password: c.Password,
		DB:       c.DB,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("redis 连接失败: %v", err)
	}

	return redisClient, nil
}
