package sqlx

import (
	"fmt"
	"strings"
)

// 分页
type PageLimit struct {
	Page     int `json:"page" form:"page"`           // 页码
	PageSize int `json:"page_size" form:"page_size"` // 每页大小
}

func (page *PageLimit) IsValid() bool {
	return page.Page >= 0 && page.PageSize > 0
}

func (page *PageLimit) Limit() int {
	if page.PageSize == 0 {
		page.PageSize = 10
	}
	return page.PageSize
}

func (page *PageLimit) Offset() int {
	return (page.Page - 1) * page.Limit()
}

// 排序语句
type Sort struct {
	Field string `json:"field"`                  // 表字段
	Order string `json:"order" enums:"asc,desc"` // 排序规则 asc|desc
}

// 排序语句
func OrderClause(sorts []*Sort) string {
	if len(sorts) == 0 {
		return ""
	}

	var query string
	var flag string
	for i, order := range sorts {
		if i == 0 {
			flag = ""
		} else {
			flag = ","
		}
		query += fmt.Sprintf("%s `%s` %s", flag, order.Field, order.Order)
	}

	return query
}

// 查询条件
type Condition struct {
	Field string      `json:"field"`                      // 表字段
	Value interface{} `json:"value"`                      // 值
	Rule  string      `json:"rule" enums:"=,like,in,<,>"` // 比较运算符（Comparison Operators）。规则 =,like,in,<,>
	Flag  string      `json:"flag" enums:"and,or"`        // 逻辑运算符（Logical Operators）。标识 and、or,默认and
}

func (condition *Condition) Clause() (string, interface{}) {
	var query string
	var arg interface{}

	switch condition.Rule {
	case "like":
		query += fmt.Sprintf("%s %s ? ", condition.Field, condition.Rule)
		arg = "%" + condition.Value.(string) + "%"
	case "in":
		query += fmt.Sprintf("%s %s (?) ", condition.Field, condition.Rule)
		arg = condition.Value
	default:
		query += fmt.Sprintf("%s %s ? ", condition.Field, condition.Rule)
		arg = condition.Value
	}

	return query, arg
}

// "`id` = ?" , 1
func NewCondition(cond string, arg interface{}) *Condition {
	key := strings.Split(cond, " ")
	if len(key) != 3 {
		return nil
	}
	return &Condition{
		Field: key[0],
		Value: arg,
		Rule:  key[1],
	}
}

// 转换条件语句
func ConditionClause(conditions []*Condition) (string, []interface{}) {
	if len(conditions) == 0 {
		return "", nil
	}

	var query string
	var args []interface{}
	var flag string
	for i, condition := range conditions {
		if i == 0 {
			flag = ""
		} else {
			flag = condition.Flag
			if flag == "" {
				flag = "and"
			}
		}

		switch condition.Rule {
		case "in":
			query += fmt.Sprintf("%s %s %s (?) ", flag, condition.Field, condition.Rule)
			args = append(args, condition.Value)
		case "like":
			query += fmt.Sprintf("%s %s %s ? ", flag, condition.Field, condition.Rule)
			args = append(args, "%"+condition.Value.(string)+"%")
		case "not like":
			query += fmt.Sprintf("%s %s %s ? ", flag, condition.Field, condition.Rule)
			args = append(args, condition.Value.(string)+"%")
		default:
			query += fmt.Sprintf("%s %s %s ? ", flag, condition.Field, condition.Rule)
			args = append(args, condition.Value)
		}
	}

	return query, args
}

// 查询是否含有某条件
func FindCondition(conditions []*Condition, name string) *Condition {
	for _, condition := range conditions {
		if condition.Field == name {
			return condition
		}
	}
	return nil
}
