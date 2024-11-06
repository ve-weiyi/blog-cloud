package helper

type (
	// Table describes a mysql table
	Table struct {
		Name        string
		Db          string
		PrimaryKey  Primary
		UniqueIndex map[string][]*Field
		Fields      []*Field
	}

	// Primary describes a primary key
	Primary struct {
		Field
		AutoIncrement bool
	}

	// Field describes a table field
	Field struct {
		// NameOriginal    string
		Name string
		// ThirdPkg        string
		DataType        string
		Comment         string
		SeqInIndex      int
		OrdinalPosition int
		// ContainsPQ      bool
	}

	// KeyType types alias of int
	KeyType int
)
