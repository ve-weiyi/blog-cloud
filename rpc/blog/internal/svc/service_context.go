package svc

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/config"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/adapter/gormx/gormlogx"
	"github.com/ve-weiyi/vkit/adapter/mail"
	"github.com/ve-weiyi/vkit/adapter/oauthx"
	"github.com/ve-weiyi/vkit/adapter/smsx"
	"github.com/ve-weiyi/vkit/adapter/storex/codestore"

	"github.com/ve-weiyi/blog-cloud/infra/dbnotify"
)

type ServiceContext struct {
	Config config.Config

	GormDB *gorm.DB
	Redis  *redis.Client

	EmailDeliver mail.IEmailDeliver
	SmsProvider  smsx.SmsProvider

	// OAuth 服务提供商 map[platform]OAuthProvider
	OAuthProviders map[string]oauthx.OAuthProvider

	CodeStore *codestore.CodeStore

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
	TMessageModel model.TMessageModel

	// website models
	TConfigModel model.TConfigModel
	TAlbumModel  model.TAlbumModel
	TPhotoModel  model.TPhotoModel
	TFriendModel model.TFriendModel
	TTalkModel   model.TTalkModel
	TPageModel   model.TPageModel
	TGuestModel  model.TGuestModel

	// notice models
	TNotifyTemplateModel model.TNotifyTemplateModel
	TNotifyMessageModel  model.TNotifyMessageModel
	TNotifyRecordModel   model.TNotifyRecordModel

	// stats/log models
	TDailyStatsModel   model.TDailyStatsModel
	TVisitLogModel     model.TVisitLogModel
	TLoginLogModel     model.TLoginLogModel
	TOperationLogModel model.TOperationLogModel
	TUploadLogModel    model.TUploadLogModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := ConnectGorm(c.MysqlConf, c.Log)
	if err != nil {
		panic(err)
	}

	// 注意：生产环境不建议使用 AutoMigrate，可能会导致数据丢失。建议使用专门的数据库迁移工具（如 golang-migrate）来管理数据库 schema 变更。
	//if err := AutoMigrate(db); err != nil {
	//	panic(fmt.Errorf("AutoMigrate 失败: %v", err))
	//}

	rds, err := ConnectRedis(c.RedisConf)
	if err != nil {
		panic(err)
	}

	dbnotify.Register(db, rds)

	emailDeliver := NewEmailDeliver(c.EmailConf)
	smsProvider := NewSmsProvider(c.SmsConf)

	return &ServiceContext{
		Config:         c,
		GormDB:         db,
		Redis:          rds,
		EmailDeliver:   emailDeliver,
		SmsProvider:    smsProvider,
		OAuthProviders: NewOAuthProviders(c.AppOAuthConf),
		CodeStore:      nil,
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
		TMessageModel: model.NewTMessageModel(db),
		// website models
		TConfigModel: model.NewTConfigModel(db),
		TAlbumModel:  model.NewTAlbumModel(db),
		TPhotoModel:  model.NewTPhotoModel(db),
		TFriendModel: model.NewTFriendModel(db),
		TTalkModel:   model.NewTTalkModel(db),
		TPageModel:   model.NewTPageModel(db),
		TGuestModel:  model.NewTGuestModel(db),
		// notice models
		TNotifyTemplateModel: model.NewTNotifyTemplateModel(db),
		TNotifyMessageModel:  model.NewTNotifyMessageModel(db),
		TNotifyRecordModel:   model.NewTNotifyRecordModel(db),
		// stats/log models
		TDailyStatsModel:   model.NewTDailyStatsModel(db),
		TVisitLogModel:     model.NewTVisitLogModel(db),
		TLoginLogModel:     model.NewTLoginLogModel(db),
		TOperationLogModel: model.NewTOperationLogModel(db),
		TUploadLogModel:    model.NewTUploadLogModel(db),
	}
}

func ConnectGorm(c config.MysqlConf, l logx.LogConf) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", c.Username, c.Password, c.Host, c.Port, c.Dbname, c.Config)

	var lg logger.Interface
	if l.Mode == "console" && l.Encoding == "plain" {
		lg = logger.New(
			gormlogx.NewGormWriter(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				gormlogx.Config{
					Skip:         2,
					SkipKeywords: []string{"/model", "/gorm"},
				}),
			logger.Config{
				SlowThreshold:             500 * time.Millisecond,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
				ParameterizedQueries:      false,
			})
	} else {
		lg = gormlogx.New(
			logger.Config{
				SlowThreshold:             500 * time.Millisecond,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
				ParameterizedQueries:      false,
			},
		)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
		},
		Logger: lg,
	})
	if err != nil {
		return nil, fmt.Errorf("GORM 数据库连接失败: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("SQL 数据库连接失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(64)
	sqlDB.SetMaxOpenConns(64)
	sqlDB.SetConnMaxLifetime(time.Minute)

	return db, nil
}

func ConnectRedis(c config.RedisConf) (*redis.Client, error) {
	address := c.Host + ":" + c.Port
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: c.Password,
		DB:       c.DB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("redis 连接失败: %v", err)
	}

	return client, nil
}

func NewEmailDeliver(c config.EmailConf) mail.IEmailDeliver {
	ec := &mail.EmailConfig{
		Host:     c.Host,
		Port:     c.Port,
		Username: c.Username,
		Password: c.Password,
		Nickname: c.Nickname,
		BCC:      c.BCC,
	}

	return mail.NewEmailDeliver(ec)
}

func NewSmsProvider(c config.SmsConf) smsx.SmsProvider {
	sc := &smsx.SmsConfig{
		Provider:  c.Provider,
		AccessKey: c.AccessKey,
		SecretKey: c.SecretKey,
		SignName:  c.SignName,
		Region:    c.Region,
		SdkAppId:  c.SdkAppId,
		Templates: c.Templates,
	}

	return smsx.NewSmsProvider(sc)
}

func NewOAuthProviders(confs map[string]config.OAuthConf) map[string]oauthx.OAuthProvider {
	providers := make(map[string]oauthx.OAuthProvider)
	for platform, c := range confs {
		providers[platform] = oauthx.NewOAuthProvider(&oauthx.OAuthConfig{
			Platform:     platform,
			ClientId:     c.ClientId,
			ClientSecret: c.ClientSecret,
			RedirectUri:  c.RedirectUri,
		})
	}
	return providers
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.TUser{},
		&model.TUserOauth{},
		&model.TRole{},
		&model.TApi{},
		&model.TMenu{},
		&model.TUserRole{},
		&model.TRoleApi{},
		&model.TRoleMenu{},
		&model.TArticle{},
		&model.TCategory{},
		&model.TTag{},
		&model.TArticleTag{},
		&model.TChat{},
		&model.TComment{},
		&model.TMessage{},
		&model.TConfig{},
		&model.TAlbum{},
		&model.TPhoto{},
		&model.TFriend{},
		&model.TTalk{},
		&model.TPage{},
		&model.TGuest{},
		&model.TNotifyTemplate{},
		&model.TNotifyMessage{},
		&model.TNotifyRecord{},
		&model.TDailyStats{},
		&model.TVisitLog{},
		&model.TLoginLog{},
		&model.TOperationLog{},
		&model.TUploadLog{},
	)
}
