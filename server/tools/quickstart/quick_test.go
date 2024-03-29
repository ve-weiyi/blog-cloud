package quickstart

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/ve-weiyi/ve-blog-golang/server/global"
	"github.com/ve-weiyi/ve-blog-golang/server/tools/quickstart/invent"
	"github.com/ve-weiyi/ve-blog-golang/server/tools/quickstart/invent/model"
	"github.com/ve-weiyi/ve-blog-golang/server/tools/quickstart/tmpl"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/jsonconv"
)

// GEN 自动生成 GORM 模型结构体文件及使用示例 https://blog.csdn.net/Jeffid/article/details/126898000
const dsn = "root:mysql7914@(veweiyi.cn:3306)/blog-veweiyi?charset=utf8mb4&parseTime=True&loc=Local"

var db *gorm.DB

func Init() {
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

func TestCodeStarter(t *testing.T) {
	Init()
	out := path.Join(global.GetRuntimeRoot(), "server/api")
	//out := path.Join("./autocode_template", "test")

	typeInt := "int"
	// 自定义字段的数据类型
	// 统一数字类型为int64,兼容protobuf
	dataMap := map[string]func(columnType gorm.ColumnType) (dataType string){
		"tinyint":   func(columnType gorm.ColumnType) (dataType string) { return typeInt },
		"smallint":  func(columnType gorm.ColumnType) (dataType string) { return typeInt },
		"mediumint": func(columnType gorm.ColumnType) (dataType string) { return typeInt },
		"bigint":    func(columnType gorm.ColumnType) (dataType string) { return typeInt },
		"int":       func(columnType gorm.ColumnType) (dataType string) { return typeInt },
		//"datetime":  func(columnType gorm.ColumnType) (dataType string) { return "*time.Time" },
	}

	cfg := Config{
		DbEngin:     db,
		ReplaceMode: invent.ModeCreateOrReplace,
		OutPath:     out,
		OutFileNS: func(tableName string) (fileName string) {
			return fmt.Sprintf("bs_%v", tableName)
		},
		FieldNameNS: func(column string) string {
			return jsonconv.Case2Camel(column)
		},
		FieldJsonNS: func(column string) string {
			return jsonconv.Camel2Case(column)
		},
		FieldValueNS: func(columnName string) (valueName string) {
			if columnName == "id" {
				return "id"
			}
			return jsonconv.Case2CamelLowerStart(columnName)
		},
		IsIgnoreKey: func(key string) bool {
			return key != tmpl.KeyModel
		},
		FieldConfig: model.FieldConfig{
			DataTypeMap: dataMap,
			// 表字段可为 null 值时, 对应结体字段使用指针类型
			FieldNullable: false, // generate pointer when field is nullable
			// 表字段有 默认值时, 对应结体字段使用指针类型
			FieldCoverable: false, // generate pointer when field has default value, to fix problem zero value cannot be assign: https://gorm.io/docs/create.html#Default-Values
			// 模型结构体字段的数字类型的符号表示是否与表字段的一致, `false`指示都用有符号类型
			FieldSignable: false, // detect integer field's unsigned type, adjust generated data type
			// 生成 gorm 标签的字段索引属性
			FieldWithIndexTag: true, // generate with gorm index tag
			// 生成 gorm 标签的字段类型属性
			FieldWithTypeTag: true, // generate with gorm column type tag
			FieldJSONTagNS: func(column string) string {
				return jsonconv.Camel2Case(column)
			},
		},
	}

	parser := NewTableParser(cfg)
	gen := NewCodeStarter(cfg)

	gen.AddInventMetas(parser.GenerateInventMetas(parser.ParseModelFromTable("menu"))...)
	//models, _ := parser.ParseModelFromSchema()
	//gen.AddInventMetas(parser.GenerateInventMetas(models...)...)

	err := gen.Execute()
	t.Log(err)
	//gen.InitPackage("hello")
	//gen.ApplyMetas(gen.GenerateMetasFromSchema())

	//gen.ApplyMetas(gen.GenerateMetasFromTable("api", "接口"))
	//gen.ApplyMetas(gen.GenerateMetasFromTable("article", "文章"))
	//gen.ApplyMetas(gen.GenerateMetasFromTable("category", "文章分类"))
	//gen.ApplyMetas(gen.GenerateMetasFromTable("friend_link", "友链"))
	//gen.ApplyMetas(gen.GenerateMetasFromTable("menu", "菜单"))
	//gen.ApplyMetas(gen.GenerateMetasFromTable("operation_log", "操作记录"))
	//gen.ApplyMetas(gen.GenerateMetasFromTable("page", "页面"))
	//gen.ApplyMetas(gen.GenerateMetasFromTable("photo", "相片"))
	//gen.ApplyMetas(gen.GenerateMetasFromTable("photo_album", "相册"))
	//gen.ApplyMetas(gen.GenerateMetasFromTable("remark", "留言"))
	//gen.ApplyMetas(gen.GenerateMetasFromTable("role", "角色"))
	//gen.ApplyMetas(gen.GenerateMetasFromTable("tag", "文章标签"))
	//gen.ApplyMetas(gen.GenerateMetasFromTable("talk", "说说"))
	//gen.ApplyMetas(gen.GenerateMetasFromTable("chat_record", "聊天记录"))

	// 不能覆盖的

	//gen.ApplyMetas(gen.GenerateMetasFromTable("comment", "评论"))
	//gen.ApplyMetas(gen.GenerateMetasFromTable("user_account", "用户账号信息"))
	//gen.ApplyMetas(gen.GenerateMetasFromTable("user_information", "用户信息"))
	//gen.ApplyMetas(gen.GenerateMetasFromTable("user_login_history", "用户登录历史"))
	//gen.ApplyMetas(gen.GenerateMetasFromTable("upload_record", "文件上传"))
	//gen.ApplyMetas(gen.GenerateMetasFromTable("unique_view", "页面访问数量"))
	//gen.ApplyMetas(gen.GenerateMetasFromTable("website_config", "网站设置"))
	//gen.GenerateCommonFile("upload", "文件上传")

	//gen.RollBack()
	//gen.Execute()
}

func TestVisitFile(t *testing.T) {
	root := path.Join(global.GetRuntimeRoot(), "server/api", "model/entity")
	err := filepath.Walk(root, visitFile)
	t.Log(err)
}

// bs_ -> base 基础文件
// ex_ -> extend 扩展文件
// sp_ -> special 特殊文件
func visitFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !info.IsDir() {
		fmt.Println("File:", path)
		// 获取原始文件名
		oldName := info.Name()

		//if strings.HasPrefix(oldName, "gen_") {
		//	return nil
		//}

		// 添加前缀 "gen_" 到文件名
		newName := strings.Replace(oldName, "ex_", "ex_", 1)

		// 修改文件名
		err := os.Rename(path, filepath.Join(filepath.Dir(path), newName))
		if err != nil {
			fmt.Println("Error renaming file:", err)
		} else {
			fmt.Println("Renamed file:", newName)
		}
	}

	return nil
}
