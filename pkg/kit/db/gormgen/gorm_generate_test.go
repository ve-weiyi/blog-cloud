package gormgen

import (
	"fmt"
	"log"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/jsonconv"
)

// GEN 自动生成 GORM 模型结构体文件及使用示例 https://blog.csdn.net/Jeffid/article/details/126898000
const dsn = "root:mysql7914@(127.0.0.1:3306)/blog-veweiyi?charset=utf8mb4&parseTime=True&loc=Local"

var db *gorm.DB

func init() {
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

func TestGenerator(t *testing.T) {

	out := "./blog"
	//path := "./"
	// 生成实例
	// 指定生成代码的具体(相对)目录，默认为：./dao
	// 默认情况下需要使用WithContext之后才可以查询，但可以通过设置gen.WithoutContext避免这个操作
	g := gen.NewGenerator(gen.Config{
		// 相对执行`go run`时的路径, 会自动创建目录
		OutPath:      out + "/dao",
		OutFile:      out + "/dao/dao.go",
		ModelPkgPath: out + "/entity", // 默认情况下会跟随OutPath参数，在同目录下生成model目录
		WithUnitTest: true,
		// WithDefaultQuery 生成默认查询结构体(作为全局变量使用), 即`Q`结构体和其字段(各表模型)
		// WithoutContext 生成没有context调用限制的代码供查询
		// WithQueryInterface 生成interface形式的查询代码(可导出), 如`Where()`方法返回的就是一个可导出的接口类型
		//Mode: gogger.WithoutContext | gogger.WithDefaultQuery | gogger.WithQueryInterface,
		Mode: gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		// 表字段可为 null 值时, 对应结体字段使用指针类型
		FieldNullable: false, // generate pointer when field is nullable
		// 表字段有 默认值时, 对应结体字段使用指针类型
		FieldCoverable: false, // generate pointer when field has default value, to fix problem zero value cannot be assign: https://gorm.io/docs/create.html#Default-Values
		// 模型结构体字段的数字类型的符号表示是否与表字段的一致, `false`指示都用有符号类型
		FieldSignable: true, // detect integer field's unsigned type, adjust generated data type
		// 生成 gorm 标签的字段索引属性
		FieldWithIndexTag: true, // generate with gorm index tag
		// 生成 gorm 标签的字段类型属性
		FieldWithTypeTag: true, // generate with gorm column type tag

	})
	// 设置目标 db
	g.UseDB(db)

	typeInt := "int"
	//typeTime := "sql.NullTime"

	// 自定义字段的数据类型
	// 统一数字类型为int64,兼容protobuf
	dataMap := map[string]func(columnType gorm.ColumnType) (dataType string){
		"tinyint":   func(columnType gorm.ColumnType) (dataType string) { return typeInt },
		"smallint":  func(columnType gorm.ColumnType) (dataType string) { return typeInt },
		"mediumint": func(columnType gorm.ColumnType) (dataType string) { return typeInt },
		"bigint":    func(columnType gorm.ColumnType) (dataType string) { return typeInt },
		"int":       func(columnType gorm.ColumnType) (dataType string) { return typeInt },
		//"time":      func(columnType gorm.ColumnType) (dataType string) { return typeTime },
		//"date":      func(columnType gorm.ColumnType) (dataType string) { return typeTime },
		//"datetime":  func(columnType gorm.ColumnType) (dataType string) { return typeTime },
		//"timestamp": func(columnType gorm.ColumnType) (dataType string) { return typeTime },
		//"json":      func(columnType gorm.ColumnType) (dataType string) { return "datatypes.JSON" },
	}
	// 要先于`ApplyBasic`执行
	g.WithDataTypeMap(dataMap)

	// 自定义模型结体字段的标签
	// 将特定字段名的 json 标签加上`string`属性,即 MarshalJSON 时该字段由数字类型转成字符串类型
	jsonField := gen.FieldJSONTagWithNS(func(columnName string) (tagContent string) {
		//toStringField := "time"
		//if strings.Contains(columnName, toStringField) {
		// return columnName + "\" example:\"2022-11-16T16:00:00.000Z"
		//}
		return columnName
	})

	// 将非默认字段名的字段定义为自动时间戳和软删除字段;
	// 自动时间戳默认字段名为:`updated_at`、`created_at, 表字段数据类型为: INT 或 DATETIME
	autoUpdateTimeField := gen.FieldGORMTag("update_time", func(tag field.GormTag) field.GormTag {
		tag.Set("column", "update_time")
		tag.Set("type", "timestamp")
		tag.Set("autoUpdateTime")
		return tag
	})
	autoCreateTimeField := gen.FieldGORMTag("create_time", func(tag field.GormTag) field.GormTag {
		tag.Set("column", "update_time")
		tag.Set("type", "timestamp")
		tag.Set("autoUpdateTime")
		return tag
	})
	// 软删除默认字段名为:`deleted_at`, 表字段数据类型为: DATETIME
	softDeleteField := gen.FieldType("delete_time", "soft_delete.DeletedAt")
	// 模型自定义选项组
	fieldOpts := []gen.ModelOpt{gen.FieldTrimPrefix("tb_"), jsonField, autoCreateTimeField, autoUpdateTimeField, softDeleteField}

	// 创建模型的结构体,生成文件在 blog 目录; 先创建的结果会被后面创建的覆盖
	// 创建全部模型文件, 并覆盖前面创建的同名模型
	g.ApplyBasic(g.GenerateAllTable(fieldOpts...)...)

	//g.ApplyBasic(g.GenerateModel("article", fieldOpts...))
	g.Execute()
}

func TestMigrate(t *testing.T) {
	err := db.AutoMigrate()
	if err != nil {
		return
	}
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

		columns, err := db.Migrator().ColumnTypes(table)
		if err != nil {
			log.Println("-->", err)
		}
		log.Println("columns-->", jsonconv.AnyToJsonIndent(columns))

		indexes, err := db.Migrator().GetIndexes(table)
		if err != nil {
			log.Println("-->", err)
		}
		log.Println("indexes-->", jsonconv.AnyToJsonIndent(indexes))
	}
}
