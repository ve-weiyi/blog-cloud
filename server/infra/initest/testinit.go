package initest

import (
	"context"
	"fmt"
	"log"
	"path"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/ve-weiyi/ve-blog-golang/server/global"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/glog"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/glog/zaplog"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/jjwt"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/copyutil"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/files"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

// @title						Swagger Example API
// @version					0.0.1
// @description				This is a sample Server pets
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						x-token
// @BasePath					/
func Init(configPath ...string) {
	InitConfig(configPath...)
	// 初始化zap日志库
	Zap()
	// 初始化gorm数据库
	Gorm()
	// 初始化redis服务
	Redis()
	// 初始化jwt服务
	JwtToken()
}

func InitConfig(configPath ...string) {
	log.Println("let's go")
	var filepath string
	if len(configPath) > 1 {
		filepath = configPath[0]
	} else {
		filepath = path.Join(global.GetRuntimeRoot() + "server/config.yaml")
	}
	global.VP = Viper(filepath) // 初始化Viper
}

func Viper(config string) *viper.Viper {

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		log.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&global.CONFIG); err != nil {
			log.Println(err)
		}
	})
	if err = v.Unmarshal(&global.CONFIG); err != nil {
		log.Println(err)
	}

	// root 适配性 根据root位置去找到对应迁移位置,保证root路径有效
	//global.CONFIG.AutoCode.Root = global.GetRuntimeRoot()
	return v
}

func Zap() {
	err := files.MkDir(global.CONFIG.Zap.CacheDir)
	if err != nil {
		log.Println(err)
	}
	cfg := zaplog.ZapConfig{}

	copyutil.DeepCopyByJson(global.CONFIG.Zap, &cfg)

	//glog.ReplaceDefaultLogger(cfg)
	global.LOG = glog.NewGlogger(1, cfg)

	global.LOG.Infof("日志组件初始化成功！")
	return
}

func Gorm() {

	cfg := &global.CONFIG.Mysql
	dsn := cfg.Dsn()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//PrepareStmt:            true, // 缓存预编译语句
		// 外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 禁用默认事务（提高运行速度）
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			// 表前缀
			TablePrefix: cfg.Prefix,
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
			SingularTable: true,
		},
		// gorm日志模式
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("GORM 数据库连接失败: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("SQL 数据库连接失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	global.DB = db

	log.Printf("Mysql 数据库连接成功！%s", cfg.Dsn())
}

func Redis() {
	redisCfg := global.CONFIG.Redis
	address := redisCfg.Host + ":" + redisCfg.Port
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DB,       // use default DB
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("Redis 连接失败, err:%v", err)
		return
	}
	client.Set(context.Background(), "connect", time.Now().String(), -1)
	global.REDIS = client

	log.Printf("Redis 连接成功%v! address:%v db:%v", pong, address, redisCfg.DB)
}

func JwtToken() {
	global.JWT = &jjwt.JwtToken{
		SigningKey:  []byte(global.CONFIG.JWT.SigningKey),
		TokenPrefix: "",
		Issuer:      "blog",
	}
}
