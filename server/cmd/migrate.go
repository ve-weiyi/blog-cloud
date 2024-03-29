/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
	"github.com/ve-weiyi/ve-blog-golang/server/config/properties"
	"github.com/ve-weiyi/ve-blog-golang/server/global"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/initest"
)

// migrateCmd represents the migrate command
type MigrateCmd struct {
	cmd      *cobra.Command
	Mysql    properties.Mysql
	createDB bool
	action   string
	sqlFile  string
}

func NewMigrateCmd() *MigrateCmd {
	migrateCmd := &MigrateCmd{}
	migrateCmd.cmd = &cobra.Command{
		Use:   "migrate",
		Short: "合并数据库表",
		Long:  `合并数据库表，更新到最新版本`,
		Run: func(cmd *cobra.Command, args []string) {
			migrateCmd.MigrateDB()
		},
	}
	migrateCmd.init()
	return migrateCmd
}

func (s *MigrateCmd) init() {
	s.cmd.PersistentFlags().StringVarP(&s.action, "action", "a", "migrate", "migrate|create|reset")
	s.cmd.PersistentFlags().StringVarP(&s.sqlFile, "file", "", "blog-veweiyi.sql", "数据库sql文件")

	s.cmd.PersistentFlags().StringVarP(&s.Mysql.Username, "username", "", "root", "账号")
	s.cmd.PersistentFlags().StringVarP(&s.Mysql.Password, "password", "", "123456", "密码")
	s.cmd.PersistentFlags().StringVarP(&s.Mysql.Host, "host", "", "localhost", "数据库ip")
	s.cmd.PersistentFlags().StringVarP(&s.Mysql.Port, "port", "", "3306", "数据库端口")
	s.cmd.PersistentFlags().StringVarP(&s.Mysql.Dbname, "name", "", "blog", "数据库名称")
	s.cmd.PersistentFlags().StringVarP(&s.Mysql.Config, "config", "", "charset=utf8mb4&parseTime=True&loc=Local", "数据库配置")
}

func (s *MigrateCmd) MigrateDB() {
	//m := s.Mysql
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", m.Username, m.Password, m.Host, m.Port, "", m.Config)
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("connect to ", dsn)

	initest.Init()
	db := global.DB
	switch s.action {
	case "create":
		s.InitDatabase(db, "ve-weiyi")
	case "reset":
		s.ResetDatabase(db)
	case "migrate":
		s.MigrateDatabase(db)
	default:
		fmt.Println("action not support")
	}
}

// 迁移数据库
func (s *MigrateCmd) MigrateDatabase(db *gorm.DB) {
	// 读取 SQL 文件内容
	content, err := os.ReadFile(s.sqlFile)
	if err != nil {
		log.Fatal(err)
	}

	// 将文件内容转为字符串
	queries := string(content)

	// 分割 SQL 命令
	queryList := strings.Split(queries, ";")

	// 执行每个 SQL 命令
	for _, query := range queryList {
		query = strings.TrimSpace(query)
		if query != "" {
			db.Exec(query)
		}
	}

	fmt.Println("Database initialized successfully.")
}

// 创建数据库
func (s *MigrateCmd) InitDatabase(db *gorm.DB, dbName string) {
	// 创建数据库（如果不存在）
	result := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	// 关闭连接
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.Close()
}

// 重置数据库
func (s *MigrateCmd) ResetDatabase(db *gorm.DB) {
	var err error
	// ****重置登录历史表****
	err = ClearTable(db, entity.TableNameUserLoginHistory)
	if err != nil {
		log.Fatal(err)
	}

	// ****重置角色菜单表****
	err = ClearTable(db, entity.TableNameRoleMenu)
	if err != nil {
		log.Fatal(err)
	}

	// ****重置角色接口表****
	err = ClearTable(db, entity.TableNameRoleApi)
	if err != nil {
		log.Fatal(err)
	}

	// ****重置菜单表****
	err = ClearTable(db, entity.TableNameMenu)
	if err != nil {
		log.Fatal(err)
	}

	// ****重置接口表****
	err = ClearTable(db, entity.TableNameApi)
	if err != nil {
		log.Fatal(err)
	}

	// ****重置操作记录表****
	err = ClearTable(db, entity.TableNameOperationLog)
	if err != nil {
		log.Fatal(err)
	}

	// ****重置上传记录表****
	err = ClearTable(db, entity.TableNameUploadRecord)
	if err != nil {
		log.Fatal(err)
	}

	// ****重置casbin记录表****
	err = ClearTable(db, entity.TableNameCasbinRule)
	if err != nil {
		log.Fatal(err)
	}
}

func ClearTable(db *gorm.DB, tableName string) (err error) {
	// 清空表的数据
	err = db.Exec(fmt.Sprintf("DELETE FROM `%v`", tableName)).Error
	if err != nil {
		return err
	}
	// 重置 AUTO_INCREMENT 值为 1
	err = db.Exec(fmt.Sprintf("ALTER TABLE `%v` AUTO_INCREMENT = 1", tableName)).Error
	if err != nil {
		return err
	}
	return nil
}
