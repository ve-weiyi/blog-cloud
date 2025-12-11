package dbparser

import (
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/model/sql/parser"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/filex"
)

// 从sql文件解析Table
func ParseTableFromSql(sql string) (list []*Table, err error) {
	n := strings.TrimRight(sql, ".sql")

	f := filex.ToAbs(sql)
	tables, err := parser.Parse(f, n, true)
	if err != nil {
		return nil, err
	}

	for _, table := range tables {

		fs := make([]*Field, 0)
		for _, field := range table.Fields {
			f := ConvertFieldToField(field)
			fs = append(fs, &f)
		}

		ufs := make(map[string][]*Field)
		for k, index := range table.UniqueIndex {
			uf := make([]*Field, 0)
			for _, field := range index {
				f := ConvertFieldToField(field)
				uf = append(uf, &f)
			}
			ufs[k] = uf
		}

		v := &Table{
			Name: table.Name.Source(),
			Db:   table.Db.Source(),
			PrimaryKey: Primary{
				AutoIncrement: table.PrimaryKey.AutoIncrement,
				Field:         ConvertFieldToField(&table.PrimaryKey.Field),
			},
			UniqueIndex: ufs,
			Fields:      fs,
		}

		list = append(list, v)

	}

	return list, nil
}

func ConvertFieldToField(col *parser.Field) Field {
	f := Field{
		Name:            col.Name.Source(),
		DataType:        col.DataType,
		Comment:         col.Comment,
		SeqInIndex:      col.SeqInIndex,
		OrdinalPosition: col.OrdinalPosition,
	}

	return f
}
