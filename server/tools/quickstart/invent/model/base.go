package model

import (
	"strings"
)

const (
	// DefaultModelPkg ...
	DefaultModelPkg = "model"
)

// Status sql status
type Status int

const (
	// UNKNOWN ...
	UNKNOWN Status = iota
	// SQL ...
	SQL
	// DATA ...
	DATA
	// VARIABLE ...
	VARIABLE
	// IF ...
	IF
	// ELSE ...
	ELSE
	// WHERE ...
	WHERE
	// SET ...
	SET
	// FOR ...
	FOR
	// END ...
	END
	// TRIM ...
	TRIM
)

// SourceCode source code
type SourceCode int

const (
	// Struct ...
	Struct SourceCode = iota
	// Table ...
	Table_
	// Object ...
	Object
)

// GormKeywords ...
var GormKeywords = KeyWord{
	words: []string{
		"UnderlyingDB", "UseDB", "UseModel", "UseTable", "Quote", "Debug", "TableName", "WithContext",
		"As", "Not", "Or", "Build", "Columns", "Hints",
		"Distinct", "Omit",
		"Select", "Where", "OrderKey", "Group", "Having", "Limit", "Offset",
		"Join", "LeftJoin", "RightJoin",
		"Save", "Create", "CreateInBatches",
		"Update", "Updates", "UpdateColumn", "UpdateColumns",
		"Find", "FindInBatches", "First", "Take", "Last", "Pluck", "Count",
		"Scan", "ScanRows", "Row", "Rows",
		"Delete", "Unscoped",
		"Scopes",
	},
}

// DOKeywords ...
var DOKeywords = KeyWord{
	words: []string{
		"Alias", "TableName", "WithContext",
	},
}

// GenKeywords ...
var GenKeywords = KeyWord{
	words: []string{
		"generateSQL", "whereClause", "setClause",
	},
}

// KeyWord ...
type KeyWord struct {
	words []string
}

// FullMatch full match
func (g *KeyWord) FullMatch(word string) bool {
	for _, item := range g.words {
		if word == item {
			return true
		}
	}
	return false
}

// Contain contain
func (g *KeyWord) Contain(text string) bool {
	for _, item := range g.words {
		if strings.Contains(text, item) {
			return true
		}
	}
	return false
}

var (
	defaultDataType             = "string"
	defaultIntType              = "int"
	dataType        dataTypeMap = map[string]dataTypeMapping{
		"numeric":    func(string) string { return defaultIntType },
		"integer":    func(string) string { return defaultIntType },
		"int":        func(string) string { return defaultIntType },
		"smallint":   func(string) string { return defaultIntType },
		"mediumint":  func(string) string { return defaultIntType },
		"bigint":     func(string) string { return "int64" },
		"float":      func(string) string { return "float32" },
		"real":       func(string) string { return "float64" },
		"double":     func(string) string { return "float64" },
		"decimal":    func(string) string { return "float64" },
		"char":       func(string) string { return "string" },
		"varchar":    func(string) string { return "string" },
		"tinytext":   func(string) string { return "string" },
		"mediumtext": func(string) string { return "string" },
		"longtext":   func(string) string { return "string" },
		"binary":     func(string) string { return "[]byte" },
		"varbinary":  func(string) string { return "[]byte" },
		"tinyblob":   func(string) string { return "[]byte" },
		"blob":       func(string) string { return "[]byte" },
		"mediumblob": func(string) string { return "[]byte" },
		"longblob":   func(string) string { return "[]byte" },
		"text":       func(string) string { return "string" },
		"json":       func(string) string { return "string" },
		"enum":       func(string) string { return "string" },
		"time":       func(string) string { return "time.Time" },
		"date":       func(string) string { return "time.Time" },
		"datetime":   func(string) string { return "time.Time" },
		"timestamp":  func(string) string { return "time.Time" },
		"year":       func(string) string { return defaultIntType },
		"bit":        func(string) string { return "[]uint8" },
		"boolean":    func(string) string { return "bool" },
		"tinyint": func(detailType string) string {
			if strings.HasPrefix(strings.TrimSpace(detailType), "tinyint(1)") {
				return "bool"
			}
			return defaultIntType
		},
	}
)

type dataTypeMapping func(detailType string) (finalType string)

type dataTypeMap map[string]dataTypeMapping

func (m dataTypeMap) Get(dataType, detailType string) string {
	//先判断varchar，再判断 varchar(11)
	if convert, ok := m[strings.ToLower(dataType)]; ok {
		return convert(detailType)
	}
	return defaultDataType
}
