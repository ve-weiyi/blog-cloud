package quickstart

import (
	"gorm.io/gorm"

	"github.com/ve-weiyi/ve-blog-golang/server/tools/quickstart/invent/model"
)

type Config struct {
	DbEngin *gorm.DB

	model.FieldConfig

	ReplaceMode int //是否替换文件 0:创建或替换 1:只创建 2:只替换

	OutPath   string                                   // 输出路径
	OutFileNS func(tableName string) (fileName string) // 输出文件名称

	FieldNameNS  func(columnName string) (fieldName string)
	FieldJsonNS  func(columnName string) (jsonName string)
	FieldValueNS func(columnName string) (valueName string)

	IsIgnoreKey func(key string) bool
}
