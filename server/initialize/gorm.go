package initialize

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/gormlogger"
	"github.com/ve-weiyi/ve-blog-golang/server/config"
)

func ConnectGorm(c config.MysqlConf) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", c.Username, c.Password, c.Host, c.Port, c.Dbname, c.Config)

	var lg logger.Interface
	// 跟随gorm的日志输出格式
	lg = logger.New(
		gormlogger.NewGormWriter(),
		logger.Config{
			SlowThreshold:             500 * time.Millisecond, // 慢 SQL 阈值，超过会提前结束
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false, // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,  // 彩色打印
			ParameterizedQueries:      false, // 使用参数化查询 (true时，会将参数值替换为?)
		})
	// 跟随go-zero的日志输出格式
	//lg = gormlogx.New(
	//	logger.Config{
	//		SlowThreshold:             500 * time.Millisecond, // 慢 SQL 阈值，超过会提前结束
	//		LogLevel:                  logger.Info,
	//		IgnoreRecordNotFoundError: false, // 忽略ErrRecordNotFound（记录未找到）错误
	//		Colorful:                  true,  // 彩色打印
	//		ParameterizedQueries:      false, // 使用参数化查询 (true时，会将参数值替换为?)
	//	},
	//)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//PrepareStmt:            true, // 缓存预编译语句
		// 禁用外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
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
		//Logger: logger.Default,
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
