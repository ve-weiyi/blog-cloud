package svctx

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/orca-zhang/ecache"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/config"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/middleware"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/initialize"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/jwtx"
)

// 注册需要用到的gorm、redis、model
type ServiceContext struct {
	Config *config.Config

	DbEngin    *gorm.DB
	RedisEngin *redis.Client
	LocalCache *ecache.Cache
	Token      *jwtx.JwtInstance

	TerminalToken gin.HandlerFunc
	UserToken     gin.HandlerFunc
	AdminToken    gin.HandlerFunc
	Operation     gin.HandlerFunc
	Permission    gin.HandlerFunc
	OperationLog  gin.HandlerFunc
	VisitLog      gin.HandlerFunc
}

func NewServiceContext(c *config.Config) *ServiceContext {
	db, err := initialize.ConnectGorm(c.Mysql)
	if err != nil {
		panic(err)
	}

	rdb, err := initialize.ConnectRedis(c.Redis)
	if err != nil {
		panic(err)
	}

	cache := ecache.NewLRUCache(16, 200, 10*time.Second).LRU2(1024)

	tk := jwtx.NewJwtInstance([]byte(c.JWT.Secret))

	return &ServiceContext{
		Config:        c,
		DbEngin:       db,
		RedisEngin:    rdb,
		LocalCache:    cache,
		Token:         tk,
		TerminalToken: middleware.TerminalToken(),
		UserToken:     middleware.UserToken(),
		AdminToken:    middleware.AdminToken(tk),
		Operation:     middleware.GinLogger(),
		Permission:    middleware.GinRecovery(true),
		OperationLog:  middleware.GinRecovery(true),
		VisitLog:      middleware.GinRecovery(true),
	}
}
