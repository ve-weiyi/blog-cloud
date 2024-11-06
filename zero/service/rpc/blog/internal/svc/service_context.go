package svc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/captcha"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/constant"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/gormlogger"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/mail"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/mq/rabbitmqx"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oauth"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oauth/feishu"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oauth/gitee"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oauth/github"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oauth/qq"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oauth/weibo"

	"github.com/ve-weiyi/ve-blog-golang/zero/internal/gormlogx"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/config"
)

type ServiceContext struct {
	Config        config.Config
	Gorm          *gorm.DB
	Redis         *redis.Client
	LocalCache    *collection.Cache
	EmailDeliver  *mail.MqEmailDeliver
	CaptchaHolder *captcha.CaptchaHolder
	Oauth         map[string]oauth.Oauth

	TUserModel             model.TUserModel
	TUserOauthModel        model.TUserOauthModel
	TUserLoginHistoryModel model.TUserLoginHistoryModel
	TRoleModel             model.TRoleModel
	TApiModel              model.TApiModel
	TMenuModel             model.TMenuModel
	TUserRoleModel         model.TUserRoleModel
	TRoleApiModel          model.TRoleApiModel
	TRoleMenuModel         model.TRoleMenuModel

	// blog models
	TWebsiteConfigModel model.TWebsiteConfigModel
	TArticleModel       model.TArticleModel
	TCategoryModel      model.TCategoryModel
	TTagModel           model.TTagModel
	TArticleTagModel    model.TArticleTagModel

	TCommentModel      model.TCommentModel
	TRemarkModel       model.TRemarkModel
	TFriendModel       model.TFriendModel
	TTalkModel         model.TTalkModel
	TPhotoModel        model.TPhotoModel
	TAlbumModel        model.TAlbumModel
	TBannerModel       model.TBannerModel
	TVisitHistoryModel model.TVisitHistoryModel

	TOperationLogModel model.TOperationLogModel
	TChatMessageModel  model.TChatMessageModel

	TFileFolderModel model.TFileFolderModel
	TFileUploadModel model.TFileUploadModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := ConnectGorm(c.MysqlConf, c.Log)
	if err != nil {
		panic(err)
	}

	rds, err := ConnectRedis(c.RedisConf)
	if err != nil {
		panic(err)
	}

	cache, err := collection.NewCache(60 * time.Minute)
	if err != nil {
		panic(err)
	}

	deliver, err := InitEmailDeliver(c)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:                 c,
		Gorm:                   db,
		Redis:                  rds,
		LocalCache:             cache,
		EmailDeliver:           deliver,
		CaptchaHolder:          captcha.NewCaptchaHolder(captcha.WithRedisStore(rds)),
		Oauth:                  InitOauth(c.OauthConfList),
		TUserModel:             model.NewTUserModel(db, rds),
		TUserOauthModel:        model.NewTUserOauthModel(db, rds),
		TUserLoginHistoryModel: model.NewTUserLoginHistoryModel(db, rds),
		TRoleModel:             model.NewTRoleModel(db, rds),
		TApiModel:              model.NewTApiModel(db, rds),
		TMenuModel:             model.NewTMenuModel(db, rds),
		TUserRoleModel:         model.NewTUserRoleModel(db, rds),
		TRoleApiModel:          model.NewTRoleApiModel(db, rds),
		TRoleMenuModel:         model.NewTRoleMenuModel(db, rds),

		// blog models
		TWebsiteConfigModel: model.NewTWebsiteConfigModel(db, rds),
		TArticleModel:       model.NewTArticleModel(db, rds),
		TCategoryModel:      model.NewTCategoryModel(db, rds),
		TTagModel:           model.NewTTagModel(db, rds),
		TArticleTagModel:    model.NewTArticleTagModel(db, rds),

		TCommentModel:      model.NewTCommentModel(db, rds),
		TRemarkModel:       model.NewTRemarkModel(db, rds),
		TFriendModel:       model.NewTFriendModel(db, rds),
		TTalkModel:         model.NewTTalkModel(db, rds),
		TPhotoModel:        model.NewTPhotoModel(db, rds),
		TAlbumModel:        model.NewTAlbumModel(db, rds),
		TBannerModel:       model.NewTBannerModel(db, rds),
		TVisitHistoryModel: model.NewTVisitHistoryModel(db, rds),

		TOperationLogModel: model.NewTOperationLogModel(db, rds),
		TChatMessageModel:  model.NewTChatMessageModel(db, rds),

		TFileFolderModel: model.NewTFileFolderModel(db, rds),
		TFileUploadModel: model.NewTFileUploadModel(db, rds),
	}
}

func ConnectGorm(c config.MysqlConf, l logx.LogConf) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", c.Username, c.Password, c.Host, c.Port, c.Dbname, c.Config)

	var lg logger.Interface
	if l.Mode == "console" && l.Encoding == "plain" {
		// 跟随gorm的日志输出格式
		lg = logger.New(
			gormlogger.NewGormWriter(gormlogger.AddSkip(1)),
			logger.Config{
				SlowThreshold:             500 * time.Millisecond, // 慢 SQL 阈值，超过会提前结束
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: false, // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  true,  // 彩色打印
				ParameterizedQueries:      false, // 使用参数化查询 (true时，会将参数值替换为?)
			})
	} else {
		// 跟随go-zero的日志输出格式
		lg = gormlogx.New(
			logger.Config{
				SlowThreshold:             500 * time.Millisecond, // 慢 SQL 阈值，超过会提前结束
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: false, // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  true,  // 彩色打印
				ParameterizedQueries:      false, // 使用参数化查询 (true时，会将参数值替换为?)
			},
		)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// PrepareStmt:            true, // 缓存预编译语句
		// 外键约束
		// DisableForeignKeyConstraintWhenMigrating: true,
		// 禁用默认事务（提高运行速度）
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			// 表前缀
			TablePrefix: "",
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
			SingularTable: true,
		},
		// gorm日志模式
		Logger: lg,
		// Logger: logger.Default,
	})

	if err != nil {
		return nil, fmt.Errorf("GORM 数据库连接失败: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("SQL 数据库连接失败: %v", err)
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(64)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(64)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Minute)

	return db, nil
}

func ConnectRedis(c config.RedisConf) (*redis.Client, error) {
	address := c.Host + ":" + c.Port
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Username: "",
		Password: c.Password, // no password set
		DB:       c.DB,       // use default DB
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("redis 连接失败: %v", err)
	}

	client.Set(context.Background(), fmt.Sprintf("redis:rpc:%s", pong), time.Now().String(), -1)
	return client, nil
}

func InitEmailDeliver(c config.Config) (*mail.MqEmailDeliver, error) {
	e := c.EmailConf
	emailSender := mail.NewEmailDeliver(
		mail.WithHost(e.Host),
		mail.WithPort(e.Port),
		mail.WithUsername(e.Username),
		mail.WithPassword(e.Password),
		mail.WithNickname(e.Nickname),
		mail.WithDeliver(e.Deliver),
	)

	r := c.RabbitMQConf
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", r.Username, r.Password, r.Host, r.Port)
	// 创建连接
	conn, err := rabbitmqx.NewRabbitmqConn(url, nil)
	if err != nil {
		log.Fatal("rabbitmq 初始化失败!", err)
	}

	queue := &rabbitmqx.QueueOptions{
		Name:    constant.EmailQueue,
		Durable: true, // 是否持久化
	}

	exchange := &rabbitmqx.ExchangeOptions{
		Name:    constant.EmailExchange,
		Kind:    rabbitmqx.ExchangeTypeFanout,
		Durable: true, // 是否持久化
	}

	binding := &rabbitmqx.BindingOptions{
		RoutingKey: "",
	}

	err = conn.Declare(queue, exchange, binding)
	if err != nil {
		log.Fatal(err)
	}

	// pub/sub模式 消息发布者只需要声明交换机
	pb := rabbitmqx.NewRabbitmqProducer(conn,
		rabbitmqx.WithPublisherExchange(exchange.Name),
		rabbitmqx.WithPublisherMandatory(true),
	)

	// pub/sub模式 消息订阅者需要声明交换机和队列
	sb := rabbitmqx.NewRabbitmqConsumer(
		conn,
		rabbitmqx.WithConsumerQueue(queue.Name),
		rabbitmqx.WithConsumerAutoAck(true),
	)

	deliver := mail.NewMqEmailDeliver(emailSender, pb, sb)
	// 订阅消息
	go deliver.SubscribeEmail()

	return deliver, nil
}

func InitOauth(c map[string]config.OauthConf) map[string]oauth.Oauth {
	var om = make(map[string]oauth.Oauth)

	for k, v := range c {
		conf := &oauth.AuthConfig{
			ClientId:     v.ClientId,
			ClientSecret: v.ClientSecret,
			RedirectUri:  v.RedirectUri,
		}
		switch k {
		case "qq":
			auth := qq.NewAuthQq(conf)
			om["qq"] = auth
		case "weibo":
			auth := weibo.NewAuthWb(conf)
			om["weibo"] = auth
		case "feishu":
			auth := feishu.NewAuthFeishu(conf)
			om["feishu"] = auth
		case "github":
			auth := github.NewAuthGithub(conf)
			om["github"] = auth
		case "gitee":
			auth := gitee.NewAuthGitee(conf)
			om["gitee"] = auth
		}
	}
	return om
}
