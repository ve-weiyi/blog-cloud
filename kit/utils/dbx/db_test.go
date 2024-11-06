package dbx

import (
	"fmt"
	"log"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
)

// GEN 自动生成 GORM 模型结构体文件及使用示例 https://blog.csdn.net/Jeffid/article/details/126898000
const dsn = "root:mysql7914@(veweiyi.cn:3306)/blog-veweiyi?charset=utf8mb4&parseTime=True&loc=Local"

var db *gorm.DB

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	var err error
	// 连接数据库
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix: "tb_",
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
			SingularTable: true,
		},
	})
	if err != nil {
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}
	log.Println("mysql connection done")
}

func TestDBService(t *testing.T) {
	var driver DBDriver
	driver = &MysqlDriver{DB: db}
	//dbList, err := driver.GetDB()
	//if err != nil {
	//	log.Println("-->", err)
	//	return
	//}
	//log.Println("dbList-->", jsonconv.AnyToJsonIndent(dbList))
	//
	//tabelList, err := driver.GetTables("blog-veweiyi")
	//if err != nil {
	//	log.Println("-->", err)
	//	return
	//}
	//log.Println("tabelList-->", jsonconv.AnyToJsonIndent(tabelList))

	columnList, err := driver.GetTableColumns("blog-veweiyi", "user_account")
	if err != nil {
		log.Println("-->", err)
		return
	}
	log.Println("columnList-->", db.Migrator().CurrentDatabase(), jsonconv.AnyToJsonIndent(columnList))
}

func TestMigrator(t *testing.T) {
	dbName := db.Migrator().CurrentDatabase()
	log.Println("dbName-->", dbName)

	tables, err := db.Migrator().GetTables()
	if err != nil {
		log.Println("-->", err)
	}

	for _, table := range tables {
		tableType, err := db.Migrator().TableType(table)
		if err != nil {
			log.Println("-->", err)
		}
		log.Println("table-->", jsonconv.AnyToJsonIndent(tableType))

		//columns, err := db.Migrator().ColumnTypes(table)
		//if err != nil {
		//	log.Println("-->", err)
		//}
		//log.Println("columns-->", columns)

		//indexes, err := db.Migrator().GetIndexes(table)
		//if err != nil {
		//	log.Println("-->", err)
		//}
		//log.Println("indexes-->", jsonconv.AnyToJsonIndent(indexes))
	}

}
