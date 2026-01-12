package gormdriver

type GormDriver interface {
	GetSchemas() (data []Db, err error)
	GetTables(dbName string) (data []Table, err error)
	GetTableColumns(dbName string, tableName string) (data []Column, err error)
}

// 数据库信息
type Db struct {
	SchemaName string `json:"schemaName" gorm:"column:SCHEMA_NAME"`
}

// 表信息
type Table struct {
	TableName    string `json:"tableName" gorm:"column:TABLE_NAME"`
	TableComment string `json:"tableComment" gorm:"column:TABLE_COMMENT"`
}

// 列信息
type Column struct {
	TableCatalog           string `gorm:"column:TABLE_CATALOG" json:"tableCatalog"`                      // 列所属的数据库名称。
	TableSchema            string `gorm:"column:TABLE_SCHEMA" json:"tableSchema"`                        // 列所属的模式名称。
	TableName              string `gorm:"column:TABLE_NAME" json:"tableName"`                            // 列所属的表名称。
	OrdinalPosition        int    `gorm:"column:ORDINAL_POSITION" json:"ordinalPosition"`                // 列在表中的位置。
	ColumnName             string `gorm:"column:COLUMN_NAME" json:"columnName"`                          // 列的名称。
	ColumnType             string `gorm:"column:COLUMN_TYPE" json:"columnType"`                          // 列的类型和长度。
	ColumnDefault          string `gorm:"column:COLUMN_DEFAULT" json:"columnDefault"`                    // 列的默认值。
	ColumnComment          string `gorm:"column:COLUMN_COMMENT" json:"columnComment"`                    // 列的注释信息。
	ColumnKey              string `gorm:"column:COLUMN_KEY" json:"columnKey"`                            // 列是否为主键或唯一键的一部分。
	IsNullable             string `gorm:"column:IS_NULLABLE" json:"isNullable"`                          // 列是否允许为空。
	DataType               string `gorm:"column:DATA_TYPE" json:"dataType"`                              // 列的数据类型。
	CharacterMaximumLength int    `gorm:"column:CHARACTER_MAXIMUM_LENGTH" json:"characterMaximumLength"` // 字符类型列的最大长度。
	NumericPrecision       int    `gorm:"column:NUMERIC_PRECISION" json:"numericPrecision"`              // 数值类型列的精度。
	NumericScale           int    `gorm:"column:NUMERIC_SCALE" json:"numericScale"`                      // 数值类型列的小数位数。
	DatetimePrecision      int    `gorm:"column:DATETIME_PRECISION" json:"datetimePrecision"`            // 日期时间类型列的精度。
	CharacterSetName       string `gorm:"column:CHARACTER_SET_NAME" json:"characterSetName"`             // 字符类型列的字符集名称。
	CollationName          string `gorm:"column:COLLATION_NAME" json:"collationName"`                    // 字符类型列的排序规则名称。
	Extra                  string `gorm:"column:EXTRA" json:"extra"`                                     // 列是否具有附加属性，如自动递增。
	Privileges             string `gorm:"column:PRIVILEGES" json:"privileges"`                           // 与列相关的权限信息。
}
