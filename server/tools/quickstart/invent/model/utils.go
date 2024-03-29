package model

import (
	"gorm.io/gorm"
)

// 获取数据库表信息
func GetTable(t *gorm.DB, tableName string) (data *Table, err error) {
	types, err := t.Migrator().TableType(tableName)
	if err != nil {
		return nil, err
	}

	var out Table
	out.TableType = types
	out.TableName = types.Name()
	out.TableComment, _ = types.Comment()
	out.Type = types.Type()
	out.Columns, err = GetTableColumns(t, tableName)

	return &out, nil
}

// 获取表字段信息
func GetTableColumns(t *gorm.DB, tableName string) (data []*Column, err error) {
	var entities []*Column

	var mapIndex map[string][]*Index

	types, err := t.Migrator().ColumnTypes(tableName)
	if err != nil {
		return nil, err
	}

	indexes, err := t.Migrator().GetIndexes(tableName)
	if err != nil {
		return nil, err
	}

	mapIndex = GroupByColumn(indexes)
	for _, entity := range types {
		col := &Column{}
		col.ColumnType = entity

		col.ColumnName = entity.Name()
		col.Indexes = mapIndex[entity.Name()]
		col.ColumnFiledType, _ = entity.ColumnType()
		col.ColumnDefault, col.HasDefault = entity.DefaultValue()
		col.ColumnComment, _ = entity.Comment()
		col.DataType = entity.DatabaseTypeName()
		col.DataTypeLong, _ = entity.Length()
		col.IsNullable, _ = entity.Nullable()
		col.IsPrimaryKey, _ = entity.PrimaryKey()
		col.IsUnique, _ = entity.Unique()
		col.IsAutoIncrement, _ = entity.AutoIncrement()
		entities = append(entities, col)
	}

	return entities, nil
}

// 获取表索引信息
func GetTableIndex(t *gorm.DB, tableName string) (indexes []gorm.Index, err error) {
	return t.Migrator().GetIndexes(tableName)
}

// GroupByColumn group columns
func GroupByColumn(indexList []gorm.Index) map[string][]*Index {
	columnIndexMap := make(map[string][]*Index, len(indexList))
	if len(indexList) == 0 {
		return columnIndexMap
	}

	for _, idx := range indexList {
		if idx == nil {
			continue
		}
		for i, col := range idx.Columns() {
			columnIndexMap[col] = append(columnIndexMap[col], &Index{
				Index:    idx,
				Priority: int32(i + 1),
			})
		}
	}
	return columnIndexMap
}
