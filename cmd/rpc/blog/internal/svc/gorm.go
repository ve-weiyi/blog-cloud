package svc

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/config"
	"github.com/ve-weiyi/blog-cloud/infra/store/gormlogx"
)

// ConnectGorm 建立 MySQL 连接并配置连接池。
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
		return nil, fmt.Errorf("initialize: connect mysql: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("initialize: get sql db: %w", err)
	}

	sqlDB.SetMaxIdleConns(64)
	sqlDB.SetMaxOpenConns(64)
	sqlDB.SetConnMaxLifetime(time.Minute)

	return db, nil
}
