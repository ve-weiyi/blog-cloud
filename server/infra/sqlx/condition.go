package sqlx

import (
	"strings"

	"gorm.io/gorm"
)

const (
	QueryLike         = "LIKE"
	QueryPrefixLike   = "PREFIX LIKE"
	QueryNotLike      = "NOT LIKE"
	QueryEqual        = "="
	QueryIn           = "IN"
	QueryAND          = "AND"
	BiggerOrQueryAND  = "<="
	SmallerOrQueryAND = ">="
	QueryIsNull       = "IS NULL"
	QueryIsNotNull    = "IS NOT NULL"
)

const (
	OrderAsc  = "asc"
	OrderDesc = "desc"
)

const (
	StateNormal  = -1
	StateDeleted = 1
)

/*
*
查询过滤条件
*/
func DbFilter(db *gorm.DB, query []Condition) *gorm.DB {

	for _, v := range query {

		if checkNull(v.Value) {
			continue
		}

		switch strings.ToUpper(v.Rule) {
		case QueryLike:
			db = db.Where(v.Field+" LIKE ? ", v.Value.(string)+"%")
			break
		case QueryPrefixLike:
			db = db.Where(v.Field+"	LIKE ? ", "%"+v.Value.(string)+"%")
			break
		case QueryNotLike:
			db = db.Where(v.Field+" NOT LIKE ? ", v.Value.(string)+"%")
			break
		case QueryIn:
			db = db.Where(v.Field+" "+v.Rule+"(?)", v.Value)
			break
		case QueryAND:
			db = db.Where(v.Field+" "+v.Rule+" ? ", v.Value)
			break
		case QueryIsNull:
			db = db.Where(v.Field + QueryIsNull)
			break
		case QueryIsNotNull:
			db = db.Where(v.Field + QueryIsNotNull)
			break

		default:
			db = db.Where(v.Field+" "+v.Rule+" ? ", v.Value)
		}
	}

	return db
}

/*
*
过滤查询条件
*/
func checkNull(v interface{}) bool {

	switch v.(type) {
	case string:
		if v == "" {
			return true
		}
	case int:
		if v == 0 {
			return true
		}
	case int64:
		a, _ := v.(int64)
		if a == 0 {
			return true
		}
	case []string:
		if v == nil {
			return true
		}
		value := v.([]string)
		if len(value) == 0 {
			return true
		}
	}

	return false
}
