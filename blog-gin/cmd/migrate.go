package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type MigrateFlags struct {
	initFile string
	dataFile string

	Host     string
	Port     string
	Username string
	Password string
	Dbname   string
	Config   string
}

var migrateFlags = MigrateFlags{
	initFile: "blog-veweiyi-init.sql",
	dataFile: "blog-veweiyi-data.sql",
	Host:     "localhost",
	Port:     "3306",
	Username: "root",
	Password: "123456",
	Dbname:   "blog-veweiyi",
	Config:   "charset=utf8mb4&parseTime=True&loc=Local",
}

func NewMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "初始化数据库表",
		Long:  `初始化数据库表，支持自定义数据库配置和sql文件`,
		RunE:  runMigrate,
	}

	cmd.Flags().StringVarP(&migrateFlags.initFile, "file", "i", migrateFlags.initFile, "数据库结构sql文件")
	cmd.Flags().StringVarP(&migrateFlags.dataFile, "data", "d", migrateFlags.dataFile, "数据库初始数据sql文件")
	cmd.Flags().StringVar(&migrateFlags.Host, "host", migrateFlags.Host, "数据库ip")
	cmd.Flags().StringVar(&migrateFlags.Port, "port", migrateFlags.Port, "数据库端口")
	cmd.Flags().StringVar(&migrateFlags.Username, "username", migrateFlags.Username, "账号")
	cmd.Flags().StringVar(&migrateFlags.Password, "password", migrateFlags.Password, "密码")
	cmd.Flags().StringVar(&migrateFlags.Dbname, "name", migrateFlags.Dbname, "数据库名称")
	cmd.Flags().StringVar(&migrateFlags.Config, "config", migrateFlags.Config, "数据库配置")

	return cmd
}

func runMigrate(cmd *cobra.Command, args []string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", migrateFlags.Username, migrateFlags.Password, migrateFlags.Host, migrateFlags.Port, "", migrateFlags.Config)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//PrepareStmt:            true, // 缓存预编译语句
		// 外键约束
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
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("connected to", dsn)

	// 清空数据库
	if err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", migrateFlags.Dbname)).Error; err != nil {
		return err
	}

	// 创建数据库（如果不存在）
	if err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", migrateFlags.Dbname)).Error; err != nil {
		return err
	}

	if err = db.Exec(fmt.Sprintf("USE `%s`", migrateFlags.Dbname)).Error; err != nil {
		return err
	}

	// 初始化表
	if err = execSqlFile(db, migrateFlags.initFile); err != nil {
		return err
	}

	// 导入数据
	if err = execSqlFile(db, migrateFlags.dataFile); err != nil {
		return err
	}

	log.Println("database initialized successfully")
	return nil
}

func execSqlFile(db *gorm.DB, sqlFile string) error {
	if sqlFile == "" {
		return nil
	}

	// 读取 SQL 文件内容
	content, err := os.ReadFile(sqlFile)
	if err != nil {
		return err
	}

	// 使用strings.NewReader创建一个读取器，并传递给bufio.NewScanner
	scanner := bufio.NewScanner(strings.NewReader(string(content)))

	// 运行语句
	var query string
	// 按行读取内容
	for scanner.Scan() {
		line := scanner.Text()

		// 忽略注释
		if strings.HasPrefix(line, "--") {
			continue
		}

		// 拼接 SQL 语句
		query += "\n" + line
		if strings.HasSuffix(query, ";") {
			query = strings.TrimSpace(query)
			if err = db.Exec(query).Error; err != nil {
				log.Fatal(err)
			}
			query = ""
		} else {
			continue
		}
	}

	// 检查是否有错误发生（例如：文件是否读取完整）
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
