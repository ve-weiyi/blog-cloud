package helper

import (
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/model/sql/parser"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/convertx"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/files"
)

// 从sql文件解析Table
func ParseTableFromSql(sql string) (list []*Table, err error) {
	n := strings.TrimRight(sql, ".sql")

	f := files.ToAbs(sql)
	tables, err := parser.Parse(f, n, false)
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

// 从数据库中解析Table
func ParseTableFromDsn(dsn string) (list []*Table, err error) {
	db, err := gorm.Open(mysql.Open(dsn), nil)
	if err != nil {
		return nil, err
	}

	dbName := db.Migrator().CurrentDatabase()
	tables, err := db.Migrator().GetTables()
	for _, tableName := range tables {
		types, err := db.Migrator().ColumnTypes(tableName)
		if err != nil {
			return nil, err
		}

		indexes, err := db.Migrator().GetIndexes(tableName)
		if err != nil {
			return nil, err
		}

		pm := Primary{}
		for _, entity := range types {
			is, ok := entity.PrimaryKey()
			if ok && is {
				pm.Field = ConvertColumnToField(entity)
				pm.AutoIncrement, _ = entity.AutoIncrement()
			}
		}

		fs := make([]*Field, 0)
		for _, entity := range types {
			f := ConvertColumnToField(entity)
			fs = append(fs, &f)
		}

		ufs := make(map[string][]*Field)
		for k, index := range GroupByColumn(indexes) {
			uf := make([]*Field, 0)
			for _, field := range index {
				for _, entity := range types {
					if entity.Name() == field {
						f := ConvertColumnToField(entity)
						uf = append(uf, &f)
					}
				}
			}
			ufs[k] = uf
		}

		v := &Table{
			Name:        tableName,
			Db:          dbName,
			PrimaryKey:  pm,
			UniqueIndex: ufs,
			Fields:      fs,
		}

		list = append(list, v)
	}
	return list, nil
}

func ConvertColumnToField(col gorm.ColumnType) Field {
	f := Field{}

	f.Name = col.Name()
	f.Comment, _ = col.Comment()
	f.DataType = convertx.ConvertMysqlToGoType(col.DatabaseTypeName())

	// col.DatabaseTypeName() int
	// col.ColumnType() int unsigned
	return f
}

// group columns
func GroupByColumn(indexList []gorm.Index) map[string][]string {

	ufs := make(map[string][]string)
	if len(indexList) == 0 {
		return ufs
	}

	for _, idx := range indexList {
		if idx == nil {
			continue
		}
		is, ok := idx.PrimaryKey()
		if ok && is {
			continue
		}

		is, ok = idx.Unique()
		if ok && is {
			name := idx.Name()
			for _, col := range idx.Columns() {
				ufs[name] = append(ufs[name], col)
			}
		}
	}
	return ufs
}
