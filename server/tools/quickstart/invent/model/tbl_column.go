package model

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"

	"github.com/ve-weiyi/ve-blog-golang/server/tools/quickstart/invent/field"
)

// 获取数据库表信息
type Schema struct {
	SchemaName string   `json:"schema_name"` // 数据库名
	Tables     []*Table `json:"tables"`      // 表
}

// 表
type Table struct {
	gorm.TableType

	SchemaName   string `json:"schema_name"`   // 数据库名
	TableName    string `json:"table_name"`    // 表名
	TableComment string `json:"table_comment"` // 表注释
	Type         string `json:"type"`          // 表类型

	Indexes []*Index  `json:"indexes"` // 表索引
	Columns []*Column `json:"columns"` // 表字段
}

// 索引
type Index struct {
	gorm.Index
	Priority int32 `gorm:"column:SEQ_IN_INDEX"`
}

// 需要的数据
type Column struct {
	gorm.ColumnType
	Indexes []*Index `json:"indexes"` // 字段索引

	ColumnName      string `json:"column_name"`       //列名
	ColumnFiledType string `json:"column_filed_type"` //字段类型 varchar(11)
	ColumnDefault   string `json:"column_default"`    //默认值
	ColumnComment   string `json:"column_comment"`    //备注
	DataType        string `json:"data_type"`         //数据类型 varchar
	DataTypeLong    int64  `json:"data_type_long"`    //数据长度

	IsNullable      bool `json:"is_nullable"`       //是否可空
	IsPrimaryKey    bool `json:"is_primary_key"`    //数据是否是主键
	IsUnique        bool `json:"is_unique"`         //数据是否是唯一的
	IsAutoIncrement bool `json:"is_auto_increment"` //数据是否是自动递增
	HasDefault      bool `json:"has_default"`       //数据是否有默认值
}

// ToField convert to field
func (c *Column) FiledType(cfg *FieldConfig) string {
	fieldType := c.getDataType(cfg)
	if cfg.FieldSignable && strings.Contains(c.columnType(), "unsigned") && strings.HasPrefix(fieldType, "int") {
		fieldType = "u" + fieldType
	}
	switch {
	case c.Name() == "deleted_at" && fieldType == "time.Time":
		fieldType = "gorm.DeletedAt"
	case cfg.FieldCoverable && c.needDefaultTag(c.defaultTagValue()):
		fieldType = "*" + fieldType
	case cfg.FieldNullable:
		if n, ok := c.Nullable(); ok && n {
			fieldType = "*" + fieldType
		}
	}

	return fieldType
}

func (c *Column) BuildGormTag() field.GormTag {
	tag := field.GormTag{
		field.TagKeyGormColumn: []string{c.Name()},
		field.TagKeyGormType:   []string{c.columnType()},
	}
	isPriKey, ok := c.PrimaryKey()
	isValidPriKey := ok && isPriKey
	if isValidPriKey {
		tag.Set(field.TagKeyGormPrimaryKey, "")
		if at, ok := c.AutoIncrement(); ok {
			tag.Set(field.TagKeyGormAutoIncrement, fmt.Sprintf("%t", at))
		}
	} else if n, ok := c.Nullable(); ok && !n {
		tag.Set(field.TagKeyGormNotNull, "")
	}

	for _, idx := range c.Indexes {
		if idx == nil {
			continue
		}
		if pk, _ := idx.PrimaryKey(); pk { //ignore PrimaryKey
			continue
		}
		if uniq, _ := idx.Unique(); uniq {
			tag.Append(field.TagKeyGormUniqueIndex, fmt.Sprintf("%s,priority:%d", idx.Name(), idx.Priority))
		} else {
			tag.Append(field.TagKeyGormIndex, fmt.Sprintf("%s,priority:%d", idx.Name(), idx.Priority))
		}
	}

	if dtValue := c.defaultTagValue(); c.needDefaultTag(dtValue) { // cannot set default tag for primary key
		tag.Set(field.TagKeyGormDefault, dtValue)
	}
	if comment, ok := c.Comment(); ok && comment != "" {
		if c.multilineComment() {
			comment = strings.ReplaceAll(comment, "\n", "\\n")
		}
		tag.Set(field.TagKeyGormComment, comment)
	}
	return tag
}

// GetDataType get data type
func (c *Column) getDataType(cfg *FieldConfig) (fieldtype string) {
	if mapping, ok := cfg.DataTypeMap[c.DatabaseTypeName()]; ok {
		return mapping(c.ColumnType)
	}
	//if c.UseScanType && c.ScanType() != nil {
	//	return c.ScanType().String()
	//}
	return dataType.Get(c.DatabaseTypeName(), c.columnType())
}

func (c *Column) multilineComment() bool {
	cm, ok := c.Comment()
	return ok && strings.Contains(cm, "\n")
}

// needDefaultTag check if default tag needed
func (c *Column) needDefaultTag(defaultTagValue string) bool {
	if defaultTagValue == "" {
		return false
	}
	switch c.ScanType().Kind() {
	case reflect.Bool:
		return defaultTagValue != "false"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return defaultTagValue != "0"
	case reflect.String:
		return defaultTagValue != ""
	case reflect.Struct:
		return strings.Trim(defaultTagValue, "'0:- ") != ""
	}
	return c.Name() != "created_at" && c.Name() != "updated_at"
}

// defaultTagValue return gorm default tag's value
func (c *Column) defaultTagValue() string {
	value, ok := c.DefaultValue()
	if !ok {
		return ""
	}
	if value != "" && strings.TrimSpace(value) == "" {
		return "'" + value + "'"
	}
	return value
}

func (c *Column) columnType() (v string) {
	if cl, ok := c.ColumnType.ColumnType(); ok {
		return cl
	}
	return c.DatabaseTypeName()
}
