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

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/config"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/likes"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq/outbox"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/notification"
	"github.com/ve-weiyi/blog-cloud/model"

	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/store/gormlogx"
	"github.com/ve-weiyi/blog-cloud/infra/store/otpx"
	"github.com/ve-weiyi/blog-cloud/infra/store/storex"
	"github.com/ve-weiyi/vkit/adapter/emailx"
	"github.com/ve-weiyi/vkit/adapter/oauthx"
	"github.com/ve-weiyi/vkit/adapter/smsx"

	"github.com/ve-weiyi/blog-cloud/infra/messaging/dbnotify"
	"github.com/ve-weiyi/blog-cloud/infra/messaging/notifyx"
)

type ServiceContext struct {
	Config config.Config

	GormDB *gorm.DB
	Redis  *redis.Client
	MQ     *mq.Client
	Outbox *outbox.GORMStore

	EmailDeliver emailx.IEmailDeliver
	SmsProvider  smsx.SmsProvider

	// OAuth 服务提供商 map[platform]OAuthProvider
	OAuthProviders map[string]oauthx.OAuthProvider

	OTPStore *otpx.Store

	// Notify 通知域服务：生产事件、扇出、读写收件箱
	Notify *notification.Service

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
	TLikeModel       model.TLikeModel

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
	// notification models
	TNotificationEventModel      model.TNotificationEventModel
	TNotificationModel           model.TNotificationModel
	TNotificationDeliveryModel   model.TNotificationDeliveryModel
	TNotificationPreferenceModel model.TNotificationPreferenceModel
	TAnnouncementModel           model.TAnnouncementModel

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

	emailDeliver, err := NewEmailDeliver(c.EmailConf)
	if err != nil {
		panic(err)
	}

	smsProvider, err := NewSmsProvider(c.SmsConf)
	if err != nil {
		panic(err)
	}

	oAuthProviders, err := NewOAuthProviders(c.AppOAuthConf)
	if err != nil {
		panic(err)
	}

	mqClient, err := mq.NewClient(c.RabbitMQConf)
	if err != nil {
		panic(err)
	}

	svc := &ServiceContext{
		Config:         c,
		GormDB:         db,
		Redis:          rds,
		MQ:             mqClient,
		Outbox:         outbox.NewGORMStore(db),
		EmailDeliver:   emailDeliver,
		SmsProvider:    smsProvider,
		OAuthProviders: oAuthProviders,
		OTPStore:       otpx.New(storex.NewRedisStore(rds)),
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
		TLikeModel:       model.NewTLikeModel(db),
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
		// notification models
		TNotificationEventModel:      model.NewTNotificationEventModel(db),
		TNotificationModel:           model.NewTNotificationModel(db),
		TNotificationDeliveryModel:   model.NewTNotificationDeliveryModel(db),
		TNotificationPreferenceModel: model.NewTNotificationPreferenceModel(db),
		TAnnouncementModel:           model.NewTAnnouncementModel(db),
		// stats/log models
		TDailyStatsModel:   model.NewTDailyStatsModel(db),
		TVisitLogModel:     model.NewTVisitLogModel(db),
		TLoginLogModel:     model.NewTLoginLogModel(db),
		TOperationLogModel: model.NewTOperationLogModel(db),
		TUploadLogModel:    model.NewTUploadLogModel(db),
	}

	svc.Notify = notification.NewService(notification.Deps{
		DB:                   db,
		Events:               svc.TNotificationEventModel,
		Notifications:        svc.TNotificationModel,
		Deliveries:           svc.TNotificationDeliveryModel,
		Preferences:          svc.TNotificationPreferenceModel,
		Announcements:        svc.TAnnouncementModel,
		Users:                svc.TUserModel,
		Comments:             svc.TCommentModel,
		Outbox:               svc.Outbox,
		Notifier:             notifySignaler{rdb: rds},
		AnnouncementNotifier: notifySignaler{rdb: rds},
		UnsubscribeSecret:    c.NotifyConf.UnsubscribeSecret,
		SiteURL:              c.NotifyConf.SiteURL,
	})
	return svc
}

// LikeDeps 返回点赞模块的依赖组合。
func (s *ServiceContext) LikeDeps() likes.Deps {
	return likes.Deps{DB: s.GormDB, Likes: s.TLikeModel}
}

// Close 按依赖顺序释放资源：先停 MQ（等待在途消费完成），再关 Redis 与数据库。
// 幂等：重复调用安全。
func (s *ServiceContext) Close() {
	if s == nil {
		return
	}
	if s.MQ != nil {
		s.MQ.Close()
	}
	if s.Redis != nil {
		_ = s.Redis.Close()
	}
	if s.GormDB != nil {
		if db, err := s.GormDB.DB(); err == nil {
			_ = db.Close()
		}
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
				ParameterizedQueries:      true,
			})
	} else {
		lg = gormlogx.New(
			logger.Config{
				SlowThreshold:             500 * time.Millisecond,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
				ParameterizedQueries:      true,
			},
		)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		TranslateError:         true,
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

func NewEmailDeliver(c config.EmailConf) (emailx.IEmailDeliver, error) {
	if c.Port == 465 && !c.IsSsl {
		return nil, fmt.Errorf("smtp port 465 requires is_ssl=true")
	}

	return emailx.NewEmailDeliver(toEmailConfig(c)), nil
}

func toEmailConfig(c config.EmailConf) *emailx.EmailConfig {
	return &emailx.EmailConfig{
		Host:     c.Host,
		Port:     c.Port,
		Username: c.Username,
		Password: c.Password,
		Nickname: c.Nickname,
		SSL:      c.IsSsl,
		BCC:      c.BCC,
	}
}

func NewSmsProvider(c config.SmsConf) (smsx.SmsProvider, error) {
	if c.Provider != "mock" {
		for _, scene := range []string{enums.CodeSceneMobileLogin, enums.CodeSceneResetPassword, enums.CodeSceneBind} {
			if c.Templates[scene] == "" {
				return nil, fmt.Errorf("sms template missing for scene %s", scene)
			}
		}
	}

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

func NewOAuthProviders(confs map[string]config.OAuthConf) (map[string]oauthx.OAuthProvider, error) {
	providers := make(map[string]oauthx.OAuthProvider)
	for platform, c := range confs {
		p, err := oauthx.NewOAuthProvider(&oauthx.OAuthConfig{
			Platform:     platform,
			ClientId:     c.ClientId,
			ClientSecret: c.ClientSecret,
			RedirectUri:  c.RedirectUri,
		})
		if err != nil {
			return nil, err
		}
		providers[platform] = p
	}
	return providers, nil
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
		&model.TNotificationEvent{},
		&model.TNotification{},
		&model.TNotificationDelivery{},
		&model.TNotificationPreference{},
		&model.TAnnouncement{},
		&model.TDailyStats{},
		&model.TVisitLog{},
		&model.TLoginLog{},
		&model.TOperationLog{},
		&model.TUploadLog{},
	)
}

// notifySignaler 把扇出写入的收件人变成实时信号。
//
// 按收件人逐条发布：信号只到达持有该用户的实例，不向全站在线连接广播。
// 发布失败只记日志——推送是加速器，丢了由客户端定时对账兜底。
type notifySignaler struct {
	rdb *redis.Client
}

func (n notifySignaler) NotifyChanged(ctx context.Context, recipients []string) {
	for _, uid := range recipients {
		if err := notifyx.PublishToUser(ctx, n.rdb, uid); err != nil {
			logx.Errorf("通知信号发送失败: user=%s err=%v", uid, err)
		}
	}
}

// NotifyAnnouncement 广播公告状态变化，管理端各实例订阅后刷新列表。
func (n notifySignaler) NotifyAnnouncement(ctx context.Context, published bool, id int64, title, category, level string, at time.Time) {
	eventType := notifyx.EventAnnouncementPublished
	if !published {
		eventType = notifyx.EventAnnouncementRevoked
	}

	if err := notifyx.Publish(ctx, n.rdb, notifyx.Event{
		Type:           eventType,
		AnnouncementId: id,
		Title:          title,
		Category:       category,
		Level:          level,
		PublishedAt:    at.Unix(),
	}); err != nil {
		logx.Errorf("公告状态广播失败: id=%d published=%v err=%v", id, published, err)
	}
}

// notifyEmailSender 把邮件通道适配到通知投递。
type notifyEmailSender struct {
	deliver emailx.IEmailDeliver
}

func (n notifyEmailSender) SendEmail(_ context.Context, to, subject, content string) error {
	if n.deliver == nil {
		return fmt.Errorf("notification: email channel is not configured")
	}
	return n.deliver.DeliveryEmail(&emailx.EmailMessage{
		To:      []string{to},
		CC:      []string{},
		Subject: subject,
		Content: content,
	})
}

// NotifyEmailSender 返回通知投递使用的邮件通道。
func (s *ServiceContext) NotifyEmailSender() notification.EmailSender {
	return notifyEmailSender{deliver: s.EmailDeliver}
}
