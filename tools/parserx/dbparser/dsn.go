package dbparser

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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
	f.DataType = ConvertMysqlToGoType(col.DatabaseTypeName())

	nullable, ok := col.Nullable()
	if ok && nullable && f.DataType == "time.Time" {
		f.DataType = "sql.NullTime"
	}

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
