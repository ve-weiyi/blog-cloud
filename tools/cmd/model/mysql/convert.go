package mysql

import (
	"fmt"
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
	"github.com/ve-weiyi/ve-blog-golang/tools/parserx/dbparser"
)

func ConvertTableToData(table *dbparser.Table) *ModelData {

	var fs []*ModelField
	for _, e := range table.Fields {
		fs = append(fs, ConvertField(e))
	}

	var ufs [][]*ModelField
	for _, es := range table.UniqueIndex {
		var u []*ModelField
		for _, e := range es {
			u = append(u, ConvertField(e))
		}
		ufs = append(ufs, u)
	}

	data := &ModelData{
		TableName:           table.Name,
		UpperStartCamelName: jsonconv.Case2Camel(table.Name),
		LowerStartCamelName: jsonconv.FirstLower(jsonconv.Case2Camel(table.Name)),
		SnakeName:           jsonconv.Case2Snake(table.Name),
		Fields:              fs,
		UniqueFields:        ufs,
	}

	return data
}

func ConvertField(f *dbparser.Field) *ModelField {
	return &ModelField{
		Name:    jsonconv.Case2Camel(f.Name),
		Type:    strings.TrimPrefix(f.DataType, "u"),
		Tag:     fmt.Sprintf(`json:"%v" gorm:"column:%v"`, f.Name, f.Name),
		Comment: f.Comment,
	}
}
