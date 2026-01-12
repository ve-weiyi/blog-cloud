package svc

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/constant"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/infra/gormlogx"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/online"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/config"
	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/captcha"
	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/mail"
	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/mq/rabbitmqx"
)

type ServiceContext struct {
	Config        config.Config
	Gorm          *gorm.DB
	Redis         *redis.Client
	LocalCache    *collection.Cache
	EmailDeliver  mail.IEmailDeliver
	CaptchaHolder *captcha.CaptchaHolder

	OnlineUserService *online.OnlineUserService

	// account models
	TUserModel      model.TUserModel
	TUserOauthModel model.TUserOauthModel
	TRoleModel      model.TRoleModel
	TApiModel       model.TApiModel
	TMenuModel      model.TMenuModel
	TUserRoleModel  model.TUserRoleModel
	TRoleApiModel   model.TRoleApiModel
	TRoleMenuModel  model.TRoleMenuModel

	// blog models
	TArticleModel    model.TArticleModel
	TCategoryModel   model.TCategoryModel
	TTagModel        model.TTagModel
	TArticleTagModel model.TArticleTagModel

	// message models
	TChatModel    model.TChatModel
	TCommentModel model.TCommentModel
	TRemarkModel  model.TRemarkModel

	// website models
	TWebsiteConfigModel   model.TWebsiteConfigModel
	TAlbumModel           model.TAlbumModel
	TPhotoModel           model.TPhotoModel
	TFriendModel          model.TFriendModel
	TTalkModel            model.TTalkModel
	TPageModel            model.TPageModel
	TVisitDailyStatsModel model.TVisitDailyStatsModel
	TVisitorModel         model.TVisitorModel

	TVisitLogModel     model.TVisitLogModel
	TLoginLogModel     model.TLoginLogModel
	TOperationLogModel model.TOperationLogModel
	TFileLogModel      model.TFileLogModel
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
		Config:        c,
		Gorm:          db,
		Redis:         rds,
		LocalCache:    cache,
		EmailDeliver:  deliver,
		CaptchaHolder: captcha.NewCaptchaHolder(captcha.WithRedisStore(rds)),

		OnlineUserService: online.NewOnlineUserService(rds, 3600*24),
		// account models
		TUserModel:      model.NewTUserModel(db),
		TUserOauthModel: model.NewTUserOauthModel(db),
		TRoleModel:      model.NewTRoleModel(db),
		TApiModel:       model.NewTApiModel(db),
		TMenuModel:      model.NewTMenuModel(db),
		TUserRoleModel:  model.NewTUserRoleModel(db),
		TRoleApiModel:   model.NewTRoleApiModel(db),
		TRoleMenuModel:  model.NewTRoleMenuModel(db),
		// blog models
		TArticleModel:    model.NewTArticleModel(db),
		TCategoryModel:   model.NewTCategoryModel(db),
		TTagModel:        model.NewTTagModel(db),
		TArticleTagModel: model.NewTArticleTagModel(db),
		// message models
		TChatModel:    model.NewTChatModel(db),
		TCommentModel: model.NewTCommentModel(db),
		TRemarkModel:  model.NewTRemarkModel(db),
		// website models
		TWebsiteConfigModel:   model.NewTWebsiteConfigModel(db),
		TAlbumModel:           model.NewTAlbumModel(db),
		TPhotoModel:           model.NewTPhotoModel(db),
		TFriendModel:          model.NewTFriendModel(db),
		TTalkModel:            model.NewTTalkModel(db),
		TPageModel:            model.NewTPageModel(db),
		TVisitDailyStatsModel: model.NewTVisitDailyStatsModel(db),
		TVisitorModel:         model.NewTVisitorModel(db),
		TVisitLogModel:        model.NewTVisitLogModel(db),
		TLoginLogModel:        model.NewTLoginLogModel(db),
		TOperationLogModel:    model.NewTOperationLogModel(db),
		TFileLogModel:         model.NewTFileLogModel(db),
	}
}

func ConnectGorm(c config.MysqlConf, l logx.LogConf) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", c.Username, c.Password, c.Host, c.Port, c.Dbname, c.Config)

	var lg logger.Interface
	if l.Mode == "console" && l.Encoding == "plain" {
		// 跟随gorm的日志输出格式
		lg = logger.New(
			gormlogx.NewGormWriter(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				gormlogx.Config{
					Skip: 2,
					// 跳过 model/ 仓库层的行号，使用业务层的行号
					SkipKeywords: []string{"/model", "/gorm"},
				}),
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

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("redis 连接失败: %v", err)
	}

	return client, nil
}

func InitEmailDeliver(c config.Config) (mail.IEmailDeliver, error) {
	e := &mail.EmailConfig{
		Host:     c.EmailConf.Host,
		Port:     c.EmailConf.Port,
		Username: c.EmailConf.Username,
		Password: c.EmailConf.Password,
		Nickname: c.EmailConf.Nickname,
		BCC:      c.EmailConf.BCC,
	}

	// 如果不使用Rabbitmq
	if c.RabbitMQConf.Host == "" {
		return mail.NewEmailDeliver(e), nil
	}

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

	// 注册队列、交换机、绑定关系
	err = conn.Declare(queue, exchange, binding)
	if err != nil {
		log.Fatal(err)
	}

	// pub/sub模式 消息发布者只需要声明交换机
	pb := rabbitmqx.NewRabbitmqProducer(conn,
		rabbitmqx.WithPublisherExchange(constant.EmailExchange),
		rabbitmqx.WithPublisherMandatory(true),
	)

	// pub/sub模式 消息订阅者需要声明交换机和队列
	sb := rabbitmqx.NewRabbitmqConsumer(
		conn,
		rabbitmqx.WithConsumerQueue(constant.EmailQueue),
		rabbitmqx.WithConsumerAutoAck(true),
	)

	// 使用消息队列投递邮件
	deliver := mail.NewMqEmailDeliver(e, pb, sb)
	// 订阅消息
	go deliver.SubscribeEmail()

	return deliver, nil
}
