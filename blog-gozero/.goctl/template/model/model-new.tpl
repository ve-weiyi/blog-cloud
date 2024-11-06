func new{{.upperStartCamelObject}}Model(db *gorm.DB{{if .withCache}}, cache *redis.Client{{end}}) *default{{.upperStartCamelObject}}Model {
	return &default{{.upperStartCamelObject}}Model{
	    DbEngin:    db,
	    {{if .withCache}}CacheEngin: cache,{{end}}
		table: {{.table}},
	}
}

