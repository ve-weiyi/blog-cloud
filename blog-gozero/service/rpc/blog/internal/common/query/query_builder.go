package query

import "strings"

// QueryBuilder 查询条件构建器
type QueryBuilder struct {
	page       int
	size       int
	sorts      string
	conditions []string
	params     []any
}

// Option 查询选项函数
type Option func(*QueryBuilder)

// WithPage 设置页码
func WithPage(page int) Option {
	return func(qb *QueryBuilder) {
		qb.page = page
	}
}

// WithSize 设置每页大小
func WithSize(size int) Option {
	return func(qb *QueryBuilder) {
		qb.size = size
	}
}

// WithSorts 设置排序
func WithSorts(sorts ...string) Option {
	return func(qb *QueryBuilder) {
		if len(sorts) > 0 {
			qb.sorts = strings.Join(sorts, ",")
		}
		if qb.sorts == "" {
			qb.sorts = "id desc"
		}
	}
}

// WithCondition 添加查询条件
func WithCondition(condition string, value any) Option {
	return func(qb *QueryBuilder) {
		qb.conditions = append(qb.conditions, condition)
		qb.params = append(qb.params, value)
	}
}

// NewQueryBuilder 创建查询构建器
func NewQueryBuilder(opts ...Option) *QueryBuilder {
	qb := &QueryBuilder{
		page:       1,
		size:       10,
		sorts:      "id desc",
		conditions: make([]string, 0),
		params:     make([]any, 0),
	}

	for _, opt := range opts {
		opt(qb)
	}

	return qb
}

// Build 构建查询条件和参数
func (qb *QueryBuilder) Build() (page int, size int, sorts string, conditions string, params []any) {
	return qb.page, qb.size, qb.sorts, strings.Join(qb.conditions, " and "), qb.params
}
